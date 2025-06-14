-- Map your generic payment methods to this partner’s types + any method‑specific flags
CREATE TABLE payment_methods (
    id                   BINARY(16)   PRIMARY KEY,
    payment_method_code  VARCHAR(70)  NOT NULL,   -- matches Payments.payment_method_code
    partner_pm_type      VARCHAR(100) NOT NULL,   -- e.g. "card", "alipay", "ideal"
    display_name         VARCHAR(255) NULL,
    config               JSON         NULL,       -- partner tweaks: {"capture_method":"manual"}
    status               ENUM('active','inactive')
                            NOT NULL DEFAULT 'active',
    created_at           DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           DATETIME     NOT NULL
                            DEFAULT CURRENT_TIMESTAMP
                            ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE(payment_method_code, partner_pm_type)
);
