# shopxpress-app

Application monorepo for ShopXpress Lab A. Contains 3 Go microservices representing 3 patterns: BFF gateway, read-heavy products, write-heavy orders.

## Layout

```
shopxpress-app/
├── pkg/                  # shared code
│   ├── logger/           # JSON structured logger
│   └── httpx/            # /healthz + /version handlers
└── services/
    ├── gateway/          # :8080 — reverse proxy + BFF
    ├── products/         # :8081 — read-only catalog (mock)
    └── orders/           # :8082 — write transactions (in-memory → RDS later)
```

## Quick start (local)

```bash
docker compose up --build
curl localhost:8080/api/products
```
