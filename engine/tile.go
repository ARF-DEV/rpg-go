package engine

import (
	"bufio"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
)

type TileType string

var textureIndex map[TileType][]image.Rectangle

const (
	Blank     TileType = " " // 00000000
	StoneWall TileType = "#" // 00000010
	Door      TileType = "|" // 00000100
	SandFloor TileType = "." // 00001000

	TILE_SIZE uint32 = 32
)

type Tile struct {
	Pos        mgl32.Vec2
	IsWalkable bool
	Type       TileType
}

func GetTextureIndex() map[TileType][]image.Rectangle {
	return textureIndex
}

func InitTextureIndex(file string, texture *Texture) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	// imgF, err := os.Open(file)
	// if err != nil {
	// 	return err
	// }
	// img, _, err := image.Decode(imgF)
	// if err != nil {
	// 	return err
	// }
	img := texture.TexBound

	textureIndex = map[TileType][]image.Rectangle{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		splits := strings.Split(line, " ")

		texCoords := []image.Rectangle{}
		identifiers, coords := splits[0], splits[1]
		coordsArr := strings.Split(coords, ",")

		fmt.Println(identifiers, coordsArr)

		xStart, err := strconv.Atoi(string(coordsArr[0]))
		if err != nil {
			return err
		}
		yStart, err := strconv.Atoi(string(coordsArr[1]))
		if err != nil {
			return err
		}
		variants, err := strconv.Atoi(string(coordsArr[2]))
		if err != nil {
			return err
		}

		for i := 0; i < variants; i++ {
			// x := xStart + (int(TILE_SIZE) * (i + 1))
			// y := yStart * (int(TILE_SIZE))
			x := xStart*int(TILE_SIZE) + (int(TILE_SIZE) * i)
			y := yStart * int(TILE_SIZE)
			if x > img.Bounds().Dx() {
				xOver := x - img.Bounds().Dx()
				// xOver div by img width in this case always > 1
				linePass := int(math.Ceil(float64(xOver) / float64(img.Bounds().Dx())))
				x = x - (linePass * img.Bounds().Dx())
				y = y + (linePass * int(TILE_SIZE))
			}
			texCoords = append(texCoords, image.Rect(x, y, x+int(TILE_SIZE), y+int(TILE_SIZE)))
		}
		textureIndex[TileType(identifiers)] = texCoords
	}
	fmt.Println(textureIndex)
	return nil
}
