# Day 3 exercise starter — `urlshorten`

Empty scaffold — `go.mod` + `main.go` เท่านั้น ที่เหลือเขียนเองทั้งหมด ทีละ step
(ดู `AGENTS.md`) ตาม `day3-rest-api-clean-arch-progressive.md`

## Roadmap — ทำทีละ step ห้ามข้าม

| Step | เนื้อหา | Endpoint / deliverable |
|---|---|---|
| 2 | `net/http` handler แรก | `GET /ping`, `GET /version`, `GET /shorten/{code}` |
| 2 | in-memory shorten | `POST /shorten`, `GET /{code}` (302), 404 |
| 3 | ย้ายไป Gin | endpoint เดิม, behavior เดิม |
| 4 | persist ลง SQLite (`database/sql` + `modernc.org/sqlite`) | `SaveLink`, `FindLink` |
| 5 | `.env` config (`godotenv`) | `POST /shorten` คืน full URL (`BASE_URL` + code) |
| 6 | JWT + `AuthMiddleware` | `POST /shorten` ต้อง auth, `GET /:code` ยัง public |
| 7 | cookies (`HttpOnly`/`Secure`/`SameSite`) + XSS fix (`html/template`) | หน้า reflect input ต้อง escape |
| 8 | health probes | `GET /health`, `GET /health/ready` (ping DB จริง) |
| 9 | structured logging (`slog`) | logging middleware ทุก request |
| 10 | rate limiting (`golang.org/x/time/rate`) | `POST /shorten` เกิน burst -> `429` |
| 11 | graceful shutdown (`signal.NotifyContext` + `srv.Shutdown`) | in-flight request รอดตอน shutdown |

จบแล้ว refactor ทั้งหมดเข้า Clean Architecture (ดูท้าย deck):

```
cmd/urlshorten/main.go   <- DI wiring
link/link.go             <- Entity + Repository interface
link/service.go          <- Use Case
sqlite/link.go           <- Repository implementation
httpapi/handler.go       <- Gin handlers
```

## Cross Functional Requirements checklist

- [ ] configuration in `.env` file
- [ ] graceful shutdown
- [ ] liveness and readiness probes
- [ ] structured logging
- [ ] rate limiting

## Dependencies

ไม่ได้ pre-install ไว้ — `go get` เพิ่มเองตอนถึง step ที่ต้องใช้:

```sh
go get github.com/gin-gonic/gin
go get modernc.org/sqlite
go get github.com/joho/godotenv
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/time/rate
```

## Workshop

หลังจบ `urlshorten` ทำ Workshop: **Product Catalog API** — เอา pattern เดียวกันทั้งหมด
ไปใช้กับ domain ใหม่ (เริ่ม `mkdir` โปรเจกต์ใหม่เองเหมือน Step 1) — ดูท้าย deck

```
go build ./...
go test ./...
```
