package main

import (
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func detectPM() string {
	for _, p := range []string{"pacman", "apt", "dnf", "zypper", "apk", "emerge"} {
		if _, err := exec.LookPath(p); err == nil {
			return p
		}
	}
	return "bilinmiyor"
}

func pmCmd(pm, op string, args ...string) string {
	var base []string
	switch pm {
	case "pacman":
		switch op {
		case "refresh":
			base = []string{"pacman", "-Sy"}
		case "update":
			base = []string{"pacman", "-Syu", "--noconfirm"}
		case "install":
			base = []string{"pacman", "-S", "--noconfirm"}
		case "remove":
			base = []string{"pacman", "-Rns", "--noconfirm"}
		case "search":
			base = []string{"pacman", "-Ss"}
		}
	case "apt":
		switch op {
		case "refresh":
			base = []string{"apt", "update"}
		case "update":
			base = []string{"apt", "full-upgrade", "-y"}
		case "install":
			base = []string{"apt", "install", "-y"}
		case "remove":
			base = []string{"apt", "purge", "-y"}
		case "search":
			base = []string{"apt", "search"}
		}
	case "dnf":
		switch op {
		case "refresh":
			base = []string{"dnf", "makecache"}
		case "update":
			base = []string{"dnf", "upgrade", "-y"}
		case "install":
			base = []string{"dnf", "install", "-y"}
		case "remove":
			base = []string{"dnf", "remove", "-y"}
		case "search":
			base = []string{"dnf", "search"}
		}
	case "zypper":
		switch op {
		case "refresh":
			base = []string{"zypper", "refresh"}
		case "update":
			base = []string{"zypper", "dup", "-y"}
		case "install":
			base = []string{"zypper", "install", "-y"}
		case "remove":
			base = []string{"zypper", "remove", "-y"}
		case "search":
			base = []string{"zypper", "search"}
		}
	case "apk":
		switch op {
		case "refresh":
			base = []string{"apk", "update"}
		case "update":
			base = []string{"apk", "upgrade"}
		case "install":
			base = []string{"apk", "add"}
		case "remove":
			base = []string{"apk", "del"}
		case "search":
			base = []string{"apk", "search"}
		}
	}
	if len(base) == 0 {
		return ""
	}
	return strings.Join(append(base, args...), " ")
}

func packagesModule() fyne.CanvasObject {
	pm := detectPM()
	distro := sh("grep PRETTY_NAME /etc/os-release | cut -d= -f2 | tr -d '\"'")

	info := widget.NewLabel(t("lbl.pmInfo") + ": " + pm + "  (" + distro + ")")
	info.TextStyle = fyne.TextStyle{Bold: true}
	pkg := widget.NewEntry()
	pkg.SetPlaceHolder(t("ph.pkg"))

	rt, sc := outBox()
	busy := widget.NewActivity()
	busy.Hide()

	setBusy := func(b bool) {
		doUI(func() {
			if b {
				busy.Start()
				busy.Show()
			} else {
				busy.Stop()
				busy.Hide()
			}
		})
	}

	run := func(cmd string, priv bool) {
		if cmd == "" {
			return
		}
		setBusy(true)
		runSafe(0, func() string {
			if priv {
				return shPriv(cmd)
			}
			return sh(cmd)
		}, func(out string) {
			setOut(rt, "> "+cmd+"\n\n"+out)
			setBusy(false)
		})
	}

	updateBtn := widget.NewButtonWithIcon(t("btn.update"), theme.ViewRefreshIcon(), func() {
		run(pmCmd(pm, "update"), true)
	})
	refreshBtn := widget.NewButtonWithIcon(t("btn.refreshRepo"), theme.DownloadIcon(), func() {
		run(pmCmd(pm, "refresh"), true)
	})
	installBtn := widget.NewButtonWithIcon(t("btn.install"), theme.ConfirmIcon(), func() {
		if strings.TrimSpace(pkg.Text) == "" {
			dialog.ShowInformation(t("dlg.enterPkg"), t("dlg.enterPkg"), fyne.CurrentApp().Driver().AllWindows()[0])
			return
		}
		run(pmCmd(pm, "install", pkg.Text), true)
	})
	removeBtn := widget.NewButtonWithIcon(t("btn.remove"), theme.DeleteIcon(), func() {
		if strings.TrimSpace(pkg.Text) == "" {
			dialog.ShowInformation(t("dlg.enterPkg"), t("dlg.enterPkg"), fyne.CurrentApp().Driver().AllWindows()[0])
			return
		}
		dialog.ShowConfirm(t("dlg.removeTitle"), pkg.Text+t("dlg.removeMsg"), func(ok bool) {
			if ok {
				run(pmCmd(pm, "remove", pkg.Text), true)
			}
		}, fyne.CurrentApp().Driver().AllWindows()[0])
	})
	searchBtn := widget.NewButtonWithIcon(t("btn.search"), theme.SearchIcon(), func() {
		if strings.TrimSpace(pkg.Text) == "" {
			dialog.ShowInformation(t("dlg.enterPkg"), t("dlg.enterPkg"), fyne.CurrentApp().Driver().AllWindows()[0])
			return
		}
		run(pmCmd(pm, "search", pkg.Text), false)
	})

	row := container.NewHBox(pkg, installBtn, removeBtn, searchBtn)
	row2 := container.NewHBox(updateBtn, refreshBtn, busy)
	note := widget.NewLabel(t("lbl.pmNote"))
	note.Wrapping = fyne.TextWrapWord

	body := container.NewBorder(
		container.NewVBox(info, row, row2, note),
		nil, nil, nil,
		sc,
	)
	return body
}