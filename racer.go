package main

import (
    "math"
	"errors"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/mixer"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"image"
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
        car.Update(t)
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

type Wheel struct {
   torque   float32
   speed  float32
   inertia    float32
   radius   float32
   forwardAxis Vector
   sideAxis Vector
   position Vector
}

func NewWheel(position Vector, radius float32) *Wheel {
    wheel := &Wheel{
        torque:0,
        speed:0,
        radius:radius,
        inertia:radius*radius,//fake
        position:position,
       }
    wheel.SetSteeringAngle(0)
    return wheel
}

func (w *Wheel) SetSteeringAngle(newAngle float32) {
    forward := Vector{0,1}
    side := Vector{-1,0}

    w.forwardAxis = forward.Rotate(newAngle)
    w.sideAxis = side.Rotate(newAngle)
}

func (w *Wheel) AddTransmissionTorque(newValue float32) {
    w.torque += newValue
}

func (w *Wheel) CalculateForce(relativeGroundSpeed Vector, tDur time.Duration) Vector {
    t := float32(tDur) * timeFactor

    patchSpeed := w.forwardAxis.MulScalar(-w.speed*w.radius)
    velDiff := relativeGroundSpeed.Add(patchSpeed)
    sideVel, _ := velDiff.Project(w.sideAxis)

    forwardVel, forwardMag := velDiff.Project(w.forwardAxis)

    responseForce := sideVel.MulScalar(-2)
    responseForce = responseForce.Add(forwardVel.MulScalar(-1))

    w.torque += forwardMag * w.radius
    w.speed += w.torque / w.inertia * t
    w.torque = 0

    return responseForce
}

type Car struct {
	maxVelocity float32
	zLevel      int
	layer       int
	owner       *Player
	sprite      *Sprite
	steerValue  float32

    mass        float32
    inertia     float32

    force       Vector
    velocity    Vector
    position    Vector

    torque      float32
    angularVelocity   float32
    angle       float32 

    wheels [4]*Wheel
}

func (car *Car) AddForce(force Vector, relOffset Vector) {
    car.force = car.force.Add(force)
    car.torque += car.force.CrossProd(relOffset)
}

func (car *Car)RelativeToWorld(relative Vector) Vector {
    return relative.Rotate(car.angle)
}

func (car *Car)WorldToRelative(relative Vector) Vector {
    return relative.Rotate(-car.angle)
}

func (car *Car)PointVel(offset Vector) Vector {
    tangent := Vector{-offset.y, offset.x}
    return tangent.MulScalar(car.angularVelocity).Add(car.velocity)
}

func (car *Car) Update(time time.Duration) {
    t := float32(time)*timeFactor

    car.SetThrottle(car.owner.JoystickY*2-1, false)

    car.Steer(car.owner.JoystickX*2-1)
//    car.Steer(car.owner.JoystickX*2-1, time)

    for _, wheel := range(car.wheels) {
        worldWheelOffset := car.RelativeToWorld(wheel.position)
        worldWheelGroundVel := car.PointVel(worldWheelOffset)
        relGroundSpeed := car.WorldToRelative(worldWheelGroundVel)
        relResponseForce := wheel.CalculateForce(relGroundSpeed, time)
        worldResponseForce := car.RelativeToWorld(relResponseForce)

        car.AddForce(worldResponseForce, worldWheelOffset)
    }

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
	car.sprite.Draw(float32(car.position.x), float32(car.position.y),
		float32(car.angle), heightMod)
}

func NewCar(owner *Player, sprite *Sprite) *Car {
	return &Car{
		position:    Vector{0, 0},
		velocity:   Vector{0, 0},
		maxVelocity: 100,
		zLevel:      0,
		layer:       0,
		owner:       owner,
		sprite:      sprite,
        force:Vector{0,0},
        torque:0,
        angularVelocity:0,
        angle:0,
        mass:6,
        inertia:1/24.0*20*8*8*24*24,
        wheels:[4]*Wheel{
            NewWheel(Vector{-8,24}, 2),
            NewWheel(Vector{8,24}, 2),
            NewWheel(Vector{8,-24}, 2),
            NewWheel(Vector{-8,-24}, 2),
        },
	}
}

var timeFactor float32 = 0.00000001

func (car *Car) Steer(steering float32) {
    steeringLock := float32(math.Pi/2)
    car.wheels[0].SetSteeringAngle(-steering*steeringLock)
    car.wheels[1].SetSteeringAngle(-steering*steeringLock)
}

func (car *Car) SetThrottle(throttle float32, allWheel bool) {
    torque := float32(40)

    if allWheel {
        car.wheels[0].AddTransmissionTorque(throttle*torque)
        car.wheels[1].AddTransmissionTorque(throttle*torque)
    }

    car.wheels[2].AddTransmissionTorque(throttle*torque)
    car.wheels[3].AddTransmissionTorque(throttle*torque)
}

func (car *Car) SetBrakes(brakes float32) {
    brakeTorque := float32(4.0)
    for _, wheel := range(car.wheels) {
        wheelVel := wheel.speed
        wheel.AddTransmissionTorque(-wheelVel * brakeTorque * brakes)
    }
}
