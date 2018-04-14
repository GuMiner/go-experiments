#version 400 core

in vec3 fs_color;

in vec3 fs_lightNormalVector;
in vec3 fs_lightVector;

out vec4 color;

uniform vec3 ambient = vec3(0.1, 0.1, 0.1);
uniform vec3 diffuseAlbedo = vec3(0.6, 0.5, 0.8);
uniform vec3 shadingColor;

// Pass-thru the color
void main(void)
{
    vec3 N = normalize(fs_lightNormalVector);
    vec3 L = normalize(fs_lightVector);

    vec3 diffuse = max(abs(dot(N, L)), 0.0) * diffuseAlbedo;

    color = vec4(fs_color.xyz * (ambient * 5 + diffuse), 1.0f) * vec4(shadingColor, 1.0f);
}