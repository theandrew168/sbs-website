---
date: 2024-03-30
title: "A Multi-Platform Modern WebGL Demo"
slug: "a-multi-platform-modern-webgl-demo"
tags: ["TypeScript", "WebGL", "Graphics"]
draft: true
---

Nearly four years ago I wrote a [blog post](/posts/a-multi-platform-modern-opengl-demo-with-sdl2/) about a cross-platform, native OpenGL demo that I'd written.
That post (which was only my second ever) was actually a response to _another_ [blog post](https://nullprogram.com/blog/2015/06/06/) written by Chris Wellons (aka [null program](https://nullprogram.com/)) on the same topic.
Chris' post use [GLFW3](https://www.glfw.org/) for platform compatibility while I used [SDL2](https://www.libsdl.org/).
Back then, I spent many hours just trying to figure how I could write (and compile) a single C program that would run on all three major platforms: Windows, macOS, and Linux.

<div style="display:flex;justify-content:center">
	<img src="/images/sdl2-opengl-demo.webp" alt="SDL2 OpenGL Demo">
</div>

# Enter WebGL

Lately, I've been quite interested in [WebGL](https://developer.mozilla.org/en-US/docs/Web/API/WebGL_API) and how it trades top-end performance for the complete elimination of _all_ difficulties related to distribution.
I touched on this topic in one of my [previous posts](/posts/why-write/#webgl-rocks).
As a brief summary, however: WebGL allows you to build graphical applications (and games) that are writte in JavaScript and executed directly in the browser.
You don't need to install anything, download any files, or tell your computer to trust a suspicious executable (even my friends were a bit suspicious when I said "Oh yeah, just ignore the virus warnings").

These days, simple graphics programs are even easier to distribute (to all platforms) by means of the modern browser and WebGL.
You can write your awesome demo / website / game, host it on a web server, and have it been instantly accessible by anyone with an internet connection.
People from anywhere in the world (and on any operating system) can simply enter your URL into their browser and experience the content!
In fact, by just visiting this page, the demo has been distributed to your machine!
Look forward to that at the end of this post.

# Code Structure

This is a pretty bare-bones WebGL demo without too many abstractions in place.
I did write a few helpers for things like resizing the canvas / viewport, and compiling / linking shader programs.
The demo uses a single shader that rotates a static square and colors it red.
A single buffer is used to hold all four of the square's points.
After building the shader and buffering the square's geometry, the program starts the draw loop.

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

The WebGL-based draw loop utilizes [requestAnimationFrame](https://developer.mozilla.org/en-US/docs/Web/API/window/requestAnimationFrame) to let the browser decide _when_ to draw a frame.
It'll usually match your screen's refresh rate but could change depending on how the browser manages its resources.
From the docs:

> The frequency of calls to the callback function will generally match the display refresh rate. The most common refresh rate is 60hz, (60 cycles/frames per second), though 75hz, 120hz, and 144hz are also widely used. requestAnimationFrame() calls are paused in most browsers when running in background tabs or hidden iframe's, in order to improve performance and battery life.

To start the loop, you call this function and pass it your drawing function.
Then, your drawing function must be sure to call `requestAnimationFrame` again at the end to queue up the next frame and continue the loop.
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

# Demo

Since this blog post is written in markdown, I included some inline HTML with the `canvas` and `script` tags.
By default, however, Hugo won't render arbitrary HTML snippets (for security reasons).
I think this is a sane default, though, since many folks _won't_ need to embed arbitrary HTML / JS into their blog posts.
Since I've written this code myself, I'll go ahead and disable this safety feature in my Hugo config file:

```
[markup.goldmark.renderer]
unsafe = true
```

Enough talk, let's see this thing!
Isn't it beautiful: a simple graphics demo distributed to your machine with the same ease and convenience of reading this blog post.
If you want to see the full source code, just view this page's source!
The JavaScript code will be there in its unaltered, hand-written form.
Thanks for reading!

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
