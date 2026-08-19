# Changelog

All notable changes to this project are documented in this file. The format is based on [Keep a Changelog](https://keepachangelog.com/).

## [1.0.0] - 2026-08-19

### Added
- Single-window GUI with 10 modules + About, built with Go and Fyne.
- System Info module: hardware, kernel, RAM, disk, GPU, IP summary.
- Command Guide module: search 40+ essential commands and natural-language → command translator.
- Package Manager module: install/remove/search/update for pacman, apt, dnf, zypper.
- Troubleshooting module: step-by-step guides for sound, Wi-Fi, freezes, printers.
- Disk & Boot module: `lsblk` and bootloader info (read-only).
- Windows Data Migration module: copy data from NTFS partitions (rsync/cp).
- Setup Wizard module: doc, media, dev, gaming, system package groups.
- Learning Tracker module: command quiz with locally saved score.
- Dual-Boot module: RTC time sync and Windows partition mounting.
- Security & Service module: ufw firewall and systemd service management.
- About module with developer, version, license and repository info.
- Interface translated into 11 languages: tr, en, es, ja, zh, ko, ar, la, de, fr, it.
- Command guide data (42 commands) fully translated into all 11 languages.
- Light and dark themes, remembered between runs via Fyne preferences.
- Crash-safe diagnostics: all system checks run under timeouts with panic recovery.
- Embedded DejaVu fonts to prevent font rendering crashes.
- Turkish launcher script `calistir.sh`.

### Changed
- Refactored command data into per-language i18n tables (`cmd_i18n.go`).
- Replaced `mod_windows.go` with `winmigrate.go` to avoid Go's platform build-tag suffix.
- Migrated all UI updates to Fyne's `fyne.Do` threading model.

### Fixed
- Font crash caused by system Noto font GSUB table (fixed by embedding DejaVu fonts).
- Stack overflow on language switch (recursive `buildUI` guarded by a building flag).
- Literal `**` markdown artifacts in the System Info module (now rendered properly).
- Tiny search boxes and barely visible text (theme text size raised).
- Invalid Fyne API calls (`SetPlaceHolder`, `SunriseIcon`/`MoonIcon`, `Entry.SetMinSize`).
- `runSafe` immediate timeout when timeout was zero.
- Untranslated source/destination labels in Windows Data Migration.
- Missing `time` import in the System Info module.

### Security
- Admin operations use polkit (`pkexec`) with a password dialog; no root password is embedded.
- No commands are executed with hardcoded credentials.