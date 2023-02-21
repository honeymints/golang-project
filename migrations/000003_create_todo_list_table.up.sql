CREATE TABLE IF NOT EXISTS lists (
ID bigserial PRIMARY KEY,
user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
title text NOT NULL
);