package imgutil

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

// The PNGImage represents a structure for working with PNG images.
type PNGImage struct {
	in        image.Image
	out       *image.RGBA64
	kernelDim int
}

// Load returns a PNGImage that was loaded based on the filePath parameter
func Load(filePath string) (*PNGImage, error) {

	inReader, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}
	defer inReader.Close()

	inImg, err := png.Decode(inReader)

	if err != nil {
		return nil, err
	}

	inBounds := inImg.Bounds()

	outImg := image.NewRGBA64(inBounds)

	return &PNGImage{inImg, outImg, 3}, nil
}

// Save saves the image to the given file
func (img *PNGImage) Save(filePath string) error {

	outWriter, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outWriter.Close()

	err = png.Encode(outWriter, img.out)
	if err != nil {
		return err
	}
	return nil
}

//clamp will clamp the comp parameter to zero if it is less than zero or to 65535 if the comp parameter
// is greater than 65535.
func clamp(comp float64) uint16 {
	return uint16(math.Min(65535, math.Max(0, comp)))
}

// Grayscale applies a grayscale filtering effect to the image
func (img *PNGImage) Grayscale() {

	bounds := img.out.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.in.At(x, y).RGBA()
			greyC := clamp(float64(r+g+b) / 3)
			img.out.Set(x, y, color.RGBA64{greyC, greyC, greyC, uint16(a)})
		}
	}
}

// Blur applies a blur filtering effect to the image
func (img *PNGImage) Blur() {
	kernel := [][]float64{
		{1. / 9, 1. / 9, 1. / 9},
		{1. / 9, 1. / 9, 1. / 9},
		{1. / 9, 1. / 9, 1. / 9}}

	img.Convolution(kernel)
}

// Convolution performs image convolution given a kernel of specified dimension.
func (img *PNGImage) Convolution(kernel [][]float64) {
	bounds := img.out.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			neighbors := img.neighbors(x, y)
			var r, g, b, a uint16
			for i := 0; i < img.kernelDim; i++ {
				for j := 0; j < img.kernelDim; j++ {
					neighbor := neighbors[i][j]
					if neighbor == nil {
						continue
					}
					nr, ng, nb, na := neighbor.RGBA()
					r += clamp(kernel[i][j] * float64(nr))
					g += clamp(kernel[i][j] * float64(ng))
					b += clamp(kernel[i][j] * float64(nb))
					a += clamp(kernel[i][j] * float64(na))
				}
			}
			img.out.Set(x, y, color.RGBA64{r, g, b, a})
		}
	}
}

func (img *PNGImage) neighbors(x, y int) [][]color.Color {
	bounds := img.out.Bounds()

	neighbors := make([][]color.Color, img.kernelDim)
	for i := range neighbors {
		neighbors[i] = make([]color.Color, img.kernelDim)
	}

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			// Edge detection
			if x+i < bounds.Min.X || x+i >= bounds.Max.X || y+j < bounds.Min.Y || y+j >= bounds.Max.Y {
				continue
			}
			color := img.in.At(x+i, y+j)
			// fmt.Println("Color at", x+i, ",", y+j, "is", color)
			neighbors[i+1][j+1] = color
		}
	}
	return neighbors
}
