
BEGIN;

-- Reverse dependency order: children before parents.
DROP TABLE IF EXISTS payment_events;
DROP TABLE IF EXISTS refunds;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS vehicles;
DROP TABLE IF EXISTS users;

-- Safe now: every trigger that used it was dropped with its table.
DROP FUNCTION IF EXISTS set_updated_at();

COMMIT;