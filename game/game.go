package game

import (
	"github.com/ARF-DEV/rpg-go/engine"
	"github.com/go-gl/gl/v3.3-core/gl"
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

func (g *Game) Update(window *glfw.Window) {
	glfw.PollEvents()
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
}

func (g *Game) init() {

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	g.TextureMap = map[string]*engine.Texture{}
}
func (g *Game) Start() {
	g.init()

	var err error
	if err := g.loadTexture("level_texture_atlas", "assets/atlas/default.png"); err != nil {
		panic(err)
	}

	if err := engine.InitTextureIndex("assets/atlas/atlas-index.txt", g.TextureMap["level_texture_atlas"]); err != nil {
		panic(err)
	}

	g.CurrentLevel, err = LoadLevelFromFile("assets/maps/level1.map", &g.Player, g.TextureMap["level_texture_atlas"])
	if err != nil {
		panic(err)
	}
}

func (g *Game) Draw(window *glfw.Window, sr *engine.SpriteRenderer, shader *engine.Shader) {
	// fmt.Println(g.Player.Posiiton)
	g.CurrentLevel.Draw(sr, shader)
	sr.Draw(shader, g.TextureMap["level_texture_atlas"], mgl32.Vec2{21 * 32, 59 * 32}, mgl32.Vec2{32, 32}, g.Player.Posiiton.Mul(32), mgl32.Vec2{32, 32}, 0, mgl32.Vec4{1, 1, 1, 1})
	window.SwapBuffers()

	gl.ClearColor(0.2, 0.5, 0.1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}
