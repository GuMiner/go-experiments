#version 400 core

layout (location = 0) in vec2 position;

// Lines are purely 2D pass-thru items.
void main(void)
{
    gl_Position = vec4(position, 1.0f, 1.0f);
}
