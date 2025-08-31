-- name: ChangeBalance :exec
INSERT INTO balance_register (user_guid, operation_ref, amount, created_at ) VALUES ($1, $2, $3, now());

-- name: GetUserBalance :one
SELECT user_guid, SUM(amount)::FLOAT FROM balance_register WHERE user_guid = @user_guid GROUP BY user_guid;
