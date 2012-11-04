// Copyright (c) 2012 by Lecture Hall Games Authors.
// All source files are distributed under the Simplified BSD License.

package main

import (
	"errors"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/mixer"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/ttf"
	"image"
	"image/color"
	"path/filepath"
	"time"
)

type Racer struct {
	cars []*Car

	carSize float32

	obstaclemap *image.Gray
	heightmap   *image.Gray

	spriteCarFG *Sprite
	spriteCarBG *Sprite

	spriteForeground *Sprite
	spriteBackground *Sprite
	spriteWaiting    *Sprite

	running   bool
	showNames bool

	music *mixer.Music
	font  *ttf.Font
}

func NewRacer(levelDir string) (*Racer, error) {
	r := &Racer{cars: make([]*Car, 0)}

	var err error
	if r.obstaclemap, err = LoadImageGray(filepath.Join(levelDir, "velocity.png")); err != nil {
		return nil, err
	}
	if r.heightmap, err = LoadImageGray(filepath.Join(levelDir, "z.png")); err != nil {
		return nil, err
	}

	carSize := 0.02 * float32(screenWidth)
	if r.spriteCarFG, err = NewSprite("data/cars/car1/fg.png", carSize, carSize); err != nil {
		return nil, err
	}
	if r.spriteCarBG, err = NewSprite("data/cars/car1/bg.png", carSize, carSize); err != nil {
		return nil, err
	}

	r.carSize = carSize

	if r.spriteForeground, err = NewSprite(filepath.Join(levelDir, "foreground.png"), screenWidth, screenHeight); err != nil {
		return nil, err
	}
	if r.spriteBackground, err = NewSprite(filepath.Join(levelDir, "background.png"), screenWidth, screenHeight); err != nil {
		return nil, err
	}

	if r.music = mixer.LoadMUS("data/music.ogg"); r.music == nil {
		return nil, errors.New(sdl.GetError())
	}

	if r.font = ttf.OpenFont("data/font.otf", 32); r.font == nil {
		return nil, errors.New(sdl.GetError())
	}

	textWaiting := ttf.RenderUTF8_Blended(r.font, "Waiting for other players. Press space to start....", sdl.Color{0, 0, 255, 0})
	r.spriteWaiting = NewSpriteFromSurface(textWaiting)

	return r, nil
}

func (r *Racer) Update(t time.Duration) {
	if !r.running {
		return
	}

	r.HandleCollisions()

	for _, car := range r.cars {
		car.Update(t, r)
	}
}

func (r *Racer) Render(screen *sdl.Surface) {
	r.spriteBackground.Draw(screenWidth/2, screenHeight/2, 0, 1, false)

	for _, car := range r.cars {
		size := (1 - 0.3*valueAt(r.heightmap, car.position.x, car.position.y))
		car.Draw(size)
	}

	if r.showNames {
		x := float32(screenWidth) * 0.02
		y := float32(screenHeight) * 0.02
		for _, car := range r.cars {
			car.spriteNick.Draw(x+.25*car.spriteNick.width, y+.25*car.spriteNick.height, 0, .5, true)
			y += .5 * 1.2 * car.spriteNick.height
		}
	}

	r.spriteForeground.Draw(screenWidth/2, screenHeight/2, 0, 1, true)

	if r.running != true {
		r.spriteWaiting.Draw(screenWidth/2, screenHeight/5, 0, 1, true)
	}
}

func (r *Racer) Join(player *Player, x, y float32) {
	if len(r.cars) == 0 {
		mixer.ResumeMusic()
		r.music.PlayMusic(-1)
	}
	car := NewCar(player, r.spriteCarFG, r.spriteCarBG, r.carSize, r.font)
	car.position.x = x
	car.position.y = y
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

type Wheel struct {
	torque      float32
	speed       float32
	inertia     float32
	radius      float32
	forwardAxis Vector
	sideAxis    Vector
	position    Vector
}

func NewWheel(position Vector, radius float32) *Wheel {
	wheel := &Wheel{
		torque:   0,
		speed:    0,
		radius:   radius,
		inertia:  radius * radius, //fake
		position: position,
	}
	wheel.SetSteeringAngle(0)
	return wheel
}

func (w *Wheel) SetSteeringAngle(newAngle float32) {
	forward := Vector{0, 1}
	side := Vector{-1, 0}

	w.forwardAxis = forward.Rotate(newAngle)
	w.sideAxis = side.Rotate(newAngle)
}

func (w *Wheel) AddTransmissionTorque(newValue float32) {
	w.torque += newValue
}

func (w *Wheel) CalculateForce(relativeGroundSpeed Vector, tDur time.Duration) Vector {
	t := float32(tDur) * timeFactor

	patchSpeed := w.forwardAxis.MulScalar(-w.speed * w.radius)
	velDiff := relativeGroundSpeed.Add(patchSpeed)
	sideVel, _ := velDiff.Project(w.sideAxis)

	forwardVel, forwardMag := velDiff.Project(w.forwardAxis)

	responseForce := sideVel.MulScalar(-1)
	responseForce = responseForce.Add(forwardVel.MulScalar(-1))

	w.torque += forwardMag * w.radius
	w.speed += w.torque / w.inertia * t
	w.torque = 0

	return responseForce
}

func valueAt(img *image.Gray, x, y float32) float32 {
	dx, dy := x/float32(screenWidth), y/float32(screenHeight)
	b := img.Bounds().Max
	px, py := int(dx*float32(b.X)), int(dy*float32(b.Y))
	v := float32(img.At(px, py).(color.Gray).Y) / 255
	return v
}

func (r *Racer) KeyPressed(input sdl.Keysym) {
	if input.Sym == sdl.K_SPACE {
		r.running = true
	}
	if input.Sym == sdl.K_TAB {
		r.showNames = !r.showNames
	}
}

type Car struct {
	maxVelocity float32
	zLevel      int
	layer       int
	owner       *Player
	steerValue  float32

	mass    float32
	inertia float32

	force    Vector
	velocity Vector
	position Vector

	torque          float32
	angularVelocity float32
	angle           float32

	wheels [2]*Wheel

	spriteBG   *Sprite
	spriteFG   *Sprite
	size       float32
	width      float32
	height     float32
	spriteNick *Sprite
}

func (car *Car) AddForce(force Vector, relOffset Vector) {
	car.force = car.force.Add(force)
	car.torque += relOffset.CrossProd(force)
}

func (car *Car) RelativeToWorld(relative Vector) Vector {
	return relative.Rotate(car.angle)
}

func (car *Car) WorldToRelative(relative Vector) Vector {
	return relative.Rotate(-car.angle)
}

func (car *Car) PointVel(offset Vector) Vector {
	tangent := Vector{-offset.y, offset.x}
	return tangent.MulScalar(car.angularVelocity).Add(car.velocity)
}

func (car *Car) Update(time time.Duration, r *Racer) {
	t := float32(time) * timeFactor
	if car.owner.ButtonA {
		car.SetThrottle(1, false)
	} else {
		car.SetThrottle(0, false)
	}
	if car.owner.ButtonB {
		car.SetThrottle(-0.8, false)
	} else {
		car.SetBrakes(0)
	}

	steer := car.owner.JoystickX
	if !car.owner.ButtonA {
		steer *= 1.1
	}
	car.Steer(steer)
	for _, wheel := range car.wheels {
		worldWheelOffset := car.RelativeToWorld(wheel.position)
		worldWheelGroundVel := car.PointVel(worldWheelOffset)
		relGroundSpeed := car.WorldToRelative(worldWheelGroundVel)
		relResponseForce := wheel.CalculateForce(relGroundSpeed, time)
		worldResponseForce := car.RelativeToWorld(relResponseForce)

		car.AddForce(worldResponseForce, worldWheelOffset)
	}

	const cRolling = 12.8 * .2
	const cDrag = 0.4257 * .3

	fDrag := car.velocity.MulScalar(-cDrag * car.velocity.Length())
	fRolling := car.velocity.MulScalar(-cRolling)

	terrain := (1 - valueAt(r.obstaclemap, car.position.x, car.position.y)) * 6
	fRolling = fRolling.MulScalar(terrain)

	car.AddForce(fDrag, Vector{0, 0})
	car.AddForce(fRolling, Vector{0, 0})

	acceleration := car.force.DivScalar(car.mass)
	car.velocity = car.velocity.Add(acceleration.MulScalar(t))
	car.position = car.position.Add(car.velocity.MulScalar(t))

	angAcc := car.torque / car.inertia
	car.angularVelocity += angAcc * t
	car.angle += car.angularVelocity * t

	car.force = Vector{0, 0}
	car.torque = 0
}

func (car *Car) Draw(heightMod float32) {
	car.spriteBG.Draw(float32(car.position.x), float32(car.position.y),
		float32(car.angle), heightMod, true)
	car.spriteFG.Draw(float32(car.position.x), float32(car.position.y),
		float32(car.angle), heightMod, true)
}

func NewCar(owner *Player, spriteFG, spriteBG *Sprite, carSize float32, font *ttf.Font) *Car {
	textNick := ttf.RenderUTF8_Blended(font, owner.Nick, sdl.Color{255, 255, 255, 0})
	spriteNick := NewSpriteFromSurface(textNick)
	return &Car{
		position:        Vector{0, 0},
		velocity:        Vector{0, 0},
		maxVelocity:     100,
		zLevel:          0,
		layer:           0,
		owner:           owner,
		force:           Vector{0, 0},
		torque:          0,
		angularVelocity: 0,
		angle:           0,
		mass:            5,
		inertia:         200,
		wheels: [2]*Wheel{
			NewWheel(Vector{0, carSize / 2.0}, 4),
			NewWheel(Vector{0, -carSize / 2.0}, 4),
		},
		spriteFG:   spriteFG,
		spriteBG:   spriteBG,
		size:       carSize,
		width:      carSize * 18 / 32.0,
		height:     carSize * 1,
		spriteNick: spriteNick,
	}
}

var timeFactor float32 = 0.00000001

func (car *Car) Steer(steering float32) {
	steeringLock := float32(0.5)
	car.wheels[1].SetSteeringAngle(-steering * steeringLock)
}

func (car *Car) SetThrottle(throttle float32, allWheel bool) {
	torque := float32(100)

	if allWheel {
		car.wheels[1].AddTransmissionTorque(throttle * torque)
	}

	car.wheels[0].AddTransmissionTorque(throttle * torque)
}

func (car *Car) SetBrakes(brakes float32) {
	brakeTorque := float32(5)
	for _, wheel := range car.wheels {
		wheelVel := wheel.speed
		wheel.AddTransmissionTorque(-wheelVel * brakeTorque * brakes)
	}
}
