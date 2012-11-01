package main

import (
	"encoding/binary"
	"fmt"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"github.com/banthar/gl"
	"io"
	"log"
	"math"
	"math/rand"
	"net"
	"runtime"
	"sync"
	"time"
)

type Game interface {
	Run()
}

const (
	ButtonA = 1
	ButtonB = 2
)

var mu sync.Mutex

func handleConnection(conn net.Conn) {
	car := NewCar(Player{"blub"}, sprite)
	defer func() {
		log.Println("client disconnected")
		mu.Lock()
		for i := range racer.cars {
			if racer.cars[i] == car {
				racer.cars = append(racer.cars[:i], racer.cars[i+1:]...)
			}
		}
		mu.Unlock()
		conn.Close()
	}()
	log.Println("client connected")
	car.position.x = 100
	car.position.y = 100

	mu.Lock()
	racer.cars = append(racer.cars, car)
	mu.Unlock()

	buf := make([]byte, 12)
	for {
		if _, err := io.ReadFull(conn, buf); err != nil {
			log.Println(err)
			return
		}
		joyX := math.Float32frombits(binary.BigEndian.Uint32(buf))
		joyY := math.Float32frombits(binary.BigEndian.Uint32(buf[4:]))
		buttons := binary.BigEndian.Uint32(buf[8:])
		fmt.Println(joyX, joyY, buttons)

		if joyX < 0 {
			joyX = 0
		} else if joyX > 1 {
			joyX = 1
		}
		mu.Lock()
		car.steerValue = joyX*2 - 1
		mu.Unlock()
	}
}

func drawTexture() {

}

var racer = NewRacer()
var sprite *Sprite

func main() {
	log.SetFlags(0)
	runtime.LockOSThread()

	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		log.Fatal(sdl.GetError())
	}

	var screen = sdl.SetVideoMode(800, 600, 32, sdl.OPENGL)
	if screen == nil {
		log.Fatal(sdl.GetError())
	}

	sdl.WM_SetCaption("Lecture Hall Games 0.1", "")
	sdl.EnableUNICODE(1)

	if gl.Init() != 0 {
		panic("gl error")

	}

	gl.Viewport(0, 0, int(screen.W), int(screen.H))
	gl.ClearColor(1, 1, 1, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(screen.W), float64(screen.H), 0, -1.0, 1.0)

	rand.Seed(time.Now().UnixNano())

	sdl.EnableUNICODE(1)

	go func() {
		ln, err := net.Listen("tcp", ":8001")
		if err != nil {
			log.Fatal("Server failed")
		}

		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Fatal("Client failed")
				continue
			}

			go handleConnection(conn)
		}
	}()

	running := true
	last := time.Now()

	width := 16
	sprite = NewSprite("artwork/auto.png", width, width*3)
	background := NewSprite("artwork/background.png", 800, 600)

	for running {
		// move objects
		current := time.Now()
		t := current.Sub(last)
		last = current

		// process events
		select {
		case event := <-sdl.Events:
			switch e := event.(type) {
			case sdl.QuitEvent:
				running = false
			case sdl.ResizeEvent:
				screen = sdl.SetVideoMode(int(e.W), int(e.H), 32, sdl.RESIZABLE)
			case sdl.KeyboardEvent:
			}
		default:
		}

		mu.Lock()
		for i := range racer.cars {
			racer.cars[i].steer(racer.cars[i].steerValue, t)
		}
		mu.Unlock()

		gl.ClearColor(1, 1, 1, 0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		racer.Update(t)
		racer.Render()

		background.Draw(400, 300, 0, 1)

		sdl.GL_SwapBuffers()
	}

	sdl.Quit()
}
