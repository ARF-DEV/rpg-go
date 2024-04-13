package engine

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type shaderID uint32
type programID uint32

type Shader struct {
	ID programID
}

func CreateShader(vertexPath string, fragmentPath string) (Shader, error) {
	vShader, err := loadShader(vertexPath, gl.VERTEX_SHADER)
	if err != nil {
		return Shader{}, fmt.Errorf("error when loading vertex shader: %w", err)
	}

	fShader, err := loadShader(fragmentPath, gl.FRAGMENT_SHADER)
	if err != nil {
		return Shader{}, fmt.Errorf("error when loading fragment shader: %w", err)
	}

	// TODO LINK PROGRAM
	id, err := createProgram(vShader, fShader)
	if err != nil {
		return Shader{}, err
	}

	return Shader{
		ID: id,
	}, nil
}

func (s *Shader) Use() {
	gl.UseProgram(uint32(s.ID))
}

func (s *Shader) SetUniform4f(name string, x, y, z, w float32) {
	uniformLocation := gl.GetUniformLocation(uint32(s.ID), gl.Str(fmt.Sprint(name, "\x00")))
	gl.Uniform4f(uniformLocation, x, y, z, w)
}

func (s *Shader) SetMatrix4f(name string, mat mgl32.Mat4) {
	uniformLocation := gl.GetUniformLocation(uint32(s.ID), gl.Str(fmt.Sprint(name, "\x00")))
	gl.UniformMatrix4fv(uniformLocation, 1, false, &mat[0])
}

func (s *Shader) SetVector4f(name string, vec mgl32.Vec4) {
	uniformLocation := gl.GetUniformLocation(uint32(s.ID), gl.Str(fmt.Sprint(name, "\x00")))
	gl.Uniform4fv(uniformLocation, 1, &vec[0])

}

func (s *Shader) SetVector2f(name string, vec mgl32.Vec2) {
	uniformLocation := gl.GetUniformLocation(uint32(s.ID), gl.Str(fmt.Sprint(name, "\x00")))
	gl.Uniform2fv(uniformLocation, 1, &vec[0])

}

func (s *Shader) SetInt(name string, val int32) {
	uniformLocation := gl.GetUniformLocation(uint32(s.ID), gl.Str(fmt.Sprint(name, "\x00")))
	gl.Uniform1i(uniformLocation, val)

}
func loadShader(path string, shaderType uint32) (shaderID, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("failed to read shader file: %v", err)
	}

	cSrc, free := gl.Strs(string(source))
	defer free()

	shader := gl.CreateShader(shaderType)
	gl.ShaderSource(shader, 1, cSrc, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var infoLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &infoLength)

		logInfo := strings.Repeat("\x00", int(infoLength)+1)
		gl.GetShaderInfoLog(shader, infoLength, nil, gl.Str(logInfo))

		return 0, fmt.Errorf("failed to compile shader: %s", logInfo)
	}

	return shaderID(shader), nil
}

func createProgram(vShader, fShader shaderID) (programID, error) {
	program := gl.CreateProgram()
	gl.AttachShader(program, uint32(vShader))
	gl.AttachShader(program, uint32(fShader))
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		logInfo := strings.Repeat("\x00", int(logLength)+1)
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(logInfo))
		return 0, fmt.Errorf("failed to link program: %s", logInfo)
	}

	return programID(program), nil
}
