//go:build unit

package provider

import (
	"regexp"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/payment"
	"github.com/stretchr/testify/require"
)

func TestCiabraBuildInvoicePayloadMatchesDocumentedShape(t *testing.T) {
	prov, err := NewCiabra("inst_1", map[string]string{
		"publicKey":        "pk",
		"secretKey":        "sk",
		"apiBase":          "https://api.az.center",
		"paymentTypes":     "PIX",
		"dueDays":          "1",
		"webhookPublicURL": "https://example.test/api/v1/payment/webhook/ciabra",
	})
	require.NoError(t, err)

	payload := prov.buildInvoicePayload("cust_123", payment.CreatePaymentRequest{
		OrderID:   "sub2_test_order",
		Subject:   "Route-Cortexx payment",
		ReturnURL: "https://app.example.com/payment/result?order=123",
	}, 1000) // R$10.00

	require.Equal(t, "cust_123", payload["customerId"])
	require.Equal(t, "sub2_test_order", payload["externalId"])
	require.Equal(t, 10.0, payload["price"])
	require.Equal(t, []string{"PIX"}, payload["paymentTypes"])
	require.Equal(t, 1, payload["installmentCount"])
	require.Equal(t, "SINGLE", payload["invoiceType"])
	require.Equal(t, "https://app.example.com/payment/result?order=123", payload["redirectTo"])

	// Items must be empty — Ciabra rejects free-form items with `.map()` crash.
	items, ok := payload["items"].([]map[string]any)
	require.True(t, ok)
	require.Empty(t, items)

	// Compatibility-superset fields removed.
	_, hasAmount := payload["amount"]
	_, hasTotal := payload["total"]
	_, hasValue := payload["value"]
	_, hasPaymentMethods := payload["paymentMethods"]
	_, hasInvoiceItems := payload["invoiceItems"]
	require.False(t, hasAmount, "amount must not be sent (not documented)")
	require.False(t, hasTotal, "total must not be sent (not documented)")
	require.False(t, hasValue, "value must not be sent (not documented)")
	require.False(t, hasPaymentMethods, "paymentMethods must not be sent (not documented)")
	require.False(t, hasInvoiceItems, "invoiceItems must not be sent (docs use items)")

	hooks, ok := payload["webhooks"].([]map[string]any)
	require.True(t, ok, "webhooks array must be present when webhookPublicURL is set")
	require.Len(t, hooks, 4)
	gotEvents := map[string]bool{}
	for _, h := range hooks {
		require.Equal(t, "https://example.test/api/v1/payment/webhook/ciabra", h["url"])
		gotEvents[h["hookType"].(string)] = true
	}
	require.True(t, gotEvents["PAYMENT_RECEIVED"])
	require.True(t, gotEvents["PAYMENT_CONFIRMED"])
	require.True(t, gotEvents["PAYMENT_REVERSED"])
	require.True(t, gotEvents["PAYMENT_GENERATED"])

	dueDate, ok := payload["dueDate"].(string)
	require.True(t, ok)
	// Ciabra examples use millisecond precision: 2024-08-21T03:00:00.000Z
	matched, err := regexp.MatchString(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$`, dueDate)
	require.NoError(t, err)
	require.True(t, matched, "dueDate must use millisecond precision: %q", dueDate)
}

func TestCiabraBuildInvoicePayloadSendsEmptyArraysWhenUnconfigured(t *testing.T) {
	prov, err := NewCiabra("inst_2", map[string]string{
		"publicKey": "pk",
		"secretKey": "sk",
	})
	require.NoError(t, err)

	payload := prov.buildInvoicePayload("cust_x", payment.CreatePaymentRequest{
		OrderID: "ord_x",
		Subject: "Single",
	}, 1000)

	// Even without a webhook URL the keys must exist as empty arrays — the
	// Ciabra API crashes with `.map()` on undefined when they are omitted.
	hooks, ok := payload["webhooks"].([]map[string]any)
	require.True(t, ok)
	require.Empty(t, hooks)

	notifications, ok := payload["notifications"].([]map[string]any)
	require.True(t, ok)
	require.Empty(t, notifications)

	_, hasRedirect := payload["redirectTo"]
	require.False(t, hasRedirect)
}

func TestCiabraEventToProviderStatus(t *testing.T) {
	cases := map[string]string{
		"PAYMENT_RECEIVED":  payment.ProviderStatusPaid,
		"PAYMENT_CONFIRMED": payment.ProviderStatusPaid,
		"PAYMENT_REVERSED":  payment.ProviderStatusRefunded,
		"PAYMENT_GENERATED": payment.ProviderStatusPending,
		"INVOICE_CREATED":   payment.ProviderStatusPending,
		"":                  "",
		"WHATEVER":          "",
	}
	for in, want := range cases {
		got := ciabraEventToProviderStatus(in)
		require.Equal(t, want, got, "input=%q", in)
	}
}

func TestNormalizeBrazilianDocument(t *testing.T) {
	require.Equal(t, "12345678901", normalizeBrazilianDocument("123.456.789-01"))
	require.Equal(t, "12345678000190", normalizeBrazilianDocument("12.345.678/0001-90"))
	require.Equal(t, "", normalizeBrazilianDocument("123"))
	require.Equal(t, "", normalizeBrazilianDocument(""))
}

func TestNormalizeCiabraAPIBaseStripsLegacyInvoicesSuffix(t *testing.T) {
	got, err := normalizeCiabraAPIBase("https://api.az.center/invoices")
	require.NoError(t, err)
	require.Equal(t, "https://api.az.center", got)

	got, err = normalizeCiabraAPIBase("")
	require.NoError(t, err)
	require.Equal(t, "https://api.az.center", got)

	_, err = normalizeCiabraAPIBase("http://insecure.example.com")
	require.Error(t, err)
}
