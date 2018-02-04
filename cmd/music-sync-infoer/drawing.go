package main

import (
	"github.com/gdamore/tcell"
	"math"
)

func drawString(x, y int, style tcell.Style, str string, screen tcell.Screen) {
	for i, r := range str {
		screen.SetContent(x+i, y, r, nil, style)
	}
}

func drawProgress(x, y int, style tcell.Style, length int, progress float64, screen tcell.Screen) {
	head := int(math.Floor(float64(length) * progress))
	_, headProgress := math.Modf(float64(length) * progress)

	filledRune := '█'
	emptyRune := ' '
	var headRune rune
	if headProgress < 0.3 {
		headRune = emptyRune
	} else if headProgress < 0.7 {
		headRune = '▌'
	} else {
		headRune = filledRune
	}
	if head == length-1 {
		headRune = filledRune
	}

	for i := 0; i < length; i++ {
		var r rune
		if i < head {
			r = filledRune
		} else if head == i {
			r = headRune
		} else {
			r = emptyRune
		}
		screen.SetContent(x+i, y, r, nil, style)
	}
}

func drawBox(x, y, w, h int, style tcell.Style, screen tcell.Screen) {
	for i := x; i < x+w; i++ {
		screen.SetContent(i, y+h-1, '═', nil, style)
		screen.SetContent(i, y, '═', nil, style)
	}

	for j := y; j < y+h; j++ {
		screen.SetContent(x, j, '║', nil, style)
		screen.SetContent(x+w-1, j, '║', nil, style)
	}
	screen.SetContent(x, y, '╔', nil, style)
	screen.SetContent(x+w-1, y, '╗', nil, style)
	screen.SetContent(x, y+h-1, '╚', nil, style)
	screen.SetContent(x+w-1, y+h-1, '╝', nil, style)
}
