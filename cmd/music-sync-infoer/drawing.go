package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"math"
)

type drawer struct {
	tcell.Screen
	w, h int
}

func (d *drawer) eventLoop() {
	for {
		ev := d.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyCtrlC {
				return
			}
		case *tcell.EventResize:
			d.Sync()
			d.w, d.h = d.Size()
		default:
			panic(fmt.Sprintf("%T", ev))
		}
	}
}

func (d *drawer) drawString(x, y int, style tcell.Style, str string) {
	for i, r := range str {
		d.SetContent(x+i, y, r, nil, style)
	}
}

func (d *drawer) drawProgress(x, y int, style tcell.Style, length int, progress float64) {
	head, filledRune, headRune, emptyRune := d.progressHead(length, progress)

	for i := 0; i < length; i++ {
		var r rune
		if i < head {
			r = filledRune
		} else if head == i {
			r = headRune
		} else {
			r = emptyRune
		}
		d.SetContent(x+i, y, r, nil, style)
	}
}

func (d *drawer) progressHead(length int, progress float64) (head int, filledRune, headRune, emptyRune rune) {
	head = int(math.Floor(float64(length) * progress))
	_, headProgress := math.Modf(float64(length) * progress)

	filledRune = '█'
	emptyRune = ' '

	if head == length-1 {
		headRune = filledRune
	} else if headProgress < 0.3 {
		headRune = emptyRune
	} else if headProgress < 0.7 {
		headRune = '▌'
	} else {
		headRune = filledRune
	}
	return
}

func (d *drawer) drawBox(x, y, w, h int, style tcell.Style) {
	for i := x; i < x+w; i++ {
		d.SetContent(i, y+h-1, '═', nil, style)
		d.SetContent(i, y, '═', nil, style)
	}

	for j := y; j < y+h; j++ {
		d.SetContent(x, j, '║', nil, style)
		d.SetContent(x+w-1, j, '║', nil, style)
	}
	d.SetContent(x, y, '╔', nil, style)
	d.SetContent(x+w-1, y, '╗', nil, style)
	d.SetContent(x, y+h-1, '╚', nil, style)
	d.SetContent(x+w-1, y+h-1, '╝', nil, style)
}
