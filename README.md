# JobGen: AI Remote Job Finder & CV Optimizer

> Empowering African (starting with Ethiopian) professionals to polish their CVs and discover high‑quality remote tech roles worldwide through an AI chat experience.  
> Stack: Go (backend APIs + orchestration), React/TypeScript (web frontend), optional mobile wrapper later. GenAI for CV critique, rewriting, job summarization & matching.

---

## Table of Contents
1. Vision & Problem
2. Core Features
3. Architecture Overview
4. Tech Stack & Rationale
5. High-Level Flow
6. Directory Structure (Proposed)
7. Data & Domain Models
8. API Design (v1 Draft)
9. AI Prompting Strategy
10. Matching Algorithm (MVP vs Enhanced)
11. Resume Text Extraction & Parsing
12. Frontend UX Notes
13. Environment Variables & Configuration
14. Local Development Setup
15. Running the Stack (Dev)
16. Testing Strategy
17. Logging, Monitoring & Observability
18. Security & Privacy
19. Performance & Scalability Considerations
20. Deployment Strategy
21. Roadmap (3 Weeks + Beyond)
22. Risks & Mitigations
23. Contribution Guidelines
24. License
25. Acknowledgements / Inspiration

---

## 1. Vision & Problem
Youth unemployment and underemployment are high in Ethiopia and across Africa. Talented graduates and self‑taught developers struggle to:
- Discover relevant, legitimate remote job opportunities.
- Optimize CVs for international applicant tracking systems (ATS).
- Understand how their experience maps to global market expectations.

**JobGen** acts as an AI career co‑pilot:
1. Ingests a CV (file or text) and produces targeted, actionable improvement suggestions.
2. Lets the user iteratively refine bullet points or sections via conversational rewrites.
3. Searches aggregated remote job sources (e.g., RemoteOK, curated feeds / APIs), ranks & summarizes top matches.
4. (Future) Offers interview prep, skill gap insights, and upskilling suggestions.

---

## 2. Core Features

| Feature | Description | Status (MVP Plan) |
|---------|-------------|-------------------|
| CV Upload & Analysis | Upload PDF/DOCX or paste text; AI returns prioritized suggestions | Week 1 |
| Interactive Rewrites | User asks for rewrites (“rewrite my database bullet with metrics”) | Week 2 |
| Job Search & Match | Query remote job APIs / feeds, compute match %, summarize | Week 2 |
| Match Explanation | Show why a job matches / missing keywords | Week 2–3 |
| Session Memory | Retain parsed CV & improvements for subsequent queries | Week 2 |
| Privacy Controls | Ephemeral storage / user “Clear Data” action | Week 3 |
| Interview Prep (Optional) | Generate probable questions per job | Stretch |
| Basic Analytics | Count analyses, job searches (internal metrics) | Week 3 |
| Localization (Light) | Simple explanatory English (Amharic help TBD) | Post-MVP |

---

## 3. Architecture Overview

```
              ┌────────────────────────────┐
              │          Frontend          │
              │ React/TypeScript (Web App) │
              │ - Chat UI                  │
              │ - CV Upload                │
              │ - Job Cards                │
              └─────────────┬──────────────┘
                            HTTPS (JSON/REST)
┌──────────────────────────────────────────────────────────────┐
│                        Go Backend                           │
│  cmd/api/main.go                                             │
│  ├── Auth (simple session / token)                           │
│  ├── CV Service:                                             │
│  │     - File ingestion & text extract                       │
│  │     - AI analysis & rewrite orchestration                 │
│  ├── Job Service:                                            │
│  │     - RemoteOK / future sources fetch                     │
│  │     - Keyword scoring + (optional) AI scoring             │
│  ├── Conversation / Context Store                            │
│  ├── Prompt Manager (templates)                              │
│  ├── Rate Limiter                                            │
│  └── Logging / Metrics                                       │
│                  │
│          External Services
│          ├── GenAI Provider (OpenAI, etc.)
│          ├── Job APIs / Feeds (RemoteOK JSON, RSS)
│          └── (Optional) PDF text extraction libs
└──────────────────────────────────────────────────────────────┘

           Storage Layer (MVP)
           - In-memory / Redis (sessions, CV cache)
           - (Future) Postgres for persistence
```

---

## 4. Tech Stack & Rationale
- **Backend:** Go (concurrency, fast IO for external API calls, strong ecosystem).
- **Frontend:** React + TypeScript (developer velocity, rich component ecosystem).
- **AI Integration:** OpenAI GPT-4/4o or fallback to a provider via a pluggable interface.
- **File Parsing:** Go packages (e.g., `rsc.io/pdf`) or sending extracted text via frontend using JS libs if simpler; DOCX via `unidoc` or convert client-side.
- **Job Sources (MVP):** RemoteOK (public JSON), curated static JSON fallback. Add more sources later.
- **Session / Cache:** In-memory map -> switch to Redis for horizontal scaling.
- **Deployment:** Docker containers; optional Compose for dev; future K8s cluster or Fly.io/Render for simplicity.

---

## 5. High-Level Flow

1. User uploads CV -> frontend sends raw text (or file) to `/api/v1/cv/analyze`.
2. Backend extracts text (if file) -> constructs prompt -> AI returns structured suggestions (prompt instructs numbered list).
3. Suggestions stored in session (ID in cookie or header token).
4. User requests rewrite -> `/api/v1/cv/rewrite` with snippet + guidance.
5. User clicks “Find Jobs” -> `/api/v1/jobs/search` with query (or derived from CV).
6. Backend fetches job feed, scores matches (keyword overlap + optional AI refinement), returns top N.
7. User optionally asks follow-up (“why is match low?”) -> backend compares missing keywords.
8. User can clear session -> remove CV data.

---

## 6. Directory Structure (Proposed)

```
jobgen/
├─ backend/
│  ├─ cmd/
│  │  └─ api/
│  │      └─ main.go
│  ├─ internal/
│  │  ├─ http/
│  │  │  ├─ middleware/
│  │  │  └─ handlers/
│  │  ├─ cv/
│  │  ├─ jobs/
│  │  ├─ ai/
│  │  ├─ matching/
│  │  ├─ session/
│  │  ├─ config/
│  │  ├─ logging/
│  │  └─ util/
│  ├─ pkg/
│  ├─ go.mod
│  └─ go.sum
├─ frontend/
│  ├─ src/
│  │  ├─ components/
│  │  ├─ pages/
│  │  ├─ hooks/
│  │  ├─ services/
│  │  ├─ context/
│  │  ├─ utils/
│  │  └─ types/
│  ├─ public/
│  ├─ package.json
│  └─ tsconfig.json
├─ deploy/
│  ├─ docker/
│  │  ├─ backend.Dockerfile
│  │  └─ frontend.Dockerfile
│  ├─ docker-compose.yml
│  └─ k8s/
├─ docs/
│  ├─ prompts/
│  └─ api/
└─ README.md
```

---

## 7. Data & Domain Models (Draft)

### Session
```go
type Session struct {
    ID            string
    CreatedAt     time.Time
    CVText        string
    ExtractedSkills []string
    Suggestions   []CVSuggestion
    LastJobQuery  string
}
```

### CVSuggestion
```go
type CVSuggestion struct {
    ID          string
    Category    string
    Text        string
    Priority    int
    Tags        []string
}
```

### JobListing
```go
type JobListing struct {
    ID           string
    Source       string
    Title        string
    Company      string
    Location     string
    Remote       bool
    URL          string
    Description  string
    RawJSON      []byte
    RetrievedAt  time.Time
}
```

### MatchResult
```go
type MatchResult struct {
    JobID         string
    MatchScore    float64
    MatchedSkills []string
    MissingSkills []string
    AIExplanation string
}
```

---

## 8. API Design (v1)

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | /api/v1/cv/analyze | Analyze CV (multipart file or JSON {text}) | Session cookie |
| POST | /api/v1/cv/rewrite | Rewrite snippet `{original, instruction}` | Session |
| GET  | /api/v1/cv/skills | Return extracted skills | Session |
| POST | /api/v1/jobs/search | Body: `{ query?: string, limit?: number }` | Session |
| GET  | /api/v1/jobs/last | Return last job results (cached) | Session |
| POST | /api/v1/session/clear | Clear stored data | Session |
| GET  | /healthz | Liveness | Public |
| GET  | /readyz | Readiness (checks AI & job source connectivity) | Public |

### Request: Analyze CV
```
POST /api/v1/cv/analyze
Content-Type: multipart/form-data OR application/json

multipart: file=<cv.pdf>
json: {"text": "paste of CV ..."}
```

### Response:
```json
{
  "sessionId": "abc123",
  "suggestions": [
    {"id":"s1","category":"Summary","priority":1,"text":"Add a 2-3 line professional summary.","tags":["summary","impact"]},
    {"id":"s2","category":"Experience","priority":1,"text":"Quantify database optimization results.","tags":["metrics","experience"]}
  ],
  "extractedSkills": ["Python","React","SQL","Docker"]
}
```

---

## 9. AI Prompting Strategy

Store prompt templates under `docs/prompts/` and load at startup.

CV Analysis Prompt:
```
You are a career coach analyzing an African tech professional's CV.
Return a JSON array of suggestions. Each suggestion:
- category (Summary, Experience, Skills, Education, Formatting, Keywords)
- priority (1=High,2=Medium,3=Low)
- text (actionable, specific)
- tags (array of lowercase keywords)
CV TEXT:
{{CV_TEXT}}
```

Rewrite Prompt:
```
Rewrite the resume bullet for stronger impact and metrics if possible.
Original Bullet:
"{{ORIGINAL}}"
User Instruction: "{{INSTRUCTION}}"
Return ONLY improved bullet(s).
```

Job Match Explanation Prompt:
```
Given CV (short summary + skills) and job description, rate match 0-100 and list matched_skills & missing_skills.
Return JSON {score, matched_skills, missing_skills, explanation}
```

Token Optimization:
- Truncate CV to max characters.
- Extract skills first (regex / dictionary), feed summarized skill list into subsequent prompts.

---

## 10. Matching Algorithm

MVP:
1. Extract candidate skills (normalized).
2. For each job:
   - Overlap count (weighted).
   - BaseScore = (matched / totalCandidateSkills) * 70.
   - Title keyword bonus.
   - Cap at 95.

Enhanced:
- AI scoring for top 10.
- Seniority alignment.
- Negative penalties for missing critical explicit skills.

---

## 11. Resume Text Extraction & Parsing

- PDF: try Go library; fallback ask for paste.
- Section detection via regex headings.
- Skill extraction: curated dictionary + pattern capture.
- Optionally AI refine skill list.

---

## 12. Frontend UX Notes

Components:
- ChatWindow
- CVUploader
- SuggestionList (group by category)
- JobCard (title/company/match bar)
- RewriteModal

Flow:
1. Prompt to upload.
2. Show suggestions with quick actions:
   - [Find Jobs] [Rewrite a Bullet] [Show Extracted Skills]
3. Loading indicators & accessible design.

---

## 13. Environment Variables

Backend `.env.example`:
```
PORT=8080
AI_PROVIDER=openai
OPENAI_API_KEY=sk-...
JOB_SOURCES=remoteok
REMOTEOK_API_URL=https://remoteok.com/api
SESSION_TTL_MINUTES=60
MAX_CV_CHARS=20000
LOG_LEVEL=info
RATE_LIMIT_RPS=3
ALLOWED_ORIGINS=http://localhost:3000
```

Frontend `.env`:
```
REACT_APP_API_BASE=http://localhost:8080
REACT_APP_BUILD_ENV=development
```

---

## 14. Local Development Setup

Prerequisites: Go 1.22+, Node 20+, Docker (optional).

```
git clone https://github.com/naolaboma/JobGen.git
cd JobGen
cp backend/.env.example backend/.env
# Add AI key
cd backend && go mod tidy && go run cmd/api/main.go
# new terminal
cd frontend && npm install && npm run dev
```

Visit `http://localhost:3000`.

---

## 15. Running the Stack (Docker)

```
docker compose up --build
```

Services: backend, frontend (and future redis).

---

## 16. Testing Strategy

Backend:
- Unit: skill extraction, scoring, prompt builder.
- Integration: analyze & job search with AI mock.
Frontend:
- Component tests (React Testing Library).
- E2E (Playwright) main flow.
CI: GitHub Actions; run `go test`, `npm test`, lint.

---

## 17. Logging, Monitoring & Observability

- Structured JSON logs (`zerolog` or `zap`).
- Correlation IDs middleware.
- Prometheus metrics endpoint (future).
- Key metrics: AI latency, job fetch count, match compute time.

---

## 18. Security & Privacy

- Ephemeral CV storage in memory.
- Clear Data endpoint.
- No full CV in logs.
- Rate limit & CORS restrictions.
- Future: encryption at rest if persistence added.

---

## 19. Performance & Scalability

- Cache job feeds (TTL).
- Parallel scoring with goroutines.
- AI scoring only for top subset.
- Horizontal scaling with shared Redis session store.

---

## 20. Deployment Strategy

- Phase 1: Single container on Render/Fly.io.
- GH Actions build & push Docker images (GHCR).
- Domain + TLS (Caddy/NGINX).
- Future: K8s manifests.

---

## 21. Roadmap

Week 1: CV analysis, suggestions, skill extraction.  
Week 2: Rewrites, job search integration, scoring, chat polish.  
Week 3: Match explanations, privacy polish, tests, performance tune.  
Post-MVP: Interview prep, more sources, persistent profiles, CV export, skill gap analysis, localization.

---

## 22. Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| AI cost | Budget | Rate limit, cache |
| PDF parse failures | Poor suggestions | Paste fallback |
| Single job source | Limited results | Add sources, cached fallback |
| Generic advice | Low value | Prompt tuning, require specifics |
| Latency | UX degrade | Parallel & cache |
| Privacy concern | Trust | Ephemeral storage & disclosure |

---

## 23. Contribution Guidelines

1. Branch naming: `feature/<slug>`, `fix/<slug>`.
2. Conventional commits.
3. Run tests + lint before PR.
4. PR Template (description, screenshots, tests).
5. Go: `gofmt`, `golangci-lint`; TS: ESLint + Prettier.
6. No secrets committed.
7. Document new endpoints.

---

## 24. License

```
MIT License
```

---

## 25. Acknowledgements / Inspiration
- RemoteOK
- OpenAI & LLM ecosystem
- African tech talent communities (A2SV, local hubs)

---

## Quick Start (TL;DR)

Backend:
```
cd backend
cp .env.example .env
go run cmd/api/main.go
```

Frontend:
```
cd frontend
npm install
npm run dev
```

Navigate to `http://localhost:3000`.

---

## Status Disclaimer
This README is forward-looking; update sections as actual implementations are completed.

---

Happy building! 🎉
