package main

import (
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type scoreRec struct {
	Correct int `json:"correct"`
	Wrong   int `json:"wrong"`
}

func scorePath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = "."
	}
	return filepath.Join(dir, "EasyLinux", "skor.json")
}

func loadScore() scoreRec {
	var s scoreRec
	data, err := os.ReadFile(scorePath())
	if err == nil {
		json.Unmarshal(data, &s)
	}
	return s
}

func saveScore(s scoreRec) {
	os.MkdirAll(filepath.Dir(scorePath()), 0o755)
	data, _ := json.Marshal(s)
	os.WriteFile(scorePath(), data, 0o644)
}

func learnModule(w fyne.Window) fyne.CanvasObject {
	rand.Seed(42)
	var order []cmdInfo
	idx := 0
	seen := 0
	score := loadScore()

	question := widget.NewRichTextFromMarkdown(t("lbl.startQ"))
	question.Wrapping = fyne.TextWrapWord
	answer := widget.NewEntry()
	answer.SetPlaceHolder(t("ph.answer"))
	answer.Disable()
	result := widget.NewLabel("")
	scoreLbl := widget.NewLabel("")
	bar := widget.NewProgressBar()
	bar.SetValue(0)

	refreshScore := func() {
		scoreLbl.SetText(t("lbl.score") + itoa(score.Correct) + t("lbl.scoreCor") + itoa(score.Wrong) + t("lbl.scoreWro"))
	}

	nextQ := func() {
		if idx >= len(order) {
			question.ParseMarkdown(t("res.done"))
			answer.Disable()
			answer.SetText("")
			result.SetText(t("res.done"))
			bar.SetValue(1)
			return
		}
		c := order[idx]
		question.ParseMarkdown("**" + t("lbl.qnum") + itoa(seen+1) + "/" + itoa(len(order)) + "**\n\n" + t("lbl.question") + "**" + cmdAct(c) + "**\n\n> " + cmdNe(c))
		answer.Enable()
		answer.SetText("")
		result.SetText("")
	}

	start := func() {
		order = order[:0]
		for _, c := range allCommands {
			order = append(order, c)
		}
		rand.Shuffle(len(order), func(i, j int) { order[i], order[j] = order[j], order[i] })
		idx = 0
		seen = 0
		bar.SetValue(0)
		nextQ()
	}

	check := func() {
		if idx >= len(order) {
			return
		}
		got := strings.ToLower(strings.TrimSpace(answer.Text))
		ok := false
		c := order[idx]
		for _, a := range c.ans {
			if strings.ToLower(a) == got {
				ok = true
				break
			}
		}
		if ok {
			score.Correct++
			result.SetText(t("res.correct"))
			result.TextStyle = fyne.TextStyle{Bold: true}
		} else {
			score.Wrong++
			result.SetText(t("res.wrong") + c.cmd)
			result.TextStyle = fyne.TextStyle{Bold: true}
		}
		saveScore(score)
		refreshScore()
		seen++
		idx++
		bar.SetValue(float64(seen) / float64(len(order)))
		answer.Disable()
	}

	checkBtn := widget.NewButtonWithIcon(t("btn.check"), theme.ConfirmIcon(), check)
	nextBtn := widget.NewButtonWithIcon(t("btn.next"), theme.MediaPlayIcon(), nextQ)
	startBtn := widget.NewButtonWithIcon(t("btn.startRound"), theme.ViewRefreshIcon(), start)
	resetBtn := widget.NewButtonWithIcon(t("btn.resetScore"), theme.DeleteIcon(), func() {
		score = scoreRec{}
		saveScore(score)
		refreshScore()
	})

	top := container.NewVBox(
		widget.NewLabel(t("lbl.learnDesc")),
		scoreLbl,
		bar,
	)
	mid := container.NewVBox(
		question,
		container.NewHBox(answer, checkBtn, nextBtn),
		container.NewHBox(startBtn, resetBtn),
		result,
	)
	refreshScore()
	return container.NewBorder(top, nil, nil, nil, container.NewVBox(mid, widget.NewLabel("")))
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	if neg {
		b = append([]byte{'-'}, b...)
	}
	return string(b)
}