package systems

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

type TextureStore struct {
	textures map[string]uint32
}

func NewTextureStore() *TextureStore {
	return &TextureStore{
		textures: make(map[string]uint32),
	}
}

// GetTexture ensures the texture is loaded only once and reused thereafter.
func (ts *TextureStore) GetTexture(texturePath string) (uint32, error) {
	if texture, exists := ts.textures[texturePath]; exists {
		return texture, nil
	}

	newTexture, err := prepareTexture(texturePath)
	if err != nil {
		return 0, err
	}

	ts.textures[texturePath] = newTexture
	return newTexture, nil
}

func prepareTexture(texturePath string) (uint32, error) {
	img, err := loadImage(texturePath)
	if err != nil {
		return 0, err
	}

	rgbaFlipped := imageToRGBA(img)
	rgba := flipImageVertically(rgbaFlipped)

	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D, textureID)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	width, height := rgba.Rect.Size().X, rgba.Rect.Size().Y
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(width), int32(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	return textureID, nil
}

func loadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if strings.HasSuffix(strings.ToLower(filename), ".jpg") || strings.HasSuffix(strings.ToLower(filename), ".jpeg") {
		return jpeg.Decode(file)
	} else if strings.HasSuffix(strings.ToLower(filename), ".png") {
		return png.Decode(file)
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
