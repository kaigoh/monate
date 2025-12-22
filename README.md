# Monate

Monate is a GraphQL backend for a Monero POS / donation experience powered by
[MoneroPay](https://moneropay.eu/). It stores stores/invoices locally (via GORM)
so Houdini/SvelteKit front-ends can control layout, theming, and live payment
updates through GraphQL subscriptions.

## Features

- GraphQL schema for stores, invoices, and a `invoiceStatus` subscription that
  streams webhook/manual polling updates.
- Store-level theme system that keeps end-user customisation in the database.
- MoneroPay client + webhook handler that forward status transitions back to
  the GraphQL layer.
- Manual fallback worker that polls MoneroPay every minute to catch missed
  callbacks (fully configurable with env vars).

## Configuration

| Variable | Description | Default |
| --- | --- | --- |
| `PORT` | HTTP port for the server | `8080` |
| `MONATE_DB_DRIVER` | Database driver (`sqlite` supported today) | `sqlite` |
| `MONATE_DB_DSN` | GORM DSN. For sqlite this is the filename. | `monate.db` |
| `MONATE_PUBLIC_BASE_URL` | Public URL used to build webhook callback URLs. | `http://localhost:${PORT}` |
| `MONATE_MANUAL_CHECK_INTERVAL` | How frequently to run the MoneroPay manual check worker. | `1m` |
| `MONATE_MANUAL_CHECK_BATCH` | Number of invoices to poll per interval. | `25` |

Each store record stores the MoneroPay API endpoint/API key pair that should be
used when creating invoices.

## Running the API

```bash
# Build/embed the UI once (requires Bun)
pushd gui
bun install
bun run generate
bun run build
popd
go generate ./web
# Run the server (defaults to sqlite database)
go run ./...
```

> Requires [Bun](https://bun.sh/) v1.1+ for UI builds.

The GraphQL Playground lives at http://localhost:8080/ and the HTTP endpoint is
`POST /query`. The Houdini/SvelteKit UI is served from `/`, and
you can run a hot-reload UI dev server (Vite) and let Go proxy to it by setting
`MONATE_UI_DEV_SERVER_URL=http://localhost:5173`.

Webhooks are exposed at
`POST /webhooks/moneropay/:secret?invoice=<uuid>`; this URL is generated
automatically for each invoice so you only need to expose the backend publicly.

If you need to initialise a clean database or run migrations from scratch, just
delete the sqlite file (defaults to `monate.db`) before booting again.

## Docker

A production-like stack (Monate, MoneroPay, and `monero-wallet-rpc`) is defined
in `docker-compose.yml`. Before starting it, copy the environment template and
fill in the secrets/endpoints:

```bash
cp .env.example .env
# edit .env with your monerod endpoint, wallet credentials, etc.
docker compose build monate
docker compose up
```

`MONEROD_RPC_ENDPOINT` must point at your external monerod node since the stack
does not run one locally. The `.env` file also defines the wallet RPC bind port,
MoneroPay API token, and Monate’s base URL so all services share the same
configuration source.

## Testing

Run the unit test suite with:

```bash
GOCACHE=$(pwd)/.gocache go test ./...
```

On sandboxes that block writes to `$HOME/.cache`, the `GOCACHE` override keeps
the build cache inside the repo.

## Next Steps

- Wire the Houdini/SvelteKit client to the `invoiceStatus` subscription.
- Extend the schema with donation campaigns, content blocks, and per-store POS
  templates.
- Implement MoneroPay webhook signature verification once the upstream exposes
  it, or add custom HMAC headers for extra security.
