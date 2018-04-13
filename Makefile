SHELL:=/bin/bash

.DEFAULT_GOAL := default

bin_folder = bin

dependencies:
	dep ensure

format:
	go fmt ./...

jenkins-status: format
	go build -o $(bin_folder)/jenkinsStatus

default: jenkins-status

clean:
	rm $(bin_folder)/*