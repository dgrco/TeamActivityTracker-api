# TeamActivityTracker API

## Get Started
### Prerequisites
* PostgreSQL server running
* Go

### Build
```bash
go build ./cmd/server/
```

### Run
* Before running the binary, you must have certain environment variables set (see `.env.example`)
    - **IMPORTANT:** Use a secure generator for your JWT secret
        - **OpenSSL** (Recommended): `openssl rand -base64 <n>` where `<n>` should be 32 (bytes) or more.
