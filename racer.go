package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"github.com/banthar/gl"
	"image"
	"image/png"
	"math"
	"os"
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
	sprite      *Sprite
	steerValue  float32
}

func (car *Car) Draw() {
	angle := math.Pi/2.0 + math.Atan2(float64(car.direction.y), float64(car.direction.x))
	car.sprite.Draw(float32(car.position.x), float32(car.position.y),
		float32(angle))
}

type Sprite struct {
	filename string
	width    int
	height   int
	texture  gl.Texture
}

func NewSprite(filename string, width, height int) *Sprite {
	fi, _ := os.Open(filename)
	img, _ := png.Decode(fi)

	gl.Enable(gl.TEXTURE_2D)
	texture := gl.GenTexture()
	texture.Bind(gl.TEXTURE_2D)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, 4, img.Bounds().Max.X, img.Bounds().Max.Y, 0, gl.RGBA, gl.UNSIGNED_BYTE,
		(img.(*image.RGBA)).Pix)
	gl.Disable(gl.TEXTURE_2D)

	return &Sprite{
		filename: filename,
		width:    width,
		height:   height,
		texture:  texture,
	}
}

func (sprite *Sprite) Draw(x, y, angle float32) {
	gl.Enable(gl.TEXTURE_2D)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Translatef(x, y, 0)
	gl.Rotatef(angle*360/(2*math.Pi), 0, 0, 1)
	sprite.texture.Bind(gl.TEXTURE_2D)
	gl.Begin(gl.QUADS)
	gl.Color3f(1, 1, 1)
	gl.TexCoord2d(0, 0)
	gl.Vertex3f(-float32(sprite.width/2), -float32(sprite.height/2), 0)
	gl.TexCoord2d(1, 0)
	gl.Vertex2f(float32(sprite.width/2), -float32(sprite.height/2))
	gl.TexCoord2d(1, 1)
	gl.Vertex3f(float32(sprite.width/2), float32(sprite.height/2), 0)
	gl.TexCoord2d(0, 1)
	gl.Vertex3f(-float32(sprite.width/2), float32(sprite.height/2), 0)
	gl.End()
	gl.Disable(gl.TEXTURE_2D)
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

func NewCar(owner Player, sprite *Sprite) *Car {
	return &Car{
		position:    Vector{0, 0},
		direction:   Vector{0, 1},
		maxVelocity: 10,
		velocity:    10,
		zLevel:      0,
		layer:       0,
		owner:       owner,
		sprite:      sprite,
	}
}

func (car *Car) steer(power float32, elapsedTime time.Duration) {
	max_angle := float32(math.Pi / 4.0)
	car.direction.x =
		car.direction.x*float32(math.Cos(float64(max_angle*power*float32(elapsedTime)*0.00000001))) -
			car.direction.y*float32(math.Sin(float64(max_angle*power*float32(elapsedTime)*0.00000001)))

	car.direction.y =
		car.direction.x*float32(math.Sin(float64(max_angle*power*float32(elapsedTime)*0.00000001))) +
			car.direction.y*float32(math.Cos(float64(max_angle*power*float32(elapsedTime)*0.00000001)))
}

func (racer *Racer) Update(elapsedTime time.Duration) {
	for _, car := range racer.cars {
		car.velocity =
			car.maxVelocity * racer.velocityGreymap.Modifier(car.position)
		car.position =
			car.position.Add(car.direction.MulScalar(car.velocity * (float32(elapsedTime)) * 0.00000001))
	}
}

func (racer *Racer) Render() {
	for _, car := range racer.cars {
		car.Draw()
	}
	// background layer
	// map layers
	// cars
}
