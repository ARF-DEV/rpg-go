package game

import (
	"github.com/ARF-DEV/rpg-go/engine"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Player struct {
	Position mgl32.Vec2
	tex      *engine.Texture
}

func CreatePlayer(pos mgl32.Vec2, tex *engine.Texture) Player {
	return Player{
		Position: pos,
		tex:      tex,
	}
}

func (p *Player) Draw(sr *engine.SpriteRenderer, shader *engine.Shader) {
	sr.Draw(shader, p.tex, mgl32.Vec2{21 * 32, 59 * 32}, mgl32.Vec2{32, 32}, p.Position.Mul(32), mgl32.Vec2{32, 32}, 0, mgl32.Vec4{1, 1, 1, 1})
}

func (p *Player) Update(in *engine.Input) {
}

func (p *Player) UpdateOnInput(in *engine.Input, lvl *Level) {
	proposedPos := p.Position
	if in.Keys[glfw.KeyA] {
		proposedPos[0] -= 1
	}
	if in.Keys[glfw.KeyD] {
		proposedPos[0] += 1
	}
	if in.Keys[glfw.KeyW] {
		proposedPos[1] -= 1
	}
	if in.Keys[glfw.KeyS] {
		proposedPos[1] += 1
	}

	// TODO CHECK IF WALKABLE
	if lvl.GetTIle(int(proposedPos[0]), int(proposedPos[1])).IsWalkable {
		p.Position = proposedPos
	}
}
