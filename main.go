package main

import (
	"encoding/binary"
	"fmt"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"github.com/banthar/gl"
	"log"
	"math"
	"math/rand"
	"net"
	"runtime"
	"time"
)

type Game interface {
	Run()
}

const (
	ButtonA = 1
	ButtonB = 2
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Println("client connected")

	buf := make([]byte, 12)
	for {
		if _, err := conn.Read(buf); err != nil {
			log.Println(err)
			return
		}
		joyX := math.Float32frombits(binary.BigEndian.Uint32(buf))
		if joyX > 1.0 {
			joyX = 1.0
		} else if joyX < -1.0 {
			joyX = -1.0
		}
		joyY := math.Float32frombits(binary.BigEndian.Uint32(buf[4:]))
		if joyY > 1.0 {
			joyY = 1.0
		} else if joyY < -1.0 {
			joyY = -1.0
		}
		buttons := binary.BigEndian.Uint32(buf[8:])
		fmt.Println(joyX, joyY, buttons)
	}
}

func drawTexture() {

}

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
	sprite := NewSprite("artwork/auto.png", width, width*3)

	car := NewCar(Player{"Marco"}, sprite)
	car.position.x = 100
	car.position.y = 100

	car2 := NewCar(Player{"Christoph"}, sprite)
	car2.position.x = 200
	car2.position.y = 100

	racer := NewRacer()
	racer.cars = append(racer.cars, car)
	racer.cars = append(racer.cars, car2)

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
				if e.Keysym.Sym == sdl.K_LEFT {
					car.steer(-1, t)
				} else if e.Keysym.Sym == sdl.K_RIGHT {
					car.steer(+1, t)
				}
			}
		default:
		}

		//		fmt.Println(t)

		gl.ClearColor(1, 1, 1, 0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		racer.Update(t)
		racer.Render()

		sdl.GL_SwapBuffers()
	}

	sdl.Quit()
}
