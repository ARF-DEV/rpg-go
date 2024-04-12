#version 330 core

out vec4 FragColor;
in vec2 TexCord;

uniform vec4 color;
uniform sampler2D texture0;

void main()
{
    FragColor = color * texture(texture0, TexCord);
}