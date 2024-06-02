package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ARF-DEV/rpg-go/engine"
	"github.com/go-gl/mathgl/mgl32"
)

type Pos struct {
	X, Y int32
}

func (p Pos) GetX() int32 {
	return p.X
}
func (p Pos) GetY() int32 {
	return p.Y
}

type Level struct {
	Map          [][]engine.Tile
	textureAtlas *engine.Texture
	DebugMap     map[Pos]bool
}

func (l *Level) GetTile(x, y int) *engine.Tile {
	return &l.Map[y][x]
}

func (l *Level) Draw(renderer engine.Renderer, shader *engine.Shader) {
	texture := l.textureAtlas
	if l.textureAtlas == nil {
		panic(fmt.Errorf("error texture atlas not set"))
	}

	r := rand.New(rand.NewSource(1))
	for y, tileLine := range l.Map {
		for x, tile := range tileLine {
			if tile != engine.Blank {
				src := engine.GetTextureIndex()[tile]
				idx := r.Intn(len(src))
				srcPos := mgl32.Vec2{float32(src[idx].Min.X), float32(src[idx].Min.Y)}
				srcSize := mgl32.Vec2{float32(engine.TILE_SIZE), float32(engine.TILE_SIZE)}
				renderer.Draw(shader, texture, srcPos, srcSize, mgl32.Vec2{float32((int(-cam[0]) + x) * 32), float32((int(-cam[1]) + y) * 32)}, mgl32.Vec2{32, 32}, 0, mgl32.Vec4{1, 1, 1, 1})
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
	level.DebugMap = map[Pos]bool{}
	return level, nil
}

func bfs(shader *engine.Shader, sr engine.Renderer, level *Level, startPos Pos) {
	bfsQ := Queue[Pos]{}
	visitedMap := map[Pos]bool{}
	bfsQ.Put(startPos)
	for bfsQ.Len() > 0 {
		// fmt.Println(bfsQ.Len())
		curPos := bfsQ.Pop()
		visitedMap[curPos] = true
		// fmt.Println(bfsQ.Len())
		// sr.Bind()
		sr.DebugDraw(shader, float32(engine.TILE_SIZE*uint32(curPos.X)), float32(engine.TILE_SIZE*uint32(curPos.Y)), float32(engine.TILE_SIZE), float32(engine.TILE_SIZE), mgl32.Vec4{1, 0, 0, 1})
		// sr.UnBind()
		// sr.Clear()
		// sr.Present()

		// sr.CopyDraw(, shader)
		// glfw.GetCurrentContext().SwapBuffers()
		time.Sleep(time.Millisecond * 50)
		for _, neighbour := range getNeighbors(curPos, level) {
			if !visitedMap[neighbour] {
				bfsQ.Put(neighbour)
				// level.DebugMap[Pos{neighbour.X, neighbour.Y}] = true
			}
		}

	}
	fmt.Println("DONE")

}

func getNeighbors(curPos Pos, level *Level) []Pos {
	neighbours := []Pos{}
	left := Pos{curPos.X - 1, curPos.Y}
	right := Pos{curPos.X + 1, curPos.Y}
	up := Pos{curPos.X, curPos.Y - 1}
	down := Pos{curPos.X, curPos.Y + 1}

	if level.GetTile(int(left.X), int(left.Y)).CanWalk() {
		neighbours = append(neighbours, left)

	}
	if level.GetTile(int(right.X), int(right.Y)).CanWalk() {
		neighbours = append(neighbours, right)
	}

	if level.GetTile(int(up.X), int(up.Y)).CanWalk() {
		neighbours = append(neighbours, up)
	}
	if level.GetTile(int(down.X), int(down.Y)).CanWalk() {
		neighbours = append(neighbours, down)
	}

	return neighbours
}
