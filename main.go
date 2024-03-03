package main

import (
	"net/url"
	"slices"

	ui "github.com/libui-ng/golang-ui"
	"github.com/skratchdot/open-golang/open"
)

var (
	초성 = []rune{'ㄱ', 'ㄲ', 'ㄴ', 'ㄷ',
		'ㄸ', 'ㄹ', 'ㅁ', 'ㅂ',
		'ㅃ', 'ㅅ', 'ㅆ', 'ㅇ',
		'ㅈ', 'ㅉ', 'ㅊ', 'ㅋ',
		'ㅌ', 'ㅍ', 'ㅎ'}
	중성 = []rune{'ㅏ', 'ㅐ', 'ㅑ', 'ㅒ',
		'ㅓ', 'ㅔ', 'ㅕ', 'ㅖ',
		'ㅗ', 'ㅘ', 'ㅙ', 'ㅚ',
		'ㅛ', 'ㅜ', 'ㅝ', 'ㅞ',
		'ㅟ', 'ㅠ', 'ㅡ', 'ㅢ',
		'ㅣ'}
	종성 = []rune{
		'ㄱ', 'ㄲ', 'ㄳ',
		'ㄴ', 'ㄵ', 'ㄶ', 'ㄷ',
		'ㄹ', 'ㄺ', 'ㄻ', 'ㄼ',
		'ㄽ', 'ㄾ', 'ㄿ', 'ㅀ',
		'ㅁ', 'ㅂ', 'ㅄ', 'ㅅ',
		'ㅆ', 'ㅇ', 'ㅈ', 'ㅊ',
		'ㅋ', 'ㅌ', 'ㅍ', 'ㅎ'}
)

func IsSyllable(c rune) bool {
	return c >= 0xAC00 && c <= 0xD7A3
}

func IsJamo(c rune) bool {
	return (c >= 0x1100 && c <= 0x11FF) || (c >= 0x3130 && c <= 0x318F)
}

func onClick(button *ui.Button, entry *ui.Entry) {
	key := button.Text()
	cur := []rune(key)[0]
	text := []rune(entry.Text())
	if len(text) == 0 {
		entry.SetText(key)
		return
	}
	last := text[len(text)-1]
	if IsJamo(last) {
		if last == cur {
			text = text[:len(text)-1]
			entry.SetText(string(append(text, last+1)))
			return
		}
		text = text[:len(text)-1]
		syll := toSyllable(last, cur, 0)
		entry.SetText(string(append(text, syll)))
		return
	}
	if !IsSyllable(last) {
		entry.SetText(string(text) + key)
		return
	}

	first, second, third := toJamo(last)
	if third > 0 {
		switch cur {
		case 'ㅏ', 'ㅐ', 'ㅑ', 'ㅒ', 'ㅓ', 'ㅔ', 'ㅕ', 'ㅖ', 'ㅗ', 'ㅛ', 'ㅜ', 'ㅠ', 'ㅡ', 'ㅣ':
			text = text[:len(text)-1]
			type Composite struct {
				Is            bool
				First, Second rune
			}
			var composite Composite
			switch third {
			case 'ㄲ':
				composite = Composite{true, 'ㄱ', 'ㄱ'}
			case 'ㄳ':
				composite = Composite{true, 'ㄱ', 'ㅅ'}
			case 'ㄵ':
				composite = Composite{true, 'ㄴ', 'ㅈ'}
			case 'ㄶ':
				composite = Composite{true, 'ㄴ', 'ㅎ'}
			case 'ㄺ':
				composite = Composite{true, 'ㄹ', 'ㄱ'}
			case 'ㄻ':
				composite = Composite{true, 'ㄹ', 'ㅁ'}
			case 'ㄼ':
				composite = Composite{true, 'ㄹ', 'ㅂ'}
			case 'ㄽ':
				composite = Composite{true, 'ㄹ', 'ㅅ'}
			case 'ㄾ':
				composite = Composite{true, 'ㄹ', 'ㅌ'}
			case 'ㄿ':
				composite = Composite{true, 'ㄹ', 'ㅍ'}
			case 'ㅀ':
				composite = Composite{true, 'ㄹ', 'ㅎ'}
			case 'ㅄ':
				composite = Composite{true, 'ㅂ', 'ㅅ'}
			case 'ㅆ':
				composite = Composite{true, 'ㅅ', 'ㅅ'}
			}
			var prev, add rune
			if composite.Is {
				prev = toSyllable(first, second, composite.First)
				add = toSyllable(composite.Second, cur, 0)
			} else {
				prev = toSyllable(first, second, 0)
				add = toSyllable(third, cur, 0)
			}
			entry.SetText(string(append(append(text, prev), add)))
			return
		}
		var composite rune
		switch third {
		case 'ㄱ':
			switch cur {
			case 'ㄱ':
				composite = 'ㄲ'
			case 'ㅅ':
				composite = 'ㄳ'
			}
		case 'ㄴ':
			switch cur {
			case 'ㅈ':
				composite = 'ㄵ'
			}
		case 'ㄹ':
			switch cur {
			case 'ㄱ':
				composite = 'ㄺ'
			case 'ㅁ':
				composite = 'ㄻ'
			case 'ㅂ':
				composite = 'ㄼ'
			case 'ㅅ':
				composite = 'ㄽ'
			case 'ㅌ':
				composite = 'ㄾ'
			case 'ㅍ':
				composite = 'ㄿ'
			case 'ㅎ':
				composite = 'ㅀ'
			}
		case 'ㅂ':
			switch cur {
			case 'ㅅ':
				composite = 'ㅄ'
			}
		case 'ㅅ':
			switch cur {
			case 'ㅅ':
				composite = 'ㅆ'
			}
		}
		if composite > 0 {
			text = text[:len(text)-1]
			syll := toSyllable(first, second, composite)
			entry.SetText(string(append(text, syll)))
			return
		}
		entry.SetText(string(append(text, cur)))
		return
	}
	text = text[:len(text)-1]
	var composite rune
	switch second {
	case 'ㅗ':
		switch cur {
		case 'ㅏ':
			composite = 'ㅘ'
		case 'ㅐ':
			composite = 'ㅙ'
		case 'ㅣ':
			composite = 'ㅚ'
		}
	case 'ㅜ':
		switch cur {
		case 'ㅓ':
			composite = 'ㅝ'
		case 'ㅔ':
			composite = 'ㅞ'
		case 'ㅣ':
			composite = 'ㅟ'
		}
	case 'ㅡ':
		switch cur {
		case 'ㅣ':
			composite = 'ㅢ'
		}
	}
	var syll rune
	if composite > 0 {
		syll = toSyllable(first, composite, 0)
	} else {
		syll = toSyllable(first, second, cur)
	}
	entry.SetText(string(append(text, syll)))
}

func start() {
	win := ui.NewWindow("자판 - 字板", 450, 120, true)
	win.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		win.Destroy()
		return true
	})
	layout := [][]rune{
		[]rune("ㅂㅈㄷㄱㅅㅛㅕㅑㅐㅔ"),
		[]rune("ㅁㄴㅇㄹㅎㅗㅓㅏㅣㅒ"),
		[]rune("ㅋㅌㅊㅍㅠㅜㅡㅖ"),
	}
	tb := ui.NewEntry()
	kb := ui.NewVerticalBox()
	for _, r := range layout {
		row := ui.NewHorizontalBox()
		row.Append(ui.NewHorizontalSeparator(), true)
		for _, c := range r {
			b := ui.NewButton(string(c))
			b.OnClicked(func(button *ui.Button) {
				onClick(button, tb)
			})
			row.Append(b, false)
		}
		row.Append(ui.NewHorizontalSeparator(), true)
		kb.Append(row, true)
	}
	kb.Append(tb, true)
	bus := ui.NewHorizontalBox()
	bus.Append(ui.NewHorizontalSeparator(), true)
	{
		bu := ui.NewButton("Hanja")
		bu.OnClicked(func(button *ui.Button) {
			text := tb.Text()
			if len(text) > 0 {
				hangul := url.PathEscape(text)
				open.Start(`https://koreanhanja.app/` + hangul)
			}
		})
		bus.Append(bu, false)
	}
	{
		bu := ui.NewButton("English")
		bu.OnClicked(func(button *ui.Button) {
			text := tb.Text()
			if len(text) > 0 {
				params := url.Values{
					"source": []string{"osdd"},
					"sl":     []string{"ko"},
					"tl":     []string{"en"},
					"text":   []string{text},
					"op":     []string{"translate"},
				}
				open.Start(`https://translate.google.de/?` + params.Encode())
			}
		})
		bus.Append(bu, false)
	}
	bus.Append(ui.NewHorizontalSeparator(), true)
	kb.Append(bus, true)
	win.SetChild(kb)
	SetIconFromFile(win, "./Pictogrammers-Material-Syllabary-hangul.128.png")
	win.Show()
}

const GA = '가'

func toSyllable(first, second, third rune) rune {
	cho := slices.Index(초성, first)
	if cho == -1 {
		return GA
	}
	jung := slices.Index(중성, second)
	if jung == -1 {
		return GA
	}
	if third >= 0 {
		jong := slices.Index(종성, third)
		if jong >= 0 {
			return rune(GA + 588*cho + 28*jung + jong + 1)
		}
	}
	return rune(GA + 588*cho + 28*jung)
}

func toJamo(syllable rune) (rune, rune, rune) {
	cho := (syllable - GA) / 21 / 28
	jung := ((syllable - GA) - (cho * 21 * 28)) / 28
	jong := (syllable - GA) - (cho * 21 * 28) - (jung * 28)
	if jong > 0 {
		return 초성[cho], 중성[jung], 종성[jong-1]
	}
	return 초성[cho], 중성[jung], 0
}

func main() {
	ui.Main(start)
}
