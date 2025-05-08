---
date: 2024-03-31
title: "A Multi-Platform Modern WebGL Demo"
slug: "a-multi-platform-modern-webgl-demo"
tags: ["TypeScript", "WebGL", "Graphics"]
---

Nearly four years ago I wrote a [blog post](/posts/a-multi-platform-modern-opengl-demo-with-sdl2/) about a native, cross-platform OpenGL demo that I'd written.
That post (which was only my second ever) was actually a response to _another_ [blog post](https://nullprogram.com/blog/2015/06/06/) written by Chris Wellons (aka [null program](https://nullprogram.com/)) on the same topic.
Back then, I spent **hours and hours** just trying to figure out how I could write (and compile) a single C program that would run on all three major platforms: Windows, macOS, and Linux.
Chris used [GLFW3](https://www.glfw.org/) to solve this problem while mine used [SDL2](https://www.libsdl.org/) (both are valid options).
Here is a screenshot of the demo for reference:

<div style="display:flex;justify-content:center">
	<img src="/images/sdl2-opengl-demo.webp" alt="SDL2 OpenGL Demo">
</div>

While not particularly fancy, [the code](https://github.com/theandrew168/sdl2-opengl-demo) serves as a decent foundation for writing this type of program.
I read through it the other day and wondered if I could **rebuild the same application in a way that was easier to share** with others...

## Enter WebGL

Lately, I've been quite interested in [WebGL](https://developer.mozilla.org/en-US/docs/Web/API/WebGL_API) and how it trades top-end performance for the complete elimination of _all_ distribution woes.
I even touched on this topic in a [previous post](/posts/why-write/#webgl-rocks).
In summary: WebGL allows you to build graphical applications (such as games) that are written in JavaScript and are executed directly in the user's browser.
They don't need to download any files, install anything, or tell their computer to trust an untrusted executable (everyone gets suspicous when they hear: "Oh yeah, just ignore the virus warnings!").

These days, simple graphics programs are incredibly easy to create and distribute via the modern browser.
You can throw some JS+WebGL onto an HTML page, host it on a web server, and have it be instantly accessible by anyone with an internet connection.
People from anywhere in the world (and on any operating system) can **simply click on a link and experience the content!**
In fact, just by visiting this page, the new version of the demo has already been delivered to your machine!
Look forward to seeing that at the end of this post.

## Code Structure

This is a pretty bare-bones WebGL demo without too many abstractions in place.
I did write a few helpers for things like resizing the canvas and building shader programs, however.
The demo uses a single shader to draw and rotate a red square.
A single buffer is used to hold all four of the square's vertices (2D points, one for each corner).
The program starts by querying for an [HTML canvas element](https://developer.mozilla.org/en-US/docs/Web/API/Canvas_API) that is included just before the script.
Then, it asks the canvas (and the browser) for a WebGL context and gets to work!

```jsx
<canvas id="glcanvas" width="640" height="640"></canvas>;
<script>
function main() {
  const canvas = document.querySelector("#glcanvas");
  const gl = canvas.getContext("webgl2");

  // load shaders
  // buffer geometry
  // start the draw loop
}

main();
</script>
```

After building the shader and buffering the square's geometry, the program starts the draw loop.
A WebGL-based draw loop utilizes [requestAnimationFrame](https://developer.mozilla.org/en-US/docs/Web/API/window/requestAnimationFrame) to let the browser decide _when_ to draw a frame.
It'll usually match your screen's refresh rate but could change depending on how the browser chooses to allocate its resources.
From the docs:

> The frequency of calls to the callback function will generally match the display refresh rate. The most common refresh rate is 60hz, (60 cycles/frames per second), though 75hz, 120hz, and 144hz are also widely used. requestAnimationFrame() calls are paused in most browsers when running in background tabs or hidden iframe's, in order to improve performance and battery life.

To start the loop, you call this function with a function of your own that is responsible for drawing the scene.
Additionally, your drawing function must be sure to call `requestAnimationFrame` again at the end to queue up the next frame and continue the loop.
Essentially, your draw function and `requestAnimationFrame` form an infinite loop by continuously calling each other.
Pretty neat, huh?
In short:

```js
function draw() {
  // activate our shader and buffer
  // update the square's rotation value
  // clear the canvas
  // draw the square
  // tell the browser to call our draw function again
  requestAnimationFrame(draw);
}
// kick off the loop by telling the browser to call our draw function
requestAnimationFrame(draw);
```

## The Demo

Enough talk, let's see this thing!

<!-- WebGL Demo Code Starts Here! -->
<div style="display:flex;justify-content:center">
	<canvas id="glcanvas" width="640" height="640" style="max-width:100%"></canvas>
</div>

<script>
	function resizeGL(gl) {
		const canvas = gl.canvas;
		const width = canvas.clientWidth;
		const height = canvas.clientHeight;
		if (gl.canvas.width != width || gl.canvas.height != height) {
			gl.canvas.width = width;
			gl.canvas.height = height;
			gl.viewport(0, 0, width, height);
		}
	}

	function compileShader(gl, shader, source) {
		gl.shaderSource(shader, source);

		gl.compileShader(shader);
		if (!gl.getShaderParameter(shader, gl.COMPILE_STATUS)) {
			throw new Error(gl.getShaderInfoLog(shader));
		}
	}

	function linkProgram(gl, program, vertShader, fragShader) {
		gl.attachShader(program, vertShader);
		gl.attachShader(program, fragShader);

		gl.linkProgram(program);
		if (!gl.getProgramParameter(program, gl.LINK_STATUS)) {
			throw new Error(gl.getProgramInfoLog(program));
		}

		gl.detachShader(program, vertShader);
		gl.detachShader(program, fragShader);
	}

	function compileAndLinkShader(gl, vertSource, fragSource) {
		const vertShader = gl.createShader(gl.VERTEX_SHADER);
		compileShader(gl, vertShader, vertSource);

		const fragShader = gl.createShader(gl.FRAGMENT_SHADER);
		compileShader(gl, fragShader, fragSource);

		const program = gl.createProgram();
		linkProgram(gl, program, vertShader, fragShader);

		gl.deleteShader(vertShader);
		gl.deleteShader(fragShader);

		return program
	}

	function main() {
		const canvas = document.querySelector("#glcanvas");

		const gl = canvas.getContext("webgl2");
		if (gl === null) {
			const msg = "Unable to initialize WebGL2. Your browser or machine may not support it.";
			throw new Error(msg);
		}

		console.log("WebGL Vendor:   %s\n", gl.getParameter(gl.VENDOR));
		console.log("WebGL Renderer: %s\n", gl.getParameter(gl.RENDERER));
		console.log("WebGL Version:  %s\n", gl.getParameter(gl.VERSION));
		console.log("GLSL Version:   %s\n", gl.getParameter(gl.SHADING_LANGUAGE_VERSION));

		const vertSource = `
		    #version 300 es

			in vec2 aPosition;

			uniform float uAngle;

			void main() {
				mat2 rotate = mat2(cos(uAngle), -sin(uAngle), sin(uAngle), cos(uAngle));
				gl_Position = vec4(0.75 * rotate * aPosition, 0.0, 1.0);
			}
		`;
		const fragSource = `
		    #version 300 es
			precision highp float;

			out vec4 vColor;

			void main() {
			    vColor = vec4(1, 0.15, 0.15, 1);
			}
		`;
		const shader = compileAndLinkShader(gl, vertSource.trim(), fragSource.trim());

		const uAngle = gl.getUniformLocation(shader, "uAngle");

		const square = [
			-1.0,  1.0,
			-1.0, -1.0,
			 1.0,  1.0,
			 1.0, -1.0,
		];

		const vao = gl.createVertexArray();
		gl.bindVertexArray(vao);

		const vbo = gl.createBuffer();
		gl.bindBuffer(gl.ARRAY_BUFFER, vbo);
		gl.bufferData(gl.ARRAY_BUFFER, new Float32Array(square), gl.STATIC_DRAW);

		gl.enableVertexAttribArray(0);
		gl.vertexAttribPointer(0, 2, gl.FLOAT, false, 8, 0);

		gl.bindVertexArray(null);

		function draw(time) {
			resizeGL(gl);
			gl.clearColor(0.15, 0.15, 0.15, 1.0);
			gl.clear(gl.COLOR_BUFFER_BIT);

			const angle = time / 1000;
			gl.useProgram(shader);
			gl.uniform1f(uAngle, angle);
			gl.bindVertexArray(vao);
			gl.drawArrays(gl.TRIANGLE_STRIP, 0, 4);

			requestAnimationFrame(draw);
		}

		requestAnimationFrame(draw);
	}

	main();
</script>
<!-- WebGL Demo Code Ends Here -->

Isn't it beautiful?
A simple graphics demo was distributed to your machine with the same ease and convenience as reading this blog post.
To reiterate the sentiment from my previous post about WebGL:

> While the upper bounds of performance are certainly lower when using a browser-based, JavaScript-based graphics programming environment, I don't think that it impacts me very much.
> I don't plan on creating anything incredibly large or high fidelity: just basic games and examples.
> If completely solving the problem of distribution means sacrificing a bit of performance, then I'm satisfied.

Thanks for reading!
