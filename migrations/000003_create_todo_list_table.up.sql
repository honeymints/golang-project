CREATE TABLE IF NOT EXISTS lists (
ID bigserial PRIMARY KEY,
user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
title text NOT NULL,
created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);