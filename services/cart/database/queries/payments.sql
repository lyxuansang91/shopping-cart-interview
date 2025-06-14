-- name: GetPaymentMethodByCode :one
SELECT
  id,
  payment_method_code,
  partner_pm_type,
  display_name,
  config,
  status,
  created_at,
  updated_at
FROM payment_methods
WHERE payment_method_code = ?;

-- name: ListPaymentMethods :many
SELECT
  id,
  payment_method_code,
  partner_pm_type,
  display_name,
  config,
  status,
  created_at,
  updated_at
FROM payment_methods
ORDER BY created_at DESC;

-- name: EnablePaymentMethodByCode :exec
UPDATE payment_methods
SET
  status = 'active',
  updated_at = CURRENT_TIMESTAMP
WHERE payment_method_code = ?;

-- name: DisablePaymentMethodByCode :exec
UPDATE payment_methods
SET
  status = 'inactive',
  updated_at = CURRENT_TIMESTAMP
WHERE payment_method_code = ?;

-- name: DeletePaymentMethodByCode :exec
DELETE FROM payment_methods
WHERE payment_method_code = ?;
