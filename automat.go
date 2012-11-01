package main

import (
//    "fmt"
    "github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
    "log"
    "net"
//    "math"
//    "os"
//    "strings"
    "time"
    "math/rand"
)

type Game interface {
    Run()
}

func handleConnection(conn Conn) {
    log.Info("Client joined, %s, %s", conn.RemoteAddr().String())
}

func main () {
    log.SetFlags(0)

    if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
        log.Fatal(sdl.GetError())
    }

    var screen = sdl.SetVideoMode(800, 600, 32, sdl.RESIZABLE)

    if screen == nil {
        log.Fatal(sdl.GetError())
    }

    rand.Seed(time.Now().UnixNano())

    sdl.EnableUNICODE(1)

    ticker := time.NewTicker(time.Second / 50)

    go func () {
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
    }

    go func () {
        for {
            select {
            case <- ticker.C:
                screen.FillRect(nil, 0x302019)
            loop:
                for {
                    select {
                    default:
                        break loop
                    }
                }
                screen.Flip()
            }
        }
    }()

    running := true
    for running {
        select {
        case _event := <-sdl.Events:
            switch e := _event.(type) {
            case sdl.QuitEvent:
                running = false
            case sdl.ResizeEvent:
                screen = sdl.SetVideoMode(int(e.W), int(e.H), 32,
                    sdl.RESIZABLE)

                if screen == nil {
                    log.Fatal(sdl.GetError())
                }
            case sdl.KeyboardEvent:
            }
        }
    }

    sdl.Quit()
}


