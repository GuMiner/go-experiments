#version 400 core

layout (location = 0) in vec2 position;

uniform float zOrder;
out vec2 fs_pos;

// Pass-Thru to do all rendering in the fragment shader
void main(void)
{
    fs_pos = position / 2 + vec2(0.5, 0.5);
    gl_Position = vec4(position, zOrder, 1.0f);
}
