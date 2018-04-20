#version 400 core

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

layout (location = 0) in vec3 position;

// Basic renderer for voxels (no shading)
void main(void)
{
    gl_Position = projection * camera * model * vec4(position, 1.0f);
}
