# HR Policy Assistant — Frontend

React + Vite boilerplate for the HR Policy Assistant UI.

## Stack

- React 19
- Vite 8
- Dev proxy to Go API on `http://localhost:8081`

## Project structure

```
frontend/
├── public/
├── src/
│   ├── components/layout/   # Header, Footer, Layout
│   ├── pages/               # Route pages
│   ├── services/            # API client
│   ├── App.jsx
│   └── main.jsx
├── vite.config.js
└── package.json
```

## Setup

```bash
cd frontend
npm install
cp .env.example .env
```

## Development

Start the Go API first (from repo root):

```bash
go run ./cmd/hrpolicy
```

Then start the frontend:

```bash
cd frontend
npm run dev
```

Open [http://localhost:5173](http://localhost:5173).

## Scripts

| Command         | Description              |
| --------------- | ------------------------ |
| `npm run dev`   | Start dev server         |
| `npm run build` | Production build         |
| `npm run preview` | Preview production build |
| `npm run lint`  | Run oxlint               |

## API proxy

During development, Vite proxies:

- `/api/*` → `http://localhost:8081`
- `/health` → `http://localhost:8081`

Use `src/services/api.js` for backend calls.
