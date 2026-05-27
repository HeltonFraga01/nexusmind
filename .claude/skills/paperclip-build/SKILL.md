---
name: paperclip-build
description: |
  Automatiza o monitoramento, patching e build da imagem Docker do Paperclip AI
  a partir do repositório upstream. Detecta novas versões, aplica patches de URL
  de API, constrói imagem multi-arch e cria issue de notificação no GitHub.
license: MIT
metadata:
  version: "1.0.0"
  author: "Antigravity + Helton Fraga"
  emoji: "paperclip"
  philosophy: "DevOps"
  triggers:
    - "atualizar paperclip"
    - "build paperclip"
    - "checar versão do paperclip"
    - "verificar paperclip upstream"
    - "paperclip patch"
    - "paperclip deploy"
  dependencies:
    - scripts/patch_and_build_paperclip.py
    - .github/workflows/build-paperclip.yml
---

# Skill: Paperclip Build & Patch Automation

## Quando Usar

Use esta skill sempre que precisar:
- Verificar se existe uma nova versão do **Paperclip** no upstream (`paperclipai/paperclip`)
- Aplicar patches e construir a imagem `heltonfraga/paperclipai` no Docker Hub
- Disparar manualmente o workflow de build via GitHub Actions

---

## Tasks Disponíveis

### 1. Verificar versão rodando no Portainer
Consulta o Portainer e retorna a versão atual da imagem `heltonfraga/paperclipai` em execução no stack.

```bash
python3 scripts/patch_and_build_paperclip.py --dry-run
```
> O modo `--dry-run` executa toda a lógica de comparação de versões sem fazer push de imagem ou criar issues.

---

### 2. Executar pipeline completo (checar → patch → build → push → issue)
Detecta nova versão no upstream, clona o repositório, aplica patches, constrói imagem multi-arch e cria issue de notificação.

```bash
python3 scripts/patch_and_build_paperclip.py
```

**Variáveis de ambiente obrigatórias para build e push:**
| Variável | Descrição |
|---|---|
| `GITHUB_TOKEN` | Token de acesso ao GitHub (para criar issues) |
| `DOCKERHUB_USERNAME` | Usuário do Docker Hub |
| `DOCKERHUB_TOKEN` | Token do Docker Hub |
| `PORTAINER_URL` | URL do Portainer (ex: `https://painel.wasend.com.br`) |
| `PORTAINER_USERNAME` | Usuário do Portainer |
| `PORTAINER_PASSWORD` | Senha do Portainer |

---

### 3. Pipeline completo em repositório específico
Define o repositório GitHub onde a issue de notificação será criada.

```bash
python3 scripts/patch_and_build_paperclip.py --repo HeltonFraga01/nexusmind
```

---

### 4. Disparar workflow via GitHub Actions (manual)
O workflow `build-paperclip.yml` pode ser disparado manualmente via `workflow_dispatch` no GitHub ou agendado diariamente às **12:00 UTC (09:00 BRT)**.

**Para disparar manualmente:**
1. Acesse o repositório no GitHub.
2. Vá em **Actions → Patch and Build Paperclip**.
3. Clique em **Run workflow**.

---

## O que os Patches Fazem

O script `patch_and_build_paperclip.py` aplica automaticamente as seguintes modificações no código do upstream:

| Patch | Descrição |
|---|---|
| **URLs Anthropic/OpenAI** | Substitui URLs hardcoded por variáveis de ambiente (`ANTHROPIC_BASE_URL`, `OPENAI_BASE_URL`) |
| **codex-home.ts** | Injeta gerador dinâmico de `config.toml` com suporte a WebSockets e roteamento NexusMind |
| **runtime-config.ts** | Injeta configuração dinâmica de provedores OpenAI com modelos NexusMind (gpt-5.x) |
| **Dockerfile** | Instala `hermes-agent` (Python CLI) e `hermes-paperclip-adapter` (npm) |

---

## Imagens Docker Geradas

| Tag | Descrição |
|---|---|
| `heltonfraga/paperclipai:<versão>-fixed` | Versão específica patcheada |
| `heltonfraga/paperclipai:latest` | Última versão patcheada |

---

## Fluxo de Atualização Pós-Build

Após o build bem-sucedido, uma issue é criada automaticamente no repositório `HeltonFraga01/nexusmind` com instruções para:
1. Abrir o Portainer.
2. Navegar até o stack **paperclip**.
3. Atualizar a imagem para `heltonfraga/paperclipai:<versão>-fixed`.
4. Clicar em **Update the stack**.
