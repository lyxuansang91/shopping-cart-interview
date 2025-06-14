-- Payment Methods Seed Data for Stripe Adapter
INSERT INTO payment_methods (
  id,
  payment_method_code,
  partner_pm_type,
  display_name,
  config,
  status
) VALUES
  (UUID_TO_BIN(UUID(), 1), '00_card_visa',      'card', 'Visa',             '{"capture_method":"automatic"}', 'active'),
  (UUID_TO_BIN(UUID(), 1), '00_card_mastercard','card', 'Mastercard',       '{"capture_method":"automatic"}', 'active'),
  (UUID_TO_BIN(UUID(), 1), '00_card_amex',      'card', 'American Express', '{"capture_method":"automatic"}', 'active');
