package engine

import "github.com/go-gl/glfw/v3.3/glfw"

type Timer struct {
	latestTime float64
	DeltaTime  float64
}

func (t *Timer) Update() {
	tempTime := t.latestTime
	t.latestTime = glfw.GetTime()

	t.DeltaTime = t.latestTime - tempTime

}
