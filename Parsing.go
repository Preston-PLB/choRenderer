package choRenderer

import (
	"github.com/tdewolff/canvas"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

//TODO: Make file path and song more fleshed out
type Song struct {
	sections   []Section
	Name       string
	PathToFile string

	Resolution Rect

	FontPath string

	fontFamily *canvas.FontFamily
}

type SongSettings struct {
	Name       string
	PathToFile string
	FontPath   string

	height string
	width  string
}

func (song *Song) LoadSettings(settings *SongSettings) {
	song.Name = settings.Name
	song.PathToFile = settings.PathToFile
	song.FontPath = settings.FontPath

	height_float, err := strconv.ParseFloat(settings.height, 64)
	if err != nil {
		log.Fatal("Couldn't convert string to float 64")
	}
	width_float, err := strconv.ParseFloat(settings.width, 64)
	if err != nil {
		log.Fatal("Couldn't convert string to float 64")
	}
	song.Resolution = Rect{height_float, width_float}
}

type Rect struct {
	H float64
	W float64
}

type Section struct {
	lines []Line
	tags  map[string]string
}

func (section *Section) initSection() {
	section.tags = make(map[string]string)
}

type Line struct {
	lyrics string
	chords []*Chord
}

type Chord struct {
	name        string
	charOffset  int
	pixelOffset float64
}

type Tag struct {
	name  string
	value string
}

//
// Function
//

func getName(stringPath string) (name string) {
	lastSlash := int(math.Max(0, float64(strings.LastIndex(stringPath, "/"))))
	lastDot := strings.LastIndex(stringPath, ".cho")

	return stringPath[lastSlash+1 : lastDot]
}

func parseTag(byteLine string) (key string, value string) {
	raw := strings.Split(byteLine, ": ")
	return raw[0][1:], raw[1][0 : len(raw[1])-1]
}

func parseLine(byteLine string) (line Line) {
	var lyricRaw []byte
	for i, k := 0, 0; i < len(byteLine); i++ {
		if byteLine[i] == '[' {
			var chordName []byte
			for j := i + 1; j < len(byteLine); j++ {
				if byteLine[j] != ']' {
					chordName = append(chordName, byteLine[j])
				} else {
					i = j
					break
				}
			}
			chord := Chord{string(chordName), k, 0.0}
			line.chords = append(line.chords, new(Chord))
			line.chords[len(line.chords)-1] = &chord
		} else {
			if byteLine[i] != '\r' {
				lyricRaw = append(lyricRaw, byteLine[i])
			}
			k++
		}
	}

	line.lyrics = string(lyricRaw)

	return line
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func (song *Song) getOutputPath() (path string) {
	return trimDirectoryPath(song.PathToFile) + string(os.PathSeparator) + song.Name //TODO: fix this, make it more smooth
}

func trimDirectoryPath(path string) (newPath string) {
	if strings.Contains(path, string(os.PathSeparator)) {
		return path[0:strings.LastIndex(path, string(os.PathSeparator))]
	} else {
		return path
	}
}
