-- RentGo baseline schema (PostgreSQL 13+)
-- Money: BIGINT rupiah. IDs: UUID.

BEGIN;

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$;

-- ---------------------------------------------------------------- users
CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         TEXT        NOT NULL,
    full_name     TEXT        NOT NULL,
    password_hash TEXT        NOT NULL,
    phone_number  TEXT        NOT NULL,
    address       TEXT        NOT NULL,
    role          TEXT        NOT NULL DEFAULT 'user'
                  CONSTRAINT users_role_check CHECK (role IN ('user','admin')),
    is_active     BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX ux_users_email ON users (lower(email));

-- -------------------------------------------------------------- vehicles
CREATE TABLE vehicles (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT        NOT NULL,
    description TEXT        NOT NULL DEFAULT '',
    category    TEXT        NOT NULL,
    quantity    INTEGER     NOT NULL DEFAULT 0
                CONSTRAINT vehicles_quantity_check CHECK (quantity >= 0),
    daily_rate  BIGINT      NOT NULL
                CONSTRAINT vehicles_daily_rate_check CHECK (daily_rate > 0),
    is_active   BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ix_vehicles_active_category
    ON vehicles (category)
    WHERE is_active;

-- -------------------------------------------------------------- bookings
CREATE TABLE bookings (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID        NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    vehicle_id      UUID        NOT NULL REFERENCES vehicles(id) ON DELETE RESTRICT,
    booking_start   TIMESTAMPTZ NOT NULL,
    booking_end     TIMESTAMPTZ NOT NULL,
    total_days      INTEGER     NOT NULL
                    CONSTRAINT bookings_total_days_check CHECK (total_days > 0),
    daily_rate      BIGINT      NOT NULL
                    CONSTRAINT bookings_daily_rate_check CHECK (daily_rate > 0),
    total_price     BIGINT GENERATED ALWAYS AS (daily_rate * total_days) STORED,
    status          TEXT        NOT NULL DEFAULT 'pending_payment'
                    CONSTRAINT bookings_status_check
                    CHECK (status IN ('pending_payment','confirmed','completed','overdue','cancelled')),
    hold_expires_at TIMESTAMPTZ,
    returned_at     TIMESTAMPTZ,
    returned_by     UUID REFERENCES users(id) ON DELETE RESTRICT,
    idempotency_key TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT bookings_dates_check
        CHECK (booking_end > booking_start),
    CONSTRAINT bookings_end_matches_days
        CHECK (booking_end = booking_start + make_interval(days => total_days)),
    CONSTRAINT bookings_hold_check
        CHECK (status <> 'pending_payment' OR hold_expires_at IS NOT NULL),

    CONSTRAINT bookings_completed_check
        CHECK (status <> 'completed' OR (returned_at IS NOT NULL AND returned_by IS NOT NULL)),
    CONSTRAINT bookings_returned_only_when_completed
        CHECK (returned_at IS NULL OR status = 'completed'),
    CONSTRAINT bookings_returned_pair
        CHECK ((returned_at IS NULL) = (returned_by IS NULL)),
    CONSTRAINT bookings_returned_after_start
        CHECK (returned_at IS NULL OR returned_at >= booking_start)
);

CREATE INDEX ix_bookings_user_created
    ON bookings (user_id, created_at DESC);
CREATE INDEX ix_bookings_active_window
    ON bookings (vehicle_id, booking_start, booking_end)
    WHERE status IN ('pending_payment','confirmed');
CREATE INDEX ix_bookings_due
    ON bookings (booking_end)
    WHERE status = 'confirmed' AND returned_at IS NULL;
CREATE UNIQUE INDEX ux_bookings_idem
    ON bookings (user_id, idempotency_key)
    WHERE idempotency_key IS NOT NULL;

-- -------------------------------------------------------------- payments
CREATE TABLE payments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id      UUID        NOT NULL REFERENCES bookings(id) ON DELETE RESTRICT,
    user_id         UUID        NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    amount          BIGINT      NOT NULL
                    CONSTRAINT payments_amount_check CHECK (amount > 0),
    currency        VARCHAR(3)  NOT NULL DEFAULT 'IDR'
                    CONSTRAINT payments_currency_check CHECK (currency = upper(currency)),
    status          TEXT        NOT NULL DEFAULT 'pending'
                    CONSTRAINT payments_status_check
                    CHECK (status IN ('pending','paid','expired','failed','refunded','partially_refunded')),
    provider        TEXT        NOT NULL,
    provider_ref    TEXT,
    idempotency_key TEXT,
    expires_at      TIMESTAMPTZ NOT NULL,
    paid_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT payments_paid_at_check
        CHECK (status IN ('pending','expired','failed') OR paid_at IS NOT NULL)
);

CREATE UNIQUE INDEX ux_payments_active
    ON payments (booking_id)
    WHERE status = 'pending';
CREATE UNIQUE INDEX ux_payments_provider_ref
    ON payments (provider, provider_ref)
    WHERE provider_ref IS NOT NULL;
CREATE UNIQUE INDEX ux_payments_idem
    ON payments (user_id, idempotency_key)
    WHERE idempotency_key IS NOT NULL;
CREATE INDEX ix_payments_booking ON payments (booking_id);
CREATE INDEX ix_payments_expiry
    ON payments (expires_at)
    WHERE status = 'pending';

-- --------------------------------------------------------------- refunds
CREATE TABLE refunds (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id      UUID        NOT NULL REFERENCES payments(id) ON DELETE RESTRICT,
    amount          BIGINT      NOT NULL
                    CONSTRAINT refunds_amount_check CHECK (amount > 0),
    reason          TEXT,
    status          TEXT        NOT NULL DEFAULT 'pending'
                    CONSTRAINT refunds_status_check
                    CHECK (status IN ('pending','succeeded','failed')),
    provider_ref    TEXT,
    idempotency_key TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT refunds_succeeded_ref_check
        CHECK (status <> 'succeeded' OR provider_ref IS NOT NULL)
);

CREATE INDEX ix_refunds_payment ON refunds (payment_id);
CREATE UNIQUE INDEX ux_refunds_idem
    ON refunds (payment_id, idempotency_key)
    WHERE idempotency_key IS NOT NULL;
CREATE UNIQUE INDEX ux_refunds_provider_ref
    ON refunds (provider_ref)
    WHERE provider_ref IS NOT NULL;

-- -------------------------------------------------------- payment_events
CREATE TABLE payment_events (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider     TEXT        NOT NULL,
    event_id     TEXT        NOT NULL,
    event_type   TEXT,
    payload      JSONB       NOT NULL,
    received_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    processed_at TIMESTAMPTZ,

    CONSTRAINT ux_payment_events_event UNIQUE (provider, event_id)
);

CREATE INDEX ix_payment_events_unprocessed
    ON payment_events (received_at)
    WHERE processed_at IS NULL;

-- -------------------------------------------------------------- triggers
CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON users    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_vehicles_updated_at
    BEFORE UPDATE ON vehicles FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_bookings_updated_at
    BEFORE UPDATE ON bookings FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_payments_updated_at
    BEFORE UPDATE ON payments FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_refunds_updated_at
    BEFORE UPDATE ON refunds  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMIT;