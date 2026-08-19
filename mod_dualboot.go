package main

import (
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func dualbootModule(w fyne.Window) fyne.CanvasObject {
	rt, sc := outBox()
	devEntry := widget.NewEntry()
	devEntry.SetPlaceHolder(t("ph.device"))
	devEntry.SetText(sh("lsblk -nlo NAME,FSTYPE | awk '$2==\"ntfs\"{print \"/dev/\"$1; exit}'"))

	refresh := func() {
		runSafe(15*time.Second, func() string {
			return sh("timedatectl 2>/dev/null | grep -iE 'RTC in local|Time zone'; echo '---'; grubinfo=$(grep GRUB_DEFAULT /etc/default/grub 2>/dev/null || echo 'GRUB yok'); echo $grubinfo; echo '---'; lsblk -o NAME,SIZE,FSTYPE,MOUNTPOINTS | grep -iE 'ntfs'")
		}, func(out string) { setOut(rt, out) })
	}

	top := container.NewVBox(
		widget.NewLabel(t("lbl.dualbootDesc")),
		container.NewHBox(
			widget.NewButtonWithIcon(t("btn.status"), theme.InfoIcon(), refresh),
			widget.NewButtonWithIcon(t("btn.rtcLocal"), theme.ConfirmIcon(), func() {
				runSafe(30*time.Second, func() string {
					return shPriv("timedatectl set-local-rtc 1") + "\n\n" + sh("timedatectl | grep -i rtc")
				}, func(out string) { setOut(rt, out) })
			}),
			widget.NewButtonWithIcon(t("btn.rtcUTC"), theme.CancelIcon(), func() {
				runSafe(30*time.Second, func() string {
					return shPriv("timedatectl set-local-rtc 0") + "\n\n" + sh("timedatectl | grep -i rtc")
				}, func(out string) { setOut(rt, out) })
			}),
		),
		container.NewHBox(
			devEntry,
			widget.NewButtonWithIcon(t("btn.mount"), theme.FolderOpenIcon(), func() {
				dev := strings.TrimSpace(devEntry.Text)
				if dev == "" {
					dev = "/dev/sda1"
				}
				runSafe(30*time.Second, func() string {
					return shPriv("mkdir -p /mnt/windows && mount -t ntfs-3g " + dev + " /mnt/windows 2>&1 && echo '/mnt/windows'")
				}, func(out string) { setOut(rt, out) })
			}),
			widget.NewButtonWithIcon(t("btn.umount"), theme.LogoutIcon(), func() {
				runSafe(30*time.Second, func() string {
					return shPriv("umount /mnt/windows 2>&1 && echo 'umounted'")
				}, func(out string) { setOut(rt, out) })
			}),
		),
	)
	body := container.NewBorder(top, nil, nil, nil, sc)
	refresh()
	return body
}