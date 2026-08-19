package main

import (
	"fmt"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func windowsModule(w fyne.Window) fyne.CanvasObject {
	rt, sc := outBox()
	srcLbl := widget.NewLabel(t("lbl.src"))
	dstLbl := widget.NewLabel(t("lbl.dst"))
	var src, dst string

	detect := func() {
		runSafe(15*time.Second, func() string {
			return sh("lsblk -o NAME,SIZE,FSTYPE,MOUNTPOINTS | grep -iE 'ntfs|mount'; echo '---'; findmnt -rno TARGET,FSTYPE 2>/dev/null | grep -i ntfs; echo '---'; ls /mnt 2>/dev/null")
		}, func(out string) { setOut(rt, out) })
	}

	pickSrc := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
		if err == nil && uri != nil {
			src = uri.Path()
			srcLbl.SetText(t("lbl.src") + " " + src)
		}
	}, w)
	pickDst := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
		if err == nil && uri != nil {
			dst = uri.Path()
			dstLbl.SetText(t("lbl.dst") + " " + dst)
		}
	}, w)

	copyBtn := widget.NewButtonWithIcon(t("btn.copyStart"), theme.MediaPlayIcon(), func() {
		if src == "" || dst == "" {
			dialog.ShowInformation(t("dlg.missingSel"), t("dlg.missingSelMsg"), w)
			return
		}
		dialog.ShowConfirm(t("dlg.copyConfirm"), t("dlg.copyMsg")+fmt.Sprintf("%s\n→ %s\n\n%s?", src, dst, t("dlg.copyConfirm")), func(ok bool) {
			if !ok {
				return
			}
			runSafe(0, func() string {
				cmd := ""
				if _, err := exec.LookPath("rsync"); err == nil {
					cmd = fmt.Sprintf("rsync -a --info=progress2 %q/ %q/ 2>&1", src, dst)
				} else {
					cmd = fmt.Sprintf("cp -a %q/. %q/ 2>&1", src, dst)
				}
				return sh(cmd)
			}, func(out string) {
				setOut(rt, out+"\n\n"+t("dlg.copyDone"))
				fyne.CurrentApp().SendNotification(fyne.NewNotification(t("dlg.copyDone"), src+" → "+dst))
			})
		}, w)
	})

	top := container.NewVBox(
		widget.NewLabel(t("lbl.winDesc")),
		container.NewHBox(
			widget.NewButtonWithIcon(t("btn.findWin"), theme.SearchIcon(), detect),
			widget.NewButtonWithIcon(t("btn.openWin"), theme.FolderOpenIcon(), func() {
				runSafe(15*time.Second, func() string {
					return sh("ntf=$(lsblk -o NAME,FSTYPE,MOUNTPOINTS | awk '/ntfs/{m=$3} END{print m}'); echo 'NTFS: '${ntf:-?}; [ -n \"$ntf\" ] && xdg-open \"$ntf\" 2>&1")
				}, func(out string) { setOut(rt, out) })
			}),
		),
		srcLbl,
		container.NewHBox(
			widget.NewButton(t("btn.pickSrc"), func() { pickSrc.Show() }),
			widget.NewButton(t("btn.pickDst"), func() { pickDst.Show() }),
		),
		dstLbl,
		copyBtn,
		widget.NewLabel(t("lbl.winHint")),
	)

	topBox := container.NewBorder(top, nil, nil, nil, sc)
	detect()
	return topBox
}