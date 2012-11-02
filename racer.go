package main

import (
	"errors"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/mixer"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/ttf"
	"image"
	"image/color"
	"math"
	"time"
"fmt"
)

type Racer struct {
	cars []*Car

	obstaclemap *image.Gray
	heightmap   *image.Gray

	spriteCarFG *Sprite
	spriteCarBG *Sprite

	spriteForeground *Sprite
	spriteBackground *Sprite
	spriteWaiting    *Sprite

	running bool

	music *mixer.Music
	font  *ttf.Font
}

func NewRacer() (*Racer, error) {
	r := &Racer{cars: make([]*Car, 0)}

	var err error
	if r.obstaclemap, err = LoadImageGray("data/levels/demolevel3/velocity.png"); err != nil {
		return nil, err
	}
	if r.heightmap, err = LoadImageGray("data/levels/demolevel3/z.png"); err != nil {
		return nil, err
	}

	carSize := 0.04 * float32(screenWidth)
	if r.spriteCarFG, err = NewSprite("data/cars/car1/fg.png", carSize, carSize); err != nil {
		return nil, err
	}
	if r.spriteCarBG, err = NewSprite("data/cars/car1/bg.png", carSize, carSize); err != nil {
		return nil, err
	}

	if r.spriteForeground, err = NewSprite("data/levels/demolevel3/foreground.png", screenWidth, screenHeight); err != nil {
		return nil, err
	}
	if r.spriteBackground, err = NewSprite("data/levels/demolevel3/background.png", screenWidth, screenHeight); err != nil {
		return nil, err
	}

	if r.music = mixer.LoadMUS("data/music.ogg"); r.music == nil {
		return nil, errors.New(sdl.GetError())
	}

	if r.font = ttf.OpenFont("data/font.otf", 72); r.font == nil {
		return nil, errors.New(sdl.GetError())
	}

	textWaiting := ttf.RenderUTF8_Blended(r.font, "X", sdl.Color{255, 0, 0, 0})
	r.spriteWaiting = NewSpriteFromSurface(textWaiting)

	return r, nil
}

func (r *Racer) Update(t time.Duration) {
	if !r.running {
		return
	}
	for _, car := range r.cars {
		car.velocity = 100 * valueAt(r.obstaclemap, car.position.x, car.position.y)
		car.position = car.position.Add(car.direction.MulScalar(car.velocity * float32(t.Seconds())))

		car.steer(car.owner.JoystickX, t)
	}
}

func (r *Racer) Render(screen *sdl.Surface) {
	r.spriteBackground.Draw(screenWidth/2, screenHeight/2, 0, 1, false)

	for _, car := range r.cars {
		// heightMod := 1/racer.heightGraymap.Modifier(car.position)

		size := float32(1.0)
		size *= (1 - 0.2*valueAt(r.heightmap, car.position.x, car.position.y))
		if car.owner.ButtonA {
			size *= 1.5
		} else if car.owner.ButtonB {
			size *= 0.75
		}
		car.Draw(size)
	}

	r.spriteForeground.Draw(screenWidth/2, screenHeight/2, 0, 1, true)

}

func (r *Racer) Join(player *Player) {
	if len(r.cars) == 0 {
		mixer.ResumeMusic()
		r.music.PlayMusic(-1)
	}
	car := NewCar(player, r.spriteCarFG, r.spriteCarBG)
	car.position.x = 200
	car.position.y = 200
	r.cars = append(r.cars, car)
}

func (r *Racer) Leave(player *Player) {
	for i := range r.cars {
		if r.cars[i].owner == player {
			r.cars[i] = r.cars[len(r.cars)-1]
			r.cars = r.cars[:len(r.cars)-1]
			break
		}
	}
	if len(r.cars) == 0 {
		r.running = false
		mixer.PauseMusic()
	}
}


func valueAt(img *image.Gray, x, y float32) float32 {
	dx, dy := x/float32(screenWidth), y/float32(screenHeight)
	b := img.Bounds().Max
	px, py := int(dx*float32(b.X)), int(dy*float32(b.Y))
	v := float32(img.At(px, py).(color.Gray).Y) / 255
	return v

func (r *Racer) KeyPressed(input sdl.Keysym) {

	fmt.Printf("%d pressed\n",input)
	if input.Sym == sdl.K_SPACE {
		r.running = true
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

	spriteBG   *Sprite
	spriteFG   *Sprite
	steerValue float32
}

func (car *Car) Draw(heightMod float32) {
	angle := math.Pi/2.0 + math.Atan2(float64(car.direction.y), float64(car.direction.x))
	car.spriteBG.Draw(float32(car.position.x), float32(car.position.y),
		float32(angle), heightMod, true)
	car.spriteFG.Draw(float32(car.position.x), float32(car.position.y),
		float32(angle), heightMod, true)
}

func NewCar(owner *Player, spriteFG, spriteBG *Sprite) *Car {
	return &Car{
		position:    Vector{0, 0},
		direction:   Vector{0, 1},
		maxVelocity: 100,
		velocity:    10,
		zLevel:      0,
		layer:       0,
		owner:       owner,
		spriteFG:    spriteFG,
		spriteBG:    spriteBG,
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
