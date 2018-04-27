#version 400 core

in vec2 fs_pos;

uniform sampler2D overlayImage;
uniform vec2 offset;
uniform vec2 scale;

out vec4 color;

// Render a Julia fractal on the background
void main(void) {
    const float minBounds = 0;
    const float maxBounds = 1;

    vec2 texturePos = fs_pos; // * scale + offset;
    if (texturePos.x < minBounds || texturePos.x > maxBounds || 
        texturePos.y < minBounds || texturePos.y > maxBounds) {
        color = vec4(1, 0, 0, 0.3); // Fully-transparent, green for debugging
    } else {
        color = vec4(texture(overlayImage, texturePos).xyz, 1.0f);
    }
}