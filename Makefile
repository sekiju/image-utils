build:
	@go build -o bin/image_utils cmd/image_utils/main.go

buildw:
	@go build -o bin/image_utils.exe cmd/image_utils/main.go

run:
	@go run cmd/image_utils/main.go