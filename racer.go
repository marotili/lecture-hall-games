package main

import (
	"errors"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/mixer"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"image"
	"image/color"
	"math"
	"time"
)

type Racer struct {
	cars []*Car

	obstaclemap *image.Gray
	heightmap   *image.Gray

	spriteCar        *Sprite
	spriteBackground *Sprite

	music *mixer.Music
}

func NewRacer() (*Racer, error) {

	r := &Racer{cars: make([]*Car, 0)}

	var err error
	if r.obstaclemap, err = LoadImageGray("data/velocity.png"); err != nil {
		return nil, err
	}
	if r.heightmap, err = LoadImageGray("data/velocity.png"); err != nil {
		return nil, err
	}

	if r.spriteCar, err = NewSprite("data/car.png", 16, 48); err != nil {
		return nil, err
	}

	if r.spriteBackground, err = NewSprite("data/background.png", 800, 600); err != nil {
		return nil, err
	}

	if r.music = mixer.LoadMUS("data/music.ogg"); r.music == nil {
		return nil, errors.New(sdl.GetError())
	}

	return r, nil
}

func (r *Racer) Update(t time.Duration) {
	for _, car := range r.cars {
		car.velocity = car.maxVelocity * float32(r.obstaclemap.At(int(car.position.x), int(car.position.y)).(color.Gray).Y) / 255
		car.position =
			car.position.Add(car.direction.MulScalar(car.velocity * float32(t.Seconds())))

		car.steer(car.owner.JoystickX*2-1, t)
	}
}

func (r *Racer) Render() {
	r.spriteBackground.Draw(400, 300, 0, 1)

	for _, car := range r.cars {
		// heightMod := 1/racer.heightGraymap.Modifier(car.position)
		car.Draw(1)
	}
}

func (r *Racer) Join(player *Player) {
	if len(r.cars) == 0 {
		mixer.ResumeMusic()
		r.music.PlayMusic(-1)
	}
	car := NewCar(player, r.spriteCar)
	car.position.x = 200
	car.position.y = 200
	r.cars = append(r.cars, car)
}

func (r *Racer) Leave(player *Player) {
	for i := range r.cars {
		if r.cars[i].owner == player {
			if i < len(r.cars) {
				r.cars = append(r.cars[:i], r.cars[i+1:]...)
			} else {
				r.cars = r.cars[:i-1]
			}
		}
	}
	if len(r.cars) == 0 {
		mixer.PauseMusic()
	}
}

type Car struct {
	position    Vector
	direction   Vector
	maxVelocity float32
	velocity    float32
	zLevel      int
	layer       int
	owner       *Player
	sprite      *Sprite
	steerValue  float32
}

func (car *Car) Draw(heightMod float32) {
	angle := math.Pi/2.0 + math.Atan2(float64(car.direction.y), float64(car.direction.x))
	car.sprite.Draw(float32(car.position.x), float32(car.position.y),
		float32(angle), heightMod)
}

func NewCar(owner *Player, sprite *Sprite) *Car {
	return &Car{
		position:    Vector{0, 0},
		direction:   Vector{0, 1},
		maxVelocity: 100,
		velocity:    10,
		zLevel:      0,
		layer:       0,
		owner:       owner,
		sprite:      sprite,
	}
}

var timeFactor float32 = 0.00000001

func (car *Car) steer(power float32, elapsedTime time.Duration) {
	max_angle := float32(math.Pi / 8.0)
	car.direction.x =
		car.direction.x*float32(math.Cos(float64(max_angle*power*float32(elapsedTime)*timeFactor))) -
			car.direction.y*float32(math.Sin(float64(max_angle*power*float32(elapsedTime)*timeFactor)))

	car.direction.y =
		car.direction.x*float32(math.Sin(float64(max_angle*power*float32(elapsedTime)*timeFactor))) +
			car.direction.y*float32(math.Cos(float64(max_angle*power*float32(elapsedTime)*timeFactor)))

	car.direction.Normalize()
}
