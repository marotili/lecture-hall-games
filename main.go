package main

import (
	"encoding/binary"
	"fmt"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
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

func main() {
	runtime.LockOSThread()
	log.SetFlags(0)

	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		log.Fatal(sdl.GetError())
	}

	sdl.WM_SetCaption("Lecture Hall Games 0.1", "")
	sdl.EnableUNICODE(1)

	screen := sdl.SetVideoMode(800, 600, 32, sdl.RESIZABLE)
	if screen == nil {
		log.Fatal(sdl.GetError())
	}

	rand.Seed(time.Now().UnixNano())

	go func() {
		ln, err := net.Listen("tcp", ":8080")
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

		// render screen
		screen.FillRect(nil, 0x302019)
	}

	sdl.Quit()
}
