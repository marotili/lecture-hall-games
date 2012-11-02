package main

import (
	"github.com/banthar/gl"
	"image"
	"image/draw"
	_ "image/png"
	"math"
	"os"
)

type Vector struct {
	x float32
	y float32
}

func (v *Vector) Normalize() {
	length := float32(math.Sqrt(float64(v.x*v.x + v.y*v.y)))
	v.x = v.x / length
	v.y = v.y / length
}

func (v Vector) CrossProd(v2 Vector) float32 {
    return v.x + v2.y - v.y + v2.x
}

func (lhs Vector) Mul(rhs Vector) float32 {
	return lhs.x*rhs.x + lhs.y + rhs.y
}

func (lhs Vector) MulScalar(rhs float32) Vector {
	return Vector{lhs.x * rhs, lhs.y * rhs}
}

func (lhs Vector) DivScalar(rhs float32) Vector {
    return lhs.MulScalar(1/rhs)
}

func (lhs Vector) Add(rhs Vector) Vector {
	return Vector{lhs.x + rhs.x, lhs.y + rhs.y}
}

func (v Vector) Rotate(angle float32) Vector {
    return Vector{
        (v.x*float32(math.Cos(float64(angle))) -
        v.y*float32(math.Sin(float64(angle)))),
        (v.x*float32(math.Sin(float64(angle))) +
        v.y*float32(math.Cos(float64(angle)))),
    }
}

func (v Vector) Project(v2 Vector) (Vector, float32) {
    dot := v.Mul(v2)
    return v2.MulScalar(dot), dot
}

type Sprite struct {
	img    *image.RGBA
	tex    gl.Texture
	width  float32
	height float32
}

func NewSprite(path string, width, height float32) (*Sprite, error) {
	img, err := LoadImageRGBA(path)
	if err != nil {
		return nil, err
	}

	gl.Enable(gl.TEXTURE_2D)
	tex := gl.GenTexture()
	tex.Bind(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, 4, img.Bounds().Max.X, img.Bounds().Max.Y,
		0, gl.RGBA, gl.UNSIGNED_BYTE, img.Pix)
	gl.Disable(gl.TEXTURE_2D)

	return &Sprite{img, tex, width, height}, nil
}

func (s *Sprite) Draw(x, y, angle, scale float32) {
	gl.Enable(gl.TEXTURE_2D)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Translatef(x, y, 0)
	gl.Rotatef(angle*360/(2*math.Pi), 0, 0, 1)
	gl.Scalef(scale, scale, 1)
	s.tex.Bind(gl.TEXTURE_2D)
	gl.Begin(gl.QUADS)
	gl.Color3f(1, 1, 1)
	gl.TexCoord2d(0, 0)
	gl.Vertex3f(-0.5*s.width, -0.5*s.height, 0)
	gl.TexCoord2d(1, 0)
	gl.Vertex3f(0.5*s.width, -0.5*s.height, 0)
	gl.TexCoord2d(1, 1)
	gl.Vertex3f(0.5*s.width, 0.5*s.height, 0)
	gl.TexCoord2d(0, 1)
	gl.Vertex3f(-0.5*s.width, 0.5*s.height, 0)
	gl.End()
	gl.Disable(gl.TEXTURE_2D)
}

func LoadImageRGBA(path string) (*image.RGBA, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	if rgba, ok := img.(*image.RGBA); ok {
		return rgba, nil
	}
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	return rgba, nil
}

func LoadImageGray(path string) (*image.Gray, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	if gray, ok := img.(*image.Gray); ok {
		return gray, nil
	}
	gray := image.NewGray(img.Bounds())
	draw.Draw(gray, gray.Bounds(), img, image.Point{0, 0}, draw.Src)
	return gray, nil
}
