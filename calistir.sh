#!/usr/bin/env bash
# Linux Yardım Merkezi başlatıcı
set -e
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$DIR"
exec ./linux-yardim "$@"