package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	svg "github.com/ajstarks/svgo"
)

// Number of notes
const notes int = 12

// 5-pentatonic, 7-natural
const gamaSize int = 3

// Directory of SVG
var DrawDir string = "gama" + strconv.Itoa(gamaSize)
var GenesisNotes []int
var GamaCounter int
var Wg sync.WaitGroup
var SVGConf = SvgConf{
	Width:  640,
	Height: 480,
}

type SvgConf struct {
	Width       int
	Height      int
	WhiteWidth  int
	WhiteHeight int
	BlackWidth  int
	BlackHeight int
}

func init() {
	for i := range notes {
		GenesisNotes = append(GenesisNotes, i+1)
	}
	fmt.Println("Number of notes", notes)
	fmt.Println("Gama size", gamaSize)
	os.Mkdir(DrawDir, 0755)
	SVGConf.WhiteWidth = SVGConf.Width / 7
	SVGConf.WhiteHeight = SVGConf.Height
	SVGConf.BlackWidth = 53
	SVGConf.BlackHeight = 320
	fmt.Println("Starting ...")
}

func main() {
	start := time.Now()
	shuffleNotes(GenesisNotes, [gamaSize]int{}, 0)
	//Wg.Wait()
	fmt.Println(GamaCounter)
	fmt.Println("Total:", GamaCounter, "gama's")
	fmt.Println("Runtime is", time.Since(start))
}

func shuffleNotes(notesSet []int, gama [gamaSize]int, pos int) {
	for _, v := range notesSet {
		var gama2 [gamaSize]int
		gama2 = gama
		notesSet = notesSet[1:]
		if pos < gamaSize {
			gama2[pos] = v
		}
		if pos < gamaSize-1 {
			shuffleNotes(notesSet, gama2, pos+1)
		}
		if pos == gamaSize-1 {
			GamaCounter++
			drawGama(gama2)
			fmt.Println(gama2)
		}
	}
}

func drawGama(gama [gamaSize]int) {
	file, err := os.Create(DrawDir + "/" + strconv.Itoa(GamaCounter) + ".svg")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	canvas := svg.New(file)
	canvas.Start(640, 480)

	//White keys
	for i, n := range [7]int{1, 3, 5, 6, 8, 10, 12} {
		if noteInGama(gama, n) {
			canvas.Rect(i*SVGConf.WhiteWidth, 0, SVGConf.WhiteWidth, SVGConf.WhiteHeight, "style=\"fill:pink;stroke:black\"")
		} else {
			canvas.Rect(i*SVGConf.WhiteWidth, 0, SVGConf.WhiteWidth, SVGConf.WhiteHeight, "style=\"fill:white;stroke:black\"")
		}
	}

	//Black keys
	for _, n := range [5]int{2, 4, 7, 9, 11} {
		if noteInGama(gama, n) {
			canvas.Rect((n-1)*SVGConf.BlackWidth, 0, SVGConf.BlackWidth, SVGConf.BlackHeight, "style=\"fill:pink;stroke:black\"")
		} else {
			canvas.Rect((n-1)*SVGConf.BlackWidth, 0, SVGConf.BlackWidth, SVGConf.BlackHeight, "style=\"fill:black;stroke:black\"")
		}
	}
	canvas.End()
}

func noteInGama(gama [gamaSize]int, note int) bool {
	for _, n := range gama {
		if note == n {
			return true
		}
	}
	return false
}
