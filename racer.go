package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"log"
	"math"
	"time"
)

type Vector struct {
	x float32
	y float32
}

func (lhs Vector) Mul(rhs Vector) float32 {
	return lhs.x*rhs.x + lhs.y + rhs.y
}

func (lhs Vector) MulScalar(rhs float32) Vector {
	return Vector{lhs.x * rhs, lhs.y * rhs}
}

func (lhs Vector) Add(rhs Vector) Vector {
	return Vector{lhs.x + rhs.x, lhs.y + rhs.y}
}

type Player struct {
	name string
}

type Car struct {
	position    Vector
	direction   Vector
	maxVelocity float32
	velocity    float32
	zLevel      int
	layer       int
	owner       Player
	surface     *sdl.Surface
}

type Greymap struct {
	data *sdl.Surface
}

func (greymap Greymap) Modifier(pos Vector) float32 {
	value := getpixel(greymap.data, uint16(pos.x), uint16(pos.y))
	return float32(255) / float32(value)
}

type Racer struct {
	cars            []*Car
	boundingRects   []Vector // stub
	velocityGreymap Greymap
}

func NewRacer() *Racer {
	return &Racer{
		cars:            nil,
		boundingRects:   nil,
		velocityGreymap: Greymap{sdl.Load("artwork/velocity.png")},
	}
}

func NewCar(owner Player) *Car {
	return &Car{
		position:    Vector{0, 0},
		direction:   Vector{1, 0},
		maxVelocity: 10,
		velocity:    10,
		zLevel:      0,
		layer:       0,
		owner:       owner,
	}
}

func (car *Car) steer(power float32) {
	max_angle := float32(math.Pi / 4.0)
	car.direction.x = car.direction.x*float32(math.Cos(float64(max_angle*power))) -
		car.direction.y*float32(math.Sin(float64(max_angle*power)))

	car.direction.y = car.direction.x*float32(math.Sin(float64(max_angle*power))) +
		car.direction.y*float32(math.Cos(float64(max_angle*power)))
}

func (racer *Racer) Update(elapsedTime time.Duration) {
	for _, car := range racer.cars {
		car.velocity =
			car.maxVelocity * racer.velocityGreymap.Modifier(car.position)
		car.position = car.position.Add(car.direction.MulScalar(car.velocity))
	}
}

func (racer *Racer) Render(screen *sdl.Surface) {
	screen.FillRect(nil, 0x000000)
	for _, car := range racer.cars {
		log.Fatal("%f", car.velocity)
	}
	// background layer
	// map layers
	// cars
}
