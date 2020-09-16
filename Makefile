.PHONY : prepare build run

$(eval $(service):;@:)
ifndef service
$(error "service" is not set)
endif

$(eval $(gomod):;@:)
ifndef gomod
$(error "go modules" is not set)
endif

init:
	go run cmd/scaffold_maker/*.go --scope=initservice --servicename=$(service) --modules=$(modules) --gomod=$(gomod)
	# @$(MAKE) -f $(lastword $(MAKEFILE_LIST)) proto

add-module:
	go run cmd/scaffold_maker/*.go --scope=addmodule --servicename=$(service) --modules=$(modules) --gomod=$(gomod)
	# @$(MAKE) -f $(lastword $(MAKEFILE_LIST)) proto

prepare:
	@if [ ! -d "cmd/$(service)" ]; then  echo "ERROR: service '$(service)' undefined"; exit 1; fi
	@ln -sf cmd/$(service)/main.go main_service.go

build: prepare
	go build -o bin

run: build
	./bin

# proto:
# 	$(foreach proto_file, $(shell find api/$(service)/proto -name '*.proto'),\
# 	protoc -I . $(proto_file) --go_out=plugins=grpc:. --go_opt=paths=source_relative;)

# docker: prepare
# 	docker build --build-arg SERVICE_NAME=$(service) -t $(service):latest .

# run-container:
# 	docker run --name=$(service) --network="host" -d $(service)

# generate-rsa-key:
# 	sh scripts/generate_rsa_key.sh

clear:
	rm main_service.go bin backend-service
