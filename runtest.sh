#!/bin/bash
go test ./... -race -coverprofile=coverage.txt -covermode=atomic
