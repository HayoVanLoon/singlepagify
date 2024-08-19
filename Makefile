GO := $(shell which go)

run:
ifndef INPUT
	$(error missing INPUT)
endif
ifndef OUTPUT
	$(error missing OUTPUT)
endif
	$(GO) run cmd/*.go '$(INPUT)' '$(OUTPUT)'
