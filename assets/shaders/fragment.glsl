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

// Point light uniform
struct PointLight {
    vec3 position;
    vec3 color;
    float intensity;
    float constant;
    float linear;
    float quadratic;
};
uniform PointLight pointLight;

vec3 calculatePointLight(PointLight light, vec3 fragPos, vec3 normal) {
    vec3 lightDir = normalize(light.position - fragPos);
    float dist = length(light.position - fragPos);
    float attenuation = 1.0 / (light.constant + light.linear * dist + light.quadratic * dist * dist);

    float diff = max(dot(normal, lightDir), 0.0);
    vec3 diffuse = light.color * light.intensity * diff * attenuation;

    return diffuse;
}

void main() {
    // Calculate ambient light
    vec3 ambient = ambientLightColor * ambientLightIntensity;

    // Directional lighting
    vec3 norm = normalize(Normal);
    vec3 lightDir = normalize(-directionalLightDirection);
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuseDir = diff * directionalLightColor * directionalLightIntensity;

    // Calculate point light
    vec3 diffusePoint = calculatePointLight(pointLight, FragPos, norm);

    // Combine the lighting components
    vec3 result = ambient + diffuseDir + diffusePoint;

    FragColor = vec4(result, 1.0) * texture(texture1, TexCoord);
}

