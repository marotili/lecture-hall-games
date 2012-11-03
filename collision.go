package main

import (
	"math"
)

func lines(car *Car, newAngle float32) [4][2]float32 {
    relx1 := car.width/2.0
    rely1 := float32(0)

    relx3 := float32(0)
    rely3 := car.height/2.0

    relx2 := -relx1
    rely2 := -rely1

    relx4 := -relx3
    rely4 := -rely3

    vectors := [4]*Vector{
        {relx1, rely1},
        {relx2, rely2},
        {relx3, rely3},
        {relx4, rely4},
    }

    for _, vector := range(vectors) {
        *vector = ((*vector).Rotate(car.angle))
        *vector = ((*vector).Add(car.position))
    }

    angle := car.angle + newAngle

    k1 := float32(math.Tan(math.Pi/2 + float64(angle)))
    d1 := -k1*vectors[0].x + vectors[0].y

	k2 := float32(math.Tan(math.Pi/2 + float64(angle)))
    d2 := -k2*vectors[1].x + vectors[1].y

	k3 := float32(math.Tan(float64(angle)))
    d3 := -k3*vectors[2].x + vectors[2].y

	k4 := float32(math.Tan(float64(angle)))
    d4 := -k4*vectors[3].x + vectors[3].y

    return [4][2]float32{
        {k1, d1}, 
        {k2, d2}, 
        {k3, d3}, 
        {k4, d4},
    }
}

func col(car *Car, car2 *Car) {

    newAngle :=float32(0)
    car1Angle := car.angle
    car2Angle := car2.angle

    for car1Angle > float32(math.Pi/2.0) {
        car1Angle -= float32(math.Pi/2.0)
    }

    for car2Angle > float32(math.Pi/2.0) {
        car2Angle -= float32(math.Pi/2.0)
    }

    if car1Angle + car2Angle < math.Pi/16.0 {
        newAngle = float32(45)
    } else if car1Angle < math.Pi/32.0 || car2Angle < math.Pi/32.0 {
        newAngle = (car1Angle + car2Angle)/2.0
    }

    c1Lines := lines(car, newAngle)
    c2Lines := lines(car2, newAngle)

    linePairs := [4][2][2]float32{
        {c2Lines[0], c2Lines[2]},
        {c2Lines[2], c2Lines[1]},
        {c2Lines[1], c2Lines[3]},
        {c2Lines[3], c2Lines[0]},
    }

    for _, c1Line := range(c1Lines) {
        for _, linePair := range(linePairs) {
            c2Line := linePair[0]

            if math.Abs(float64(c1Line[0] - c2Line[0])) < 0.001 {
                continue
            }

            x1 := (c2Line[1] - c1Line[1])/(c1Line[0] - c2Line[0])
            y1 := c1Line[0]*x1 + c1Line[1]

            c2Line = linePair[1]

            if math.Abs(float64(c1Line[0] - c2Line[0])) < 0.001 {
                continue
            }

            x2 := (c2Line[1] - c1Line[1])/(c1Line[0] - c2Line[0])
            y2 := c1Line[0]*x2 + c1Line[1]

            dx := math.Abs(float64(x1-x2))
            dy := math.Abs(float64(y1-y2))

            d := math.Sqrt(dx*dx + dy*dy)

            if d < 2{
                car.owner.Vibrate()
                car2.owner.Vibrate()

                carAbs := Vector{(x1 + x2)/2.0, (y1 + y2)/2.0}
                carRel := carAbs.Sub(car.position)
                car2Rel := carAbs.Sub(car2.position)

                carRel = carRel.Rotate(-newAngle)
                car2Rel = car2Rel.Rotate(-newAngle)

                car1Vel :=
                car.velocity
                car2Vel := car2.velocity
                car1AngVel := car.angularVelocity
                car2AngVel := car2.angularVelocity

                car.force = car.force.MulScalar(-2)
                car2.force = car.force.MulScalar(-2)
                car.velocity = car2Vel
                car2.velocity = car1Vel

                car.torque *= -1
                car2.torque *=-1
                car.angularVelocity = car2AngVel
                car2.angularVelocity = car1AngVel
                
                forceOn1 := car2.force
                forceOn2 := car.force

                car.AddForce(forceOn1, car.RelativeToWorld(carRel))
                car2.AddForce(forceOn2, car2.RelativeToWorld(car2Rel))
                return 
            }
        }
    }
}

func (r *Racer) HandleCollisions() {
	collisionHandled := [][2]*Car{}
	for _, car := range r.cars {
		for _, car2 := range r.cars {
			if car == car2 {
				continue
			}

			for _, collisions := range collisionHandled {
				if (car == collisions[0] && car2 == collisions[1]) || (car2 ==
					collisions[0] && car == collisions[1]) {
					continue
				}
			}

			if car2.position.Sub(car.position).Length() < (car.size+
				car2.size)/2 {
				collisionHandled = append(collisionHandled, [2]*Car{car,
					car2})
				col(car, car2)
			}
		}
	}
}
