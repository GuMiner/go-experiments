#version 400 core

layout (location = 0) in vec3 position;
layout (location = 1) in vec2 texPos;

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

out vec2 vs_texPos;

// Perform our position and projection transformations, and pass-through the color / texture data
void main(void)
{
    vs_texPos = texPos;
    gl_Position = projection * camera * model * vec4(position, 1);
}