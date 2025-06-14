#!/bin/bash
set -e

# Create databases
mysql -u root -p$MYSQL_ROOT_PASSWORD <<EOF
CREATE DATABASE IF NOT EXISTS temporal;
CREATE DATABASE IF NOT EXISTS temporal_visibility;
CREATE DATABASE IF NOT EXISTS cart;
EOF

# Grant permissions to cinch user
mysql -u root -p$MYSQL_ROOT_PASSWORD <<EOF
GRANT ALL PRIVILEGES ON temporal.* TO 'cinch'@'%';
GRANT ALL PRIVILEGES ON temporal_visibility.* TO 'cinch'@'%';
GRANT ALL PRIVILEGES ON cart.* TO 'cinch'@'%';
FLUSH PRIVILEGES;
EOF
