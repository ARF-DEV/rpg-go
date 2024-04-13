package engine

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/go-gl/gl/v4.3-core/gl"
)

type textureID uint32

type Texture struct {
	ID       textureID
	TexBound image.Rectangle
}

func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, uint32(t.ID))
}
func (t *Texture) UnBind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func loadImage(path string, flip bool) (*image.RGBA, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if flip {
		for x := 0; x < rgba.Bounds().Dx(); x++ {
			for y := 0; y < rgba.Bounds().Dy(); y++ {
				rgba.Set(x, y, img.At(img.Bounds().Dx()-1-x, img.Bounds().Dy()-1-y))
			}
		}
	} else {
		draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
	}

	return rgba, nil
}
func CreateTexture(path string) (Texture, error) {
	img, err := loadImage(path, false)
	if err != nil {
		return Texture{}, err
	}
	fmt.Println(img.Bounds())

	// generate texture
	// activate texture unit (optional??)
	// bind texture to texture_2d
	// set texture parameter
	// set texture data
	// gen mipmap (optional)
	var texID uint32
	gl.ActiveTexture(gl.TEXTURE0)
	gl.GenTextures(1, &texID)
	gl.BindTexture(gl.TEXTURE_2D, texID)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(img.Bounds().Dx()), int32(img.Bounds().Dy()), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	// gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.BindTexture(gl.TEXTURE_2D, 0)
	return Texture{
		ID:       textureID(texID),
		TexBound: img.Rect,
	}, nil
}
