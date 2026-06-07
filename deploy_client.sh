#!/bin/bash

TARGET=$1
AGENT_ID=$2
TOKEN=$3
SERVER_ADDR=$4
INTERVAL=$5

if [ -z "$TARGET" ] || [ -z "$AGENT_ID" ] || [ -z "$TOKEN" ] || [ -z "$SERVER_ADDR" ] || [ -z "$INTERVAL" ]; then
  echo "Uso: ./deploy.sh <usuario@ip> <nome-do-agente> <token> <endereco-servidor> <intervalo-segundos>"
  exit 1
fi

echo "Compilando o agente para Linux..."
CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o homelens-agent ./cmd/client/main.go

if [ $? -ne 0 ]; then
  echo "Erro na compilação! Corrija o código Go e tente novamente."
  exit 1
fi
echo "Build concluído com sucesso!"

echo "Transferindo executável para $TARGET..."
scp homelens-agent $TARGET:~

echo "Configurando agente remotamente via SSH..."
ssh -t $TARGET "
  sudo mv ~/homelens-agent /usr/local/bin/homelens-agent
  sudo chmod +x /usr/local/bin/homelens-agent

  sudo mkdir -p /etc/homelens

  echo 'Configurando variáveis de ambiente...'
  echo 'HOMELENS_AUTH_TOKEN=$TOKEN' | sudo tee /etc/homelens/agent.env > /dev/null
  echo 'HOMELENS_AGENT_ID=$AGENT_ID' | sudo tee -a /etc/homelens/agent.env > /dev/null
  echo 'HOMELENS_SERVER_ADDR=$SERVER_ADDR' | sudo tee -a /etc/homelens/agent.env > /dev/null
  echo 'HOMELENS_SECONDS_INTERVAL=$INTERVAL' | sudo tee -a /etc/homelens/agent.env > /dev/null

  echo 'Atualizando serviço do systemd...'
  echo '[Unit]
Description=HomeLens Agent
After=network.target

[Service]
Type=simple
EnvironmentFile=/etc/homelens/agent.env
ExecStart=/usr/local/bin/homelens-agent
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target' | sudo tee /etc/systemd/system/homelens-agent.service > /dev/null

  echo 'Reiniciando serviço...'
  sudo systemctl daemon-reload
  sudo systemctl enable homelens-agent
  sudo systemctl restart homelens-agent
  
  echo 'Status do serviço:'
  sudo systemctl status homelens-agent --no-pager | head -n 5
"

echo "Deploy concluído para $AGENT_ID."
