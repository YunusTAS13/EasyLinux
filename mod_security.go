package main

import (
	"os/exec"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func securityModule(w fyne.Window) fyne.CanvasObject {
	// ---- Güvenlik duvarı ----
	fwRT, fwSc := outBox()
	fwRT.Wrapping = fyne.TextWrapWord
	fwLabel := widget.NewLabel(t("lbl.fwTitle"))
	if _, err := exec.LookPath("ufw"); err != nil {
		fwLabel.SetText(t("lbl.fwTitle") + " (ufw yok; firewalld: systemctl status firewalld)")
	}

	fwRun := func(cmd string) {
		runSafe(30*time.Second, func() string {
			return shPriv(cmd + " 2>&1")
		}, func(out string) { setOut(fwRT, out) })
	}
	fwButtons := container.NewHBox(
		widget.NewButtonWithIcon(t("btn.status"), theme.InfoIcon(), func() { fwRun("ufw status verbose") }),
		widget.NewButtonWithIcon(t("btn.fwOn"), theme.ConfirmIcon(), func() { fwRun("ufw enable") }),
		widget.NewButtonWithIcon(t("btn.fwOff"), theme.CancelIcon(), func() { fwRun("ufw disable") }),
		widget.NewButtonWithIcon(t("btn.fwLog"), theme.DocumentIcon(), func() { fwRun("ufw logging on") }),
	)

	fwCard := container.NewBorder(
		container.NewVBox(fwLabel, fwButtons),
		nil, nil, nil,
		fwSc,
	)

	// ---- Servisler ----
	type svc struct{ name, active, sub, desc string }
	var services []svc
	var current []svc

	loadSvc := func() []svc {
		out := sh("systemctl list-units --type=service --no-pager --no-legend 2>/dev/null")
		var res []svc
		for _, line := range strings.Split(out, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			f := strings.Fields(line)
			if len(f) < 4 {
				continue
			}
			res = append(res, svc{strings.TrimSuffix(f[0], ".service"), f[2], f[3], strings.Join(f[4:], " ")})
		}
		return res
	}

	svcFilter := widget.NewEntry()
	svcFilter.SetPlaceHolder(t("ph.service"))

	var svcList *widget.List
	selectedSvc := -1
	refreshSvc := func() {
		q := strings.ToLower(svcFilter.Text)
		current = current[:0]
		for _, s := range services {
			if q == "" || strings.Contains(strings.ToLower(s.name), q) || strings.Contains(strings.ToLower(s.desc), q) {
				current = append(current, s)
			}
		}
		if svcList != nil {
			svcList.Refresh()
		}
	}

	svcList = widget.NewList(
		func() int { return len(current) },
		func() fyne.CanvasObject { return widget.NewLabel("...") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < 0 || id >= len(current) {
				return
			}
			s := current[id]
			obj.(*widget.Label).SetText(s.name + "  [" + s.active + "]  " + s.desc)
		},
	)
	svcList.OnSelected = func(id widget.ListItemID) {
		selectedSvc = id
	}
	svcFilter.OnChanged = func(string) { refreshSvc() }

	act := func(op string) {
		if selectedSvc < 0 || selectedSvc >= len(current) {
			return
		}
		name := current[selectedSvc].name
		runSafe(30*time.Second, func() string {
			return shPriv("systemctl " + op + " " + name + " 2>&1")
		}, func(out string) {
			fyne.CurrentApp().SendNotification(fyne.NewNotification(t("lbl.svcTitle"), out))
			services = loadSvc()
			doUI(refreshSvc)
		})
	}

	svcButtons := container.NewHBox(
		widget.NewButtonWithIcon(t("btn.svcStart"), theme.MediaPlayIcon(), func() { act("start") }),
		widget.NewButtonWithIcon(t("btn.svcStop"), theme.MediaStopIcon(), func() { act("stop") }),
		widget.NewButtonWithIcon(t("btn.svcRestart"), theme.ViewRefreshIcon(), func() { act("restart") }),
		widget.NewButtonWithIcon(t("btn.svcEnable"), theme.ConfirmIcon(), func() { act("enable --now") }),
		widget.NewButtonWithIcon(t("btn.svcDisable"), theme.CancelIcon(), func() { act("disable --now") }),
	)

	svcReload := widget.NewButtonWithIcon(t("btn.svcReload"), theme.ViewRefreshIcon(), func() {
		runSafe(15*time.Second, func() string {
			services = loadSvc()
			return ""
		}, func(string) { doUI(refreshSvc) })
	})

	svcTop := container.NewVBox(
		widget.NewLabel(t("lbl.svcTitle")),
		container.NewHBox(svcFilter, svcReload),
		svcButtons,
	)
	svcSc := container.NewScroll(svcList)
	svcCard := container.NewBorder(svcTop, nil, nil, nil, svcSc)

	services = loadSvc()
	refreshSvc()

	body := container.NewVSplit(fwCard, svcCard)
	return body
}