package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func aboutModule() fyne.CanvasObject {
	title := widget.NewLabel(t("ab.title"))
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	sub := widget.NewLabel(t("ab.sub"))
	sub.Alignment = fyne.TextAlignCenter

	version := widget.NewLabel(t("ab.version"))
	dev := widget.NewLabel(t("ab.dev"))
	lang := widget.NewLabel(t("ab.lang"))
	modules := widget.NewLabel(t("ab.modules"))
	license := widget.NewLabel(t("ab.license"))
	repo := widget.NewLabel(t("ab.repo"))
	repo.TextStyle = fyne.TextStyle{Italic: true}
	tip := widget.NewLabel(t("ab.tip"))
	thanks := widget.NewLabel(t("ab.thanks"))
	thanks.TextStyle = fyne.TextStyle{Italic: true}

	info := container.NewVBox(
		title, sub,
		widget.NewSeparator(),
		version,
		dev,
		lang,
		modules,
		license,
		repo,
		widget.NewSeparator(),
		tip,
		thanks,
	)

	top := container.NewHBox(
		widget.NewButtonWithIcon(t("btn.refresh"), theme.ViewRefreshIcon(), func() {
			fyne.CurrentApp().SendNotification(fyne.NewNotification(
				t("ab.title"),
				fmt.Sprintf("%s • %s", t("ab.version"), t("ab.thanks")),
			))
		}),
	)

	body := container.NewCenter(container.NewVBox(info, top))
	return container.NewBorder(nil, nil, nil, nil, body)
}