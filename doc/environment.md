# Environment Variables Reference (`.env`)

TicketApi loads environment variables automatically on startup via `github.com/joho/godotenv/autoload` (`cmd/api/main.go`).

Secrets and runtime environment flags belong in `.env`. Non-sensitive structural settings belong in `config.yaml`.

---

## 1. Security Rules for Enterprise Projects

1. **`.env` must never be committed to Git**: Ensure `.env` is listed in `.gitignore`.
2. **`.env.example` contains only placeholders**: Standard industry practice. Contains dummy/template keys with no production passwords or production API tokens.
3. **`JWT_SECRET` must be generated randomly**: Minimum 32 characters in production.
4. **`GIN_MODE` must be `release` in production**: Leaving `GIN_MODE=debug` leaks CAPTCHA answers in API responses and returns raw internal error stack traces.

---

## 2. Variables Reference

| Variable | Required | Default / Fallback | Allowed Values | Usage & Effects |
| :--- | :--- | :--- | :--- | :--- |
| `GIN_MODE` | No | `debug` | `debug`, `release`, `test` | Controls framework mode.<br>• `release`: Production mode. Strips `answer` from `/captcha/GetCaptcha/`, hides detailed internal error messages.<br>• `debug`: Development mode. Returns `answer` in CAPTCHA response, exposes full error traces. |
| `JWT_SECRET` | **Yes** (in prod) | `""` (empty string) | String (>= 32 chars recommended) | Secret key used to sign and verify HMAC-SHA256 JWT tokens (`AuthToken`, `CaptchaToken`, `OneTimeToken`). |
| `MONGODB_URI` | Yes (if `mongo.enable: true`) | `""` | Valid MongoDB connection URI | Connection string (e.g. `mongodb://user:password@localhost:27017`). |
| `MONGODB_DB` | No | `config.Get().Mongo.DBName` (`"ticket_db"`) | Database name string | Overrides target database name in MongoDB. |
| `REDIS_PASSWORD` | No | `""` | String | Redis AUTH password. Leave empty if Redis runs without password. |
| `ACCESS_KEY_MINIO` | Yes (if `minio.enable: true`) | `""` | String | MinIO / S3 access key (username). |
| `SECRET_KEY_MINIO` | Yes (if `minio.enable: true`) | `""` | String | MinIO / S3 secret key (password). |

---

## 3. Quick Setup

Copy example template to create local `.env`:

```bash
cp .env.example .env
```

Edit `.env` values for your local or production environment.
