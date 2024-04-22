package graphics

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.3-core/gl"
)

type ShaderProgram struct {
	ID uint32
}

func InitShaderProgram(vertexPath, fragmentPath string) (*ShaderProgram, error) {
	vertexShaderSource, err := os.ReadFile(vertexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read vertex shader: %w", err)
	}

	fragmentShaderSource, err := os.ReadFile(fragmentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read fragment shader: %w", err)
	}

	vertexShader, err := compileShader(string(vertexShaderSource)+"\x00", gl.VERTEX_SHADER)
	if err != nil {
		return nil, fmt.Errorf("failed to compile vertex shader: %w", err)
	}

	fragmentShader, err := compileShader(string(fragmentShaderSource)+"\x00", gl.FRAGMENT_SHADER)
	if err != nil {
		gl.DeleteShader(vertexShader) // Clean up vertex shader if fragment shader fails to compile
		return nil, fmt.Errorf("failed to compile fragment shader: %w", err)
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	// Check for linking errors
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
		return nil, fmt.Errorf("failed to link program: %s", log)
	}

	gl.DeleteShader(vertexShader)   // Don't need the shader after linking
	gl.DeleteShader(fragmentShader) // Don't need the shader after linking

	return &ShaderProgram{ID: program}, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csource, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csource, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", shaderType, log)
	}

	return shader, nil
}

func (sp *ShaderProgram) Use() {
	gl.UseProgram(sp.ID)
}
