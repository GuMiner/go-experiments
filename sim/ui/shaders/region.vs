#version 400 core

layout (location = 0) in vec2 position;

uniform vec2 offset;
uniform float scale;
uniform float orientation;

// Modify the region based on the given transformations.
void main(void)
{
    // position starts out from -0.5 to 0.5, based on our region type
    float cosOr = cos(orientation);
    float sinOr = sin(orientation);
    vec2 pos = vec2(
        position.x * cosOr - position.y * sinOr, 
        position.y * cosOr + position.x * sinOr);

    pos = offset + (pos / scale);
    gl_Position = vec4(pos, 1.0f, 1.0f);
}
