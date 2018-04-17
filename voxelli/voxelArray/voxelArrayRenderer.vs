#version 400 core

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

uniform mat4 partialShadowMatrix;

layout (location = 0) in vec3 position;
layout (location = 1) in vec3 normal;
layout (location = 2) in vec3 color;

out vec3 fs_color;

out vec3 fs_lightNormalVector;
out vec3 fs_lightVector;
out vec4 fs_shadowCoordinate;

// Basic renderer for voxels (no shading)
void main(void)
{
    fs_color = color;

    vec4 viewSpace = model * vec4(position, 1.0f);

    fs_lightNormalVector = mat3(model) * normal;
    fs_lightVector = normalize(vec3(-1, -1, -2));

    fs_shadowCoordinate = partialShadowMatrix * viewSpace;
    gl_Position = projection * camera * viewSpace;
}
