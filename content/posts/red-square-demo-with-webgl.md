---
date: 2024-03-30
title: "Red Square Demo With WebGL"
slug: "red-square-demo-with-webgl"
tags: ["TypeScript", "WebGL", "Graphics"]
draft: true
---

Link back to my [old post](/a-multi-platform-modern-opengl-demo-with-sdl2) about building a cross-platform OpenGL example with SDL2.
Show how little code it takes to achieve the same thing and how being “cross-platform” isn’t even much of a concern when building web apps.

# Overview

Back in the day, you had do quite a bit of research and structure your code very intentionally to support writing C-based OpenGL graphics for multiple platforms.
I eventually figured it out and wrote a blog post about it (using SDL2 for cross-platform compatibility).
These days, simple graphics programs are even easier to distribute (to all platforms) by means of the modern browser and WebGL.
In fact, by just visiting this page, the demo has been distributed to your machine!
Scroll down and take a look.

You don't need to install anything, download any files, or tell your computer to trust a suspicious executable.
Your browser has enough built-in safety features to trust some JavaScript-based WebGL code.

# Code Structure

This is a pretty bare-bones WebGL demo without too many abstractions in place.
I did write a few helpers for things like: resizing the canvas / viewport, and compiling / linking shader programs.
The demo uses a single shader that rotates a static square and colors it red.
A single buffer is used to hold all 8 of the square's points.
After building the shader and buffering the square's geometry, the program starts the draw loop.

The WebGL-based draw loop typically utilizes [requestAnimationFrame](https://developer.mozilla.org/en-US/docs/Web/API/window/requestAnimationFrame) to let the browser decide _when_ to draw a frame.
It'll usually match your screen's refresh rate but could change depending on how the browser manages its resources.
From the docs:

> The frequency of calls to the callback function will generally match the display refresh rate. The most common refresh rate is 60hz, (60 cycles/frames per second), though 75hz, 120hz, and 144hz are also widely used. requestAnimationFrame() calls are paused in most browsers when running in background tabs or hidden iframe's, in order to improve performance and battery life.

To start the loop, you call this function and pass it your drawing function.
Then, your drawing function must be sure to call `requestAnimationFrame` again at the end to queue up the next frame and continue the loop.
In short:

```js
function draw() {
  // perform any drawing tasks before calling requestAnimationFrame
  requestAnimationFrame(draw);
}
requestAnimationFrame(draw);
```

# Demo

How does inline HTML work?
Can I do a canvas and just embed this demo live?
In order for Hugo to render this, I had to update a setting:

```
[markup.goldmark.renderer]
unsafe = true
```

Otherwise, Hugo would strip out the following HTML / JavaScript code.
I think this is a sane default, though, since many folks _won't_ need to embed arbitrary HTML / JS into their blog posts.

<canvas id="glcanvas" width="480" height="480"></canvas>

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
			gl.clearColor(0.2, 0.3, 0.4, 1.0);
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
