package mimg

import (
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
)

type FC interface {
	RGBA() (r, g, b, a float64)
}

type FRGBA struct {
	R, G, B, A float64
}

func (f *FRGBA) RGBA() (float64, float64, float64, float64) {
	return f.R, f.G, f.B, f.A
}

type FColor struct {
	c color.Color
}

func u2F(d uint32) float64 {
	d = d >> 8
	return float64(d) / 255.0
}

func (f *FColor) RGBA() (r, g, b, a float64) {
	ri, gi, bi, ai := f.c.RGBA()
	r, g, b, a = u2F(ri), u2F(gi), u2F(bi), u2F(ai)
	return
}

type FImage struct {
	src image.Image
}

func (f *FImage) ColorModel() color.Model {
	return f.src.ColorModel()
}

func (f *FImage) Bounds() image.Rectangle {
	return f.src.Bounds()
}

func (f *FImage) At(x, y int) FC {
	return &FColor{f.src.At(x, y)}
}

type Image struct {
	src FI
}

func f2U(d float64) uint8 {
	if d < 0.0 {
		d = 0.0
	}
	if d > 1.0 {
		d = 1.0
	}
	return uint8(d * 255.0)
}

func (i *Image) At(x, y int) color.Color {
	c := i.src.At(x, y)
	r, g, b, a := c.RGBA()
	return &color.RGBA{f2U(r), f2U(g), f2U(b), f2U(a)}
}

func (i *Image) Bounds() image.Rectangle {
	return i.src.Bounds()
}

func (i *Image) ColorModel() color.Model {
	return color.RGBAModel
}

type FI interface {
	At(int, int) FC
	Bounds() image.Rectangle
	ColorModel() color.Model
}

func FI2Image(src FI) image.Image {
	return &Image{src}
}

func Save(f FI, filename string, mod int) error {
	d := FI2Image(f)
	return saveImage(d, filename, mod)
}

func NewFI(src image.Image) FI {
	return &FImage{src}
}

func savePng(src image.Image, filename string) error {
	fp, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fp.Close()
	return png.Encode(fp, src)
}

const (
	PNG  = 1
	JPEG = 2
)

func saveJpg(src image.Image, filename string) error {
	fp, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fp.Close()
	return jpeg.Encode(fp, src, nil)
}

func saveImage(src image.Image, filename string, mod int) error {
	switch mod {
	case PNG:
		return savePng(src, filename)
	case JPEG:
		return saveJpg(src, filename)
	default:
		return errors.New("unknown save mod")
	}
}

func Load(filename string) (FI, error) {
	w, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer w.Close()
	src, _, err := image.Decode(w)
	if err != nil {
		return nil, err
	}
	return NewFI(src), err
}
