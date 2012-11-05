// Copyright (c) 2012 by Lecture Hall Games Authors.
// All source files are distributed under the Simplified BSD License.

package main

import (
	"encoding/binary"
	"fmt"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/mixer"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/ttf"
	"github.com/banthar/gl"
	"go/build"
	"io"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"runtime"
	"sync"
	"time"
)

const basePkg = "github.com/fruhwirth-marco/lecture-hall-games"
const numberLevels = 2
const numberCars = 1

type Player struct {
	Conn      net.Conn
	Nick      string
	ButtonA   bool
	ButtonB   bool
	JoystickX float32
	JoystickY float32
}

const (
	screenWidth  = 1024
	screenHeight = 768
)

type Game interface {
	Update(t time.Duration)
	Render(screen *sdl.Surface)
	Join(player *Player, x, y float32)
	Leave(player *Player)
	KeyPressed(input sdl.Keysym)
}

func (p *Player) Vibrate() {
	binary.Write(p.Conn, binary.BigEndian, uint32(42))
}

func handleConnection(conn net.Conn) {
	player := &Player{Conn: conn}
	defer func() {
		log.Printf("Player %q left (%s)\n", player.Nick, conn.RemoteAddr())
		mu.Lock()
		game.Leave(player)
		mu.Unlock()
	}()

	var nickLength uint32
	binary.Read(conn, binary.BigEndian, &nickLength)
	nickBytes := make([]byte, nickLength)
	if _, err := io.ReadFull(conn, nickBytes); err != nil {
		log.Println(err)
		return
	}
	player.Nick = string(nickBytes)

	mu.Lock()
	game.Join(player, 200, 200)
	mu.Unlock()

	log.Printf("Player %q joined (%s)\n", player.Nick, conn.RemoteAddr())

	buf := make([]byte, 12)
	for {
		if _, err := io.ReadFull(conn, buf); err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			return
		}
		mu.Lock()
		player.JoystickX = math.Float32frombits(binary.BigEndian.Uint32(buf))
		player.JoystickY = math.Float32frombits(binary.BigEndian.Uint32(buf[4:]))
		buttons := binary.BigEndian.Uint32(buf[8:])
		player.ButtonA = buttons&1 != 0
		player.ButtonB = buttons&2 != 0
		if player.JoystickX < -1 {
			player.JoystickX = -1
		} else if player.JoystickX > 1 {
			player.JoystickX = 1
		}
		if player.JoystickY < -1 {
			player.JoystickY = 1
		} else if player.JoystickY > 1 {
			player.JoystickY = 1
		}
		mu.Unlock()
	}
}

var (
	game Game
	mu   sync.Mutex
)

func main() {
	runtime.LockOSThread()

	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		log.Fatal(sdl.GetError())
	}
	var screen = sdl.SetVideoMode(screenWidth, screenHeight, 32, sdl.OPENGL|sdl.HWSURFACE|sdl.GL_DOUBLEBUFFER|sdl.FULLSCREEN)
	if screen == nil {
		log.Fatal(sdl.GetError())
	}
	sdl.WM_SetCaption("Lecture Hall Games", "")
	sdl.EnableUNICODE(1)
	if gl.Init() != 0 {
		log.Fatal("could not initialize OpenGL")
	}
	gl.Viewport(0, 0, int(screen.W), int(screen.H))
	gl.ClearColor(1, 1, 1, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(screen.W), float64(screen.H), 0, -1.0, 1.0)
	gl.Disable(gl.LIGHTING)
	gl.Disable(gl.DEPTH_TEST)
	gl.TexEnvi(gl.TEXTURE_ENV, gl.TEXTURE_ENV_MODE, gl.MODULATE)

	if mixer.OpenAudio(mixer.DEFAULT_FREQUENCY, mixer.DEFAULT_FORMAT,
		mixer.DEFAULT_CHANNELS, 4096) != 0 {
		log.Fatal(sdl.GetError())
	}

	if ttf.Init() != 0 {
		log.Fatal(sdl.GetError())
	}

	if p, err := build.Default.Import(basePkg, "", build.FindOnly); err == nil {
		os.Chdir(p.Dir)
	}

	var err error

	rand.Seed(time.Now().UnixNano())
	levelDir := fmt.Sprintf("data/levels/demolevel%d", 3+rand.Intn(numberLevels))
	//carsDir := fmt.Sprintf(" data/cars/car%d/", 1+rand.Intn(numberCars))
	if game, err = NewRacer(levelDir); err != nil {
		log.Fatal(err)
	}

	go func() {
		listen, err := net.Listen("tcp", ":8001")
		if err != nil {
			log.Fatal(err)
		}
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			go handleConnection(conn)
		}
	}()

	running := true
	last := time.Now()
	for running {
		select {
		case event := <-sdl.Events:
			switch e := event.(type) {
			case sdl.QuitEvent:
				running = false
			case sdl.ResizeEvent:
				screen = sdl.SetVideoMode(int(e.W), int(e.H), 32, sdl.RESIZABLE)
			case sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN {
					if(e.Keysym.Sym == sdl.K_ESCAPE) {
						running = false;
					} else {
						game.KeyPressed(e.Keysym)
					}
				}
			}
		default:
		}

		current := time.Now()
		t := current.Sub(last)
		last = current

		mu.Lock()
		game.Update(t)
		game.Render(screen)
		mu.Unlock()

		sdl.GL_SwapBuffers()
	}

	sdl.Quit()
}
