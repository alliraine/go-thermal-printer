.PHONY: swagger dev

dev: swagger
	go run ./cmd/go-thermal-printer/main.go

swagger:
	swag init -g ./cmd/go-thermal-printer/main.go -o ./pkg/docs --parseDependency --parseInternal --useStructName
