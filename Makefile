.PHONY: prepare
prepare:
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
	go install github.com/cosmtrek/air@latest
	cd frontend/; npm ci;

	# dummy
	mkdir -p frontend/dist;
	touch frontend/dist/index.html;

.PHONY: gen
gen:
	buf generate --path=api/
	cd frontend/; npx buf generate --path=api/;

.PHONY: dev
dev: gen
	APP_ENV=development air -c .air.toml

.PHONY: build
build: gen
	cd frontend/; npm run build;
	go build main.go
