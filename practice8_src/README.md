# Practice 8 – Unit Testing in Go

## Setup (run once)

```bash
go mod tidy

# Install mockgen tool
go install go.uber.org/mock/mockgen@latest

# Add mockgen to PATH (if needed)
export PATH=$PATH:$(go env GOPATH)/bin
```

## Regenerate mock (if you change the interface)

```bash
mockgen -source=repository/user_repository.go -package=repository > repository/mock_user_repository.go
```

## Run all tests

```bash
# Task 1 – calc tests (root package)
go test -v .

# Task 2 – service tests
go test -v ./service/...

# Task 3 – exchange service tests (root package, same as task 1)
go test -v .

# ALL at once
go test -v ./...
```

## Run tests with race detector (optional)
```bash
go test -race ./...
```

## Coverage report (optional bonus)
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```
