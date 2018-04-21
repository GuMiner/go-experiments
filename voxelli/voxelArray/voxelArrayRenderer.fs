#version 400 core

uniform sampler2DShadow shadowTexture;
uniform vec3 shadingColor;

in vec3 fs_color;

in vec3 fs_lightNormalVector;
in vec3 fs_lightVector;

in vec4 fs_shadowCoordinate;

out vec4 color;

uniform vec3 ambient = vec3(0.1, 0.1, 0.1);
uniform vec3 diffuseAlbedo = vec3(0.6, 0.5, 0.8);

// Helpful data taken from https://github.com/opengl-tutorials/ogl/blob/master/tutorial16_shadowmaps/ShadowMapping.fragmentshader
vec2 poissonDisk[16] = vec2[]( 
   vec2( -0.94201624, -0.39906216 ), 
   vec2( 0.94558609, -0.76890725 ), 
   vec2( -0.094184101, -0.92938870 ), 
   vec2( 0.34495938, 0.29387760 ), 
   vec2( -0.91588581, 0.45771432 ), 
   vec2( -0.81544232, -0.87912464 ), 
   vec2( -0.38277543, 0.27676845 ), 
   vec2( 0.97484398, 0.75648379 ), 
   vec2( 0.44323325, -0.97511554 ), 
   vec2( 0.53742981, -0.47373420 ), 
   vec2( -0.26496911, -0.41893023 ), 
   vec2( 0.79197514, 0.19090188 ), 
   vec2( -0.24188840, 0.99706507 ), 
   vec2( -0.81409955, 0.91437590 ), 
   vec2( 0.19984126, 0.78641367 ), 
   vec2( 0.14383161, -0.14100790 ) 
);

// Returns a random number based on a vec3 and an int.
float random(vec3 seed, int i){
	vec4 seed4 = vec4(seed,i);
	float dot_product = dot(seed4, vec4(12.9898,78.233,45.164,94.673));
	return fract(sin(dot_product) * 43758.5453);
}

// Pass-thru the color
void main(void)
{
    vec3 N = normalize(fs_lightNormalVector);
    vec3 L = normalize(fs_lightVector);

    vec3 diffuse = max(abs(dot(N, L)), 0.0) * diffuseAlbedo;

    const float bias = 0.005;
    const int multisamplingFactor = 16;
    const float minShadowValue = 0.20;
    const float differencePerStep = (1.0 - minShadowValue) / multisamplingFactor;
    const float randomDiskSize = 750.0;
    
    float shadowValue = 1.0;
    for (int i = 0; i < multisamplingFactor; i++){
		shadowValue -= differencePerStep * (1.0 - 
            texture(shadowTexture, vec3(fs_shadowCoordinate.xy + poissonDisk[i]/randomDiskSize, 
                (fs_shadowCoordinate.z - bias)/fs_shadowCoordinate.w)));
    }
    
    color =  shadowValue *
        vec4(fs_color.xyz * (ambient * 5 + diffuse), 1.0f) * vec4(shadingColor, 1.0f);
}