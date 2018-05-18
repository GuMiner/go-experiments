#version 400 core

in vec2 fs_pos;

uniform float terrainSize;
uniform sampler2D overlayImage;

out vec4 color;

// Renders the texture directly to the screen, cutting out parts that are offscreen.
void main(void) {
    const float minBounds = 0;
    const float maxBounds = 1;

    vec2 texelPos = fs_pos * terrainSize + vec2(0.0);
    ivec2 itexelPos = ivec2(int(texelPos.x), int(texelPos.y));
    if (fs_pos.x < minBounds || fs_pos.x > maxBounds || 
        fs_pos.y < minBounds || fs_pos.y > maxBounds) {
        discard;
    } else {
        color = vec4(texelFetch(overlayImage, itexelPos, 0).xyz, 1.0f);
    }
}