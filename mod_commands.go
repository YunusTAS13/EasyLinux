package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type cmdInfo struct {
	cmd string
	ans []string
}

var allCommands = []cmdInfo{
	{"ls", []string{"ls", "ls -la"}},
	{"cd", []string{"cd"}},
	{"pwd", []string{"pwd"}},
	{"mkdir", []string{"mkdir"}},
	{"rmdir", []string{"rmdir"}},
	{"rm", []string{"rm", "rm -r", "rm -rf"}},
	{"cp", []string{"cp"}},
	{"mv", []string{"mv"}},
	{"cat", []string{"cat"}},
	{"nano", []string{"nano"}},
	{"grep", []string{"grep"}},
	{"find", []string{"find"}},
	{"chmod", []string{"chmod"}},
	{"sudo", []string{"sudo"}},
	{"apt", []string{"apt", "apt install", "sudo apt install"}},
	{"pacman", []string{"pacman", "pacman -S", "sudo pacman -S"}},
	{"df", []string{"df", "df -h"}},
	{"du", []string{"du"}},
	{"free", []string{"free", "free -h"}},
	{"top", []string{"top"}},
	{"htop", []string{"htop"}},
	{"ps", []string{"ps"}},
	{"kill", []string{"kill"}},
	{"ping", []string{"ping"}},
	{"ip", []string{"ip", "ip a"}},
	{"curl", []string{"curl"}},
	{"wget", []string{"wget"}},
	{"tar", []string{"tar"}},
	{"unzip", []string{"unzip"}},
	{"man", []string{"man"}},
	{"history", []string{"history"}},
	{"clear", []string{"clear"}},
	{"exit", []string{"exit"}},
	{"shutdown", []string{"shutdown", "shutdown -h now"}},
	{"reboot", []string{"reboot"}},
	{"ln", []string{"ln"}},
	{"whoami", []string{"whoami"}},
	{"hostname", []string{"hostname"}},
	{"date", []string{"date"}},
	{"systemctl", []string{"systemctl"}},
	{"ufw", []string{"ufw"}},
	{"xdg-open", []string{"xdg-open"}},
}

func cmdAct(c cmdInfo) string {
	act, _ := cmdTextFor(curLang, c.cmd)
	return act
}

func cmdNe(c cmdInfo) string {
	_, ne := cmdTextFor(curLang, c.cmd)
	return ne
}

func translate(s string) []cmdInfo {
	q := strings.ToLower(strings.TrimSpace(s))
	if q == "" {
		return nil
	}
	rules, ok := kwI18n[curLang]
	if !ok {
		rules = kwI18n["tr"]
	}
	var res []cmdInfo
	for _, rule := range rules {
		if strings.Contains(q, rule.key) {
for _, name := range rule.cmds {
			c := cmdByName(name)
			if c.cmd != "" {
				res = appendUnique(res, c)
			}
		}
		}
	}
	return res
}

func commandsModule() fyne.CanvasObject {
	search := widget.NewEntry()
	search.SetPlaceHolder(t("ph.searchCmd"))

	var list *widget.List
	var current []cmdInfo
	refresh := func() {
		q := strings.ToLower(strings.TrimSpace(search.Text))
		current = current[:0]
		for _, c := range allCommands {
			act := strings.ToLower(cmdAct(c))
			ne := strings.ToLower(cmdNe(c))
			if q == "" ||
				strings.Contains(strings.ToLower(c.cmd), q) ||
				strings.Contains(act, q) ||
				strings.Contains(ne, q) {
				current = append(current, c)
			}
		}
		list.Refresh()
	}

	detail := widget.NewRichTextFromMarkdown(t("lbl.selectCmd"))
	detail.Wrapping = fyne.TextWrapWord
	detailSc := container.NewScroll(detail)

	selectedID := -1

	list = widget.NewList(
		func() int { return len(current) },
		func() fyne.CanvasObject { return widget.NewLabel("...") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id < 0 || id >= len(current) {
				return
			}
			obj.(*widget.Label).SetText(current[id].cmd + "  —  " + cmdAct(current[id]))
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		selectedID = id
		if id < 0 || id >= len(current) {
			return
		}
		c := current[id]
		detail.ParseMarkdown("**" + t("cmd.name") + ":** `" + c.cmd + "`\n\n**" + t("cmd.purpose") + ":** " + cmdAct(c) + "\n\n**" + t("cmd.desc") + ":** " + cmdNe(c))
	}
	listSc := container.NewScroll(list)
	listSc.SetMinSize(fyne.NewSize(380, 0))

	search.OnChanged = func(string) { refresh() }
	refresh()

	copyBtn := widget.NewButtonWithIcon(t("btn.copyCmd"), theme.ContentCopyIcon(), func() {
		if selectedID >= 0 && selectedID < len(current) {
			fyne.CurrentApp().Clipboard().SetContent(current[selectedID].cmd)
		}
	})

	tIn := widget.NewEntry()
	tIn.SetPlaceHolder(t("ph.translate"))
	tOut := widget.NewRichTextFromMarkdown("")
	tOut.Wrapping = fyne.TextWrapWord
	tSc := container.NewScroll(tOut)
	tSc.SetMinSize(fyne.NewSize(0, 130))
	translateBtn := widget.NewButtonWithIcon(t("btn.translate"), theme.SearchIcon(), func() {
		m := translate(tIn.Text)
		if len(m) == 0 {
			tOut.ParseMarkdown("`ls`, `cd`, `mkdir`, `cp`, `mv`, `rm`, `cat`, `nano`, `sudo pacman -S <paket>`, `df -h`, `free -h`")
			return
		}
		var sb strings.Builder
		sb.WriteString("**" + t("btn.translate") + ":**\n\n")
		for _, c := range m {
			sb.WriteString("- `" + c.cmd + "` → " + cmdAct(c) + "\n")
		}
		tOut.ParseMarkdown(sb.String())
	})

	right := container.NewBorder(
		container.NewHBox(copyBtn),
		nil, nil, nil,
		container.NewVSplit(
			detailSc,
			container.NewBorder(
				container.NewVBox(
					widget.NewLabel(t("lbl.translator")),
					container.NewHBox(tIn, translateBtn),
				),
				nil, nil, nil,
				tSc,
			),
		),
	)

	left := container.NewBorder(
		container.NewVBox(search, widget.NewLabel(t("lbl.cmdList"))),
		nil, nil, nil,
		listSc,
	)

	split := container.NewHSplit(left, right)
	split.SetOffset(0.4)
	return split
}