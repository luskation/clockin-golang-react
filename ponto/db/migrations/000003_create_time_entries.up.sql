CREATE TYPE entry_type AS ENUM ('clock_in', 'clock_out');

CREATE TABLE time_entries (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type        entry_type NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    note        TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_time_entries_user_id ON time_entries(user_id);
CREATE INDEX idx_time_entries_recorded_at ON time_entries(recorded_at);