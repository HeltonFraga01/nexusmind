// nexusmind — i18n overlay for NexusMind-only features (Ciabra PIX payment).
//
// Why this file exists: upstream en.ts / zh.ts must stay untouched so that
// `git merge upstream/main` never conflicts on locale files. The keys below are
// deep-merged into every locale at load time (see i18n/index.ts). pt-BR.ts
// already carries the Portuguese translations of these keys in its own
// `overrides` object, so the pt-BR entry here is intentionally empty.

type LocaleTree = Record<string, any>

const ciabraEn: LocaleTree = {
  admin: {
    settings: {
      payment: {
        providerCiabra: 'Ciabra (AZ Center PIX)',
        field_paymentTypes: 'Ciabra Payment Types',
        field_dueDays: 'Ciabra Due Days',
        field_minAmountCents: 'Minimum Amount (cents)',
        field_webhookPublicURL: 'Webhook Public URL',
        field_webhookSignatureHeader: 'Webhook Signature Header',
        field_webhookTimestampHeader: 'Webhook Timestamp Header',
        field_webhookToleranceSeconds: 'Webhook Tolerance (seconds)',
        field_defaultDocument: 'Default CPF/CNPJ',
        field_customerEmailDomain: 'Customer Email Domain',
        field_ciabraApiBaseHint: 'Use the AZ Center invoices base URL. Default: https://api.az.center.',
        field_ciabraPaymentTypesHint: 'Comma-separated payment methods accepted by AZ Center. Default: PIX.',
        field_ciabraDueDaysHint: 'Invoice expiration in days. Default: 1.',
        field_ciabraMinAmountHint: 'Minimum order amount in BRL cents. Orders below this are rejected by AZ Center. Default: 500 (R$ 5,00).',
        field_ciabraWebhookPublicURLHint: 'Public URL registered per invoice for webhook callbacks (e.g. https://your-host/api/v1/payment/webhook/ciabra). Leave empty if you configured a global webhook in the Ciabra dashboard.',
        field_ciabraWebhookSecretHint: 'Secret used to validate the webhook signature. Must match the value configured in Ciabra/AZ Center.',
        field_ciabraDefaultDocumentHint: 'Fallback CPF/CNPJ used when the order does not carry one. Mainly useful for testing.',
        ciabraWebhookHint: 'Configure the following URL as a Webhook endpoint in Ciabra/AZ Center and keep signature/timestamp header settings aligned with your account.',
        ciabraGuideSummary: 'Ciabra flow uses AZ Center invoices: create customer, create invoice (PIX), then confirm via webhook status PAID.',
        ciabraGuideNote: 'Use BRL for Ciabra invoice currency unless your AZ Center contract explicitly supports additional currencies. Keep webhook secret and header names synchronized between NexusMind and AZ Center.',
      },
    },
  },
  payment: {
    methods: {
      ciabra: 'Ciabra',
    },
    ciabra: {
      documentLabel: 'CPF or CNPJ',
      documentPlaceholder: '000.000.000-00',
      documentHint: 'Required to issue the PIX charge via Ciabra Invoice.',
    },
  },
}

// Keyed by locale code. zh inherits the English Ciabra labels; pt-BR.ts owns its own.
const overlay: Record<string, LocaleTree> = {
  en: ciabraEn,
  zh: ciabraEn,
  'pt-BR': {},
}

export default overlay
