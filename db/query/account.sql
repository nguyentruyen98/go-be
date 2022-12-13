-- name: CreateAccount :one
INSERT INTO accounts(owner, balance, currency)
VALUES($1, $2, $3)
RETURNING *;
-- name: GetAccountForUpdate :one
SELECT *
FROM accounts
WHERE id = $1
LIMIT 1 FOR NO KEY
UPDATE;
-- ^ từ khoá UPDATE đễ chặn những transaction khác đang cố truy cập vào cùng 1 record trong bảng (phải chờ transaction này xong thì transaction kia mới vào đc, còn không sẽ bị chặn).
-- name: GetAccount :one
SELECT *
FROM accounts
WHERE id = $1
LIMIT 1;
-- name: ListAccounts :many
SELECT *
FROM accounts
ORDER BY id
LIMIT $1 OFFSET $2;
-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;
-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;