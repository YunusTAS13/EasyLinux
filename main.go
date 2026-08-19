package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var currentTheme = defaultTheme()
var buildingUI = false

func buildUI(a fyne.App, w fyne.Window) {
	buildingUI = true
	defer func() { buildingUI = false }()

	langSel := widget.NewSelect(langList(), func(name string) {
		if buildingUI {
			return
		}
		code := langNameToCode(name)
		setLang(code)
		a.Preferences().SetString("lang", code)
		buildUI(a, w)
	})
	langSel.PlaceHolder = t("top.lang")
	langSel.SetSelected(langNames[langIndex(curLang)])

	var updateThemeBtn func()
	themeBtn := widget.NewButton("", func() {
		if currentTheme.variant == theme.VariantDark {
			currentTheme = bundleTheme{variant: theme.VariantLight}
		} else {
			currentTheme = bundleTheme{variant: theme.VariantDark}
		}
		a.Preferences().SetBool("dark", currentTheme.variant == theme.VariantDark)
		a.Settings().SetTheme(currentTheme)
		updateThemeBtn()
	})
	updateThemeBtn = func() {
		if currentTheme.variant == theme.VariantDark {
			themeBtn.SetText(t("btn.light"))
			themeBtn.SetIcon(theme.ColorPaletteIcon())
		} else {
			themeBtn.SetText(t("btn.dark"))
			themeBtn.SetIcon(theme.ColorPaletteIcon())
		}
	}
	updateThemeBtn()

	top := container.NewHBox(
		widget.NewLabel(t("top.lang")+":"), langSel,
		widget.NewLabel("   "+t("top.theme")+":"), themeBtn,
	)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(t("tab.system"), theme.InfoIcon(), systemModule()),
		container.NewTabItemWithIcon(t("tab.commands"), theme.HelpIcon(), commandsModule()),
		container.NewTabItemWithIcon(t("tab.packages"), theme.DownloadIcon(), packagesModule()),
		container.NewTabItemWithIcon(t("tab.trouble"), theme.WarningIcon(), troubleModule()),
		container.NewTabItemWithIcon(t("tab.disk"), theme.MediaRecordIcon(), diskModule()),
		container.NewTabItemWithIcon(t("tab.windows"), theme.FolderOpenIcon(), windowsModule(w)),
		container.NewTabItemWithIcon(t("tab.setup"), theme.CheckButtonIcon(), setupModule(w)),
		container.NewTabItemWithIcon(t("tab.learn"), theme.ComputerIcon(), learnModule(w)),
		container.NewTabItemWithIcon(t("tab.dualboot"), theme.RadioButtonIcon(), dualbootModule(w)),
		container.NewTabItemWithIcon(t("tab.security"), theme.ConfirmIcon(), securityModule(w)),
		container.NewTabItemWithIcon(t("tab.about"), theme.InfoIcon(), aboutModule()),
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	w.SetContent(container.NewBorder(top, nil, nil, nil, tabs))
}

func langIndex(code string) int {
	for i, c := range langCodes {
		if c == code {
			return i
		}
	}
	return 0
}

func main() {
	a := app.NewWithID("tr.EasyLinux")
	prefs := a.Preferences()
	if l := prefs.String("lang"); l != "" {
		setLang(l)
	}
	if prefs.Bool("dark") {
		currentTheme = bundleTheme{variant: theme.VariantDark}
	} else {
		currentTheme = defaultTheme()
	}
	a.Settings().SetTheme(currentTheme)
	w := a.NewWindow(t("app.title"))
	w.Resize(fyne.NewSize(1100, 720))
	buildUI(a, w)
	w.ShowAndRun()
}