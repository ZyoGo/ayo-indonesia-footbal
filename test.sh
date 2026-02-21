#!/bin/bash
export PATH=$PATH:/usr/local/go/bin:/home/noodle/go/bin

cd /mnt/d/project/ayo-indonesia
go run go.uber.org/mock/mockgen -source=internal/club/domain/repository.go -destination=internal/club/mock/repository_mock.go -package=mock
go test ./internal/club/app/... -v
