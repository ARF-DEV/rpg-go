package game

import (
	"github.com/ARF-DEV/rpg-go/engine"
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

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
	glfw.PollEvents()
	g.Player.Update(in)
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

}

func (g *Game) Draw(window *glfw.Window, sr *engine.SpriteRenderer, shader *engine.Shader) {
	g.CurrentLevel.Draw(sr, shader)
	g.Player.Draw(sr, shader)
	sr.DebugDraw(shader, float32(WIDTH/2)-16, float32(HEIGHT/2)-16, 16, 16, engine.COLOR_BLACK)

	window.SwapBuffers()

	gl.ClearColor(0.2, 0.5, 0.1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}
