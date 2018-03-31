#version 400 core

in vec3 fs_color;

out vec4 color;

// Pass-thru the color
void main(void)
{
    color = vec4(fs_color.xyz, 1.0f);
}