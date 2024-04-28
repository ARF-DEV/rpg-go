package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"

	"github.com/ARF-DEV/rpg-go/engine"
	"github.com/go-gl/mathgl/mgl32"
)

type Level struct {
	Map          [][]engine.Tile
	textureAtlas *engine.Texture
}

func (l *Level) GetTile(x, y int) *engine.Tile {
	return &l.Map[y][x]
}

func (l *Level) Draw(renderer *engine.SpriteRenderer, shader *engine.Shader) {
	texture := l.textureAtlas
	if l.textureAtlas == nil {
		panic(fmt.Errorf("error texture atlas not set"))
	}

	r := rand.New(rand.NewSource(1))
	for y, tileLine := range l.Map {
		for x, tile := range tileLine {
			if tile != engine.Blank {
				src := engine.GetTextureIndex()[tile]
				// fmt.Println(src, tile)
				idx := r.Intn(len(src))
				srcPos := mgl32.Vec2{float32(src[idx].Min.X), float32(src[idx].Min.Y)}
				srcSize := mgl32.Vec2{float32(engine.TILE_SIZE), float32(engine.TILE_SIZE)}
				renderer.Draw(shader, texture, srcPos, srcSize, mgl32.Vec2{float32((int(cam[0]) + x) * 32), float32((int(cam[1]) + y) * 32)}, mgl32.Vec2{32, 32}, 0, mgl32.Vec4{1, 1, 1, 1})
			}
		}
	}
}

func LoadLevelFromFile(path string, player *Player, texAtlas *engine.Texture) (Level, error) {
	f, err := os.Open(path)
	if err != nil {
		return Level{}, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	levelLines := []string{}
	longestWidth := 0
	for scanner.Scan() {
		l := scanner.Text()
		if len(l) > longestWidth {
			longestWidth = len(l)
		}

		levelLines = append(levelLines, l)
	}

	var level Level
	var t engine.Tile
	for y, line := range levelLines {
		var tiles []engine.Tile
		for x, r := range line {
			switch r {
			case '\n', '\r', ' ', '\t':
				t = engine.Blank
			case '#':
				t = engine.StoneWall
			case '|':
				t = engine.ClosedDoor
			case '/':
				t = engine.OpenDoor
			case '.':
				t = engine.SandFloor
			case 'p':
				t = engine.SandFloor
				player.Position[0] = float32(x)
				player.Position[1] = float32(y)
				// set player pos
			default:
				return Level{}, fmt.Errorf("error: case for atlas index %c is not yet set for ", r)
			}
			tiles = append(tiles, t)
		}

		if len(line) < longestWidth {
			for i := 0; i < longestWidth-len(line); i++ {
				tiles = append(tiles, engine.Blank)
			}
		}
		level.Map = append(level.Map, tiles)
	}
	level.textureAtlas = texAtlas

	return level, nil
}
