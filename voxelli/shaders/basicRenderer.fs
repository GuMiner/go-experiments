#version 400 core

uniform float runTime;

in vec3 fs_color;

out vec4 color;

// Pass-thru the color
void main(void)
{
	// float blue = float(int((fs_color.z + runTime) * 1000) % 1000) / 1000.0f;
    color = vec4(fs_color.xyz, 1.0f);
}