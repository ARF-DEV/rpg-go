#version 330 core

layout (location = 0) in vec3 aPos;
layout (location = 1) in vec2 aTexCord;

out vec2 TexCord; 

uniform mat4 model;
uniform mat4 projection;
uniform vec2 texOffset;
uniform vec2 texSize;

void main()
{
    gl_Position = projection * model * vec4(aPos.x, aPos.y, aPos.z, 1.0);
    TexCord = aTexCord * texSize + texOffset;
}