-- Manual script: Create all custom MySQL functions required by the Cinch monorepo
-- Run this after running your main migrations (once per database)
--
-- This script is idempotent if you drop functions before creating them.

-- Drop and recreate BIN_TO_UUID_STRING for UUID formatting
DROP FUNCTION IF EXISTS BIN_TO_UUID_STRING;
CREATE FUNCTION BIN_TO_UUID_STRING(bin_uuid BINARY(16))
RETURNS VARCHAR(36)
DETERMINISTIC
RETURN LOWER(CONCAT(
    SUBSTR(HEX(bin_uuid), 1, 8), '-',
    SUBSTR(HEX(bin_uuid), 9, 4), '-',
    SUBSTR(HEX(bin_uuid), 13, 4), '-',
    SUBSTR(HEX(bin_uuid), 17, 4), '-',
    SUBSTR(HEX(bin_uuid), 21)
));

-- Add more custom functions below as needed
