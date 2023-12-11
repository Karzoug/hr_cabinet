gen-api-types:
	oapi-codegen -generate "types" -package api ./api/api.json > ./internal/server/internal/api/types.gen.go

gen-api-server:
	oapi-codegen -generate "chi-server" -package api ./api/api.json > ./internal/server/internal/api/server.gen.go

install-dependencies:
	go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest