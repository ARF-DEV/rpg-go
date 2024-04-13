package engine

import (
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

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

	shader.SetMatrix4f("model", model)
	shader.SetVector4f("color", color)
	shader.SetVector2f("texOffset", srcPos)
	shader.SetVector2f("texSize", srcSize)

	gl.BindVertexArray(r.quadVAO)
	shader.Use()
	texture.Bind()

	gl.DrawElementsWithOffset(gl.TRIANGLES, 6, gl.UNSIGNED_INT, 0)

	texture.UnBind()
}
