package window

import (
	"github.com/omar0ali/sysmontui/screentui"
	"github.com/omar0ali/sysmontui/screentui/interfaces"
)

func LineHorizontal(s interfaces.ScreenControl, length, yPos int, char rune) {
	for i := range length {
		s.SetContent(i, yPos, char)
	}
}

func Text(s interfaces.ScreenControl, p screentui.Point, txt string) {
	for i, char := range txt {
		s.SetContent(int(p.X)+i, int(p.Y), char)
	}
}

func TextWrap(s interfaces.ScreenControl, p screentui.Point, width int, txt string) {
	x, y := int(p.X), int(p.Y)
	line := ""
	for _, char := range txt {
		if len(line) >= width {
			for i, c := range line {
				s.SetContent(x+i, y, c)
			}
			line = ""
			y++
		}
		line += string(char)
	}
	// Draw remaining line
	for i, c := range line {
		s.SetContent(x+i, y, c)
	}
}

func TextCenter(s interfaces.ScreenControl, p screentui.Point, txt string) {
	width, _ := s.Size()
	startX := int(p.X) + (width-len(txt))/2
	for i, char := range txt {
		s.SetContent(startX+i, int(p.Y), char)
	}
}

func TextBoxCenter(s interfaces.ScreenControl, topLeft screentui.Point, boxWidth, boxHeight int, txt string) {
	lines := []string{}
	line := ""
	for _, c := range txt {
		if len(line) >= boxWidth {
			lines = append(lines, line)
			line = ""
		}
		line += string(c)
	}
	if line != "" {
		lines = append(lines, line)
	}

	startY := int(topLeft.Y) + (boxHeight-len(lines))/2
	for i, l := range lines {
		startX := int(topLeft.X) + (boxWidth-len(l))/2
		for j, c := range l {
			s.SetContent(startX+j, startY+i, c)
		}
	}
}

func TextScreenCenter(s interfaces.ScreenControl, txt string) {
	screenWidth, screenHeight := s.Size()
	x := (screenWidth - len(txt)) / 2
	y := screenHeight / 2
	for i, c := range txt {
		s.SetContent(x+i, y, c)
	}
}

func DrawBox(s interfaces.ScreenControl, p screentui.Point, width, height int) {
	x := int(p.X)
	y := int(p.Y)
	// Draw corners
	s.SetContent(x, y, '┌')         // top-left
	s.SetContent(x+width-1, y, '┐') // top-right
	s.SetContent(x, y+height-1, '└')
	s.SetContent(x+width-1, y+height-1, '┘')

	// Draw top and bottom borders
	for i := 1; i < width-1; i++ {
		s.SetContent(x+i, y, '─')
		s.SetContent(x+i, y+height-1, '─') // bottom
	}

	// Draw left and right borders
	for i := 1; i < height-1; i++ {
		s.SetContent(x, y+i, '│')         // left
		s.SetContent(x+width-1, y+i, '│') // right
	}
}
