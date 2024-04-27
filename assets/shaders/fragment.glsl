#version 330 core
out vec4 FragColor;

in vec2 TexCoord;

uniform sampler2D texture1;  // Texture sampler
uniform vec3 ambientLightColor;
uniform float ambientLightIntensity;

void main() {
    // Calculate ambient light
    vec3 ambient = ambientLightColor * ambientLightIntensity;

    FragColor = vec4(ambient, 1.0) * texture(texture1, TexCoord);
}

