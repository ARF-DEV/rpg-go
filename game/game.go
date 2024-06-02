package game

import (
	"fmt"
	"math"
	"sort"

	"github.com/ARF-DEV/rpg-go/engine"
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var gTimer engine.Timer

const WIDTH = 800
const HEIGHT = 600

type camera mgl32.Vec2

func (c *camera) MoveTo(pos mgl32.Vec2) {
	vecDiff := pos.Sub(c.GetCenter())
	*c = camera(mgl32.Vec2(*c).Add(vecDiff))
}

func (c *camera) Update(p *Player) {
	camCenter := c.GetCenter()
	tileDiff := camCenter.Sub(p.Position)
	if int64(mgl32.Abs(tileDiff[0])) > 3 {
		if tileDiff[0] > 0 {
			cam.MoveTo(camCenter.Add(mgl32.Vec2{-1, 0}))
		} else {
			cam.MoveTo(camCenter.Add(mgl32.Vec2{1, 0}))
		}
	}
	if int64(mgl32.Abs(tileDiff[1])) > 3 {
		if tileDiff[1] > 0 {
			cam.MoveTo(camCenter.Add(mgl32.Vec2{0, -1}))
		} else {
			cam.MoveTo(camCenter.Add(mgl32.Vec2{0, 1}))
		}
	}
}

// return coordinates in tile space
func (c *camera) GetCenter() mgl32.Vec2 {
	return mgl32.Vec2{
		float32(WIDTH / 2 / engine.TILE_SIZE),
		float32(HEIGHT / 2 / engine.TILE_SIZE),
	}.Add(mgl32.Vec2(*c))
}

var cam camera = camera{}
var traversalVis TileTraversalViz[Pos]
var aStarTravVis TileTraversalViz[PriorityPoint]

type Game struct {
	Player       Player
	CurrentLevel Level
	// State
	TextureMap map[string]*engine.Texture
}

func (g *Game) loadTexture(key, path string) error {
	tex, err := engine.CreateTexture(path)
	if err != nil {
		return err
	}
	g.TextureMap[key] = &tex
	return nil
}

func (g *Game) Update(window *glfw.Window, in *engine.Input) {
	gTimer.Update()
	glfw.PollEvents()
	g.Player.Update(in, &g.CurrentLevel)

	// traversalVis.Update(&g.CurrentLevel)
	aStarTravVis.Update(&g.CurrentLevel)

	// cam.MoveTo(g.Player.Position)
	// fmt.Println(cam)
	cam.Update(&g.Player)
}

func (g *Game) UpdateOnInput(in *engine.Input) {
	g.Player.UpdateOnInput(in, &g.CurrentLevel)

}

func (g *Game) init() {
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	g.TextureMap = map[string]*engine.Texture{}
}
func (g *Game) Start(in *engine.Input) {
	g.init()

	var err error
	if err := g.loadTexture("texture_atlas", "assets/atlas/default.png"); err != nil {
		panic(err)
	}

	g.Player = CreatePlayer(mgl32.Vec2{}, g.TextureMap["texture_atlas"])

	if err := engine.InitTextureIndex("assets/atlas/atlas-index.txt", g.TextureMap["texture_atlas"]); err != nil {
		panic(err)
	}

	g.CurrentLevel, err = LoadLevelFromFile("assets/maps/level1.map", &g.Player, g.TextureMap["texture_atlas"])
	if err != nil {
		panic(err)
	}
	in.AddSubcsriber(g)

	// gl.DrawBuffer(gl.FRONT)

	traversalVis = CreateTileTravViz(
		Pos{int32(g.Player.Position[0]), int32(g.Player.Position[1])},
		Pos{},
		func(trav *TileTraversalViz[Pos], lvl *Level) {
			if trav.visitedMap == nil {
				trav.visitedMap = map[Pos]bool{}
			}

			if trav.time > 1 {
				fmt.Println(trav.q.Len())
				if !trav.started {
					trav.q.Put(trav.start)
					trav.started = true
				}
				if trav.q.Len() > 0 {
					curPos := trav.q.Pop()
					if !trav.visitedMap[curPos] {
						trav.visitedMap[curPos] = true
						for _, neighbour := range getNeighbors(curPos, lvl) {
							if !trav.visitedMap[neighbour] {
								trav.q.Put(neighbour)
							}
						}
						trav.pVisited = curPos
						trav.time = 0
					}
				}
			}
			trav.time += gTimer.DeltaTime
		},
	)

	aStarTravVis = CreateTileTravViz[PriorityPoint](
		Pos{int32(14), int32(2)},
		Pos{int32(g.Player.Position[0]), int32(g.Player.Position[1])},
		func(trav *TileTraversalViz[PriorityPoint], lvl *Level) {
			if trav.visitedMap == nil {
				trav.visitedMap = map[Pos]bool{}
			}

			if trav.hasReach {
				for pos := trav.pVisited.Prev; pos != nil; pos = pos.Prev {
					trav.FinalPath = append(trav.FinalPath, pos.Pos)
				}
			}
			if trav.time > 0.2 {
				// fmt.Println(trav.q.Len())
				if !trav.started {
					start := PriorityPoint{
						trav.start,
						0,
						0,
						nil,
					}
					trav.q.Put(start)
					trav.started = true
				}
				if trav.q.Len() > 0 && !trav.hasReach {
					// fmt.Println(trav.q.values)
					sort.Slice(trav.q.values, func(i, j int) bool {
						return trav.q.values[i].Value < trav.q.values[j].Value
					})
					fmt.Println(trav.q.values)
					curPos := trav.q.Pop()
					if curPos.Pos == trav.end {
						trav.pVisited = curPos
						trav.hasReach = true

					}
					if !trav.visitedMap[curPos.Pos] {
						trav.visitedMap[curPos.Pos] = true
						for _, neighbour := range getNeighbors(curPos.Pos, lvl) {
							if !trav.visitedMap[neighbour] {
								diffY := int32(math.Abs(float64(trav.end.Y) - float64(neighbour.Y)))
								diffX := int32(math.Abs(float64(trav.end.X) - float64(neighbour.X)))
								fmt.Println(diffX)
								fmt.Println(diffY)
								fmt.Println(int32(math.Sqrt(float64(diffX)*float64(diffX) + float64(diffY)*float64(diffY))))
								newPriorityPoint := PriorityPoint{
									neighbour,
									(curPos.Step + 1) + diffX + diffY,
									curPos.Step + 1,
									&curPos,
								}
								// fmt.Println(newPriorityPoint)
								trav.q.Put(newPriorityPoint)
							}
						}
						trav.pVisited = curPos
						trav.time = 0
					}
				}
			}
			trav.time += gTimer.DeltaTime

		},
	)
}

func (g *Game) Draw(window *glfw.Window, sr engine.Renderer, shader *engine.Shader) {
	sr.Bind()
	sr.Clear()

	// draw start
	g.CurrentLevel.Draw(sr, shader)
	g.Player.Draw(sr, shader, &g.CurrentLevel)
	// traversalVis.Draw(sr, shader)
	aStarTravVis.Draw(sr, shader)

	sr.DebugDraw(shader, float32(WIDTH/2)-16, float32(HEIGHT/2)-16, 16, 16, engine.COLOR_BLACK)
	sr.UnBind()
	sr.Present()

	window.SwapBuffers()
}
