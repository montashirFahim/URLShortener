#!/bin/bash

# Configuration for Consul
export CONSUL_URL=localhost:8500
export CONSUL_PATH=config/backend

# Run the project using the serve command
go run main.go serve
