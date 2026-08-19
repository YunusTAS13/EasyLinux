package main

import (
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func systemModule() fyne.CanvasObject {
	rt, sc := outBox()

	refresh := func() {
		runSafe(15*time.Second, func() string {
			var b strings.Builder
			b.WriteString("## " + t("sys.title") + "\n\n")
			b.WriteString("- **" + t("sys.os") + ":** " + sh("grep PRETTY_NAME /etc/os-release | cut -d= -f2 | tr -d '\"'") + "\n")
			b.WriteString("- **" + t("sys.kernel") + ":** " + sh("uname -r") + "  (" + sh("uname -m") + ")\n")
			b.WriteString("- **" + t("sys.host") + ":** " + sh("hostname") + "  **" + t("sys.user") + ":** " + sh("whoami") + "\n")
			b.WriteString("- **" + t("sys.uptime") + ":** " + sh("uptime -p") + "\n")
			b.WriteString("- **" + t("sys.cpu") + ":** " + sh("lscpu | grep -m1 'Model name' | cut -d: -f2 | xargs") + "\n")
			b.WriteString("- **" + t("sys.ram") + ":** " + sh("free -h | awk 'NR==2{print $2\" " + t("sys.total") + ", \"$3\" " + t("sys.used") + "\"}'") + "\n")
			b.WriteString("- **" + t("sys.gpu") + ":** " + sh("lspci | grep -iE 'vga|3d|display' | head -1 | sed 's/^[^:]*: //'") + "\n")
			b.WriteString("- **" + t("sys.root") + ":** " + sh("df -h / | awk 'NR==2{print $2\" " + t("sys.total") + ", \"$4\" " + t("sys.free") + " (\"$5\" " + t("sys.used") + ")\"}'") + "\n")
			b.WriteString("- **" + t("sys.ip") + ":** " + sh("hostname -I 2>/dev/null | awk '{print $1}'") + "\n")
			b.WriteString("- **" + t("sys.desktop") + ":** " + sh("echo ${XDG_CURRENT_DESKTOP:-${DESKTOP_SESSION:-?}}") + "\n")
			return b.String()
		}, func(out string) { setMarkdown(rt, out) })
	}
	refresh()

	top := container.NewHBox(
		widget.NewButtonWithIcon(t("btn.refresh"), theme.ViewRefreshIcon(), refresh),
		widget.NewButtonWithIcon(t("btn.copy"), theme.ContentCopyIcon(), func() {
			fyne.CurrentApp().Clipboard().SetContent(rt.String())
		}),
	)
	return container.NewBorder(top, nil, nil, nil, sc)
}