-- name: CreateUser :one
INSERT INTO t1.users (
    created_at,
    created_by,
    updated_at,
    updated_by,
    deleted_at,
    deleted_by,
    username,
    type
) VALUES (
    now() AT TIME ZONE 'utc', -- created_at
    @created_by,
    now() AT TIME ZONE 'utc', -- updated_at
    @created_by, -- updated_by
    NULL, -- deleted_at
    NULL, -- deleted_by
    @username,
    @type
) RETURNING *;

-- name: UpsertUser :one
INSERT INTO t1.users (
    created_at,
    created_by,
    updated_at,
    updated_by,
    deleted_at,
    deleted_by,
    username,
    type
) VALUES (
    now() AT TIME ZONE 'utc', -- created_at
    @created_by,
    now() AT TIME ZONE 'utc', -- updated_at
    @created_by, -- updated_by
    NULL, -- deleted_at
    NULL, -- deleted_by
    @username,
    @type
) ON CONFLICT(commit_hash)
do UPDATE SET
    updated_at = now() AT TIME ZONE 'utc',
    updated_by = @created_by,
    deleted_at = NULL,
    deleted_by = NULL,
    username = @username,
    type = @type
RETURNING *;

-- name: FindUser :one
SELECT * FROM t1.users WHERE deleted_at IS NULL AND id = $1 LIMIT 1;

-- name: FindDeletedUser :one
SELECT * FROM t1.users WHERE deleted_at IS NOT NULL AND id = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE t1.users SET
    created_at = created_at,
    created_by = created_by,
    updated_at = now() AT TIME ZONE 'utc',
    updated_by = @updated_by,
    deleted_at = NULL,
    deleted_by = NULL,
    username = @username,
    type = @type
WHERE deleted_at IS NULL AND id = $1 RETURNING *;

-- name: DeleteUser :execrows
UPDATE t1.users SET deleted_at = now() AT TIME ZONE 'utc', deleted_by = $2 WHERE deleted_at IS NULL AND id = $1;

-- name: PurgeUsersOlderThan :execrows
DELETE FROM t1.users WHERE deleted_at IS NOT NULL AND deleted_at < $1;
