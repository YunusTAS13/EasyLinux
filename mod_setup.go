package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type profile struct {
	nameKey string
	descKey string
	pkgs    map[string][]string
	check   *widget.Check
}

var profiles = []*profile{
	{"pro0n", "pro0d", map[string][]string{
		"pacman": {"libreoffice-fresh", "okular", "zathura-pdf-mupdf"},
		"apt":    {"libreoffice", "okular", "zathura"},
		"dnf":    {"libreoffice", "okular", "zathura"},
	}, nil},
	{"pro1n", "pro1d", map[string][]string{
		"pacman": {"vlc", "gimp", "audacity", "ffmpeg", "rhythmbox"},
		"apt":    {"vlc", "gimp", "audacity", "ffmpeg", "rhythmbox"},
		"dnf":    {"vlc", "gimp", "audacity", "ffmpeg", "rhythmbox"},
	}, nil},
	{"pro2n", "pro2d", map[string][]string{
		"pacman": {"git", "code", "nodejs", "npm", "docker", "htop"},
		"apt":    {"git", "code", "nodejs", "npm", "docker.io", "htop"},
		"dnf":    {"git", "code", "nodejs", "npm", "docker", "htop"},
	}, nil},
	{"pro3n", "pro3d", map[string][]string{
		"pacman": {"steam", "lutris", "wine", "gamemode"},
		"apt":    {"steam", "lutris", "wine", "gamemode"},
		"dnf":    {"steam", "lutris", "wine", "gamemode"},
	}, nil},
	{"pro4n", "pro4d", map[string][]string{
		"pacman": {"htop", "tmux", "fastfetch", "curl", "wget", "unzip", "file-roller", "tree"},
		"apt":    {"htop", "tmux", "fastfetch", "curl", "wget", "unzip", "file-roller", "tree"},
		"dnf":    {"htop", "tmux", "fastfetch", "curl", "wget", "unzip", "file-roller", "tree"},
	}, nil},
}

func setupModule(w fyne.Window) fyne.CanvasObject {
	pm := detectPM()
	rt, sc := outBox()

	for _, p := range profiles {
		p := p
		p.check = widget.NewCheck(t(p.nameKey)+" — "+t(p.descKey), nil)
	}

	installBtn := widget.NewButtonWithIcon(t("btn.installSel"), theme.ConfirmIcon(), func() {
		var pkgs []string
		for _, p := range profiles {
			if p.check.Checked {
				pkgs = append(pkgs, p.pkgs[pm]...)
			}
		}
		if len(pkgs) == 0 {
			dialog.ShowInformation(t("dlg.noSel"), t("dlg.noSelMsg"), w)
			return
		}
		cmd := pmCmd(pm, "install", pkgs...)
		dialog.ShowConfirm(t("dlg.instPkgs"), t("dlg.instNote")+strings.Join(pkgs, ", ")+t("dlg.instCmd")+cmd+"\n\n"+t("dlg.copyConfirm")+"?", func(ok bool) {
			if !ok {
				return
			}
			runSafe(0, func() string {
				return shPriv(cmd)
			}, func(out string) {
				setOut(rt, "> "+cmd+"\n\n"+out)
				fyne.CurrentApp().SendNotification(fyne.NewNotification(t("dlg.installDone"), strings.Join(pkgs, ", ")))
			})
		}, w)
	})

	var checks []fyne.CanvasObject
	for _, p := range profiles {
		checks = append(checks, p.check)
	}

	top := container.NewVBox(
		widget.NewLabel(t("lbl.setupDesc")+pm),
		container.NewHBox(installBtn),
	)
	center := container.NewVBox(checks...)
	body := container.NewBorder(top, nil, nil, nil, container.NewVSplit(center, sc))
	return body
}