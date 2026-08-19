# Linux Yardım Merkezi / Linux Help Center

**Language:** [English](#english) | [Türkçe](#türkçe)

> A multi-language GUI assistant for people new to Linux — **10 modules in a single window** built with Go + Fyne.
> Linux'a yeni geçenler için çok dilli bir GUI yardımcı programı — **tek pencerede 10 modül**, Go + Fyne ile geliştirildi.

---

## English

**Linux Help Center** is a free, open-source desktop application that helps newcomers get things done on Linux without memorizing dozens of terminal commands. It offers **10 practical modules** in one window, translated into **11 languages**, with light/dark themes and safe, crash-resistant system diagnostics.

### Features

| Module | Description |
|---|---|
| System Info | Hardware, kernel, RAM, disk, GPU and IP summary |
| Command Guide | Search 40+ essential commands + natural-language → command translator |
| Package Manager | Install / remove / search / update via pacman, apt, dnf, zypper |
| Troubleshooting | Step-by-step solutions for sound, Wi-Fi, freezes, printers and more |
| Disk & Boot | `lsblk` / bootloader info (read-only) |
| Windows Data Migration | Copy data from an NTFS partition (rsync/cp) |
| Setup Wizard | Doc, Media, Dev, Gaming, System package groups |
| Learning Tracker | Command quiz with persistent score |
| Dual-Boot | Time sync (RTC) and Windows partition mounting |
| Security & Service | ufw firewall + systemd service management |
| About | Developer, version, license and project info |

### Highlights

- 🌍 **11 languages**: Turkish, English, Spanish, Japanese, Chinese, Korean, Arabic, Latin, German, French, Italian
- 🌗 **Light & dark themes** — your choice is remembered between runs
- 🛡️ **Crash-safe diagnostics** — all system checks run with timeouts and error recovery
- 🔤 **Command translator** — type in your own language and get the right terminal command
- 📚 **Learning mode** — test your command knowledge, score is saved locally
- 🖥️ **Admin tasks** use a polkit password dialog (`pkexec`) — no hardcoded root
- 📦 **Embedded fonts** — DejaVu fonts are bundled to avoid font rendering crashes

### Requirements

- Go 1.21+ (to build from source)
- gcc / CGO enabled
- GTK3 / X11 / OpenGL development libraries

### Build & Run

```bash
git clone https://github.com/YunusTAS13/EasyLinux.git
cd EasyLinux
go build -o EasyLinux .
./EasyLinux
```

Or just run the launcher script:

```bash
./calistir.sh
```

### Usage Notes

- Read-only modules (Disk & Boot) are completely safe. For partition editing, use GParted and always back up first.
- Admin operations open the polkit password dialog automatically.
- The learning score is stored in `~/.config/EasyLinux/skor.json`.

### Reporting Issues / Contributing

Found a bug or want a new feature? Open an issue at:

- **Bug reports:** https://github.com/YunusTAS13/EasyLinux/issues

> This project is developed and maintained by **Yunus Taş**. External contributors are not accepted at this time, but detailed bug reports, suggestions and translations are always welcome.

### License

This project is licensed under the **GNU General Public License v3.0** (GPL-3.0). See the [LICENSE](LICENSE) file for details.

Copyright (C) 2026 **Yunus Taş**

---

## Türkçe

**Linux Yardım Merkezi**, Linux'a yeni geçenlerin düzinelerce terminal komutu ezberlemeden işlerini halletmesini sağlayan ücretsiz, açık kaynaklı bir masaüstü uygulamasıdır. Tek pencerede **10 pratik modül** sunar; **11 dile** çevrilidir; açık/karanlık temalar ve çökme korumalı, güvenli sistem teşhisi içerir.

### Özellikler

| Modül | Açıklama |
|---|---|
| Sistem Bilgi | Donanım, çekirdek, RAM, disk, GPU ve IP özeti |
| Komut Rehberi | 40+ temel komut arama + doğal dil → komut çevirici |
| Paket Yöneticisi | pacman, apt, dnf, zypper için kur / kaldır / ara / güncelle |
| Sorun Giderme | Ses, Wi-Fi, donma, yazıcı vb. için adım adım çözüm |
| Disk & Önyükleme | `lsblk` / önyükleyici bilgisi (sadece okuma) |
| Windows Veri Taşıma | NTFS bölümünden veri kopyalama (rsync/cp) |
| Kurulum Sihirbazı | Doküman, Medya, Geliştirici, Oyun, Sistem paket grupları |
| Öğrenme Takibi | Komut bilgisi testi, skor yerel olarak saklanır |
| Dual-Boot | Zaman uyumu (RTC) ve Windows bölümü bağlama |
| Güvenlik & Servis | ufw güvenlik duvarı + systemd servis yönetimi |
| Hakkında | Geliştirici, sürüm, lisans ve proje bilgisi |

### Öne Çıkanlar

- 🌍 **11 dil**: Türkçe, İngilizce, İspanyolca, Japonca, Çince, Korece, Arapça, Latince, Almanca, Fransızca, İtalyanca
- 🌗 **Açık & karanlık tema** — seçiminiz sonraki açılışlarda da hatırlanır
- 🛡️ **Çökme korumalı teşhis** — tüm sistem kontrolleri zaman aşımı ve hata yakalama ile çalışır
- 🔤 **Komut çevirici** — kendi dilinizde yazın, doğru terminal komutunu bulun
- 📚 **Öğrenme modu** — komut bilginizi test edin, skor yerelde saklanır
- 🖥️ **Yönetici işlemleri** polkit parola penceresiyle (`pkexec`) çalışır — root şifresi gömülü değildir
- 📦 **Gömülü fontlar** — font çizim çökmelerine karşı DejaVu fontları pakete gömülüdür

### Gereksinimler

- Go 1.21+ (kaynak koddan derlemek için)
- gcc / CGO etkin
- GTK3 / X11 / OpenGL geliştirme kütüphaneleri

### Derleme & Çalıştırma

```bash
git clone https://github.com/YunusTAS13/EasyLinux.git
cd EasyLinux
go build -o EasyLinux .
./EasyLinux
```

Veya başlatıcı betiği kullanın:

```bash
./calistir.sh
```

### Kullanım Notları

- Yalnızca okuma yapan modüller (Disk & Önyükleme) tamamen güvenlidir. Bölüm düzenlemek için GParted kullanın ve önce yedek alın.
- Yönetici işlemleri polkit parola penceresini otomatik açar.
- Öğrenme skoru `~/.config/EasyLinux/skor.json` içinde saklanır.

### Sorun Bildirme / Katkı

Bir hata buldunuz veya yeni bir özellik mi istiyorsunuz? Issue açın:

- **Hata raporları:** https://github.com/YunusTAS13/EasyLinux/issues

> Bu proje **Yunus Taş** tarafından geliştirilmekte ve bakımı yapılmaktadır. Şu anda harici katkıcı kabul edilmemektedir; ancak ayrıntılı hata raporları, öneriler ve çeviriler her zaman memnuniyetle karşılanır.

### Lisans

Bu proje **GNU Genel Kamu Lisansı v3.0** (GPL-3.0) ile lisanslanmıştır. Ayrıntılar için [LICENSE](LICENSE) dosyasına bakın.

Telif Hakkı (C) 2026 **Yunus Taş**