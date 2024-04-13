package game

import (
	"github.com/ARF-DEV/rpg-go/engine"
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

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
	// fmt.Println(g.Player.Posiiton)
	g.CurrentLevel.Draw(sr, shader)
	g.Player.Draw(sr, shader)
	window.SwapBuffers()

	gl.ClearColor(0.2, 0.5, 0.1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}
