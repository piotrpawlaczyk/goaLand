# goaLand

A football results and statistics application for the Premier League, La Liga, Bundesliga, Serie A and Ligue 1.

The Go backend imports data from football-data.org, stores it in PostgreSQL and exposes a REST API for the React frontend. League standings are calculated from match results and kept in memory for API requests. Scheduled jobs update the stored data independently of requests from the browser.

## Features

- Upcoming fixtures, match results and live match indicators.
- League tables with points, goal statistics and recent form.
- Top scorers with goal and assist counts.
- Club pages with fixtures, results, standings and squad lists.
- Season selection for imported editions of a competition.
- Separate commands for initial data import and normal server startup.

The initial import includes the current season of all five leagues and the 2025/26 Premier League season. The interface is in Polish.

## Technology

| Area | Tools |
| --- | --- |
| Backend | Go, `net/http`, Chi |
| Database | PostgreSQL 17, GORM |
| Background updates | robfig/cron |
| Frontend | React, TypeScript, Vite, Tailwind CSS |
| Tests | Go testing package, Testify, GoMock |
| Local environment | Docker Compose, Make |
| Football data | [football-data.org API v4](https://www.football-data.org/documentation/quickstart) |

## Backend

- Scheduled synchronization every five minutes, staggered across leagues.
- External API client with HTTP timeouts and request throttling.
- Relational models and upserts scoped to competitions and seasons.
- Standings calculated from matches and stored in memory, protected by an `RWMutex`.
