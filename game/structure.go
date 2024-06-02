package game

import (
	"fmt"

	"github.com/ARF-DEV/rpg-go/engine"
	"github.com/go-gl/mathgl/mgl32"
)

type Traversable interface {
	GetX() int32
	GetY() int32
	comparable
}

type PriorityPoint struct {
	Pos
	Value int32
	Step  int32
	Prev  *PriorityPoint
}

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

type TileTraversalViz[T Traversable] struct {
	q          Queue[T]
	start      Pos
	end        Pos
	time       float64
	visitedMap map[Pos]bool
	started    bool
	pVisited   T
	uFunc      func(trav *TileTraversalViz[T], lvl *Level)
	hasReach   bool
	FinalPath  []Pos
}

func CreateTileTravViz[T Traversable](start Pos, end Pos, uFunc func(trav *TileTraversalViz[T], lvl *Level)) TileTraversalViz[T] {
	return TileTraversalViz[T]{
		q:          Queue[T]{},
		start:      start,
		end:        end,
		time:       0,
		visitedMap: map[Pos]bool{},
		uFunc:      uFunc,
		hasReach:   false,
	}
}

func (bfs *TileTraversalViz[T]) Update(lvl *Level) {
	bfs.uFunc(bfs, lvl)
}

func (bfs *TileTraversalViz[T]) Draw(renderer engine.Renderer, shader *engine.Shader) {
	// fmt.Println(len(bfs.visitedMap))
	for pos, _ := range bfs.visitedMap {
		renderer.DebugDraw(shader, float32(engine.TILE_SIZE*uint32(int32(-cam[0])+pos.GetX())), float32(engine.TILE_SIZE*uint32(int32(-cam[1])+pos.GetY())), float32(engine.TILE_SIZE), float32(engine.TILE_SIZE), mgl32.Vec4{1, 0, 0, 1})
	}
	renderer.DebugDraw(shader, float32(engine.TILE_SIZE*uint32(int32(-cam[0])+bfs.pVisited.GetX())), float32(engine.TILE_SIZE*uint32(int32(-cam[1])+bfs.pVisited.GetY())), float32(engine.TILE_SIZE), float32(engine.TILE_SIZE), engine.COLOR_BLUE)

	if bfs.hasReach {
		for _, pos := range bfs.FinalPath {
			renderer.DebugDraw(shader, float32(engine.TILE_SIZE*uint32(int32(-cam[0])+pos.GetX())), float32(engine.TILE_SIZE*uint32(int32(-cam[1])+pos.GetY())), float32(engine.TILE_SIZE), float32(engine.TILE_SIZE), engine.COLOR_GREEN)
		}
	}
	fmt.Println(bfs.q.values)
}
