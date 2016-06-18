LINTIGNOREDEPS='vendor/.+\.go'
KAMIMAI_ONLY_PKGS=$(shell go list ./... 2> /dev/null | grep -v "/misc/" | grep -v "/vendor/")

all:

ansible:
	ansible-playbook misc/ansible/localhost.yml

unit: lint vet test

lint:
	@echo "go lint"
	@lint=`golint ./...`; \
	lint=`echo "$$lint" | grep -E -v -e ${LINTIGNOREDEPS}`; \
	echo "$$lint"; \
	if [ "$$lint" != "" ]; then exit 1; fi

vet:
	@echo "go vet"
	@go tool vet -all -structtags -shadow $(shell ls -d */ | grep -v "misc" | grep -v "vendor")

test:
	@go test $(KAMIMAI_ONLY_PKGS)


.PHONY: all ansible lint vet
