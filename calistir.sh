#!/usr/bin/env bash
# EasyLinux başlatıcı
set -e
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$DIR"
exec ./EasyLinux "$@"