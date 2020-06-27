package choRenderer

import (
	"bufio"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/rasterizer"
	"math"
	"os"
	"strings"
)

func (song *Song) getTextBoxBounds(fontSize float64, str string, c *canvas.Canvas) canvas.Rect {
	face := song.fontFamily.Face(fontSize, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	var box = canvas.NewTextLine(face, str, canvas.Left)

	return box.Bounds()
}

func (song *Song) RenderSong() {
	file, err := os.Open(song.PathToFile)
	handle(err)
	scanner := bufio.NewScanner(file)

	var sections []Section

	var section = Section{}

	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "CCLI") {
			break
		} else if strings.HasPrefix(text, "{") {
			key, value := parseTag(text)
			section.tags[key] = value
		} else if len(text) > 0 {
			section.lines = append(section.lines, parseLine(text))
		} else {
			sections = append(sections, section)
			section = Section{}
		}

	}
	sections = append(sections, section)
	section = Section{}

	song.renderSections(sections)
}

func (song *Song) initCanvas() (c *canvas.Canvas, context *canvas.Context) {
	song.fontFamily = canvas.NewFontFamily("Ubuntu")
	song.fontFamily.Use(canvas.CommonLigatures)
	if err := song.fontFamily.LoadFontFile("/usr/share/fonts/truetype/ubuntu/Ubuntu-M.ttf", canvas.FontRegular); err != nil {
		panic(err)
	}

	c = canvas.New(song.Resolution.W, song.Resolution.H)
	context = canvas.NewContext(c)

	return c, context
}

func (song *Song) renderSections(sections []Section) {
	for _, section := range sections {
		if len(section.tags) > 0 && section.tags["comment"] != "" {
			song.renderSection(section)
		}
	}
}

func (song *Song) renderSection(section Section) {

	c, ctx := song.initCanvas()

	//setUp canvas
	ctx.SetFillColor(canvas.Black)
	fontSize, hMax, wMax := song.calcFontSize(section, c)

	song.calcPixelOffset(&section, fontSize, c)

	lineOffset := math.Max(hMax, 0)
	yOffset := math.Max((ctx.Height()-float64(len(section.lines)*2)*hMax)/2, 0)
	xOffset := math.Max((ctx.Width()-wMax)/2, 0)

	face := song.fontFamily.Face(fontSize, canvas.White, canvas.FontRegular, canvas.FontNormal)
	chordFace := song.fontFamily.Face(fontSize-40, canvas.White, canvas.FontRegular, canvas.FontNormal)

	i := 0
	for _, line := range section.lines {
		for _, chord := range line.chords {
			chordLine := canvas.NewTextLine(chordFace, chord.name, canvas.Left)
			y := (c.H - yOffset) - (lineOffset * float64(i))
			ctx.DrawText(xOffset+chord.pixelOffset, y, chordLine)
		}
		y := (c.H - yOffset) - (lineOffset * float64(i+1))
		lineBox := canvas.NewTextLine(face, line.lyrics, canvas.Left)
		ctx.DrawText(xOffset, y, lineBox)
		i += 2
	}

	name := section.tags["comment"]

	err := c.WriteFile(name+".png", rasterizer.PNGWriter(1.0))

	handle(err)
}

func (song *Song) calcFontSize(section Section, c *canvas.Canvas) (pnt, hMax, wMax float64) {

	fontSize := 12.0
	fontHeight := 0.0
	fontWidth := 0.0

	lines := section.lines

	longestLine := ""
	for _, line := range lines {
		if len(line.lyrics) > len(longestLine) {
			longestLine = line.lyrics
		}
	}

	if !strings.ContainsAny(longestLine, "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM") {
		return fontSize, 0.0, 0.0
	}

	for fontWidth < c.W && (fontHeight*2.0*float64(len(lines)) < c.H) {
		size := song.getTextBoxBounds(fontSize, longestLine, c)

		fontHeight = size.H
		fontWidth = size.W

		fontSize += 1
		//fmt.Printf("Testing font %f \n", fontSize) TODO: Remove debug code
	}
	return fontSize - 1, fontHeight, fontWidth
}

func (song *Song) calcPixelOffset(section *Section, fontSize float64, c *canvas.Canvas) {
	for _, line := range section.lines {
		for _, chord := range line.chords {
			chord.pixelOffset = song.getTextBoxBounds(fontSize, line.lyrics[0:chord.charOffset], c).W
		}
	}
}
