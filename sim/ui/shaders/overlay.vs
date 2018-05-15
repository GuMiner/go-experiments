#version 400 core

layout (location = 0) in vec2 position;

uniform vec2 offset;
uniform vec2 scale;
uniform float zOrder;
out vec2 fs_pos;

// Modify our fragment shader position to be in the right place to read from the shader.
// Pass through the actual 2D position.
void main(void)
{
    // Invert and map from 0 to 1
    fs_pos = vec2(position.x, -position.y);
    fs_pos = fs_pos / 2 + vec2(0.5, 0.5);
    fs_pos = ((fs_pos - offset) / scale);
    gl_Position = vec4(position, zOrder, 1.0f);
}
