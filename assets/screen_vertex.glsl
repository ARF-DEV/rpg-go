#version 330 core

layout (location = 0) in vec2 aPos;
layout (location = 1) in vec2 aTex;

out vec2 TexCord;

void main() {
    gl_Position = vec4(aPos, 0, 1);
    TexCord = aTex;
}