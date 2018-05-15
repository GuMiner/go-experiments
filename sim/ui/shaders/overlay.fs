#version 400 core

in vec2 fs_pos;

uniform sampler2D overlayImage;

out vec4 color;

// Renders the texture directly to the screen, cutting out parts that are offscreen.
void main(void) {
    const float minBounds = 0;
    const float maxBounds = 1;

    if (fs_pos.x < minBounds || fs_pos.x > maxBounds || 
        fs_pos.y < minBounds || fs_pos.y > maxBounds) {
        discard;
    } else {
        color = vec4(texture(overlayImage, fs_pos).xyz, 1.0f);
    }
}