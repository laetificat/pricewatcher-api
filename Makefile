.DEFAULT_GOAL := help

help:
	@echo "Available options:\n- build: Build new version\n- test: Run all tests\n- run-webserver: Run API with debug logging\n- clean: Clean project\n- check: Test and check project\n- imports: run goimports with golangci rules"

build: clean
	@printf "%s" "Building darwin/amd64..."
	@env GOOS=darwin GOARCH=amd64 go build -o builds/pricewatcher-darwin-amd64
	@printf " %s\n" "Done!"
	@printf "%s" "Building linux/amd64..."
	@env GOOS=linux GOARCH=amd64 go build -o builds/pricewatcher-linux-amd64
	@printf " %s\n" "Done!"
	@printf "%s" "Building windows/amd64..."
	@env GOOS=windows GOARCH=amd64 go build -o builds/pricewatcher-windows-amd64.exe
	@printf " %s\n" "Done!"

test:
	@echo "Running gotest..."
	@gotest ./... -coverprofile=coverage.out -count=1

run-webserver:
	@go run main.go webserver -v=debug

clean:
	@printf "%s" "Cleaning project..."
	@rm -rf builds
	@printf " %s\n" "Done!"

check: test
	@golangci-lint run

imports:
	@goimports -local github.com/golangci/golangci-lint -w .