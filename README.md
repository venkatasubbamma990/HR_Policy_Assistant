# HR Policy Assistant

A Go-based Retrieval-Augmented Generation (RAG) application for answering employee questions using company HR policy documents. Policies are stored as Markdown files, parsed into structured metadata and sections, and prepared for embedding and semantic search.

Built for **NovaTech Solutions** policy documents including leave, salary, WFH, onboarding, exit, and notice period policies.

---

## Features

- **Document ingestion pipeline** — loads and parses Markdown policy files from `documents/`
- **Structured metadata extraction** — `source`, `policy_type`, `doc_id`, `version`, `effective_date`
- **Markdown parsing** — headings, tables, lists, and paragraphs
- **Structured logging** — [Uber Zap](https://github.com/uber-go/zap) with JSON or console output
- **Docker support** — multi-stage build and Docker Compose for containerized runs
- **Standalone ingest command** — re-index policies when files are added or updated

---

## Architecture

```
documents/*.md
      │
      ▼
┌─────────────────┐
│  Ingest Pipeline │  Load → Parse metadata → Parse markdown
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  PolicyDocument  │  Structured docs with sections, tables, lists
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   RAG Engine     │  Embeddings + retrieval + LLM (planned)
└─────────────────┘
```

### Project layout

```
HRPolicyAssistant/
├── cmd/
│   ├── hrpolicy/          # Main application entry point
│   └── ingest/            # Standalone document ingestion CLI
├── internal/
│   ├── config/            # Environment-based configuration
│   ├── documents/         # Loader, metadata & markdown parser
│   ├── ingest/            # Indexing pipeline
│   ├── logger/            # Zap logger setup
│   └── rag/               # RAG engine (query handling planned)
├── documents/             # HR policy Markdown files
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── go.mod
```

---

## Prerequisites

- **Go 1.23+**
- **Docker & Docker Compose** (optional, for containerized runs)
- **Make** (optional, for convenience targets)
- **OpenAI API key** (optional for now; required when LLM/embeddings are added)

---

## Quick start

### 1. Clone and configure

```bash
git clone <repository-url>
cd HRPolicyAssistant
cp .env.example .env
```

Edit `.env` and set `OPENAI_API_KEY` when you are ready to use LLM features.

### 2. Run locally

```bash
# Build the application
go build -o bin/hrpolicy ./cmd/hrpolicy

# Start the assistant (loads and indexes policies on startup)
./bin/hrpolicy
```

### 3. Ingest documents only

Run this when you add or update policy files:

```bash
go run ./cmd/ingest
```

Or with verbose section-level debug logs:

```bash
LOG_LEVEL=debug HR_INGEST_VERBOSE=1 go run ./cmd/ingest
```

---

## Docker

```bash
# Build and start in background
make up

# View logs
docker compose logs -f

# Stop containers
make down
```

Documents are mounted read-only from `./documents`, so you can update policies without rebuilding the image.

---

## Makefile targets

| Target | Description |
|--------|-------------|
| `make build` | Build local binary → `bin/hrpolicy` |
| `make ingest` | Run document ingestion pipeline |
| `make test` | Run all Go tests |
| `make up` | Build Docker image and start container |
| `make down` | Stop and remove containers |
| `make clean` | Remove binary and local Docker artifacts |

---

## Configuration

All settings are loaded from environment variables (or a `.env` file for Docker Compose).

| Variable | Default | Description |
|----------|---------|-------------|
| `HR_DOCUMENTS_DIR` | `./documents` | Directory containing policy Markdown files |
| `OPENAI_API_KEY` | — | OpenAI API key for embeddings and chat |
| `HR_EMBED_MODEL` | `text-embedding-3-small` | Embedding model name |
| `HR_CHAT_MODEL` | `gpt-4o-mini` | Chat/completion model name |
| `LOG_LEVEL` | `info` | Log level: `debug`, `info`, `warn`, `error` |
| `LOG_FORMAT` | `console` | Log format: `console` (local) or `json` (Docker) |
| `HR_INGEST_VERBOSE` | `0` | Set to `1` for per-section debug logs during ingest |

---

## Document ingestion

Each policy file in `documents/` is parsed with the following metadata:

| Field | Source | Example |
|-------|--------|---------|
| `source` | Filename | `leave-policy.md` |
| `policy_type` | Derived from filename | `leave`, `salary`, `wfh` |
| `doc_id` | Document header | `HR-POL-LVE-002` |
| `version` | Document header | `4.1` |
| `effective_date` | Document header | `April 1, 2025` |

Policy files should follow this header format:

```markdown
# Leave Policy

**Company:** NovaTech Solutions Pvt. Ltd.
**Document ID:** HR-POL-LVE-002
**Version:** 4.1
**Effective Date:** April 1, 2025

---
```

### Included policies

| File | Policy type | Document ID |
|------|-------------|-------------|
| `leave-policy.md` | leave | HR-POL-LVE-002 |
| `salary-policy.md` | salary | HR-POL-SAL-001 |
| `wfh-policy.md` | wfh | HR-POL-WFH-003 |
| `entry-policy.md` | entry | HR-POL-ENT-004 |
| `exit-policy.md` | exit | HR-POL-EXT-005 |
| `notice-policy.md` | notice | HR-POL-NOT-006 |

---

## Logging

The application uses structured logging with [Zap](https://github.com/uber-go/zap). Loggers are scoped by component:

- `hrpolicy` — main application
- `ingest` — ingestion pipeline
- `documents` — file loading and parsing
- `rag` — RAG engine

Example console output:

```
INFO  ingest.documents  found policy files  {"count": 6, "files": ["leave-policy.md", ...]}
INFO  ingest.ingest     indexed policy document  {"source": "leave-policy.md", "doc_id": "HR-POL-LVE-002", ...}
```

In Docker, logs default to JSON for easier aggregation.

---

## Development

```bash
# Run tests
go test ./...

# Run with debug logging
LOG_LEVEL=debug go run ./cmd/hrpolicy

# Format and vet (recommended before commits)
go fmt ./...
go vet ./...
```

---

## Roadmap

- [ ] Text chunking for embedding
- [ ] Vector store integration
- [ ] OpenAI embeddings and retrieval
- [ ] HTTP API for policy Q&A
- [ ] Source citation in answers

---

## License

Internal use — NovaTech Solutions Pvt. Ltd.
