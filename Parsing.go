package choRenderer

import (
	"fmt"
	"github.com/tdewolff/canvas"
	"image/color"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

//TODO: Make file path and song more fleshed out
type Song struct {
	sections   []Section
	Name       string
	PathToFile string
	FontPath   string
	FontColor  color.RGBA

	NashvilleNumber bool

	SlideDelimiter string

	Resolution Rect

	fontFamily *canvas.FontFamily
}

type SongSettings struct {
	Name       string
	PathToFile string
	FontPath   string
	FontColor  string

	NashvilleNumber string

	SlideDelimiter string

	Height string
	Width  string
}

func (song *Song) LoadSettings(settings *SongSettings) {

	if settings.Name == "" {
		song.Name = getName(settings.PathToFile)
	} else {
		song.Name = settings.Name
	}

	hexColor, err := parseHexColor(settings.FontColor)
	song.FontColor = hexColor

	song.NashvilleNumber, err = strconv.ParseBool(settings.NashvilleNumber)
	if err != nil {
		log.Fatal("Couldn't convert string to bool")
	}

	song.SlideDelimiter = settings.SlideDelimiter

	song.PathToFile = settings.PathToFile
	song.FontPath = settings.FontPath

	height_float, err := strconv.ParseFloat(settings.Height, 64)
	if err != nil {
		log.Fatal("Couldn't convert string to float 64")
	}
	width_float, err := strconv.ParseFloat(settings.Width, 64)
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

func parseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")
	}
	return
}

var chromaticScale = map[string]int{
	"ab": 1,
	"a":  2,
	"a#": 3,
	"bb": 3,
	"b":  4,
	"b#": 5,
	"cb": 4,
	"c":  5,
	"c#": 6,
	"db": 6,
	"d":  7,
	"d#": 8,
	"eb": 8,
	"e":  9,
	"e#": 10,
	"fb": 9,
	"f":  10,
	"f#": 11,
	"gb": 11,
	"g":  12,
}

var scaleMap = map[int]int{
	0:  1,
	2:  2,
	4:  3,
	5:  4,
	7:  5,
	9:  6,
	11: 7,
}

func (song *Song) convertToNashville() {
	key := strings.ToLower(song.sections[0].tags["key"])
	keyNumber := chromaticScale[key]

	for _, section := range song.sections {
		for _, line := range section.lines {
			for i, chord := range line.chords {
				chordName := strings.ToLower(chord.name)
				var chordNumber int

				if !strings.ContainsAny(chordName, "abcdefg") {
					continue
				}

				if len(chordName) > 1 {
					match, err := regexp.MatchString(`[a-g]b|[a-g]#`, chordName[0:2])
					if err != nil {
						log.Fatal("Error parsing regex for nashville numbers")
					}
					if match {
						chordNumber = scaleMap[chromaticScale[chordName[0:2]]-keyNumber]
						line.chords[i].name = strconv.Itoa(chordNumber) + chordName[2:]
					} else {
						chordNumber = scaleMap[chromaticScale[chordName[0:1]]-keyNumber]
						line.chords[i].name = strconv.Itoa(chordNumber) + chordName[1:]
					}
				} else {
					chordNumber = scaleMap[chromaticScale[chordName]-keyNumber]
					line.chords[i].name = strconv.Itoa(chordNumber)
				}
			}
		}
	}
}
