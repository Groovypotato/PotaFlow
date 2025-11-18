## Week 1 – Skeleton & Foundations
- Create GitHub repo and add scaffolding
- Add README, .gitignore, Makefile
- Initialize Go module and dependencies
- Design initial DB schema (users, workflows, triggers, actions, workflow_runs, logs)
- Create migration 001_init.sql and run with migrate
- Implement DB connection
- Create basic API server with health endpoint

## Week 2 – Auth & Basic Workflow CRUD
- Implement Argon2 password hashing
- Implement JWT auth and middleware
- Endpoints: /auth/register, /auth/login, /me
- Implement workflow CRUD (Create, List, Get, Update, Delete)
- Test in Postman/curl

## Week 3 – Triggers, Actions & Worker Pool
- Add triggers and actions tables
- Implement CRUD for triggers/actions
- Build worker pool (goroutines + job queue)
- Implement workflow_runs table and RunWorkflow()
- Worker service loads pending runs and executes them

## Week 4 – Webhooks, Cron & Logs + Simple UI
- Add webhook trigger + POST /hooks/{trigger_id}
- Add cron trigger using robfig/cron
- Implement action execution pipeline
- Log results to workflow_run_logs
- UI: login, dashboard, workflow view, runs page

## Week 5 – Integrations, Polish & Deployment
- Implement Slack, Email, HTTP actions
- Improve error handling, config, logging
- Finalize docker-compose.yml
- Deploy to Fly.io/Railway/Render
- Validate full flow end-to-end

## Week 6 – Testing, Docs & Presentation
- Unit tests for auth, workflows, worker pool
- Polish README with diagrams/screenshots
- Document example workflows
- Prepare live demo script
