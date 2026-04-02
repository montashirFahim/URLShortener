#!/bin/bash

# Configuration YAML with corrected password 'db01'
CONFIG_YAML=$(cat <<EOF
postgres:
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "db01"
  db_name: "url_shortener"
  maxidleconn: 10
  maxopenconn: 100
  maxconnlifetime: 3600

redis:
  host: "localhost"
  port: 6379
  db: 0

jwt_secret: "eat-sleep^repeat?"
port: 8080
debug: true
EOF
)

# Seed Consul
curl --request PUT \
  --data "$CONFIG_YAML" \
  http://localhost:8500/v1/kv/config/backend
