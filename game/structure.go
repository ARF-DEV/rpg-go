package game

import (
	"fmt"

	"github.com/ARF-DEV/rpg-go/engine"
	"github.com/go-gl/mathgl/mgl32"
)

// type Node interface {
// 	GetNeighbors() B
// }

type Queue[T any] struct {
	values []T
}

func (q *Queue[T]) Put(val T) {
	q.values = append(q.values, val)
}

func (q *Queue[T]) Pop() T {
	var def T
	if len(q.values) == 0 {
		return def
	}
	val := q.values[0]
	q.values = q.values[1:]
	return val
}

func (q *Queue[T]) Len() int {
	return len(q.values)
}

type Bfs struct {
	q     Queue[Pos]
	start Pos
	// end        Pos
	time       float64
	visitedMap map[Pos]bool
	started    bool
	pVisited   Pos
}

func CreateBfs(start Pos) Bfs {
	return Bfs{
		q:          Queue[Pos]{},
		start:      start,
		time:       0,
		visitedMap: map[Pos]bool{},
	}
}

func (bfs *Bfs) Update(lvl *Level) {
	if bfs.visitedMap == nil {
		bfs.visitedMap = map[Pos]bool{}
	}

	if bfs.time > 1 {
		fmt.Println(bfs.q.Len())
		bfs.time = 0
		if !bfs.started {
			bfs.q.Put(bfs.start)
			bfs.started = true
		}
		if bfs.q.Len() > 0 {
			curPos := bfs.q.Pop()
			if !bfs.visitedMap[curPos] {
				bfs.visitedMap[curPos] = true
				for _, neighbour := range getNeighbors(curPos, lvl) {
					if !bfs.visitedMap[neighbour] {
						bfs.q.Put(neighbour)
					}
				}
				bfs.pVisited = curPos
			}
		}
	}

	bfs.time += gTimer.DeltaTime
}

func (bfs *Bfs) Draw(renderer engine.Renderer, shader *engine.Shader) {
	fmt.Println(len(bfs.visitedMap))
	for pos, _ := range bfs.visitedMap {
		renderer.DebugDraw(shader, float32(engine.TILE_SIZE*uint32(int32(-cam[0])+pos.X)), float32(engine.TILE_SIZE*uint32(int32(-cam[1])+pos.Y)), float32(engine.TILE_SIZE), float32(engine.TILE_SIZE), mgl32.Vec4{1, 0, 0, 1})
	}
	renderer.DebugDraw(shader, float32(engine.TILE_SIZE*uint32(int32(-cam[0])+bfs.pVisited.X)), float32(engine.TILE_SIZE*uint32(int32(-cam[1])+bfs.pVisited.Y)), float32(engine.TILE_SIZE), float32(engine.TILE_SIZE), engine.COLOR_BLUE)
	fmt.Println(bfs.q.values)
}
