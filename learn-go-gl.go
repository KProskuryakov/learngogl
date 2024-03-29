package main

import (
	"math"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/kproskuryakov/learngogl/shaders"
)

func init() {
	// https://github.com/go-gl/glfw
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

// window dimensions for both glfw and our framebuffers
const width int = 800
const height int = 600

func main() {
	// init glfw for window handling
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	// call terminate when we're done
	defer glfw.Terminate()

	// configure glfw to use OpenGL ver 3.3
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	// Set CORE PROFILE which disables backwards-compatible features that are unneeded
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	// May not be necessary, dunno
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Create our glfw window object
	window, err := glfw.CreateWindow(width, height, "LearnOpenGL", nil, nil)
	if err != nil {
		panic(err)
	}
	// Make it the main context
	window.MakeContextCurrent()

	// Init opengl, otherwise horrible segfaults
	if err := gl.Init(); err != nil {
		panic(err)
	}

	// Set initial viewport size for normalized coordinate calculation
	gl.Viewport(0, 0, int32(width), int32(height))

	// callback when the window is resized
	callback := func(window *glfw.Window, newWidth int, newHeight int) {
		gl.Viewport(0, 0, int32(newWidth), int32(newHeight))
	}
	window.SetFramebufferSizeCallback(callback)

	shader1 := shaders.MakeShader("shaders/vshader1.vx", "shaders/fshader1.fx")
	shader2 := shaders.MakeShader("shaders/vshader1.vx", "shaders/fshader2.fx")

	// init vao and vbo
	var vao, vbo [2]uint32
	gl.GenVertexArrays(2, &vao[0])
	gl.GenBuffers(2, &vbo[0])

	gl.BindVertexArray(vao[0])
	// put vertices in array buffer

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo[0])
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// set vertex attribute pointers
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	// set color attribute pointers
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(vao[1])

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo[1])
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices2)*4, gl.Ptr(vertices2), gl.STATIC_DRAW)

	// set vertex attribute pointers
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	// render loop
	for !window.ShouldClose() {
		// input
		processInput(window)

		// renderArrays
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		timeVal := glfw.GetTime()
		greenVal := float32((math.Sin(timeVal) / 2.0) + 0.5)

		// draw
		shader1.Use()
		gl.BindVertexArray(vao[0])
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		shader2.Use()
		shader2.SetFloat("greenVal", greenVal)
		gl.BindVertexArray(vao[1])
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		gl.BindVertexArray(0)

		// double buffer system
		window.SwapBuffers()
		// check for input events
		glfw.PollEvents()
	}
}

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
}

var vertices = []float32{
	// pos         colors
	0.5, 0.5, 0.0, 1.0, 0.0, 0.0, // top right
	0.5, -0.5, 0.0, 0.0, 1.0, 0.0, // bottom right
	-0.5, 0.5, 0.0, 0.0, 0.0, 1.0, // top left
}

var vertices2 = []float32{
	0.5, -0.5, 0.0, 1.0, 0.0, 0.0, // bottom right
	-0.5, -0.5, 0.0, 0.0, 1.0, 0.0, // bottom left
	-0.5, 0.5, 0.0, 0.0, 0.0, 1.0, // top left
}
