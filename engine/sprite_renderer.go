package engine

import (
	"fmt"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	COLOR_BLACK       mgl32.Vec4 = mgl32.Vec4{0, 0, 0, 1}
	COLOR_WHITE       mgl32.Vec4 = mgl32.Vec4{1, 1, 1, 1}
	COLOR_GREEN       mgl32.Vec4 = mgl32.Vec4{0, 1, 0, 1}
	COLOR_BLUE        mgl32.Vec4 = mgl32.Vec4{0, 0, 1, 1}
	COLOR_RED         mgl32.Vec4 = mgl32.Vec4{1, 0, 0, 1}
	COLOR_BLACK_FADED mgl32.Vec4 = mgl32.Vec4{0, 0, 0, 0.5}
	COLOR_WHITE_FADED mgl32.Vec4 = mgl32.Vec4{1, 1, 1, 0.5}
	COLOR_GREEN_FADED mgl32.Vec4 = mgl32.Vec4{0, 1, 0, 0.5}
	COLOR_BLUE_FADED  mgl32.Vec4 = mgl32.Vec4{0, 0, 1, 0.5}
	COLOR_RED_FADED   mgl32.Vec4 = mgl32.Vec4{1, 0, 0, 0.5}
)

var DebugSpriteRenderer SpriteRendererV2

type SpriteRendererV2 struct {
	quadVAO          uint32
	screenQuadVAO    uint32
	textureColorBuff uint32
	frameBuff        uint32
	rbo              uint32
	screenShader     *Shader
}

func CreateRenderer(width, height int32) (SpriteRendererV2, error) {
	r := SpriteRendererV2{}
	vertices := []float32{
		-1, -1, 0, 0,
		1, -1, 1, 0,
		-1, 1, 0, 1,
		1, 1, 1, 1,
	}

	quadVertices := []float32{
		0, 0, 0, 0,
		1, 0, 1, 0,
		0, 1, 0, 1,
		1, 1, 1, 1,
	}

	elements := []uint32{
		0, 1, 2,
		2, 3, 1,
	}

	// Generate screenQuadVAO for screen buffers
	gl.CreateVertexArrays(1, &r.screenQuadVAO)
	gl.BindVertexArray(r.screenQuadVAO)

	var VBO, EBO uint32
	gl.GenBuffers(1, &VBO)
	gl.GenBuffers(1, &EBO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)

	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(elements), gl.Ptr(elements), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 2, gl.FLOAT, false, 4*4, 0)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 4*4, 2*4)

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)

	// Generate QuadVAO for off screen buffers
	gl.CreateVertexArrays(1, &r.quadVAO)
	gl.BindVertexArray(r.quadVAO)

	gl.GenBuffers(1, &VBO)
	gl.GenBuffers(1, &EBO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)

	gl.BufferData(gl.ARRAY_BUFFER, 4*len(quadVertices), gl.Ptr(quadVertices), gl.STATIC_DRAW)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(elements), gl.Ptr(elements), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 2, gl.FLOAT, false, 4*4, 0)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 4*4, 2*4)

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)

	gl.GenFramebuffers(1, &r.frameBuff)

	gl.GenTextures(1, &r.textureColorBuff)
	gl.BindTexture(gl.TEXTURE_2D, r.textureColorBuff)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, width, height, 0, gl.RGB, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.GenRenderbuffers(1, &r.rbo)
	gl.BindRenderbuffer(gl.RENDERBUFFER, r.rbo)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH24_STENCIL8, width, height)

	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, r.frameBuff)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, r.textureColorBuff, 0)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, r.rbo)

	if gl.CheckFramebufferStatus(gl.FRAMEBUFFER) != gl.FRAMEBUFFER_COMPLETE {
		return SpriteRendererV2{}, fmt.Errorf("error: framebuffer is not complete")
	}

	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)

	// TODO: load shader
	shader, err := CreateShader("assets/screen_vertex.glsl", "assets/screen_fragment.glsl", "screenShader")
	if err != nil {
		return SpriteRendererV2{}, err
	}
	r.screenShader = &shader

	if DebugSpriteRenderer == (SpriteRendererV2{}) {
		DebugSpriteRenderer = r
	}
	return r, nil
}

func (r *SpriteRendererV2) Bind() {
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, r.frameBuff)
	// gl.Enable(gl.DEPTH_TEST)
}

func (r *SpriteRendererV2) UnBind() {
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	// gl.Disable(gl.DEPTH_TEST)
}
func (r *SpriteRendererV2) CopyDraw(obj Drawable, shader *Shader) {
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, r.frameBuff)
	gl.Enable(gl.DEPTH_TEST)

	obj.Draw(r, shader)

	gl.Disable(gl.DEPTH_TEST)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
}

func (r *SpriteRendererV2) Clear() {
	gl.ClearColor(0.2, 0.5, 0.1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (r *SpriteRendererV2) Present() {
	r.screenShader.Use()
	gl.BindVertexArray(r.screenQuadVAO)
	gl.BindTexture(gl.TEXTURE_2D, r.textureColorBuff)
	gl.DrawElementsWithOffset(gl.TRIANGLES, 6, gl.UNSIGNED_INT, 0)
	gl.BindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	r.screenShader.Unuse()

	// glfw.GetCurrentContext().SwapBuffers()
}

func (r *SpriteRendererV2) Draw(
	shader *Shader,
	texture *Texture,
	srcPos,
	srcSize,
	dstPos,
	dstSize mgl32.Vec2,
	rotate float32,
	color mgl32.Vec4,

) {
	translate := mgl32.Translate3D(dstPos[0], dstPos[1], 0)

	rot := mgl32.Translate3D(0.5, 0.5, 0).
		Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(rotate))).
		Mul4(mgl32.Translate3D(-0.5, -0.5, 0))

	scale := mgl32.Scale3D(dstSize[0], dstSize[1], 1)
	model := translate.Mul4(rot).Mul4(scale)

	if srcPos != (mgl32.Vec2{}) {
		srcPos[0] /= float32(texture.TexBound.Dx())
		srcPos[1] /= float32(texture.TexBound.Dy())
	}

	if srcSize != (mgl32.Vec2{}) {
		srcSize[0] /= float32(texture.TexBound.Dx())
		srcSize[1] /= float32(texture.TexBound.Dy())
	}
	shader.Use()
	shader.SetMatrix4f("model", model)
	shader.SetVector4f("color", color)
	shader.SetVector2f("texOffset", srcPos)
	shader.SetVector2f("texSize", srcSize)
	shader.SetInt("debug", 0)

	gl.BindVertexArray(r.quadVAO)
	texture.Bind()

	gl.DrawElementsWithOffset(gl.TRIANGLES, 6, gl.UNSIGNED_INT, 0)

	gl.BindVertexArray(0)
	texture.UnBind()
	shader.Unuse()
}

func (r *SpriteRendererV2) DebugDraw(shader *Shader, x, y float32, sizeX, sizeY float32, color mgl32.Vec4) {
	translate := mgl32.Translate3D(x, y, 1)
	scale := mgl32.Scale3D(sizeX, sizeY, 1)
	model := translate.Mul4(scale)

	shader.Use()
	shader.SetMatrix4f("model", model)
	shader.SetInt("debug", 1)
	shader.SetVector4f("color", color)

	gl.BindVertexArray(r.quadVAO)
	gl.DrawElementsWithOffset(gl.TRIANGLES, 6, gl.UNSIGNED_INT, 0)

	gl.BindVertexArray(0)
	shader.Unuse()
}

type SpriteRenderer struct {
	quadVAO uint32
}

func CreateSpriteRenderer() SpriteRenderer {
	sr := SpriteRenderer{}

	vertices := []float32{
		0, 0, 0, 0,
		1, 0, 1, 0,
		0, 1, 0, 1,
		1, 1, 1, 1,
	}

	elements := []uint32{
		0, 1, 2,
		2, 3, 1,
	}

	gl.CreateVertexArrays(1, &sr.quadVAO)
	gl.BindVertexArray(sr.quadVAO)

	var VBO, EBO uint32
	gl.GenBuffers(1, &VBO)
	gl.GenBuffers(1, &EBO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)

	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(elements), gl.Ptr(elements), gl.STATIC_DRAW)

	gl.VertexAttribPointerWithOffset(0, 2, gl.FLOAT, false, 4*4, 0)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 4*4, 2*4)

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)

	return sr
}

func (r *SpriteRenderer) Draw(
	shader *Shader,
	texture *Texture,
	srcPos,
	srcSize,
	dstPos,
	dstSize mgl32.Vec2,
	rotate float32,
	color mgl32.Vec4,

) {
	translate := mgl32.Translate3D(dstPos[0], dstPos[1], 0)

	rot := mgl32.Translate3D(0.5, 0.5, 0).
		Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(rotate))).
		Mul4(mgl32.Translate3D(-0.5, -0.5, 0))

	scale := mgl32.Scale3D(dstSize[0], dstSize[1], 1)
	model := translate.Mul4(rot).Mul4(scale)

	if srcPos != (mgl32.Vec2{}) {
		srcPos[0] /= float32(texture.TexBound.Dx())
		srcPos[1] /= float32(texture.TexBound.Dy())
	}

	if srcSize != (mgl32.Vec2{}) {
		srcSize[0] /= float32(texture.TexBound.Dx())
		srcSize[1] /= float32(texture.TexBound.Dy())
	}
	shader.Use()
	shader.SetMatrix4f("model", model)
	shader.SetVector4f("color", color)
	shader.SetVector2f("texOffset", srcPos)
	shader.SetVector2f("texSize", srcSize)
	shader.SetInt("debug", 0)

	gl.BindVertexArray(r.quadVAO)
	texture.Bind()

	gl.DrawElementsWithOffset(gl.TRIANGLES, 6, gl.UNSIGNED_INT, 0)

	gl.BindVertexArray(0)
	texture.UnBind()
	shader.Unuse()
}

func (r *SpriteRenderer) DebugDraw(shader *Shader, x, y float32, sizeX, sizeY float32, color mgl32.Vec4) {
	translate := mgl32.Translate3D(x-sizeX/2, y-sizeY/2, 1)
	scale := mgl32.Scale3D(sizeX, sizeY, 1)
	model := translate.Mul4(scale)

	shader.Use()
	shader.SetMatrix4f("model", model)
	shader.SetInt("debug", 1)
	shader.SetVector4f("color", color)

	gl.BindVertexArray(r.quadVAO)
	gl.DrawElementsWithOffset(gl.TRIANGLES, 6, gl.UNSIGNED_INT, 0)

	gl.BindVertexArray(0)
	shader.Unuse()
}
