-- name: InsertUser :exec
INSERT INTO users (guid, created_at) VALUES ($1, now()) ON CONFLICT (guid) DO NOTHING;
