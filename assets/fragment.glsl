#version 330 core

out vec4 FragColor;
in vec2 TexCord;

uniform vec4 color;
uniform sampler2D texture0;

uniform bool debug = false;

void main()
{
    if (debug) {
        FragColor = color;
    } else {
        FragColor = color * texture(texture0, TexCord);
    }
}