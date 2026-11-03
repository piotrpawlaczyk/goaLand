# goaLand

A football results and datas application for the Premier League, La Liga, Bundesliga, Serie A and Ligue 1.

The Go backend imports data from football-data.org, stores it in PostgreSQL and exposes a REST API for the React frontend. League standings are calculated from match results and kept in memory for API requests. Scheduled jobs update the stored data independently of requests from the browser.

## Features

- Upcoming fixtures, match results and live match indicators.
- League tables with points, goal statistics and recent form.
- Top scorers with goal and assist counts.
- Club pages with fixtures, results, standings and squad lists.
- Season selection for imported editions of a competition.
- Separate commands for initial data import and normal server startup.

The initial import includes the seasons from 29024 to current in top five football leagues. The interface is in Polish.

## Technology

| Area   |   Tools |
| ---    |   ---   |
| Backend | Go, `net/http`, Chi |
| Database | PostgreSQL 17, GORM |
| Background updates | robfig/cron |
| Frontend | React, TypeScript, Vite, Tailwind CSS |
| Tests | Go testing package, Testify, GoMock |
| Local environment | Docker Compose, Make |
| Football data | [football-data.org API v4](https://www.football-data.org/documentation/quickstart) |

## Backend

- Scheduled synchronization every five minutes, staggered across leagues.
- External API client with request timeouts and API rate-limit handling.
- Relational models and upserts scoped to competitions and seasons.
- Standings calculated from matches and stored in memory, protected by an `RWMutex`.

## Run locally

Requires Docker Compose v2, Node.js 22.12+, npm, Make and a Bash-compatible shell. Use WSL on Windows.

### 1. Configure the environment

```bash
git clone https://github.com/piotrpawlaczyk/goaLand
cd goaLand
make env
```

Set `FOOTBALL_DATA_API_TOKEN` in the root `.env` to your own [football-data.org token](https://www.football-data.org/client/register). It needs access to the imported competitions and seasons. Keep the other default settings for this setup.

### 2. Start the backend and database

```bash
docker compose up -d postgres
docker compose run --build --rm backend ./main --init-data
docker compose up -d backend
```

Wait for the import to finish before running the last command. Initialization creates missing tables and imports data without dropping the database. It can take several minutes due to API request limits. On later starts, only `docker compose up -d backend` is needed.

### 3. Start the frontend

In another terminal, from the repository root:

```bash
cd frontend
npm ci
npm run dev
```

Open **[localhost:5173](http://localhost:5173)**. The API is available at `http://localhost:8080`; check it with:

```bash
curl -f http://localhost:8080/health
```

Keep `VITE_API_URL=/api` in `frontend/.env`. Vite forwards `/api` requests to the backend at `http://localhost:8080`, both locally and through the demo tunnel.

## Share the demo

Install [`cloudflared`](https://developers.cloudflare.com/cloudflare-one/networks/connectors/cloudflare-tunnel/downloads/). On Windows, run it in the same WSL environment as the frontend. [Cloudflare Quick Tunnels](https://developers.cloudflare.com/tunnel/get-started/quick-tunnels/) provide a temporary public HTTPS address without an account or domain.

Keep the backend and database running. In a separate terminal, from the repository root, build the frontend:

```bash
cd frontend
npm run build
```

Start the tunnel in a separate terminal and keep it running:

```bash
cloudflared tunnel --url http://127.0.0.1:4173
```

Copy the generated `https://...trycloudflare.com` address and keep the tunnel running in this terminal. In another terminal, open the frontend directory and start the preview:
```bash
__VITE_ADDITIONAL_SERVER_ALLOWED_HOSTS=your-tunnel.trycloudflare.com \
npm run preview -- --strictPort
```

Replace `your-tunnel.trycloudflare.com` with the generated hostname, without `https://`. The preview uses the same `/api` proxy as the development server, so one tunnel serves both the frontend and API.

Open the public link on another computer. Check that competitions, matches and standings load. The preview is also available locally at `http://localhost:4173`.

The computer must stay online with the backend, frontend preview and tunnel running. Stopping the tunnel disconnects the public link; restarting it generates a new address. Rebuild the frontend after changing its code, assets or `.env` settings.

## Tests

Requires Go 1.25+. Run the commands below from the repository root.

### Unit tests

Unit tests cover synchronization with mocked dependencies, standings calculations and featured-match selection. They do not need PostgreSQL or live API access. Dummy API settings are required by the configuration test:

```bash
source ./scripts/load_env.sh
make test
```

### Integration tests

Integration tests check database operations, API requests and data initialization. They require PostgreSQL and a real football-data.org token with access to the seasons used by the tests. Set the token in the root `.env`.

**These tests drop and recreate tables. Use a separate test database.** Start PostgreSQL and create `goaland_test` once:

```bash
docker compose up -d postgres
docker compose exec postgres sh -c 'createdb -U "$POSTGRES_USER" goaland_test'
```

Load your configuration and run the unit and integration tests:

```bash
source  ./scripts/load_env.sh
POSTGRES_NAME=goaland_test make all-test
```


## Design

Initial UI designs created in Figma. The final application includes some changes.

[View the Figma designs](https://www.figma.com/design/mA9OdTxdDkIrBs33yh8OQI/goaLand?node-id=46-280&t=sGZUpzW3S0V4ddKL-1)

[Home page](docs/design/home-page.png)

[View all design screenshots](docs/design/)