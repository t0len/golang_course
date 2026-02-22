# Practice 3 (Go)

## Run
1) Create DB in pgAdmin: `mydb` (owner: postgres)
2) Set env vars (if needed):
   - PG_HOST=localhost
   - PG_PORT=5432
   - PG_USER=postgres
   - PG_PASS=<your postgres password>
   - PG_DB=mydb
   - PG_SSL=disable
   - API_KEY=dev-secret (optional)
3) Install deps:
   `go mod tidy`
4) Start:
   `go run cmd/api/main.go`

## Test
- GET /healthz
- GET /users (requires header X-API-KEY: dev-secret)
