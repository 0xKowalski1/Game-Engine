#version 330 core
out vec4 FragColor;

in vec2 TexCoord;
in vec3 FragPos;  
in vec3 Normal;

// Texture sampler
uniform sampler2D texture1;  

// Ambient Light Uniforms
uniform vec3 ambientLightColor;
uniform float ambientLightIntensity;

// Directional light uniforms
uniform vec3 directionalLightDirection;
uniform vec3 directionalLightColor;
uniform float directionalLightIntensity;

void main() {
    // Calculate ambient light
    vec3 ambient = ambientLightColor * ambientLightIntensity;

    // Directional lighting
    vec3 norm = normalize(Normal);
    vec3 lightDir = normalize(-directionalLightDirection);  // Light direction is reversed
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = diff * directionalLightColor * directionalLightIntensity;

    // Combine the lighting components
    vec3 result = ambient + diffuse;

    FragColor = vec4(result, 1.0) * texture(texture1, TexCoord);
}

