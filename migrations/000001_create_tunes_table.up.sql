CREATE TABLE IF NOT EXISTS tunes (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    styles text[] NOT NULL,
    keys text[] NOT NULL,
    time_signature text NOT NULL,
    structure text NOT NULL,
    version integer NOT NULL DEFAULT 1
);