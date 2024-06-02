package game

import (
	"fmt"

	"github.com/ARF-DEV/rpg-go/engine"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Player struct {
	Position mgl32.Vec2
	tex      *engine.Texture
	prevDir  mgl32.Vec2
	Search   bool
}

func CreatePlayer(pos mgl32.Vec2, tex *engine.Texture) Player {
	return Player{
		Position: pos,
		tex:      tex,
	}
}

func (p *Player) Draw(sr engine.Renderer, shader *engine.Shader, lvl *Level) {
	front := p.Position.Add(p.prevDir).Mul(32).Add(mgl32.Vec2{16, 16})
	sr.Draw(shader, p.tex, mgl32.Vec2{21 * 32, 59 * 32}, mgl32.Vec2{32, 32}, p.Position.Sub(mgl32.Vec2(cam)).Mul(32), mgl32.Vec2{32, 32}, 0, mgl32.Vec4{1, 1, 1, 1})

	camOffsetPix := mgl32.Vec2(cam).Mul(32)
	sr.DebugDraw(shader, -camOffsetPix[0]+front[0], -camOffsetPix[1]+front[1], 10, 10, engine.COLOR_WHITE)

	// shader := engine.ShaderMap["defaultShader"]
	if p.Search {
		bfs(shader, &engine.DebugSpriteRenderer, lvl, Pos{int32(p.Position[0]), int32(p.Position[1])})
		p.Search = false
	}
}

func (p *Player) Update(in *engine.Input, lvl *Level) {
}

// NOTE: make a global shader map
func (p *Player) UpdateOnInput(in *engine.Input, lvl *Level) {
	proposedPos := p.Position
	if in.Keys[glfw.KeyA] && !in.PrevKeys[glfw.KeyA] {
		proposedPos[0] -= 1
	} else if in.Keys[glfw.KeyD] && !in.PrevKeys[glfw.KeyD] {
		proposedPos[0] += 1
	} else if in.Keys[glfw.KeyW] && !in.PrevKeys[glfw.KeyW] {
		proposedPos[1] -= 1
	} else if in.Keys[glfw.KeyS] && !in.PrevKeys[glfw.KeyS] {
		proposedPos[1] += 1
	}

	if in.Keys[glfw.KeyE] && !in.PrevKeys[glfw.KeyE] {

		front := p.Position.Add(p.prevDir)
		frontTile := lvl.GetTile(int(front[0]), int(front[1]))

		fmt.Println(p.prevDir, p.Position, *frontTile, front)
		if frontTile.IsInteractable() {
			frontTile.Interact()
		}
	}
	if proposedPos != p.Position {
		p.moveTo(lvl, proposedPos)
	}
	if in.Keys[glfw.KeyG] && !in.PrevKeys[glfw.KeyG] {
		p.Search = true

		// shader := engine.ShaderMap["defaultShader"]
		// bfs(&shader, &engine.DebugSpriteRenderer, lvl, Pos{int32(p.Position[0]), int32(p.Position[1])})
	}
}

func (p *Player) moveTo(lvl *Level, pos mgl32.Vec2) {
	p.prevDir = pos.Sub(p.Position)
	fmt.Println(p.Position, pos, p.prevDir)
	if lvl.GetTile(int(pos[0]), int(pos[1])).CanWalk() {
		p.Position = pos
	}
}
