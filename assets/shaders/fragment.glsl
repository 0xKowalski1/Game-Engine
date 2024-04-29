#version 330 core
out vec4 FragColor;

in vec2 TexCoords;
in vec3 FragPos;  
in vec3 Normal;

// Ambient Light Uniform
struct AmbientLight {
    vec3 color;
    float intensity;
};
uniform AmbientLight ambientLight;

// Directional light uniforms
struct DirectionalLight {
    vec3 direction;
    vec3 color;
    float intensity;
};
uniform DirectionalLight directionalLight;

// Point light uniform
struct PointLight {
    vec3 position;
    vec3 color;
    float intensity;
    float constant;
    float linear;
    float quadratic;
};
#define MAX_POINT_LIGHTS 10
uniform int pointLightsCount;
uniform PointLight pointLights[MAX_POINT_LIGHTS];

// Spot light uniform
struct SpotLight {
    vec3 position;
    vec3 color;
    vec3 direction;
    float cutOff;
    float outerCutOff;
    float intensity;
    float constant;
    float linear;
    float quadratic;
};
uniform SpotLight spotLight;

// Material uniform
struct Material { 
    sampler2D diffuseMap;
    sampler2D specularMap;
    float shininess;
};
uniform Material material;

// View pos
uniform vec3 viewPos;

vec3 calculateAmbientLight(AmbientLight light, Material material) {
    vec3 ambient = light.color * vec3(texture(material.diffuseMap, TexCoords)) * light.intensity;
    
    return ambient;
}

vec3 calculateDirectionalLight(DirectionalLight light, vec3 norm, Material material){ 
    vec3 lightDir = normalize(-light.direction);
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = diff * vec3(texture(material.diffuseMap, TexCoords)) * light.color * light.intensity;

    vec3 reflectDir = reflect(-lightDir, norm);
    vec3 viewDir = normalize(viewPos - FragPos);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shininess);

    vec3 specular = spec * vec3(texture(material.specularMap, TexCoords));

    return diffuse + specular;
}

vec3 calculatePointLight(PointLight light, vec3 fragPos, vec3 normal, Material material) {
    vec3 lightDir = normalize(light.position - fragPos);
    float dist = length(light.position - fragPos);
    float attenuation = 1.0 / (light.constant + light.linear * dist + light.quadratic * (dist * dist));

    float diff = max(dot(normal, lightDir), 0.0);
    vec3 diffuse = diff * vec3(texture(material.diffuseMap, TexCoords)) * light.color * light.intensity * attenuation;

    // Specular 
    vec3 viewDir = normalize(viewPos - FragPos);
    vec3 reflectDir = reflect(-lightDir, normal); 

    float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shininess);

    vec3 specular = spec * vec3(texture(material.specularMap, TexCoords)) * light.color;

    return diffuse + specular;
}

vec3 calculateSpotLight(SpotLight light, vec3 fragPos, vec3 normal, Material material) {
    // diffuse 
    vec3 norm = normalize(Normal);
    vec3 lightDir = normalize(light.position - FragPos);
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = diff * texture(material.diffuseMap, TexCoords).rgb;  
    
    // specular
    vec3 viewDir = normalize(viewPos - FragPos);
    vec3 reflectDir = reflect(-lightDir, norm);  
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), material.shininess);
    vec3 specular = spec * texture(material.specularMap, TexCoords).rgb;  
    
    // spotlight (soft edges)
    float theta = dot(lightDir, normalize(-light.direction)); 
    float epsilon = (light.cutOff - light.outerCutOff);
    float intensity = clamp((theta - light.outerCutOff) / epsilon, 0.0, 1.0);
    diffuse  *= intensity;
    specular *= intensity;
    
    // attenuation
    float distance    = length(light.position - FragPos);
    float attenuation = 1.0 / (light.constant + light.linear * distance + light.quadratic * (distance * distance));    

    diffuse   *= attenuation;
    specular *= attenuation;

    return diffuse + specular;
}


void main() {
    // Calculate ambient light
    vec3 result = calculateAmbientLight(ambientLight, material);

    vec3 norm = normalize(Normal);

    // Directional lighting
    result += calculateDirectionalLight(directionalLight, norm, material);

    // Calculate point lights
    for(int i = 0; i < pointLightsCount; i++)
        result += calculatePointLight(pointLights[i], FragPos, norm, material);

    // Calculate spot light
    result += calculateSpotLight(spotLight, FragPos, norm, material);

    FragColor = vec4(result, 1.0);
}

