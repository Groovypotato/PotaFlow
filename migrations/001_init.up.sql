CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           TEXT NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);


CREATE TABLE workflows (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    is_enabled      BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX workflows_user_idx ON workflows(user_id);

CREATE TABLE triggers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id     UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    type            TEXT NOT NULL, -- 'webhook', 'cron'
    config          JSONB NOT NULL, -- dynamic: webhook_secret, cron_expr
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX triggers_workflow_idx ON triggers(workflow_id);
CREATE INDEX triggers_type_idx     ON triggers(type);


CREATE TABLE actions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id     UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    type            TEXT NOT NULL, -- 'slack', 'email', 'http', 'sheets'
    position        INT NOT NULL,  -- ordering for execution
    config          JSONB NOT NULL, -- dynamic based on action type
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX actions_workflow_idx ON actions(workflow_id);
CREATE INDEX actions_position_idx ON actions(workflow_id, position);


CREATE TABLE workflow_runs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id     UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    status          TEXT NOT NULL, -- 'pending', 'running', 'success', 'failed'
    trigger_type    TEXT NOT NULL, -- 'webhook', 'cron'
    started_at      TIMESTAMPTZ DEFAULT NULL,
    finished_at     TIMESTAMPTZ DEFAULT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX workflow_runs_workflow_idx ON workflow_runs(workflow_id);
CREATE INDEX workflow_runs_status_idx   ON workflow_runs(status);


CREATE TABLE workflow_run_logs (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    run_id            UUID NOT NULL REFERENCES workflow_runs(id) ON DELETE CASCADE,
    action_id         UUID NOT NULL REFERENCES actions(id) ON DELETE CASCADE,
    action_position   INT NOT NULL,
    success           BOOLEAN NOT NULL,
    message           TEXT NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX run_logs_run_idx       ON workflow_run_logs(run_id);
CREATE INDEX run_logs_action_idx    ON workflow_run_logs(action_id);
CREATE INDEX run_logs_position_idx  ON workflow_run_logs(run_id, action_position);


CREATE TABLE api_keys (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    hashed_key      TEXT NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX api_keys_user_idx ON api_keys(user_id);

