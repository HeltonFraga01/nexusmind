---
name: devops-omniroute-update
description: |
  Ensina o especialista em DevOps e infraestrutura a atualizar e homologar
  novas versões do OmniRoute no Portainer Swarm, aplicando patches de UI
  e validando a integridade das conexões.
license: MIT
metadata:
  version: "1.0.0"
  author: "Antigravity + Helton Fraga"
  emoji: "rocket"
  philosophy: "DevOps"
  triggers:
    - "atualizar o omniroute"
    - "atualizar omniroute"
    - "update omniroute"
    - "omniroute deploy"
    - "omniroute update"
  dependencies:
    - /Users/heltonfraga/Documents/Develop/NEXUSMIND/scripts/update_omniroute.py
---

# Skill: DevOps OmniRoute Update

## Quando Usar

Use esta skill sempre que houver necessidade de atualizar a versão do contêiner **OmniRoute** (`cortexx-omniroute_cortexx-omniroute`) no ambiente de produção do Portainer Swarm.

## Procedimento Padrão de Atualização

Toda a lógica de autenticação no Portainer, captura de configuração do serviço, incremento do contador de ForceUpdate, alteração da tag da imagem e monitoramento do deploy está encapsulada no script `scripts/update_omniroute.py`.

### Passo a Passo:

1. **Executar o Script de Atualização:**
   Chame o script passando a versão alvo desejada (ex: `3.8.2`):
   ```bash
   python3 scripts/update_omniroute.py --version <target_version>
   ```

   *Nota:* O script tentará ler as variáveis de ambiente `PORTAINER_URL`, `PORTAINER_USERNAME` e `PORTAINER_PASSWORD` para se autenticar. Caso não estejam definidas, usará os valores padrões configurados para o ambiente.

2. **Acompanhar Logs e Reinicialização:**
   Durante o deploy, o contêiner do OmniRoute executa o script de entrada `start-with-patch.sh`, que:
   - Aplica os patches de UI do Qdrant.
   - Executa as migrações automáticas de banco de dados SQLite.
   - Inicializa a aplicação Next.js.
   
   Você deve monitorar os logs do contêiner para garantir que não haja erros de inicialização:
   ```bash
   python3 scratch/portainer_omniroute_logs.py
   ```

3. **Verificar Portas de Escuta (API e Dashboard):**
   Execute o script de verificação interna para testar a comunicação interna nas portas `20128` (Dashboard) e `20129` (API):
   ```bash
   python3 scratch/portainer_verify_omniroute_api.py
   ```
   *Expectativa:* O retorno esperado para as chamadas sem credenciais é um status `401 Unauthorized` com JSON estruturado (indicando que a aplicação está ouvindo e rejeitando acessos não autorizados corretamente).

4. **Homologar Visualmente:**
   Navegue até a URL do painel (`https://omniroute.cortexx.online/home`) e verifique:
   - A versão exibida na barra lateral esquerda (deve corresponder à nova versão instalada).
   - A ausência do banner de notificação de atualização pendente.

## Diretrizes Críticas:

- **Preservação de Estado:** O banco de dados SQLite do OmniRoute reside em um volume persistente (`cortexx-omniroute-data`). O script de atualização do Portainer utiliza a mesma especificação de volumes, garantindo a integridade dos dados históricos, chaves e combos cadastrados.
- **Tratamento de Erros:** Se o contêiner falhar ou retornar status `unhealthy` durante o deploy, o script de atualização exibirá os logs de erro imediatamente e sairá com código de falha. Nesse caso, analise os logs em busca de erros nas migrações SQLite ou na aplicação de patches.
