SHELL:=/bin/bash

.DEFAULT_GOAL := default

bin_folder = bin

dependencies:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor

format:
	go fmt ./...

jenkins-status: dependencies format
	go build -o $(bin_folder)/jenkinsStatus

default: jenkins-status

clean:
	rm $(bin_folder)/*