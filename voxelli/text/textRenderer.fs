#version 400 core

uniform sampler2D fontImage;

uniform vec3 foregroundColor;
uniform vec3 backgroundColor;

in vec2 vs_texPos;

out vec4 color;

void main(void)
{
    vec4 lookupColor = texture2D(fontImage, vs_texPos);
    color = mix(vec4(foregroundColor, 1.0), vec4(backgroundColor, 1.0), lookupColor.x);
}