package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func diskModule() fyne.CanvasObject {
	rt, sc := outBox()

	refresh := func() {
		runSafe(15*time.Second, func() string {
			out := sh("lsblk -o NAME,SIZE,TYPE,FSTYPE,MOUNTPOINTS 2>/dev/null || lsblk")
			out += "\n\n### " + t("tab.disk") + "\n\n"
			boot := sh("(systemd-bootctl status 2>/dev/null || efibootmgr 2>/dev/null || (ls /boot/grub >/dev/null 2>&1 && echo 'GRUB') || ls /boot) | head -30")
			out += boot
			return out
		}, func(out string) { setOut(rt, out) })
	}
	refresh()

	top := container.NewHBox(
		widget.NewButtonWithIcon(t("btn.refresh"), theme.ViewRefreshIcon(), refresh),
		widget.NewButtonWithIcon(t("btn.kernel"), theme.InfoIcon(), func() {
			runSafe(15*time.Second, func() string {
				return sh("uname -a; echo '---'; lscpu | grep -E 'Model name|Architecture|Thread|Core'")
			}, func(out string) { setOut(rt, out) })
		}),
		widget.NewButtonWithIcon(t("btn.gparted"), theme.FolderOpenIcon(), func() {
			runSafe(15*time.Second, func() string {
				return sh("(gparted &) 2>&1 || echo 'GParted yok'")
			}, func(out string) { setOut(rt, out) })
		}),
	)

	warn := widget.NewLabel(t("lbl.diskWarn"))
	warn.Wrapping = fyne.TextWrapWord

	return container.NewBorder(container.NewVBox(warn, top), nil, nil, nil, sc)
}