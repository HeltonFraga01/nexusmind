# Cortexx Customizations

This fork keeps Cortexx changes small and easy to rebase on top of upstream
Sub2API releases.

## Current customizations

- Added Brazilian Portuguese locale support:
  - `frontend/src/i18n/locales/pt-BR.ts`
  - `frontend/src/i18n/index.ts`
  - `frontend/src/i18n/__tests__/usageServiceTierLocales.spec.ts`
- Made Docker frontend builds deterministic:
  - `Dockerfile` pins `PNPM_VERSION=10.32.1`
  - `frontend/package.json` explicitly allows build scripts for `esbuild` and `vue-demi`
- Added local deployment Dockerfiles:
  - `Dockerfile.cortexx` builds the backend inside Docker and embeds the
    already-built frontend from `backend/internal/web/dist`.
  - `Dockerfile.cortexx.runtime` packages a locally-built backend binary. This
    is the preferred path when Docker BuildKit stalls on local machines.

## Localization strategy

`pt-BR.ts` imports the upstream English locale and deep-merges a small override
tree on top of it. This keeps updates safe:

- new upstream English keys are inherited automatically;
- untranslated keys remain visible in English instead of breaking screens;
- our translation surface stays concentrated in one file.

When translating more UI, add only the changed keys to `overrides` in
`pt-BR.ts`. Avoid copying the full English locale unless we intentionally want
to own every string.

## Updating from upstream

```bash
cd /Users/heltonfraga/Documents/Develop/sub2api-cortexx
git fetch upstream
git switch cortexx/pt-br-localization
git rebase upstream/main
cd frontend
pnpm install --frozen-lockfile
pnpm test:run frontend/src/i18n/__tests__/usageServiceTierLocales.spec.ts
pnpm build
```

If upstream changes `frontend/src/i18n/index.ts`, resolve the conflict by
keeping the `pt-BR` locale registration and browser-language detection.
If upstream changes the Dockerfile frontend stage, keep deterministic pnpm
pinning unless upstream has replaced it with an equivalent reproducible build.

## Building the custom image

```bash
cd /Users/heltonfraga/Documents/Develop/sub2api-cortexx
/Users/heltonfraga/.npm-global/bin/pnpm --dir frontend install --frozen-lockfile
/Users/heltonfraga/.npm-global/bin/pnpm --dir frontend build
(cd backend && \
  VERSION_VALUE=$(tr -d '\r\n' < cmd/server/VERSION) && \
  DATE_VALUE=$(date -u +%Y-%m-%dT%H:%M:%SZ) && \
  CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
  -tags embed \
  -ldflags="-s -w -X main.Version=${VERSION_VALUE} -X main.Commit=cortexx-local -X main.Date=${DATE_VALUE} -X main.BuildType=release" \
  -trimpath \
  -o ../build/sub2api-linux-arm64 \
  ./cmd/server)
/Applications/Docker.app/Contents/Resources/bin/docker build \
  -f Dockerfile.cortexx.runtime \
  -t sub2api-cortexx:pt-br .
```

The full Docker build is also available:

```bash
/Applications/Docker.app/Contents/Resources/bin/docker build \
  -f Dockerfile.cortexx \
  -t sub2api-cortexx:pt-br-full .
```

Then set the deployed compose image to `sub2api-cortexx:pt-br` and recreate the
application container.
