# Hermes Agent — Arquitetura e Operação (v6)

> **Última atualização**: 2026-05-27  
> **Versão da Stack**: v6 (State-of-the-Art)  
> **Mantenedor**: HeltonFraga01/nexusmind

---

## 1. Visão Geral

O **Hermes Agent** é um agente de inteligência artificial conversacional, orquestrado via Docker Swarm e integrado ao ecossistema corporativo do Portainer. Ele é ativado pelo Telegram e utiliza o NexusMind como proxy reverso para acesso aos modelos de linguagem (LLMs) da Anthropic e da OpenAI — sem custo de API adicional além da assinatura NexusMind.

### Stack no Portainer
- **Stack ID**: `345`
- **Stack Name**: `cortexx-hermes`
- **Service Name**: `cortexx-hermes_hermes-agent`
- **Imagem Docker**: `heltonfraga/hermes-agent:amd64`

---

## 2. Princípios de Arquitetura

### 2.1 Sem Gambiarras — Configuração 100% via ENVs
**Regra fundamental**: Toda configuração do Hermes deve estar na seção `Env` da Stack do Portainer. O script de bootstrap no `Args` do Swarm lê dinamicamente essas variáveis via `${VAR}`.

- ❌ Não há hardcodes em scripts de inicialização
- ❌ Não há scripts ad-hoc executados após a subida do container
- ✅ Toda mudança de configuração é feita nas ENVs da Stack no Portainer
- ✅ Reinicializações são idempotentes e auto-suficientes

### 2.2 Fallback Automático
O Hermes opera com dois modelos:
| Camada | Provider | Modelo | Base URL |
|---|---|---|---|
| **Primário** | anthropic (via NexusMind) | `claude-sonnet-4-6` | `http://nexusmind_nexusmind:8080` |
| **Fallback** | openai (via NexusMind) | `gpt-5.4` | `http://nexusmind_nexusmind:8080/v1` |

Se o provider primário falhar, o Hermes chaveia automaticamente para o fallback **sem intervenção manual**.

---

## 3. Mapeamento de ENVs da Stack

### 3.1 ENVs Obrigatórias

| ENV | Descrição |
|---|---|
| `HERMES_MODEL_DEFAULT` | Modelo primário (ex: `claude-sonnet-4-6`) |
| `HERMES_MODEL_PROVIDER` | Provider primário (`anthropic`) |
| `HERMES_MODEL_BASE_URL` | URL interna do NexusMind (ex: `http://nexusmind_nexusmind:8080`) |
| `HERMES_MODEL_API_KEY` | Chave de API do NexusMind para o modelo primário |
| `HERMES_FALLBACK_MODEL` | Modelo de contingência (ex: `gpt-5.4`) |
| `HERMES_FALLBACK_PROVIDER` | Provider de contingência (`openai`) |
| `HERMES_FALLBACK_BASE_URL` | URL interna do NexusMind para fallback (`/v1`) |
| `HERMES_FALLBACK_API_KEY` | Chave de API do NexusMind para fallback |
| `TELEGRAM_BOT_TOKEN` | Token do bot do Telegram |
| `TELEGRAM_ALLOWED_USERS` | Lista de IDs de usuários autorizados |
| `STT_PROVIDER` | Provider de transcrição de áudio (`openai`) |
| `STT_OPENAI_BASE_URL` | URL interna do Whisper (`http://cortexx-whisper_whisper:8000/v1`) |
| `VOICE_TOOLS_OPENAI_KEY` | Chave de API para o Whisper (pode ser a do NexusMind) |
| `STT_OPENAI_MODEL` | Modelo Whisper local (`Systran/faster-whisper-base`) |

### 3.2 ENVs de Superpoderes Corporativos (v6)

| ENV | Serviço | URL Interna Swarm |
|---|---|---|
| `SEARXNG_URL` | Busca anônima local | `http://cortexx-searxng_searxng:8080` |
| `QDRANT_URL` | Vector DB para RAG | `http://cortexx-qdrant_qdrant:6333` |
| `MINIO_ENDPOINT` | Object Storage | `http://minio_minio:9000` |
| `MINIO_ACCESS_KEY` | Credencial MinIO | — |
| `MINIO_SECRET_KEY` | Credencial MinIO | — |
| `MINIO_BUCKET` | Bucket do Hermes | `hermes-artifacts` |
| `NOCODB_URL` | Planilhas inteligentes | `http://nocodb_nocodb:8080` |
| `WUZAPI_URL` | API WhatsApp | `http://wuzapi-cortexx_wuzapi-cortexx:8089` |
| `N8N_URL` | Automação de workflows | `http://n8n_editor_n8n_editor:5678` |

---

## 4. Bootstrap Script (Args do Swarm)

O container executa um script `sh -c` no `command` que realiza, em ordem:

1. **Verificação de pacotes essenciais** (`python-telegram-bot`, `fastmcp`) no virtualenv do container
2. **Instalação leve e rápida** somente se necessário (`pip install --no-cache-dir`)
3. **Instalação editable core** do `hermes-agent` sem dependências pesadas de áudio (`.[all]`)
4. **Criação de diretórios** necessários (`/opt/data/plans`, `/opt/data/skills`, etc.)
5. **Limpeza de configs antigos** para evitar estados corrompidos de reinicializações anteriores
6. **Geração do guia de auto-compreensão** em `/opt/data/skills/hermes_capabilities_guide.md`
7. **Geração dinâmica do `config.yaml`** a partir das ENVs da Stack
8. **Geração dinâmica do `.env`** a partir das ENVs da Stack
9. **Cópia dos configs** para `/root/.hermes/` (local esperado pelo agente)
10. **Inicialização do proxy FastMCP** devmemory em background (`porta 8092`)
11. **Execução do gateway do Telegram** (`gateway run`)

### Virtualenv do Container
O container utiliza um virtualenv pré-instalado com todas as dependências pesadas de IA:
```
/root/.hermes/hermes-agent/venv/bin/python
/root/.hermes/hermes-agent/venv/bin/fastmcp
```

---

## 5. Integração com Whisper (STT — Speech To Text)

### Diagnóstico da Porta (Histórico)
O serviço Whisper ativo no cluster é o `fedirz/faster-whisper-server:latest-cpu`, que escuta na porta **`8000`** — e **não** na porta `9000` (que pertencia a um container antigo descontinuado).

### Modelo Offline
O container do Whisper opera sem acesso à internet. O modelo disponível localmente no volume `cortexx-whisper_whisper-models` é:
```
Systran/faster-whisper-base
```
A configuração do `STT_OPENAI_MODEL` deve **sempre** apontar para esse modelo para evitar o erro `LocalEntryNotFoundError`.

### Configuração no `config.yaml`
```yaml
stt:
  provider: openai
  base_url: http://cortexx-whisper_whisper:8000/v1
  api_key: <VOICE_TOOLS_OPENAI_KEY>
  openai:
    model: Systran/faster-whisper-base
```

---

## 6. Recursos de Memória e Persistência

### Volumes Docker
| Volume | Ponto de Montagem | Conteúdo |
|---|---|---|
| `hermes-data` | `/opt/data` | config.yaml, .env, skills, plans, mcp |
| `hermes-workspace` | `/workspace` | workspace de trabalho do agente |

### Arquivos Chave em `/opt/data`
| Arquivo | Descrição |
|---|---|
| `config.yaml` | Configuração principal do Hermes |
| `.env` | Variáveis de ambiente para ferramentas Python |
| `skills/hermes_capabilities_guide.md` | Manual de auto-compreensão gerado no bootstrap |
| `mcp/devmemory_mcp_server.py` | Proxy FastMCP para o devmemory |

---

## 7. Recursos de Infraestrutura (Docker Swarm)

### Limites de Recursos
```yaml
deploy:
  resources:
    reservations:
      cpus: "0.1"
      memory: 256M
    limits:
      cpus: "4.0"
      memory: 8192M
```

> ⚠️ **Importante**: Reservas formais altas (ex: 2GB) podem impedir o agendamento do container em nodes com alta concorrência de memória, resultando em erro `no suitable node (insufficient resources)`. Mantenha reservas baixas.

### Política de Restart
```yaml
restart_policy:
  condition: any
  delay: 10s
  max_attempts: 0   # Sem limite de tentativas
```

### Rede
O Hermes opera na rede externa do Swarm `network_public`, que dá conectividade a todos os outros serviços do cluster.

---

## 8. Scripts de Diagnóstico (Scratch)

Os scripts abaixo existem localmente em `scratch/` (não versionados por `.gitignore`) para diagnóstico e operação:

| Script | Descrição |
|---|---|
| `redeploy_hermes_definitive_v6.py` | Deploy definitivo da Stack v6 |
| `exec_in_hermes_new.py` | Executa comandos dentro do container ativo |
| `get_hermes_logs_decoded.py` | Busca e decodifica logs do container |
| `inspect_hermes_tasks.py` | Inspeciona as tasks do Swarm do serviço |
| `find_hermes_id.py` | Encontra o ID do serviço Hermes no Swarm |

> Para rodar qualquer script: `python3 scratch/<nome>.py` a partir da raiz do projeto.

---

## 9. Histórico de Issues Resolvidas

| Issue ID | Título | Resolução |
|---|---|---|
| `FRG-418` | Primary model failed — fallback não funciona | Configurado `fallback_model` nativo via ENVs da Stack com `gpt-5.4` via OpenAI NexusMind |
| STT Whisper Port | Whisper apontando para porta errada `9000` | Corrigido para porta `8000` do container `fedirz/faster-whisper-server` |
| STT LocalEntryNotFoundError | Modelo Whisper não encontrado offline | Configurado `STT_OPENAI_MODEL=Systran/faster-whisper-base` (pré-baixado no volume) |
| Bootstrap Loop | Container travando na re-instalação pesada `.[all]` | Otimizado bootstrap para instalar apenas pacotes essenciais leves |
| Insufficient Resources | Container em `pending` por reserva de 2GB RAM | Reduzidas reservas para `256M` mantendo limite de `8192M` |
| Paperclip Loop | 23+ issues duplicadas e travadas | Cleanup cirúrgico via Postgres no container `paperclip-production_paperclip-db` |
