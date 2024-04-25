package components

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.3-core/gl"
)

type TextureComponent struct {
	TextureID uint32
}

func NewTextureComponent(imgPath string) (*TextureComponent, error) {
	// Load the image
	img, err := loadImage(imgPath)
	if err != nil {
		return nil, err
	}

	// Convert image to RGBA
	rgbaFlipped := imageToRGBA(img)

	rgba := flipImageVertically(rgbaFlipped)

	// Generate and bind a new texture
	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D, textureID)

	// Set texture parameters
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	// Upload texture data to GPU
	width, height := rgba.Rect.Size().X, rgba.Rect.Size().Y
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(width), int32(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))

	// Generate mipmaps
	gl.GenerateMipmap(gl.TEXTURE_2D)

	// Unbind the texture
	gl.BindTexture(gl.TEXTURE_2D, 0)

	return &TextureComponent{TextureID: textureID}, nil
}

func loadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Determine the file format based on the file extension
	if strings.HasSuffix(strings.ToLower(filename), ".jpg") || strings.HasSuffix(strings.ToLower(filename), ".jpeg") {
		img, err := jpeg.Decode(file)
		if err != nil {
			return nil, err
		}
		return img, nil
	} else if strings.HasSuffix(strings.ToLower(filename), ".png") {
		img, err := png.Decode(file)
		if err != nil {
			return nil, err
		}
		return img, nil
	}

	return nil, fmt.Errorf("unsupported file format for %v", filename)
}

func imageToRGBA(img image.Image) *image.RGBA {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)
	return rgba
}

func flipImageVertically(img *image.RGBA) *image.RGBA {
	src := img.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, src.Dx(), src.Dy()))
	for y := 0; y < src.Dy(); y++ {
		for x := 0; x < src.Dx(); x++ {
			dst.Set(x, src.Dy()-y-1, img.At(x, y))
		}
	}
	return dst
}
