# URL Shortener (Ruckus Networks) — README

Simple URL shortener service implemented in Go. This repo contains a small HTTP API to shorten URLs and redirect short codes to original URLs. The service uses an in-memory storage implementation and is packaged with Docker and Docker Compose.

Quick links (workspace)
- Files:
  - [docker-compose.yaml](docker-compose.yaml)
  - [Dockerfile](Dockerfile)
  - [go.mod](go.mod)
  - [Readme.md](Readme.md)
  - [cmd/server/main.go](cmd/server/main.go)
  - [internals/handler/handler.go](internals/handler/handler.go)
  - [internals/service/service.go](internals/service/service.go)
  - [internals/storage/memory.go](internals/storage/memory.go)
  - [internals/storage/storage.go](internals/storage/storage.go)
  - [test/service_test.go](test/service_test.go)
  - [test/handler_test.go](test/handler_test.go)

- Important symbols:
  - [`service.NewURLService`](internals/service/service.go)
  - [`service.URLService.ShortenURL`](internals/service/service.go)
  - [`service.URLService.GetLongURL`](internals/service/service.go)
  - [`service.URLService.GenerateShortCode`](internals/service/service.go)
  - [`service.ErrInvalidURL`](internals/service/service.go)
  - [`handler.NewHandler`](internals/handler/handler.go)
  - [`handler.Handler.ShortenURL`](internals/handler/handler.go)
  - [`handler.Handler.RedirectURL`](internals/handler/handler.go)
  - [`storage.NewMemoryStorage`](internals/storage/memory.go)
  - [`storage.Storage` interface](internals/storage/storage.go)
  - [`storage.ErrNotFound`](internals/storage/storage.go)

Requirements
- Docker & Docker Compose
- (Optional) Go 1.24 or newer for running tests locally (module set in [go.mod](go.mod))

How to run the code

1) Run with Docker Compose (recommended)
- Build and start (background):
```sh
docker compose up --build -d
```
- Check container logs:
```sh
docker compose logs -f app
```
- Stop and remove:
```sh
docker compose down
```
Notes:
- The service listens on port 8080 (host port mapped in [docker-compose.yaml](docker-compose.yaml)).
- Base URL is set via the `BASE_URL` environment variable in [docker-compose.yaml](docker-compose.yaml) and defaults to `http://localhost:8080`. The server code that reads these is in [cmd/server/main.go](cmd/server/main.go).

2) Run locally (without Docker)
- With Go installed (recommended for development):
```sh
# from repo root
go run ./cmd/server
# or build binary
go build -o urlshort ./cmd/server
./urlshort
```
- Environment variables:
  - PORT (default 8080)
  - BASE_URL (default http://localhost:8080)
The app is wired in [cmd/server/main.go](cmd/server/main.go) which creates the storage [`storage.NewMemoryStorage`](internals/storage/memory.go), the service [`service.NewURLService`](internals/service/service.go) and the handlers [`handler.NewHandler`](internals/handler/handler.go).

3) Run tests
- Locally (requires Go toolchain):
```sh
go test ./... -v
```
Tests are in [test/service_test.go](test/service_test.go) and [test/handler_test.go](test/handler_test.go).
- Inside a container (no local Go):
```sh
docker run --rm -v "$(pwd):/app" -w /app golang:1.24 go test ./... -v
```

4) Manual / API testing (curl or Postman)
- Shorten a URL (POST):
```sh
curl -s -X POST -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}' \
  http://localhost:8080/api/shorten
# -> {"short_url":"http://localhost:8080/EAaArVRs","long_url":"https://example.com"}
```
- Inspect redirect headers (HEAD):
```sh
curl -I http://localhost:8080/<shortCode>
# -> HTTP/1.1 302 Found, Location: https://example.com
```
- Follow redirect (GET):
```sh
curl -L http://localhost:8080/<shortCode>
```
Postman:
- Create environment variable `base_url = http://localhost:8080`.
- POST {{base_url}}/api/shorten with JSON body `{ "url": "https://example.com" }`.
- For redirect request GET {{base_url}}/{{short_code}} — disable auto-follow redirects to inspect the 302 Location header.

Behavior summary
- POST /api/shorten: returns JSON `{ "short_url": "...", "long_url": "..." }`. Implemented in [`handler.Handler.ShortenURL`](internals/handler/handler.go) and uses [`service.URLService.ShortenURL`](internals/service/service.go).
- GET /{shortCode}: returns HTTP 302 with Location header on success. Implemented in [`handler.Handler.RedirectURL`](internals/handler/handler.go) and resolves via [`service.URLService.GetLongURL`](internals/service/service.go).
- Storage is in-memory via [`storage.NewMemoryStorage`](internals/storage/memory.go) — restarting the app clears stored mappings.

