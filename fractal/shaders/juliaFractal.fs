#version 400 core

in vec2 fs_pos;

uniform vec2 c;
uniform float time;
uniform sampler1D fractalGradient;
uniform int maxIterations;

out vec4 color;

// Render a Julia fractal on the background
void main(void)
{
	float dist = 0.50f;
	float speed = 0.01f;
	vec2 point = vec2(cos(time*speed) * dist, sin(time * speed) * dist);
	
    // Make this a smaller quad by increasing our x and y positoins
    vec2 z = vec2(fs_pos.x * 2.5, fs_pos.y * 2.5);
    int iterations = 0;
    float thresholdSqd = 4.0f;
    while (iterations < maxIterations && dot(z, z) < thresholdSqd)
    {
        vec2 zSqd = vec2(z.x * z.x - z.y * z.y, 2.0 * z.x * z.y) + point;
        z = zSqd + c;
        ++iterations;
    }
    
    if (iterations == maxIterations || iterations < 2)
    {
        color = vec4(0, 0, 0, 1);
    }
    else
    {
        color = vec4(texture(fractalGradient, float(iterations) / float(maxIterations)).xyz, 1.0f);
    }

    // Enable this to get a rainbow texture to test the gradient
    // color = vec4(texture(fractalGradient, (float(fs_pos.x) + 1.0) / 2.0).xyz, 1.0f);
}