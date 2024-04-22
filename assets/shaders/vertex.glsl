#version 330 core
layout (location = 0) in vec3 aPos;       // Position attribute
layout (location = 1) in vec3 aColor;     // Color attribute
layout (location = 2) in vec2 aTexCoord;  // Texture coordinate attribute

out vec3 ourColor;                        // Pass color to fragment shader
out vec2 TexCoord;                        // Pass texture coordinates to fragment shader

void main() {
    gl_Position = vec4(aPos, 1.0);
    ourColor = aColor;
    TexCoord = aTexCoord;
}

