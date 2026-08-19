package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func sh(cmd string) string {
	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	s := strings.TrimSpace(string(out))
	if err != nil {
		if s != "" {
			s += "\n"
		}
		s += "[hata: " + err.Error() + "]"
	}
	return s
}

func shPriv(cmd string) string {
	if _, err := exec.LookPath("pkexec"); err == nil {
		out, err := exec.Command("pkexec", "sh", "-c", cmd).CombinedOutput()
		s := strings.TrimSpace(string(out))
		if err != nil {
			if s != "" {
				s += "\n"
			}
			s += "[hata: " + err.Error() + "]"
		}
		return s
	}
	c := exec.Command("sudo", "sh", "-c", cmd)
	done := make(chan string, 1)
	go func() {
		out, err := c.CombinedOutput()
		s := strings.TrimSpace(string(out))
		if err != nil {
			if s != "" {
				s += "\n"
			}
			s += "[hata: " + err.Error() + "]"
		}
		done <- s
	}()
	select {
	case s := <-done:
		return s
	case <-time.After(6 * time.Second):
		c.Process.Kill()
		return "Zaman aşımı: yönetici işlemi için polkit (pkexec) gerekli."
	}
}

func outBox() (*widget.RichText, *container.Scroll) {
	rt := widget.NewRichTextFromMarkdown("```\n\n```")
	rt.Wrapping = fyne.TextWrapWord
	sc := container.NewScroll(rt)
	sc.SetMinSize(fyne.NewSize(0, 260))
	return rt, sc
}

func setOut(rt *widget.RichText, txt string) {
	if strings.TrimSpace(txt) == "" {
		txt = "(boş çıktı)"
	}
	content := txt
	doUI(func() {
		rt.ParseMarkdown("```\n" + content + "\n```")
	})
}

func setMarkdown(rt *widget.RichText, md string) {
	content := md
	if strings.TrimSpace(md) == "" {
		content = "_boş_"
	}
	doUI(func() {
		rt.ParseMarkdown(content)
	})
}

func runSafe(timeout time.Duration, fn func() string, done func(string)) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				if done != nil {
					doUI(func() { done("[hata: " + fmt.Sprint(r) + "]") })
				}
			}
		}()
		ch := make(chan string, 1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					ch <- "[hata: " + fmt.Sprint(r) + "]"
				}
			}()
			ch <- fn()
		}()
		var res string
		if timeout <= 0 {
			res = <-ch
		} else {
			select {
			case res = <-ch:
			case <-time.After(timeout):
				res = "[zaman aşımı: işlem " + timeout.String() + " içinde bitmedi]"
			}
		}
		if done != nil {
			doUI(func() { done(res) })
		}
	}()
}

func appendUnique(list []cmdInfo, c cmdInfo) []cmdInfo {
	for _, x := range list {
		if x.cmd == c.cmd {
			return list
		}
	}
	return append(list, c)
}

func cmdByName(name string) cmdInfo {
	for _, c := range allCommands {
		if c.cmd == name {
			return c
		}
	}
	return cmdInfo{}
}