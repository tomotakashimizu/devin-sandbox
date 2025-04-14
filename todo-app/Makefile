.PHONY: generate-api

generate-api:
	oapi-codegen -config api/oapi-codegen-config.yaml -package api api/openapi.yaml > internal/interfaces/api/generated.go
