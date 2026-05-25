# NexusMind — Customizações do fork

Fork **estruturado** de [Wei-Shaw/sub2api](https://github.com/Wei-Shaw/sub2api).
Objetivo: trazer novas versões do upstream com `git merge` sem quebrar as melhorias.

## ⚠️ Regra de Ouro: Princípio da Alteração Mínima (Minimal Change Principle)

> [!IMPORTANT]
> **Fazer o mínimo de mudança de código possível é a nossa prioridade absoluta.**
> 
> **Por que fazemos isso?**
> Para permitir que novas versões, patches de segurança e melhorias do repositório original (*upstream* `Wei-Shaw/sub2api`) sejam integrados perfeitamente com um simples `git merge upstream/main`, sem gerar conflitos complexos de código e mantendo a stack ultra-estável e fácil de manter.
> 
> ### Diretrizes Obrigatórias:
> 1. **Home e Branding via Admin (Sem Código)**:
>    O layout, logotipo, cores, banners e todo o conteúdo da landing page (`Home Content`) são configurados dinamicamente via painel de administração em **Admin → Settings**, e **nunca** por modificação direta do código Vue.
> 2. **Somente Duas Coisas Exigem Código**:
>    - A tradução e internacionalização **pt-BR** (em arquivos isolados).
>    - O provedor de pagamento brasileiro **Ciabra PIX**.
> 3. **Marcador Obrigatório (`nexusmind`)**:
>    Toda linha de código adicionada ou modificada em arquivos originais do upstream deve, sem exceção, levar o comentário com o marcador correspondente:
>    - `// nexusmind` (Go, TypeScript, JavaScript)
>    - `/* nexusmind */` (CSS, Go, JS)
>    - `<!-- nexusmind -->` (HTML, Vue templates)
>    
>    Você pode auditar todas as customizações a qualquer momento executando:
>    ```bash
>    grep -rn "nexusmind" backend frontend
>    ```
> 4. **Surgical Patches (Patches Cirúrgicos)**:
>    Para alterações de infraestrutura e dependências externas (como o patch do container do Paperclip), priorize o uso de scripts de build ou patches dinâmicos em vez de modificar arquivos estáticos do core.


## Estrutura

| Item | Valor |
|------|-------|
| Branch única | `main` (produção, contém as customizações) |
| Remote `upstream` | `https://github.com/Wei-Shaw/sub2api.git` |
| Remote `origin` | `https://github.com/HeltonFraga01/nexusmind.git` |
| Imagem Docker | `heltonfraga/nexusmind` |

## Arquivos NOVOS (zero conflito no merge)

| Arquivo | Função |
|---------|--------|
| `frontend/src/i18n/locales/pt-BR.ts` | Locale pt-BR completo — `deepMerge(en, overrides)` |
| `frontend/src/i18n/nexusmind-i18n.ts` | Overlay de chaves NexusMind/Ciabra para `en`/`zh` |
| `backend/internal/payment/provider/ciabra.go` | Provedor Ciabra PIX (AZ Center) |
| `backend/internal/payment/provider/ciabra_payload_test.go` | Testes do provedor |
| `nexusmind/**` | Esta pasta — docs, deploy, scripts, home.html |

## Arquivos do upstream EDITADOS (touchpoints — possível conflito)

Todos com edição **aditiva** marcada `nexusmind`.

### i18n
- `frontend/src/i18n/index.ts` — registra loader pt-BR + merge do overlay
- `frontend/src/i18n/__tests__/usageServiceTierLocales.spec.ts` — teste pt-BR

### Ciabra — backend
- `backend/internal/payment/types.go` — `TypeCiabra`, campos de cliente
- `backend/internal/payment/provider/factory.go` — case `TypeCiabra`
- `backend/internal/handler/payment_handler.go` — campo `customer_document`
- `backend/internal/handler/payment_webhook_handler.go` — `CiabraWebhook` + parsing
- `backend/internal/handler/payment_webhook_handler_test.go` — testes
- `backend/internal/server/routes/payment.go` — rota `webhook/ciabra`
- `backend/internal/service/payment_service.go` — campo `CustomerDocument`
- `backend/internal/service/payment_config_providers.go` — registro do provider
- `backend/internal/service/payment_currency.go` — currency do Ciabra
- `backend/internal/service/payment_order.go` — snapshot, customer email/nome
- `backend/internal/service/payment_order_provider_snapshot.go` — validação
- `backend/internal/service/payment_*_test.go` — testes

### Ciabra — frontend
- `frontend/src/components/payment/providerConfig.ts` — campos de config Ciabra
- `frontend/src/components/payment/paymentFlow.ts` — alias/tipo + `customerDocument`
- `frontend/src/components/payment/PaymentMethodSelector.vue`
- `frontend/src/components/payment/PaymentProviderDialog.vue`
- `frontend/src/components/payment/ProviderCard.vue`
- `frontend/src/components/admin/payment/AdminOrderTable.vue`
- `frontend/src/views/admin/SettingsView.vue`
- `frontend/src/views/admin/orders/AdminOrdersView.vue`
- `frontend/src/views/user/PaymentView.vue` — campo CPF/CNPJ
- `frontend/src/types/payment.ts`
- `frontend/src/components/payment/__tests__/*.spec.ts` — testes

### Branding NexusMind (defaults/fallbacks — ver "Home e branding via admin")
- `backend/internal/service/setting_service.go` — seed `site_name = NexusMind`
- `backend/internal/service/auth_service.go`, `auth_email_binding.go`, `auth_oauth_email_flow.go`
- `backend/internal/service/balance_notify_service.go`, `content_moderation.go`
- `backend/internal/service/totp_service.go`, `user_service.go`, `payment_order.go`
- `frontend/index.html`, `frontend/src/main.ts`, `frontend/src/stores/app.ts`
- `frontend/src/router/index.ts` — `/home` público (landing)

## Home e branding via admin (sem código)

O upstream já suporta home 100% customizável: `HomeView.vue` renderiza o setting
`HomeContent` (HTML cru ou URL em iframe). Nome/logo/subtítulo também são settings.

Depois de instalar, em **Admin → Settings**:
- `Site Name` → `NexusMind`
- `Site Logo` → logo
- `Home Content` → colar o HTML de [`nexusmind/home.html`](home.html) (ou uma URL)

Assim a home muda **sem novo build**.

## Build e deploy

```bash
docker buildx build --platform linux/arm64,linux/amd64 \
  --build-arg VERSION=<versão> \
  -t heltonfraga/nexusmind:<versão> -t heltonfraga/nexusmind:latest \
  -f Dockerfile --push .
```

O `Dockerfile` do upstream **não é editado** — ele já fixa `pnpm@9` (build
determinístico). A imagem é nomeada pela flag `-t`, não pelo Dockerfile.

Deploy em Swarm ([`nexusmind/deploy/`](deploy/)):
- `stack.swarm.yml` — baseline genérico (porta publicada direta).
- `stack.swarm.nexusmind.digital.yml` — produção em **nexusmind.digital**
  (Traefik + Let's Encrypt). Deployar como stack `nexusmind`.

Domínio de produção: **https://nexusmind.digital** (deployment próprio,
independente de route.cortexx.online).

## Atualizar do upstream

```bash
bash nexusmind/sync-upstream.sh   # git fetch + merge upstream/main
# resolver conflitos só nos touchpoints acima (edições aditivas, manter ambos os lados)
bash nexusmind/verify.sh          # typecheck + build + testes + auditoria de marcadores
git push origin main
```

## Manutenção de segurança

A cada sync, rodar `govulncheck ./...` no `backend/` e bumpar dependências
sinalizadas. Mantém o `go.mod` sem divergência manual do upstream.
