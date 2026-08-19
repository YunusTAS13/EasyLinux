package main

import (
	"embed"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

//go:embed fonts/DejaVuSans.ttf fonts/DejaVuSans-Bold.ttf fonts/DejaVuSansMono.ttf
var fontFS embed.FS

var (
	resRegular = mustResource("fonts/DejaVuSans.ttf")
	resBold    = mustResource("fonts/DejaVuSans-Bold.ttf")
	resMono    = mustResource("fonts/DejaVuSansMono.ttf")
)

func mustResource(path string) fyne.Resource {
	data, err := fontFS.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return fyne.NewStaticResource(path, data)
}

type bundleTheme struct {
	variant fyne.ThemeVariant
}

func (t bundleTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace {
		return resMono
	}
	if style.Bold {
		return resBold
	}
	return resRegular
}

func (t bundleTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t bundleTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	base := theme.LightTheme()
	if t.variant == theme.VariantDark {
		base = theme.DarkTheme()
	}
	return base.Color(name, variant)
}

func (t bundleTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return 16
	}
	return theme.DefaultTheme().Size(name)
}

func defaultTheme() bundleTheme {
	return bundleTheme{variant: theme.VariantLight}
}

func doUI(f func()) {
	fyne.Do(f)
}