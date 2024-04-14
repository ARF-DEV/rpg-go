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
)

type Tile string

func (t Tile) CanWalk() bool {
	switch t {
	case Blank, StoneWall, ClosedDoor:
		return false
	case SandFloor, OpenDoor:
		return true
	}
	return false
}

func (t Tile) IsInteractable() bool {
	switch t {
	case ClosedDoor, OpenDoor:
		return true
	}
	return false
}

func (t *Tile) Interact() {
	switch *t {
	case ClosedDoor:
		*t = OpenDoor
	case OpenDoor:
		*t = ClosedDoor
	}
	fmt.Println(*t)
}

var textureIndex map[Tile][]image.Rectangle

const (
	Blank      Tile = " "
	StoneWall  Tile = "#"
	OpenDoor   Tile = "/"
	ClosedDoor Tile = "|"
	SandFloor  Tile = "."

	TILE_SIZE uint32 = 32
)

func GetTextureIndex() map[Tile][]image.Rectangle {
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

	textureIndex = map[Tile][]image.Rectangle{}
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
		t := Tile(identifiers)

		textureIndex[t] = texCoords
	}
	fmt.Println(textureIndex)
	return nil
}
