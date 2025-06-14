#!/bin/bash
set -e

# Create databases
mysql -u root -p$MYSQL_ROOT_PASSWORD <<EOF
CREATE DATABASE IF NOT EXISTS wallets;
CREATE DATABASE IF NOT EXISTS payments;
CREATE DATABASE IF NOT EXISTS bff;
CREATE DATABASE IF NOT EXISTS billing;
CREATE DATABASE IF NOT EXISTS ledgers;
CREATE DATABASE IF NOT EXISTS notifications;
CREATE DATABASE IF NOT EXISTS temporal;
CREATE DATABASE IF NOT EXISTS temporal_visibility;
CREATE DATABASE IF NOT EXISTS adapterstripe;
CREATE DATABASE IF NOT EXISTS organisations;
EOF

# Grant permissions to cinch user
mysql -u root -p$MYSQL_ROOT_PASSWORD <<EOF
GRANT ALL PRIVILEGES ON wallets.* TO 'cinch'@'%';
GRANT ALL PRIVILEGES ON payments.* TO 'cinch'@'%';
GRANT ALL PRIVILEGES ON bff.* TO 'cinch'@'%';
GRANT ALL PRIVILEGES ON billing.* TO 'cinch'@'%';
GRANT ALL PRIVILEGES ON ledgers.* TO 'cinch'@'%';
GRANT ALL PRIVILEGES ON notifications.* TO 'cinch'@'%';
GRANT ALL PRIVILEGES ON temporal.* TO 'cinch'@'%';
GRANT ALL PRIVILEGES ON temporal_visibility.* TO 'cinch'@'%';
GRANT ALL PRIVILEGES ON adapterstripe.* TO 'cinch'@'%';
GRANT ALL PRIVILEGES ON organisations.* TO 'cinch'@'%';
FLUSH PRIVILEGES;
EOF
