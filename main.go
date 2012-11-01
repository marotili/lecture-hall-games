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

func loadTexture() {
	gl.Enable(gl.TEXTURE_2D)
	gl.Disable(gl.TEXTURE_2D)
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

	gl.MatrixMode(gl.PROJECTION)
	gl.Viewport(0, 0, int(screen.W), int(screen.H))
	gl.LoadIdentity()
	gl.Ortho(0, float64(screen.W), float64(screen.H), 0, -1.0, 1.0)
	gl.ClearColor(1, 1, 1, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

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

	for running {

		// process events
		select {
		case event := <-sdl.Events:
			switch e := event.(type) {
			case sdl.QuitEvent:
				running = false
			case sdl.ResizeEvent:
				screen = sdl.SetVideoMode(int(e.W), int(e.H), 32, sdl.RESIZABLE)
			}
		default:
		}
		// move objects
		current := time.Now()
		t := current.Sub(last)
		last = current
		fmt.Println(t)

		sdl.GL_SwapBuffers()
	}

	sdl.Quit()
}
