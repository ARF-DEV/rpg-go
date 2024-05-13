package engine

import "github.com/go-gl/mathgl/mgl32"

type Drawable interface {
	Draw(sr Renderer, shader *Shader)
	// DebugDraw(shader *Shader, x, y float32, sizeX, sizeY float32, color mgl32.Vec4)
}

type Renderer interface {
	Draw(
		shader *Shader,
		texture *Texture,
		srcPos,
		srcSize,
		dstPos,
		dstSize mgl32.Vec2,
		rotate float32,
		color mgl32.Vec4,
	)
	DebugDraw(shader *Shader, x, y float32, sizeX, sizeY float32, color mgl32.Vec4)
	Bind()
	Present()
	UnBind()
	Clear()
}
