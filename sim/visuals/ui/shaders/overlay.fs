#version 400 core

in vec2 fs_pos;

uniform sampler2D overlayImage;

out vec4 color;

// Render a Julia fractal on the background
void main(void) {
    const float minBounds = 0;
    const float maxBounds = 1;

    if (fs_pos.x < minBounds || fs_pos.x > maxBounds || 
        fs_pos.y < minBounds || fs_pos.y > maxBounds) {
        // color = vec4(0, 1, 0, 0.3); // Fully-transparent, green for debugging
        discard;
    } else {
        color = vec4(texture(overlayImage, fs_pos).xyz, 1.0f);
    }
}