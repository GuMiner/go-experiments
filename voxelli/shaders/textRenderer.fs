#version 400 core

uniform sampler2D fontImage;

in vec4 vs_color;
in vec2 vs_texPos;

out vec4 color;

// Render our text with a sharp boundary.
void main(void)
{
    color = vs_color * texture2D(fontImage, vs_texPos);
    
    // Reduce font transparency to 0.
    if (color.a > 0.20)
    {
        color.a = 1.0f;
    }
    else
    {
        color.a = 0.0f;
    }
}