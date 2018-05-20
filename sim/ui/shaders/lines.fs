#version 400 core

uniform vec3 givenColor;

out vec4 color;

// Color pass-thru
void main(void) {
    color = vec4(givenColor, 1.0f);
}