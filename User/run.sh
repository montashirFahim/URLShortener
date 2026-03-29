#!/bin/bash

export CONSUL_URL=localhost:8500
export CONSUL_PATH=config/user

go run main.go user