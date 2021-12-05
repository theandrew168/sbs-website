#!/bin/sh
set -e

# Create sbs group (if it doesn't exist)
if ! getent group sbs >/dev/null; then
    groupadd --system sbs
fi

# Create sbs user (if it doesn't exist)
if ! getent passwd sbs >/dev/null; then
    useradd                        \
        --system                   \
        --gid sbs                  \
        --shell /usr/sbin/nologin  \
        --comment "sbs website"    \
        sbs
fi

# Update config file permissions (idempotent)
chown root:sbs /etc/sbs.conf
chmod 0640 /etc/sbs.conf

# Reload systemd to pickup sbs.service
systemctl daemon-reload
