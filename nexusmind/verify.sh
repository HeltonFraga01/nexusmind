#!/usr/bin/env bash
# nexusmind — verifica que as customizações sobreviveram (build + testes).
set -euo pipefail
cd "$(dirname "$0")/.."

echo "== Auditoria de marcadores 'nexusmind' =="
COUNT=$(grep -rn "nexusmind" backend frontend 2>/dev/null | wc -l | tr -d ' ')
echo "   $COUNT linhas marcadas encontradas"
if [ "$COUNT" -eq 0 ]; then
  echo "   ERRO: nenhum marcador — as customizações podem ter sido perdidas no merge."
  exit 1
fi

echo "== Backend: build =="
(cd backend && go build ./...)

echo "== Backend: testes (Ciabra + pagamento) =="
(cd backend && go test -tags=unit \
  ./internal/payment/... ./internal/handler/... ./internal/service/)

echo "== Frontend: deps + typecheck =="
(cd frontend && pnpm install --frozen-lockfile && pnpm run typecheck)

echo "== Frontend: testes (i18n + pagamento + router) =="
(cd frontend && npx vitest run \
  src/i18n/__tests__/usageServiceTierLocales.spec.ts \
  src/components/payment/__tests__/providerConfig.spec.ts \
  src/components/payment/__tests__/paymentFlow.spec.ts \
  src/router/__tests__/guards.spec.ts)

echo "== Frontend: build =="
(cd frontend && pnpm run build)

echo
echo "OK — customizações verificadas."
