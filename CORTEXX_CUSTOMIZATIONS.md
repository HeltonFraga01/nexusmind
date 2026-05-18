# Route-Cortexx — Customizações do Fork

Fork de [Wei-Shaw/sub2api](https://github.com/Wei-Shaw/sub2api) adaptado para operação brasileira.

---

## Estrutura de branches

| Branch | Finalidade |
|--------|-----------|
| `sub2api-cortexx` | **Branch de produção.** Toda mudança vai aqui. |
| `main` | Mirror do upstream (Wei-Shaw/sub2api). Não commitamos aqui. |
| `upstream/main` | Remote do repo original. Usado para sync. |

Regra simples: **um PR por feature, sempre apontando para `sub2api-cortexx`**.

---

## Customizações aplicadas

### 1. Localização pt-BR (`feat: cortexx pt-BR localization`)
- `frontend/src/i18n/locales/pt-BR.ts` — tradução completa do locale
- `frontend/src/i18n/index.ts` — registro do locale pt-BR + detecção de idioma do browser
- `frontend/src/main.ts` — lang default `pt-BR`
- Estratégia: importa o locale inglês e aplica `deepMerge` com overrides — novas chaves upstream são herdadas automaticamente em inglês

### 2. Integração Ciabra PIX (`feat(payment): add Ciabra PIX integration`)
- `backend/internal/payment/provider/ciabra.go` — provider de pagamento PIX via Ciabra
- `backend/internal/payment/provider/factory.go` — registra o provider `ciabra`
- `backend/internal/payment/types.go` — tipo `ciabra` adicionado ao enum
- `backend/internal/handler/payment_handler.go` e `payment_webhook_handler.go` — rotas de webhook
- `backend/internal/service/payment_*.go` — lógica de fulfillment e snapshots
- `backend/internal/server/routes/payment.go` — rota pública do webhook
- `frontend/src/api/payment.ts` — endpoint `getPlans()` público
- `frontend/src/components/payment/PaymentMethodSelector.vue` — filtro para exibir Ciabra
- `frontend/src/views/admin/orders/AdminOrdersView.vue` — filtro de tipo `ciabra` no admin
- `Dockerfile` — build multi-arch com `PNPM_VERSION` fixo para builds determinísticos
- `Dockerfile.cortexx` — build local (usa frontend já compilado, mais rápido)
- `Dockerfile.cortexx.runtime` — empacota binário compilado localmente

### 3. Branding Route-Cortexx (`feat(brand): apply Route-Cortexx defaults`)
- `backend/internal/service/setting_service.go` — default `site_name = Route-Cortexx`
- `backend/internal/service/auth_service.go` e relacionados — issuer/branding nos emails
- `frontend/src/stores/app.ts` — fallback `siteName = Route-Cortexx`
- `frontend/index.html` — título inicial `Route-Cortexx`
- `frontend/src/router/index.ts` — `/home` liberada como rota pública

### 4. Home Page v2 (`feat(home): redesign landing page with Route-Cortexx v2 theme`)
- `frontend/src/views/HomeView.vue` — landing page completa redesenhada
  - Design dark com grid animado e diagrama de roteamento ao vivo
  - Seção de pricing busca planos **dinamicamente** via `GET /api/v1/payment/public/plans`
  - Fallback para planos default se API falhar
  - Seção de math (comparativo de preços), providers, tutorial, FAQ, testimonials
  - Floating CTA pill, scroll reveal animations

---

## Planos de assinatura (produção)

Configurados via admin em `https://route.cortexx.online/admin`:

| # | Nome | Preço | Tokens/mês | Requisições/mês |
|---|------|-------|-----------|-----------------|
| 1 | Solo dev | R$ 79 | 10M | 1.500 |
| 2 | Pro | R$ 249 | 100M | 10.000 |
| 3 | Equipe | R$ 749 | 500M | 30.000 |

Para atualizar planos sem deploy: **admin panel → Payment → Plans** ou via API:
```bash
TOKEN=$(curl -s -X POST https://route.cortexx.online/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@route.cortexx.online","password":"<SENHA>"}' \
  | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['access_token'])")

curl -s -X PUT https://route.cortexx.online/api/v1/admin/payment/plans/<ID> \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Pro","price":249,"description":"..."}'
```

---

## Build e deploy

### Build multi-arch (arm64 + amd64) — padrão de produção

```bash
cd /Users/heltonfraga/Documents/Develop/sub2api-cortexx

docker buildx use multiarch
docker buildx build \
  --platform linux/arm64,linux/amd64 \
  --build-arg VERSION=<versão> \
  --build-arg COMMIT=$(git rev-parse --short HEAD) \
  --build-arg DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  -t heltonfraga/sub2api-cortexx:<versão> \
  -t heltonfraga/sub2api-cortexx:latest \
  -f Dockerfile \
  --push \
  .
```

Imagem publicada em: `docker.io/heltonfraga/sub2api-cortexx`

### Versões publicadas

| Versão | Commit | Descrição |
|--------|--------|-----------|
| 0.1.128 | 965f5c9b | Patch CVEs: otel 1.37→1.43, Alpine 3.21→3.23 |
| 0.1.127 | 81c4cda2 | Home v2 + branding + Ciabra PIX + pt-BR |
| 0.1.126 | — | Versão anterior (upstream) |

### Deploy no servidor (via Portainer)

1. Abre o Portainer do servidor
2. Vai em **Stacks → cortexx-app**
3. Atualiza a variável `SUB2API_IMAGE` para `heltonfraga/sub2api-cortexx:0.1.127`
4. Clica em **Update the stack**

Ou via SSH no nó manager:
```bash
docker service update --image heltonfraga/sub2api-cortexx:0.1.127 <nome-do-service>
```

---

## Atualizando do upstream

Quando o Wei-Shaw/sub2api lançar versão nova:

```bash
# 1. Busca as novidades
git fetch upstream

# 2. Rebase nossas mudanças em cima da nova versão
git checkout sub2api-cortexx
git rebase upstream/main

# 3. Resolve conflitos (os arquivos mais prováveis estão abaixo)
# Se houver conflito, prioridade:
#   - frontend/src/views/HomeView.vue → manter NOSSA versão
#   - frontend/src/i18n/locales/pt-BR.ts → manter NOSSA versão
#   - Dockerfile → verificar se upstream mudou algo crítico, depois reaplica pnpm pin
#   - backend/internal/payment/provider/factory.go → re-registrar ciabra após o merge

# 4. Verifica que tudo builda
pnpm --dir frontend run typecheck
pnpm --dir frontend run build
cd backend && go build ./...

# 5. Sobe nova versão da imagem
docker buildx build --platform linux/arm64,linux/amd64 \
  -t heltonfraga/sub2api-cortexx:<nova-versão> \
  -t heltonfraga/sub2api-cortexx:latest \
  -f Dockerfile --push .

# 6. Atualiza no servidor (Portainer ou docker service update)
```

### Arquivos mais suscetíveis a conflito

| Arquivo | Por que conflita | Como resolver |
|---------|-----------------|---------------|
| `frontend/src/views/HomeView.vue` | Upstream tem home genérica, nós temos a v2 | Manter nossa versão |
| `frontend/src/i18n/locales/pt-BR.ts` | Upstream não tem pt-BR | Manter nossa versão |
| `frontend/src/i18n/index.ts` | Upstream pode adicionar locales | Manter registro do pt-BR |
| `backend/internal/payment/provider/factory.go` | Upstream pode adicionar providers | Re-registrar `ciabra` |
| `Dockerfile` | Upstream atualiza versões de Node/Go | Verificar compatibilidade, manter `PNPM_VERSION` pin |
