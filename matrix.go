package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/nsf/termbox-go"
)

type Point struct {
	Head bool
	Char byte
	Age  int
}

var clear map[string]func()

var cols [][]Point
var width int
var height int

func main() {
	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}

	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)
	width, height = termbox.Size()

	cols = make2d(width, height)

	data := make([]byte, width*height)
	rand.Read(data)

	for i := range cols {
		for j := range cols[i] {
			cols[i][j] = Point{false, (data[i*height+j] % 94) + 33, 0}
			//cols[i][j] = byte(48 + (i*height+j)%10)
		}
	}

	data = make([]byte, 32)
	rand.Read(data)

	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()

	/*for i := 0; i < len(data)/2; i += 2 {
		a := int(data[i]) % width
		b := int(data[i+1]) % height
	}*/
	create()
loop:
	for {
		select {
		case ev := <-event_queue:
			if ev.Type == termbox.EventKey && (ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC) {
				break loop
			}
		default:
			//clearScreen()
			print(cols)
			time.Sleep(time.Millisecond * 50)
			step()
			create()
		}
	}
}

func make2d(width int, height int) [][]Point {
	arr := make([][]Point, width)
	for i := range arr {
		arr[i] = make([]Point, height)
	}
	return arr
}

func create() {
	chance := 7
	q := rand.Int() % chance
	for i := 0; i < q; i++ {
		p := rand.Int() % width
		l := rand.Int()%24 + 9
		cols[p][len(cols[p])-1].Age = l
	}
}

func step() {
	newcols := make([][]Point, width)
	for i := range newcols {
		newcols[i] = append([]Point(nil), cols[i]...)
	}
	for i := len(cols) - 1; i >= 0; i-- {
		for j := 0; j < len(cols[i]); j++ {
			if j != len(cols[i])-1 {
				if cols[i][j].Age == 0 && cols[i][j+1].Age > 0 {
					newcols[i][j].Age = cols[i][j+1].Age
					newcols[i][j].Head = true
					newcols[i][j+1].Head = false
				}
			}
			if cols[i][j].Age > 0 {
				newcols[i][j].Age--
			}
		}
	}
	cols = newcols
}

func print(data [][]Point) {
	for i := range data[0] {
		for j := range data {
			toshow := ' '

			a := j
			b := height - i - 1

			point := data[a][b]

			if data[a][b].Age > 0 {
				toshow = []rune(string(point.Char))[0]
			}
			if point.Head {
				termbox.SetCell(a, height-b-1, toshow, termbox.ColorWhite, termbox.ColorBlack)
				//fmt.Printf("\033[0;37m%s\033[0;37m", string(toshow))
			} else {
				termbox.SetCell(a, height-b-1, toshow, termbox.ColorGreen, termbox.ColorBlack)
				//fmt.Printf("\033[0;32m%s\033[0;32m", string(toshow))
			}
		}
		if i < len(data[0])-1 {
			//fmt.Printf("\n")
		}
	}
	termbox.Flush()
}

func initClear() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["darwin"] = func() {
		cmd := exec.Command("clear") //Mac example, its maybe tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func clearScreen() {
	fmt.Printf("%s\n", runtime.GOOS)
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
