package provider

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/payment"
)

// Path constants mirror the Ciabra Invoice API (https://docs.ciabra.com.br).
// All paths are relative to the configured apiBase (default https://api.az.center).
const (
	ciabraDefaultAPIBase        = "https://api.az.center"
	ciabraCustomersPath         = "/invoices/applications/customers"
	ciabraInvoicesPath          = "/invoices/applications/invoices"
	ciabraAuthCheckPath         = "/auth/applications/check"
	ciabraDefaultPaymentTypePIX = "PIX"
	ciabraDefaultInvoiceType    = "SINGLE"
	ciabraHTTPTimeout           = 15 * time.Second
	ciabraMaxResponseSize       = 1 << 20
	ciabraMaxErrorSummary       = 512
	ciabraDefaultMinAmountCents = 500 // R$5.00 — conservative default; provider may have stricter rules.
)

// Ciabra implements payment.Provider backed by the Ciabra Invoice API.
//
// Required config keys:
//   - publicKey, secretKey
//
// Optional config keys:
//   - apiBase (default https://api.az.center)
//   - paymentTypes (comma-separated, default PIX)
//   - dueDays (default 1)
//   - minAmountCents (default 500)
//   - webhookSecret  — when set, VerifyNotification enforces HMAC-SHA256
//   - webhookSignatureHeader (default x-signature)
//   - webhookTimestampHeader (default x-timestamp)
//   - webhookToleranceSeconds (default 300)
//   - webhookPublicURL  — public URL where Ciabra should POST PAYMENT_* events
//                        (typically https://<your-host>/api/v1/payment/webhook/ciabra).
//                        When empty, no per-invoice webhooks are registered; the
//                        admin is responsible for configuring them globally in
//                        the Ciabra dashboard.
//   - defaultDocument — fallback CPF/CNPJ used when the request doesn't carry one
type Ciabra struct {
	instanceID string
	config     map[string]string
	httpClient *http.Client
}

func NewCiabra(instanceID string, config map[string]string) (*Ciabra, error) {
	for _, key := range []string{"publicKey", "secretKey"} {
		if strings.TrimSpace(config[key]) == "" {
			return nil, fmt.Errorf("ciabra config missing required key: %s", key)
		}
	}

	cfg := cloneStringMap(config)
	apiBase, err := normalizeCiabraAPIBase(cfg["apiBase"])
	if err != nil {
		return nil, err
	}
	cfg["apiBase"] = apiBase

	if strings.TrimSpace(cfg["paymentTypes"]) == "" {
		cfg["paymentTypes"] = ciabraDefaultPaymentTypePIX
	}
	if strings.TrimSpace(cfg["dueDays"]) == "" {
		cfg["dueDays"] = "1"
	}
	if strings.TrimSpace(cfg["webhookSignatureHeader"]) == "" {
		cfg["webhookSignatureHeader"] = "x-signature"
	}
	if strings.TrimSpace(cfg["webhookTimestampHeader"]) == "" {
		cfg["webhookTimestampHeader"] = "x-timestamp"
	}
	if strings.TrimSpace(cfg["webhookToleranceSeconds"]) == "" {
		cfg["webhookToleranceSeconds"] = "300"
	}

	return &Ciabra{
		instanceID: instanceID,
		config:     cfg,
		httpClient: &http.Client{Timeout: ciabraHTTPTimeout},
	}, nil
}

func (c *Ciabra) Name() string        { return "Ciabra" }
func (c *Ciabra) ProviderKey() string { return payment.TypeCiabra }
func (c *Ciabra) SupportedTypes() []payment.PaymentType {
	return []payment.PaymentType{payment.TypeCiabra}
}

func (c *Ciabra) MerchantIdentityMetadata() map[string]string {
	if c == nil {
		return nil
	}
	metadata := map[string]string{}
	if origin := strings.TrimSpace(c.config["apiBase"]); origin != "" {
		metadata["api_base"] = origin
	}
	return metadata
}

// normalizeCiabraAPIBase accepts any HTTPS URL and strips any trailing
// `/invoices` legacy suffix that earlier versions of this provider required.
// Paths used by the implementation are always absolute (e.g.
// `/invoices/applications/customers`) so the base must NOT carry a path.
func normalizeCiabraAPIBase(raw string) (string, error) {
	base := strings.TrimSpace(raw)
	if base == "" {
		base = ciabraDefaultAPIBase
	}
	parsed, err := url.Parse(base)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("ciabra apiBase must be a valid URL")
	}
	if parsed.Scheme != "https" {
		return "", fmt.Errorf("ciabra apiBase must use HTTPS")
	}
	parsed.RawQuery = ""
	parsed.Fragment = ""
	parsed.RawPath = ""
	// Strip path entirely — endpoint paths are absolute from origin.
	parsed.Path = ""
	return parsed.String(), nil
}

// CheckCredentials calls /auth/applications/check to validate that the configured
// publicKey/secretKey pair is accepted by the Ciabra API. Useful for an admin
// "Test Connection" feature.
func (c *Ciabra) CheckCredentials(ctx context.Context) error {
	var out map[string]any
	return c.doJSON(ctx, http.MethodGet, ciabraAuthCheckPath, nil, &out)
}

func (c *Ciabra) CreatePayment(ctx context.Context, req payment.CreatePaymentRequest) (*payment.CreatePaymentResponse, error) {
	minorAmount, err := payment.AmountToMinorUnit(req.Amount, payment.DefaultPaymentCurrency)
	if err != nil {
		return nil, fmt.Errorf("ciabra create payment: invalid amount %s", req.Amount)
	}
	minAmount := ciabraDefaultMinAmountCents
	if raw := strings.TrimSpace(c.config["minAmountCents"]); raw != "" {
		if parsed, err := parsePositiveInt(raw); err == nil {
			minAmount = parsed
		}
	}
	if minorAmount < int64(minAmount) {
		return nil, fmt.Errorf("ciabra create payment: minimum amount is %.2f", float64(minAmount)/100)
	}

	customerID, err := c.ensureCustomer(ctx, req)
	if err != nil {
		return nil, err
	}

	invoicePayload := c.buildInvoicePayload(customerID, req, minorAmount)
	var invoiceResp map[string]any
	if err := c.doJSON(ctx, http.MethodPost, ciabraInvoicesPath, invoicePayload, &invoiceResp); err != nil {
		return nil, fmt.Errorf("ciabra create invoice: %w", err)
	}

	invoiceID := ciabraFindString(invoiceResp,
		[]string{"id"},
		[]string{"data", "id"},
		[]string{"invoice", "id"},
	)
	payURL := ciabraFindString(invoiceResp,
		[]string{"checkoutUrl"},
		[]string{"paymentUrl"},
		[]string{"url"},
		[]string{"data", "checkoutUrl"},
		[]string{"data", "paymentUrl"},
		[]string{"payment", "url"},
	)
	qrCode := ciabraFindString(invoiceResp,
		[]string{"pixQrCode"},
		[]string{"qrCode"},
		[]string{"pix", "qrCode"},
		[]string{"pix", "copyPaste"},
		[]string{"data", "pix", "qrCode"},
		[]string{"data", "pix", "copyPaste"},
	)

	if (payURL == "" || qrCode == "") && invoiceID != "" {
		if invoiceDetails, queryErr := c.fetchInvoice(ctx, invoiceID); queryErr == nil {
			if payURL == "" {
				payURL = ciabraFindString(invoiceDetails,
					[]string{"checkoutUrl"},
					[]string{"paymentUrl"},
					[]string{"url"},
					[]string{"pix", "paymentUrl"},
				)
			}
			if qrCode == "" {
				qrCode = ciabraFindString(invoiceDetails,
					[]string{"pixQrCode"},
					[]string{"qrCode"},
					[]string{"pix", "qrCode"},
					[]string{"pix", "copyPaste"},
				)
			}
		}
	}

	if payURL == "" && qrCode == "" {
		return nil, fmt.Errorf("ciabra create invoice: payment URL/QR not found in response")
	}
	if payURL == "" {
		payURL = qrCode
	}

	return &payment.CreatePaymentResponse{
		TradeNo: invoiceID,
		PayURL:  payURL,
		QRCode:  qrCode,
	}, nil
}

func (c *Ciabra) QueryOrder(ctx context.Context, tradeNo string) (*payment.QueryOrderResponse, error) {
	invoiceID := strings.TrimSpace(tradeNo)
	if invoiceID == "" {
		return nil, fmt.Errorf("ciabra query order: missing invoice id")
	}
	invoice, err := c.fetchInvoice(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("ciabra query order: %w", err)
	}

	statusRaw := ciabraFindString(invoice,
		[]string{"status"},
		[]string{"data", "status"},
		[]string{"invoice", "status"},
	)
	providerStatus := ciabraProviderStatus(statusRaw)
	amount := ciabraFindAmount(invoice,
		[]string{"price"},
		[]string{"amount"},
		[]string{"total"},
		[]string{"data", "price"},
		[]string{"data", "amount"},
		[]string{"invoice", "price"},
	)

	metadata := map[string]string{"status": statusRaw}
	if currency := strings.ToUpper(strings.TrimSpace(c.config["currency"])); currency != "" {
		metadata["currency"] = currency
	}
	for k, v := range c.MerchantIdentityMetadata() {
		metadata[k] = v
	}

	return &payment.QueryOrderResponse{
		TradeNo:  invoiceID,
		Status:   providerStatus,
		Amount:   amount,
		Metadata: metadata,
	}, nil
}

// VerifyNotification parses a Ciabra webhook event. Per the docs, events are
// identified by `hookType` (one of INVOICE_CREATED, INVOICE_DELETED,
// PAYMENT_GENERATED, PAYMENT_RECEIVED, PAYMENT_CONFIRMED, PAYMENT_REVERSED).
//
// Returns nil for events that are not relevant to order fulfillment so the
// webhook handler can acknowledge them with 200 and stop retries.
func (c *Ciabra) VerifyNotification(_ context.Context, rawBody string, headers map[string]string) (*payment.PaymentNotification, error) {
	if err := c.verifyWebhookSignature(rawBody, headers, time.Now()); err != nil {
		return nil, err
	}
	if strings.TrimSpace(rawBody) == "" {
		return nil, fmt.Errorf("ciabra webhook empty body")
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(rawBody), &payload); err != nil {
		return nil, fmt.Errorf("ciabra parse webhook: %w", err)
	}

	// Prefer `hookType` (documented contract). Fall back to legacy `status`
	// shape so we keep working with older test fixtures and edge cases.
	hookType := ciabraFindString(payload,
		[]string{"hookType"},
		[]string{"type"},
		[]string{"event"},
		[]string{"event", "type"},
		[]string{"data", "hookType"},
		[]string{"data", "type"},
	)

	statusRaw := ciabraFindString(payload,
		[]string{"status"},
		[]string{"data", "status"},
		[]string{"invoice", "status"},
		[]string{"data", "invoice", "status"},
		[]string{"event", "status"},
	)

	providerStatus := ciabraEventToProviderStatus(hookType)
	if providerStatus == "" {
		providerStatus = ciabraProviderStatus(statusRaw)
	}

	// Only PAID confirmations affect the order. Ack other events with nil so
	// the handler returns 200 and Ciabra stops retrying.
	if providerStatus != payment.ProviderStatusPaid {
		return nil, nil
	}

	orderID := ciabraFindString(payload,
		[]string{"externalId"},
		[]string{"data", "externalId"},
		[]string{"invoice", "externalId"},
		[]string{"data", "invoice", "externalId"},
		[]string{"merchantOrderId"},
	)
	tradeNo := ciabraFindString(payload,
		[]string{"id"},
		[]string{"invoiceId"},
		[]string{"data", "id"},
		[]string{"data", "invoiceId"},
		[]string{"invoice", "id"},
		[]string{"data", "invoice", "id"},
	)

	if orderID == "" || tradeNo == "" {
		return nil, fmt.Errorf("ciabra webhook missing invoice id or external id")
	}

	amount := ciabraFindAmount(payload,
		[]string{"price"},
		[]string{"amount"},
		[]string{"total"},
		[]string{"data", "price"},
		[]string{"data", "amount"},
		[]string{"invoice", "price"},
		[]string{"data", "invoice", "price"},
	)

	metadata := map[string]string{}
	if hookType != "" {
		metadata["hook_type"] = hookType
	}
	if statusRaw != "" {
		metadata["status"] = statusRaw
	}
	if currency := strings.ToUpper(strings.TrimSpace(c.config["currency"])); currency != "" {
		metadata["currency"] = currency
	}
	for k, v := range c.MerchantIdentityMetadata() {
		metadata[k] = v
	}

	return &payment.PaymentNotification{
		TradeNo:  tradeNo,
		OrderID:  orderID,
		Amount:   amount,
		Status:   payment.NotificationStatusSuccess,
		RawData:  rawBody,
		Metadata: metadata,
	}, nil
}

func (c *Ciabra) Refund(_ context.Context, _ payment.RefundRequest) (*payment.RefundResponse, error) {
	return nil, fmt.Errorf("ciabra refund is not supported yet")
}

func (c *Ciabra) ensureCustomer(ctx context.Context, req payment.CreatePaymentRequest) (string, error) {
	externalID := strings.TrimSpace(req.OrderID)
	if externalID == "" {
		externalID = fmt.Sprintf("sub2-%d", time.Now().UnixNano())
	}

	name := strings.TrimSpace(req.CustomerName)
	if len([]rune(name)) < 3 {
		name = "Route-Cortexx Customer"
	}
	email := strings.TrimSpace(req.CustomerEmail)
	if email == "" {
		return "", fmt.Errorf("ciabra create customer: customer email is required")
	}

	document := normalizeBrazilianDocument(req.CustomerDocument)
	if document == "" {
		document = normalizeBrazilianDocument(c.config["defaultDocument"])
	}
	if document == "" {
		return "", fmt.Errorf("ciabra create customer: customer CPF/CNPJ is required (set customer_document in the order request or defaultDocument in the instance config)")
	}

	payload := map[string]any{
		"fullName":   name,
		"email":      email,
		"document":   document,
		"externalId": externalID,
	}
	if business := strings.TrimSpace(req.Subject); business != "" {
		payload["business"] = business
	}

	var resp map[string]any
	if err := c.doJSON(ctx, http.MethodPost, ciabraCustomersPath, payload, &resp); err != nil {
		return "", fmt.Errorf("ciabra create customer: %w", err)
	}
	customerID := ciabraFindString(resp,
		[]string{"id"},
		[]string{"data", "id"},
		[]string{"customer", "id"},
	)
	if customerID == "" {
		return "", fmt.Errorf("ciabra create customer: response missing id")
	}
	return customerID, nil
}

// normalizeBrazilianDocument strips formatting from a CPF/CNPJ. Returns ""
// when the input doesn't contain at least 11 digits.
func normalizeBrazilianDocument(raw string) string {
	var b strings.Builder
	for _, r := range raw {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	out := b.String()
	if len(out) < 11 {
		return ""
	}
	return out
}

func (c *Ciabra) buildInvoicePayload(customerID string, req payment.CreatePaymentRequest, minorAmount int64) map[string]any {
	dueDays := 1
	if rawDueDays := strings.TrimSpace(c.config["dueDays"]); rawDueDays != "" {
		if parsed, err := parsePositiveInt(rawDueDays); err == nil {
			dueDays = parsed
		}
	}
	// Ciabra examples use millisecond precision (e.g. 2024-08-21T03:00:00.000Z).
	dueDate := time.Now().UTC().Add(time.Duration(dueDays) * 24 * time.Hour).Format("2006-01-02T15:04:05.000Z")

	paymentTypes := ciabraNormalizePaymentTypes(c.config["paymentTypes"])
	if len(paymentTypes) == 0 {
		paymentTypes = []string{ciabraDefaultPaymentTypePIX}
	}

	description := strings.TrimSpace(req.Subject)
	if description == "" {
		description = strings.TrimSpace(req.OrderID)
	}
	if description == "" {
		description = "Route-Cortexx payment"
	}

	price := payment.MinorUnitToAmount(minorAmount, payment.DefaultPaymentCurrency)

	payload := map[string]any{
		"customerId":       customerID,
		"description":      description,
		"dueDate":          dueDate,
		"externalId":       req.OrderID,
		"paymentTypes":     paymentTypes,
		"price":            price,
		"installmentCount": 1,
		"invoiceType":      ciabraDefaultInvoiceType,
		// `items` must be sent as an empty array. The Ciabra API only accepts
		// items that reference a pre-registered product UUID (free-form items
		// crash with `.map()` on undefined). For PIX checkouts the line-item
		// breakdown is irrelevant — `price` + `description` carry the total.
		"items": []map[string]any{},
		// notifications and webhooks are documented as arrays. The upstream API
		// crashes with "Cannot read properties of undefined (reading 'map')"
		// when these keys are absent — send empty arrays as the safe default.
		"notifications": []map[string]any{},
	}

	if ret := strings.TrimSpace(req.ReturnURL); ret != "" {
		payload["redirectTo"] = ret
	}

	if hooks := c.buildWebhookSpecs(); len(hooks) > 0 {
		payload["webhooks"] = hooks
	} else {
		payload["webhooks"] = []map[string]any{}
	}

	return payload
}

// buildWebhookSpecs constructs the per-invoice webhooks array Ciabra expects.
// Returns nil when no public webhook URL is configured (admin has presumably
// set global webhooks in the Ciabra dashboard instead).
func (c *Ciabra) buildWebhookSpecs() []map[string]any {
	hookURL := strings.TrimSpace(c.config["webhookPublicURL"])
	if hookURL == "" {
		return nil
	}
	events := []string{
		"PAYMENT_GENERATED",
		"PAYMENT_RECEIVED",
		"PAYMENT_CONFIRMED",
		"PAYMENT_REVERSED",
	}
	out := make([]map[string]any, 0, len(events))
	for _, ev := range events {
		out = append(out, map[string]any{
			"hookType": ev,
			"url":      hookURL,
		})
	}
	return out
}

func ciabraNormalizePaymentTypes(raw string) []string {
	parts := strings.Split(raw, ",")
	seen := make(map[string]struct{}, len(parts))
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		v := strings.ToUpper(strings.TrimSpace(part))
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

func (c *Ciabra) fetchInvoice(ctx context.Context, invoiceID string) (map[string]any, error) {
	invoiceID = strings.TrimSpace(invoiceID)
	if invoiceID == "" {
		return nil, fmt.Errorf("missing invoice id")
	}
	var resp map[string]any
	if err := c.doJSON(ctx, http.MethodGet, ciabraInvoicesPath+"/"+url.PathEscape(invoiceID), nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// verifyWebhookSignature validates an HMAC-SHA256 signature using a shared
// secret. Comparison is constant-time. When no secret is configured, the
// signature step is skipped (relying on URL secrecy + idempotent order
// matching for safety); operators should always configure webhookSecret in
// production.
func (c *Ciabra) verifyWebhookSignature(rawBody string, headers map[string]string, now time.Time) error {
	secret := strings.TrimSpace(c.config["webhookSecret"])
	if secret == "" {
		return nil
	}

	sigHeader := strings.ToLower(strings.TrimSpace(c.config["webhookSignatureHeader"]))
	timestampHeader := strings.ToLower(strings.TrimSpace(c.config["webhookTimestampHeader"]))
	if sigHeader == "" {
		sigHeader = "x-signature"
	}
	if timestampHeader == "" {
		timestampHeader = "x-timestamp"
	}

	signature := strings.TrimSpace(headers[sigHeader])
	if signature == "" {
		return fmt.Errorf("ciabra webhook missing signature header %s", sigHeader)
	}

	if timestampRaw := strings.TrimSpace(headers[timestampHeader]); timestampRaw != "" {
		toleranceSeconds := 300
		if rawTol := strings.TrimSpace(c.config["webhookToleranceSeconds"]); rawTol != "" {
			if parsed, err := parsePositiveInt(rawTol); err == nil {
				toleranceSeconds = parsed
			}
		}
		ts, err := parseUnixTimestamp(timestampRaw)
		if err != nil {
			return fmt.Errorf("ciabra webhook timestamp is invalid")
		}
		if delta := now.Sub(ts).Abs(); delta > time.Duration(toleranceSeconds)*time.Second {
			return fmt.Errorf("ciabra webhook timestamp outside tolerance")
		}
	}

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(rawBody))
	macSum := mac.Sum(nil)
	expectedHex := strings.ToLower(hex.EncodeToString(macSum))
	expectedB64 := base64.StdEncoding.EncodeToString(macSum)
	expectedB64Raw := base64.RawStdEncoding.EncodeToString(macSum)

	provided := strings.ToLower(strings.TrimSpace(signature))
	hexMatch := subtle.ConstantTimeCompare([]byte(provided), []byte(expectedHex)) == 1
	b64Match := subtle.ConstantTimeCompare([]byte(signature), []byte(expectedB64)) == 1
	rawB64Match := subtle.ConstantTimeCompare([]byte(signature), []byte(expectedB64Raw)) == 1
	if !hexMatch && !b64Match && !rawB64Match {
		return fmt.Errorf("ciabra webhook signature mismatch")
	}
	return nil
}

func parseUnixTimestamp(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, fmt.Errorf("empty timestamp")
	}
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	if len(raw) > 10 {
		return time.UnixMilli(value), nil
	}
	return time.Unix(value, 0), nil
}

func parsePositiveInt(raw string) (int, error) {
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return 0, err
	}
	if value <= 0 {
		return 0, fmt.Errorf("value must be positive")
	}
	return value, nil
}

func (c *Ciabra) doJSON(ctx context.Context, method, path string, payload any, out any) error {
	var bodyReader io.Reader
	if payload != nil {
		body, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(body)
	}

	requestURL := strings.TrimRight(c.config["apiBase"], "/") + path
	req, err := http.NewRequestWithContext(ctx, method, requestURL, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+c.basicAuth())
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	body, status, err := c.do(req)
	if err != nil {
		return err
	}
	if status < http.StatusOK || status >= http.StatusMultipleChoices {
		summary := summarizeCiabraResponse(body)
		return fmt.Errorf("HTTP %d: %s", status, summary)
	}
	if out == nil || len(bytes.TrimSpace(body)) == 0 {
		return nil
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("parse response: %w", err)
	}
	return nil
}

func (c *Ciabra) basicAuth() string {
	token := strings.TrimSpace(c.config["publicKey"]) + ":" + strings.TrimSpace(c.config["secretKey"])
	return base64.StdEncoding.EncodeToString([]byte(token))
}

func (c *Ciabra) do(req *http.Request) ([]byte, int, error) {
	client := c.httpClient
	if client == nil {
		client = &http.Client{Timeout: ciabraHTTPTimeout}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(io.LimitReader(resp.Body, ciabraMaxResponseSize))
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return body, resp.StatusCode, nil
}

func summarizeCiabraResponse(body []byte) string {
	summary := strings.Join(strings.Fields(string(body)), " ")
	if summary == "" {
		return "<empty>"
	}
	if len(summary) > ciabraMaxErrorSummary {
		return summary[:ciabraMaxErrorSummary] + "..."
	}
	return summary
}

// ciabraEventToProviderStatus maps documented Ciabra hookType events to the
// internal provider status enum. Returns "" when the event is not relevant to
// order fulfillment (caller should fall back to status fields).
func ciabraEventToProviderStatus(hookType string) string {
	switch strings.ToUpper(strings.TrimSpace(hookType)) {
	case "PAYMENT_RECEIVED", "PAYMENT_CONFIRMED":
		return payment.ProviderStatusPaid
	case "PAYMENT_REVERSED", "PAYMENT_REFUNDED":
		return payment.ProviderStatusRefunded
	case "INVOICE_CREATED", "INVOICE_CHANGED", "PAYMENT_GENERATED", "INVOICE_DELETED":
		return payment.ProviderStatusPending
	}
	return ""
}

// ciabraProviderStatus keeps the legacy `status` field handling for QueryOrder
// responses (the docs don't enumerate them but the test fixtures show shapes
// like PAID/SETTLED/CONFIRMED).
func ciabraProviderStatus(raw string) string {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "PAID", "SETTLED", "SUCCESS", "CONFIRMED", "APPROVED":
		return payment.ProviderStatusPaid
	case "FAILED", "CANCELED", "CANCELLED", "EXPIRED", "VOIDED", "REFUNDED":
		return payment.ProviderStatusFailed
	default:
		return payment.ProviderStatusPending
	}
}

func ciabraFindString(root map[string]any, paths ...[]string) string {
	for _, path := range paths {
		if value, ok := ciabraFindAny(root, path); ok {
			switch typed := value.(type) {
			case string:
				if trimmed := strings.TrimSpace(typed); trimmed != "" {
					return trimmed
				}
			}
		}
	}
	return ""
}

func ciabraFindAmount(root map[string]any, paths ...[]string) float64 {
	for _, path := range paths {
		if value, ok := ciabraFindAny(root, path); ok {
			switch typed := value.(type) {
			case float64:
				return typed
			case float32:
				return float64(typed)
			case int:
				return float64(typed)
			case int32:
				return float64(typed)
			case int64:
				return float64(typed)
			case json.Number:
				if f, err := typed.Float64(); err == nil {
					return f
				}
			case string:
				if f, err := strconv.ParseFloat(strings.TrimSpace(typed), 64); err == nil {
					return f
				}
			}
		}
	}
	return 0
}

func ciabraFindAny(root map[string]any, path []string) (any, bool) {
	if root == nil || len(path) == 0 {
		return nil, false
	}
	var current any = root
	for _, segment := range path {
		object, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}
		value, exists := object[segment]
		if !exists {
			return nil, false
		}
		current = value
	}
	return current, true
}
