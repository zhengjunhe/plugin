#!/usr/bin/env bash
/root/dplatform -f /root/dplatform.toml &
# to wait nginx start
sleep 15
/root/dplatform -f "$PARAFILE"
