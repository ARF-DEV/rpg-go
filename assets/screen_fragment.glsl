#version 330 core

in vec2 TexCord;
out vec4 FragColor;

uniform sampler2D offscreenTex;

void main () {
    FragColor = texture(offscreenTex, TexCord);
}