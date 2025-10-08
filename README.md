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
│  ├─ pkg/ (shared libs if needed)
│  ├─ go.mod
│  └─ go.sum
├─ frontend/
│  ├─ src/
│  │  ├─ components/
│  │  ├─ pages/
│  │  ├─ hooks/
│  │  ├─ services/ (API clients)
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
│  └─ k8s/ (future manifests)
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
    Category    string   // e.g. "Skills", "Experience", "Summary"
    Text        string
    Priority    int      // 1=High
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
    RawJSON      []byte // optional cache
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
    AIExplanation string // optional
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

Example: CV Analysis Prompt (pseudo-template):
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

Job Match Explanation Prompt (optional per top N):
```
Given CV (short summary + skills) and job description, rate match 0-100 and list matched_skills & missing_skills.
Return JSON {score, matched_skills, missing_skills, explanation}
```

Token Optimization:
- Truncate CV to max characters.
- Extract skills first (regex / skill dictionary), feed summarized skill list into subsequent prompts instead of full CV when possible.

---

## 10. Matching Algorithm

MVP:
1. Extract candidate skills (normalized lowercase).
2. For each job description:
   - Tokenize & count overlap with skills (weighted by importance).
   - BaseScore = (matchedSkillCount / totalCandidateSkills) * 70.
   - Bonus if title keywords align (e.g., “engineer”, “python”, “full stack”).
   - Clip 0–95 (reserve >95 for AI refined).
3. Sort descending, pick top 3–5.

Enhanced (Week 3 or Post-MVP):
- Use AI scoring for top ~10 jobs only.
- Incorporate seniority detection (e.g., “Senior”, “Lead” vs candidate experience).
- Negative penalties (e.g., missing “React” for “React Developer” role).

---

## 11. Resume Text Extraction & Parsing

Approach:
- For PDF: attempt text extraction server-side with minimal dependency (Go PDF text extractor). If extraction fails or yields low density, fallback: ask user to paste text (“Couldn’t parse – please paste”).
- Section heuristics: find headings (case-insensitive) `(?m)^(experience|education|skills|projects|summary)\b`.
- Skill extraction: maintain a curated skill dictionary (JSON) plus dynamic capture (capitalized tech words). Optionally ask AI to refine.

---

## 12. Frontend UX Notes

Components:
- `ChatWindow`: messages with role (user|ai|system).
- `CVUploader`: file + paste area.
- `SuggestionList`: collapsible categories.
- `JobCard`: title, company, match bar (progress), missing skill badges (grey), matched skill badges (green).
- `RewriteModal`: show original, new suggestion preview.

Conversation Pattern:
1. System hint: “Upload your CV to begin.”
2. After analysis: quick actions:
   - Buttons: [Find Jobs] [Rewrite a Bullet] [Show Extracted Skills]
3. Loading states: skeleton for suggestions, spinner for jobs.

Accessibility:
- Keyboard navigation for chat.
- Clear alt text for icons.

---

## 13. Environment Variables

Backend (`.env.example`):
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

Frontend (`.env`):
```
REACT_APP_API_BASE=http://localhost:8080
REACT_APP_BUILD_ENV=development
```

---

## 14. Local Development Setup

Prereqs:
- Go 1.22+
- Node 20+ / PNPM or Yarn
- Docker (optional)
- Make (optional)

Steps:
```
git clone https://github.com/naolaboma/JobGen.git
cd JobGen
cp backend/.env.example backend/.env
# Add OpenAI key
cd backend && go mod tidy && go run cmd/api/main.go
# In new terminal
cd frontend && npm install && npm run dev
```
Visit: `http://localhost:3000`

---

## 15. Running the Stack (Docker)

`docker-compose.yml` (planned) includes:
- `backend` (expose 8080)
- `frontend` (expose 3000, depends_on backend)
- `redis` (future)

Command:
```
docker compose up --build
```

---

## 16. Testing Strategy

Backend:
- Unit tests for:
  - Skill extraction
  - Match scoring
  - Prompt builder (deterministic segments)
- Integration tests:
  - `/cv/analyze` with fixture CV text (mock AI provider)
  - `/jobs/search` with mocked job feed
- Use an AI mock layer returning canned JSON.

Frontend:
- Component tests (React Testing Library): `SuggestionList`, `JobCard`.
- E2E (Playwright) for main flow (upload -> suggestions -> job search).

CI:
- GitHub Actions running `go test ./...` then `npm test`.
- Lint: `golangci-lint`, `eslint`, `prettier`.

---

## 17. Logging, Monitoring & Observability

- Structured JSON logs (`zap` or `zerolog`).
- Correlation ID per request (middleware).
- Metrics (future): Prometheus endpoint `/metrics`.
- Log events: AI call start/end (duration), job fetch count, rewrite requests, errors.

---

## 18. Security & Privacy

- Do not permanently persist CV text in MVP (store in memory session only).
- Provide “Clear Data” endpoint.
- Sanitize logs (never log full CV).
- Rate limit to mitigate abuse / cost runaway.
- CORS locked to known origins in prod.
- Future: OAuth (GitHub/Google) for persistent user profile.

---

## 19. Performance & Scalability

Anticipated Bottlenecks:
- AI latency (1–4s).
- Job feed fetch (network).
Mitigations:
- Cache job feed for X minutes (configurable).
- Parallel AI scoring for top matches (goroutines).
- Use streaming responses (future) for incremental display.

Scaling:
- Stateless backend instances behind load balancer.
- Redis session store.
- CDN for frontend static assets.

---

## 20. Deployment Strategy

Stages:
1. MVP: Single VM / PaaS (Render/Fly.io).
2. Container images published to GHCR.
3. Domain & TLS (Caddy / Nginx reverse proxy).
4. Add CD pipeline on main branch tag.

---

## 21. Roadmap

Week 1:
- CV upload & raw analysis endpoint.
- Basic chat UI scaffold.
- Skill extraction + suggestions display.

Week 2:
- Rewrite endpoint.
- Job search integration (RemoteOK).
- Match scoring MVP & job cards.

Week 3:
- Polishing (UX, error states).
- Optional AI match explanation.
- Privacy controls & simple analytics.
- README, docs, test coverage uplift.

Post-MVP:
- Interview prep mode.
- Additional job sources (LinkedIn API / Indeed integration if allowed).
- User accounts & saved profiles.
- ATS-friendly CV export.
- Skill gap analysis & learning path suggestions.
- Amharic helper mode.

---

## 22. Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| AI cost escalation | Budget overruns | Rate limiting, cache suggestions |
| Poor PDF extraction | Bad analysis | Fallback ask user to paste text |
| Unreliable job source | Missing results | Cache & multiple sources |
| Generic AI advice | Low perceived value | Prompt tuning; require actionable items |
| Latency | User frustration | Parallelization & caching |
| Data privacy concerns | Trust loss | Ephemeral storage + disclosure note |

---

## 23. Contribution Guidelines

1. Fork & branch naming: `feature/<slug>` or `fix/<slug>`.
2. Conventional commits (e.g., `feat(cv): add rewrite endpoint`).
3. Run tests & lints before PR.
4. PR template should include:
   - Description
   - Screenshots (UI changes)
   - Test coverage or manual test steps
5. Code Style:
   - Go: `gofmt`, `golangci-lint`.
   - TS: ESLint + Prettier.
6. Avoid committing secrets; use `.env` (in `.gitignore`).
7. Document new endpoints in `docs/api/`.

---

## 24. License

(Choose one—suggestion: MIT for broad adoption.)

Example:
```
MIT License
Copyright (c) 2025 ...
```

---

## 25. Acknowledgements / Inspiration
- RemoteOK job feed
- Global career coaching frameworks
- OpenAI / modern LLM capabilities inspiring conversational UX
- Communities empowering African tech talent (e.g., A2SV, local hubs)

---

## Quick Start (TL;DR)

Backend:
```
cd backend
cp .env.example .env  # add AI key
go run cmd/api/main.go
```
Frontend:
```
cd frontend
npm install
npm run dev
```
Visit: `http://localhost:3000`

---

## Status Disclaimer
This README reflects a forward-looking architecture and may include placeholders pending implementation (AI provider abstraction, multi-source job integration). Update it as actual components solidify.

---

Happy building! Feel free to open issues for clarification or improvement proposals.
