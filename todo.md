# HomeLens — Task Checklist

## 1. Authentication & Identity

- [x] Define `HOMELENS_AGENT_ID` env var on the client
- [x] Define `HOMELENS_AUTH_TOKEN` env var on the client
- [x] Define `HOMELENS_SERVER_ADDR` env var on the client
- [x] Send agent ID + token as query params on WebSocket connect (`ws://server:8080/ws?token=xxx&agent_id=ubuntu-1`)
- [x] Server validates token before accepting the connection
- [x] Server rejects and closes connection on invalid token
- [x] Server registers agent ID and associates incoming snapshots to it

## 2. Reconnection (Client)

- [x] Detect server disconnect (write/read error)
- [x] Implement reconnect loop with exponential backoff (1s, 2s, 4s, 8s... max 30s)
- [x] Reset backoff on successful reconnection
- [x] Log reconnection attempts
- [x] Continue collecting metrics during disconnect (discard or buffer)

## 3. SQLite (Server)

- [x] Add SQLite dependency (`modernc.org/sqlite` or `mattn/go-sqlite3`)
- [x] Define schema: `agents` table (id, name, last_seen)
- [x] Define schema: `snapshots` table (id, agent_id, timestamp, data JSON)
- [x] Insert snapshot on receive (throttle: 1 per 10s or 1 per minute for storage)
- [ ] Data retention: cron/goroutine to delete snapshots older than X days
- [x] Upsert agent `last_seen` on each snapshot

## 4. Server — Agent Management

- [x] Track connected agents in memory (map of agent ID → connection info)
- [x] Remove agent from map on disconnect
- [x] Endpoint or method to list all agents with status (online/offline/last_seen)
- [x] Store latest snapshot per agent in memory for instant access

## 5. REST API (Server)

- [x] `GET /api/agents` — list all agents with status
- [x] `GET /api/agents/:id` — agent detail with latest snapshot
- [x] `GET /api/agents/:id/history?from=&to=` — historical snapshots for graphs
- [x] `GET /api/stats/:id` — aggregated metrics (avg CPU, max memory over time range)

## 6. Frontend WebSocket (Server → Browser)

- [x] WebSocket endpoint for frontend clients (`/ws/live`)
- [x] On agent snapshot received, broadcast to all connected frontend clients
- [x] Send initial state (all agents + latest snapshots) on frontend connect
- [x] Handle frontend disconnect gracefully

## 7. Frontend (Web UI)

- [x] Choose framework (plain HTML+JS, React, or Templ for Go templates)
- [x] Dashboard page: overview of all agents (name, status, CPU, memory, temp)
- [x] Agent detail page: per-agent graphs (CPU, memory, network, disk over time)
- [x] Online/offline indicator per agent
- [x] Auto-update via WebSocket (live data)
- [x] Historical graphs using REST API data

## 8. Alerts

- [ ] Define alert rules (CPU > 90% for 5 min, disk > 95%, agent offline > 2 min)
- [ ] Alert engine: goroutine evaluating rules against incoming snapshots
- [ ] Alert state management (firing, resolved, cooldown to avoid spam)
- [ ] Log alerts to SQLite
- [ ] Notification channel: Telegram bot integration (webhook/API)
- [ ] Optional: email, generic webhook

## 9. Deploy & Infrastructure

- [x] Dockerfile for the server
- [x] Systemd unit file for the agent
- [x] Build script: `go build -o homelens-agent ./cmd/agent` and `go build -o homelens-server ./cmd/server`
- [x] `.env.example` for both agent and server
- [ ] README with setup instructions

## 10. Nice to Have

- [ ] Multi-path disk space monitoring (not just `/`)
- [x] Read thermal zone type from `/sys/class/thermal/thermal_zone*/type` for descriptive labels
- [x] Filter disk I/O to whole disks only (skip partitions, loop, zram, dm)
- [x] Network interface filtering (skip `lo`, optionally skip zero-traffic interfaces)
- [ ] Event log (agent connected, disconnected, alert fired/resolved)
- [ ] Server config file (YAML/TOML) for alert thresholds, retention period, port, etc.
- [ ] Agent config file as alternative to env vars
- [ ] Docker container metrics (via Docker socket)

## HomeLens — Checklist de Tarefas Atómicas

### Fase 1: Redesign do Front-end e Layout

- [x] Instanciar o medidor de Temperatura no topo da tela de detalhes
- [x] Definir `animation: false` na configuração principal do ECharts
- [x] Criar o container vazio com placeholder para a futura tabela de processos no layout
- [x] Criar o container vazio com placeholder para a futura tabela do Docker no layout

### Fase 2: Identidade e Refatoração do Back-end

- [x] Criar função no agente em Go para ler o arquivo `/etc/machine-id`
- [x] Adicionar o `machine-id` extraído como query param na string de conexão do WebSocket
- [x] Alterar o esquema da tabela `agents` no SQLite para usar o `machine-id` como chave primária
- [x] Alterar a query de inserção/upsert no servidor para usar a nova chave do `machine-id`
- [x] Refatorar a struct do payload do agente alterando o campo de temperatura para ponteiro (`*float64`)
- [x] Definir a interface `TempStrategy` com os métodos `Name()` e `Read()` no agente
- [x] Criar a struct `ThermalZoneStrategy` assinando a interface de temperatura
- [x] Implementar a leitura de `/sys/class/thermal/thermal_zone*/temp` na `ThermalZoneStrategy`
- [x] Criar a struct `HwmonStrategy` assinando a interface de temperatura
- [x] Implementar a leitura de `/sys/class/hwmon/hwmon*/temp*_input` na `HwmonStrategy`
- [x] Atualizar o construtor do coletor do agente para injetar e percorrer o array de estratégias

### Fase 4: Integração com o Docker

- [x] Criar função no agente para testar se o arquivo `/var/run/docker.sock` está acessível
- [x] Criar structs para o payload do Docker (nome, status, cpu, memória) no agente e no servidor
- [x] Implementar a listagem de containers ativos usando o método `ContainerList` do SDK
- [x] Injetar o array de containers coletados dentro da struct principal de snapshot do agente
- [x] Atualizar o loop de transmissão do servidor para repassar os dados de containers ao front-end
- [x] Substituir o placeholder do Docker no front-end por uma tabela React real usando os dados do WebSocket

### Fase 5: Leitura de Processos Nativos

- [x] Adicionar a biblioteca `gopsutil` (`github.com/shirou/gopsutil/v3/process`) no `go.mod` do agente
- [x] Criar a struct de payload para processos (pid, usuário, nome, cpu) no agente e servidor
- [x] Criar função no agente para listar todos os PIDs rodando no sistema operacional
- [x] Implementar o cálculo de porcentagem de CPU consumida por cada processo mapeado
- [x] Implementar a lógica de ordenação (sort) do array de processos do maior para o menor consumo
- [x] Fatiar o array ordenado para manter estritamente os 10 processos mais pesados
- [x] Injetar o array do Top 10 processos na struct principal de snapshot do agente
- [x] Substituir o placeholder de processos no front-end por uma tabela React real usando os dados do WebSocket

### Fase 6: Motor de Alertas e Limpeza

- [x] Criar uma goroutine com `time.Ticker` no servidor dedicada exclusivamente a rodar o motor de alertas
- [x] Criar a lógica de buffer/histórico na memória do servidor para monitorar a persistência da CPU alta
- [x] Adicionar a verificação de CPU > 90% durando mais de 5 minutos antes de mudar o estado do alerta
- [x] Criar um ticker secundário na goroutine para verificar agentes com `now() - last_seen > 2 minutos`
- [x] Implementar um mapa de controle de estado em memória para gerenciar o cooldown e evitar spam de alertas
- [x] Criar uma goroutine com ticker de 1 hora no servidor para executar `DELETE FROM snapshots WHERE timestamp < ...`

### Fase: Configurações e Motor de Alertas

**Banco de Dados (SQLite)**

- [x] Criar a tabela `alert_configs` no `schema.sql` com as colunas (id, cpu_threshold, mem_threshold, disk_threshold, offline_mins, webhook_url)
- [x] Adicionar query no `query.sql` para buscar a configuração atual (`GetAlertConfig`)
- [x] Adicionar query no `query.sql` para atualizar a configuração (`UpdateAlertConfig`) garantindo que afete apenas a linha principal
- [x] Rodar o `sqlc generate` para criar os métodos em Go

**Back-end: Memória e API**

- [x] Criar struct `AlertEngine` ou similar, contendo o cache da configuração protegido por `sync.RWMutex` e o mapa de estado `map[string]*AlertState`
- [x] Carregar a configuração do banco de dados para a memória logo no startup do servidor (`main.go`)
- [x] Criar o endpoint `GET /api/alerts/config` para retornar a configuração que está no cache de memória
- [x] Criar o endpoint `POST /api/alerts/config` para receber novos limites, salvar no SQLite e atualizar o cache em memória em seguida

**Back-end: Motor de Alertas (Goroutine)**

- [x] Instanciar uma goroutine no servidor com um `time.Ticker` (ex: a cada 15 ou 30 segundos)
- [x] Implementar a lógica de avaliação: pegar o snapshot em memória de cada agente e comparar com os limites do cache
- [x] Implementar a lógica de janela de tempo: registrar o `start_time` do pico no mapa de estados e acionar o alerta apenas se a anomalia durar mais de 5 minutos
- [x] Escrever a função que dispara o Webhook (HTTP POST) com payload JSON informando "Alerta Disparado", alterando o estado em memória para `is_firing = true`
- [x] Implementar a regra de resolução: se a métrica voltar ao normal e o `is_firing` for true, disparar o Webhook de "Alerta Resolvido" e limpar o estado do mapa
- [x] Implementar a regra de inatividade: calcular a diferença entre o tempo atual e o `last_seen` do agente, disparando alerta se passar do limite configurado

**Front-end (UI)**

- [x] Adicionar as tipagens referentes ao payload de configuração no arquivo `models.ts`
- [x] Criar as funções de fetch no client da API (`getAlertConfig` e `updateAlertConfig`)
- [x] Criar o componente da página de configurações (`/settings`) e registrar a rota lá no `App.tsx`
- [x] Adicionar um ícone de engrenagem no `Header.tsx` para permitir a navegação até a página de configurações
- [x] Montar o formulário usando o `react-hook-form`, reaproveitando o componente `TextInput` para os numéricos (CPU, Mem, Disk) e a string da URL
- [x] Fazer o binding do submit do formulário para chamar a API e disparar os toasts da biblioteca avisando se salvou com sucesso ou se deu erro

Back-end: Concorrência e Contextos

[] Refatorar a assinatura de todas as funções do agente (ex: readCPUTime, readMemoryUsage) para receberem ctx context.Context como primeiro parâmetro

[ ] Adicionar checagens de ctx.Done() nos loops mais pesados ou demorados dentro dos coletores do agente

[ ] Amarrar o ciclo de vida do WebSocket no server.go ao contexto da requisição HTTP (r.Context()), garantindo que o cancelamento se propague caso o cliente feche a aba abruptamente

Back-end: Testes com a Standard Library

[ ] Criar o arquivo client/cpu_test.go e escrever um table-driven test puro para a função getCPU, validando os cálculos de porcentagem

[ ] Criar o arquivo client/disk_test.go e client/net_test.go com table-driven tests para as funções calcDiskIOUsage e calcNetUsage

[ ] Definir uma interface Store (ou Repository) no pacote server/db ou server/api que contenha a assinatura de todos os métodos que você usa do db.Queries (ex: GetAgents, InsertSnapshot, etc.)

[ ] Refatorar a struct API e AgentServer para dependerem da interface Store em vez do ponteiro concreto \*db.Queries

[ ] Criar um arquivo server/api/mock_store_test.go e implementar a struct MockStore assinando a interface, com campos de função customizáveis (ex: ListAgentsFunc func() ([]db.Agent, error))

[ ] Criar o arquivo server/api/api_test.go usando o pacote net/http/httptest para simular requisições (ex: GET /api/agents) usando o MockStore, testando os status codes (200 OK, 500 Internal Error)

Front-end: Setup e Testes Unitários

[ ] Instalar as dependências de teste no front-end: npm i -D vitest @testing-library/react @testing-library/jest-dom jsdom

[ ] Atualizar o arquivo vite.config.ts para incluir a configuração do Vitest (adicionar test: { environment: 'jsdom', globals: true })

[ ] Criar o arquivo src/utils/index.test.ts e escrever testes para a função formatByteStr (testar conversões para KB, MB, GB, etc.)

[ ] Escrever testes unitários no mesmo arquivo para a função convertByteToMetric garantindo que a matemática bata corretamente

Front-end: Testes de Componentes Visuais

[ ] Criar um arquivo de setup para o Vitest (setupTests.ts) e importar o @testing-library/jest-dom para ter acesso aos matchers do DOM

[ ] Criar o arquivo src/components/metricBar.test.tsx

[ ] Renderizar o <MetricBar /> no teste com valor 50 (seguro) e garantir que as classes de aviso não estejam presentes

[ ] Renderizar o <MetricBar /> no teste com valor 80 e garantir que a barra recebe a classe CSS referente ao estado de warning

[ ] Renderizar o <MetricBar /> no teste com valor 95 e garantir que a barra recebe a classe CSS referente ao estado critical
