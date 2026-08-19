package main

import (
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type tcase struct {
	key   string
	steps []string
	diag  string
}

var troubles = []tcase{
	{"trb0", []string{
		"`pactl info` ile ses sunucusunun çalıştığını kontrol edin (PipeWire/PulseAudio).",
		"Yeniden başlatmayı deneyin: `systemctl --user restart pipewire` veya `pulseaudio -k`.",
		"`speaker-test -c 2 -t wav` ile test tonu çalın.",
		"`pavucontrol` kurup çıkış cihazının doğru seçildiğinden ve ses seviyesinin açık olduğundan emin olun.",
	}, "pactl info 2>/dev/null | head -12; echo '---'; aplay -l 2>/dev/null"},
	{"trb1", []string{
		"`nmcli device wifi list` ile görünür ağları listeleyin.",
		"Bağlanın: `nmcli device wifi connect 'AG_ADI' password 'sifre'`.",
		"Ağ yöneticisini yeniden başlatın: `sudo systemctl restart NetworkManager`.",
		"Ağ görünmüyorsa `sudo rfkill unblock wifi` ile engeli kaldırın.",
	}, "nmcli device status 2>/dev/null; echo '---'; nmcli device wifi list 2>/dev/null | head -10"},
	{"trb2", []string{
		"`Ctrl+Alt+F2` ile başka bir sanal konsola geçin.",
		"`ps aux | grep <donan-program>` ile programı bulup `kill -9 PID` ile kapatın.",
		"Masaüstü hâlâ yanıt vermiyorsa `Ctrl+Alt+F1` ile dönün veya `sudo systemctl restart display-manager` çalıştırın.",
	}, "ps aux --sort=-%cpu | head -8"},
	{"trb3", []string{
		"`free -h` ile RAM yeterli mi, swap var mı bakın.",
		"`systemd-analyze blame` ile hangi servisin açılışı yavaşlattığını görün.",
		"`df -h` ile kök disk dolu mu kontrol edin.",
		"`htop` ile CPU/RAM yiyen programı tespit edin.",
	}, "free -h; echo '---'; df -h /; echo '---'; systemd-analyze blame 2>/dev/null | head -12"},
	{"trb4", []string{
		"`lpstat -p` ile yazıcı tanınıyor mu bakın.",
		"`systemctl status cups` ile CUPS servisi aktif mi kontrol edin.",
		"Servisi yeniden başlatın: `sudo systemctl restart cups`.",
	}, "lpstat -p 2>&1; echo '---'; systemctl is-active cups 2>&1"},
	{"trb5", []string{
		"`lsusb` ile cihaz görünüyor mu bakın.",
		"Cihazı takıp çıkarırken `dmesg | tail` çıktısına bakın.",
		"`sudo dmesg | grep -i usb` ile usb çekirdek mesajlarını inceleyin.",
	}, "lsusb; echo '---'; dmesg 2>/dev/null | tail -15"},
	{"trb6", []string{
		"`lspci -k | grep -A3 -iE 'vga|3d'` ile hangi sürücünün yüklü olduğunu görün.",
		"NVIDIA için `nvidia-smi` çalışıyor mu test edin.",
		"AMD için `glxinfo -B` ile renderer adını kontrol edin (mesa-vulkan gerekebilir).",
	}, "lspci -k 2>/dev/null | grep -iE 'vga|3d|display' | head; echo '---'; nvidia-smi 2>/dev/null || glxinfo -B 2>/dev/null | head -6"},
	{"trb7", []string{
		"`ip a` ile ağ arayüzünün IP aldığını kontrol edin.",
		"`ping -c3 8.8.8.8` ile IP üzerinden bağlantıyı test edin.",
		"DNS için `resolvectl status` kontrol edin.",
		"`sudo systemctl restart NetworkManager` ile ağ yöneticisini yeniden başlatın.",
	}, "ip a 2>/dev/null | head -20; echo '---'; ping -c3 8.8.8.8 2>&1; echo '---'; resolvectl status 2>/dev/null | head -8"},
	{"trb8", []string{
		"`df -h` ile hangi bölümün dolu olduğunu görün.",
		"`du -xh --max-depth=1 ~ 2>/dev/null | sort -h | tail` ile ev klasöründe büyük ne olduğunu bulun.",
		"Sistem günlüklerini küçültün: `sudo journalctl --vacuum-size=100M`.",
		"Paket önbelleğini temizleyin: `sudo pacman -Sc` / `sudo apt clean`.",
	}, "df -h; echo '---'; du -xh --max-depth=1 ~ 2>/dev/null | sort -h | tail -8"},
	{"trb9", []string{
		"`Ctrl+Alt+F2` ile konsola geçin.",
		"`sudo systemctl restart display-manager` ile masaüstünü yeniden başlatmayı deneyin.",
		"`sudo dmesg | tail` ile ekran kartı hatası var mı bakın.",
		"Son çare olarak `sudo pacman -Syu` ile sistemi güncelleyip yeniden başlatın.",
	}, "systemctl status display-manager 2>&1 | head -8; echo '---'; dmesg 2>/dev/null | tail -10"},
}

func troubleModule() fyne.CanvasObject {
	var titles []string
	for _, tt := range troubles {
		titles = append(titles, t(tt.key))
	}

	stepRT, stepSc := outBox()
	diagRT, diagSc := outBox()
	busy := widget.NewActivity()
	busy.Hide()
	sel := widget.NewSelect(titles, nil)

	var current tcase
	sel.OnChanged = func(s string) {
		for _, tt := range troubles {
			if t(tt.key) == s {
				current = tt
				break
			}
		}
		var sb strings.Builder
		sb.WriteString("## " + t("lbl.steps") + "\n\n")
		for _, st := range current.steps {
			sb.WriteString("1. " + st + "\n")
		}
		setMarkdown(stepRT, sb.String())
		setMarkdown(diagRT, "_" + t("trb.diagEmpty") + "_")
	}
	sel.PlaceHolder = t("trb.choose")

	diagBtn := widget.NewButtonWithIcon(t("btn.diag"), theme.SearchIcon(), func() {
		if current.key == "" || current.diag == "" {
			return
		}
		busy.Show()
		busy.Start()
		runSafe(20*time.Second, func() string {
			return sh(current.diag)
		}, func(out string) {
			busy.Stop()
			busy.Hide()
			setMarkdown(diagRT, "```\n"+out+"\n```")
		})
	})

	stepsCard := container.NewBorder(
		container.NewVBox(widget.NewLabel(t("lbl.steps"))),
		nil, nil, nil,
		stepSc,
	)

	diagCard := container.NewBorder(
		container.NewHBox(widget.NewLabel(t("lbl.diag")), diagBtn, busy),
		nil, nil, nil,
		diagSc,
	)

	top := container.NewVBox(
		widget.NewLabel(t("lbl.troubleHint")),
		sel,
	)
	body := container.NewVSplit(stepsCard, diagCard)
	return container.NewBorder(top, nil, nil, nil, body)
}