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

## การใช้งานและการทดสอบ (Usage & Testing)

### การตั้งค่า Environment
ก๊อปปี้ไฟล์ `.env.example` ไปเป็น `.env` และตั้งค่าพอร์ตและตัวแปรต่างๆ:
```sh
cp .env.example .env
```

### วิธีการรัน Server
```sh
go run main.go
```

### วิธีการรัน Unit Test
```sh
go test -v ./...
```

### รายละเอียดการเพิ่ม Unit Test
มีการเพิ่ม Unit Test เพื่อช่วยยืนยันความถูกต้องของ Handler และ Utility Logic (โดยที่ไม่ต้องต่อฐานข้อมูลจริง):
1. **[short_test.go](file:///Users/pallat/workspace/gitlab.com/gophernment/fundamental/day3-starter/shorten/short_test.go)**: ทดสอบว่ารหัสย่อ URL (`shortenURL`) ถูกเจนขึ้นมาได้ถูกต้อง และมีความเป็น Unique (ไม่ซ้ำกันง่ายๆ จากการสุ่ม)
2. **[shorten_test.go](file:///Users/pallat/workspace/gitlab.com/gophernment/fundamental/day3-starter/shorten/shorten_test.go)**: ทดสอบการทำงานของ `Handler.Shorten` โดยการใช้ Mocking สำหรับ interfaces `shorter` และ `storer` เพื่อจำลองการบันทึกข้อมูลและตรวจสอบการตอบสนองของ HTTP responses ในกรณีต่างๆ เช่น:
   - ตรวจสอบ HTTP Method ที่ไม่อนุญาต (Method Not Allowed)
   - กรณีรันสำเร็จ (Success path) และตรวจสอบการเซฟค่า
   - กรณี Shorter/Storer ทำงานผิดพลาด (Error handling)

### รายละเอียด .gitignore
มีการเพิ่มไฟล์ `.gitignore` เพื่อแยกแยะไฟล์ที่ไม่ควรนำขึ้นระบบ Git:
- Binaries และ Executables ต่างๆ
- ไฟล์เก็บความลับเฉพาะเครื่อง เช่น `.env` (แต่ยกเว้น `.env.example`)
- ไฟล์ SQLite Database (`*.db`, `*.db-journal`, `*.db-shm`, `*.db-wal`) เพื่อป้องกันข้อมูลทับซ้อน
- การตั้งค่าเฉพาะของ IDEs (VS Code, GoLand, .idea, ฯลฯ)

