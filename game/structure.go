package game

import (
	"github.com/ARF-DEV/rpg-go/engine"
)

type Traversable interface {
	GetX() int32
	GetY() int32
	comparable
}

type TravTile struct {
	Pos
	Step int32
	Prev *TravTile
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

type TileTraversalFunc[T Traversable] func(trav *TileTraversalViz[T], lvl *Level)
type TileTraversalViz[T Traversable] struct {
	q          PrioQueue[T]
	start      Pos
	end        Pos
	time       float64
	visitedMap map[Pos]bool
	started    bool
	pVisited   T
	uFunc      TileTraversalFunc[T]
	hasReach   bool
	FinalPath  []Pos
}

func CreateTileTravViz[T Traversable](start Pos, end Pos, uFunc func(trav *TileTraversalViz[T], lvl *Level)) TileTraversalViz[T] {
	return TileTraversalViz[T]{
		q:          NewPriorityQueue[T](MinPriority),
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
		renderer.DebugDraw(shader, float32(engine.TILE_SIZE*uint32(int32(-cam[0])+pos.GetX())), float32(engine.TILE_SIZE*uint32(int32(-cam[1])+pos.GetY())), float32(engine.TILE_SIZE), float32(engine.TILE_SIZE), engine.COLOR_RED_FADED)
	}
	renderer.DebugDraw(shader, float32(engine.TILE_SIZE*uint32(int32(-cam[0])+bfs.pVisited.GetX())), float32(engine.TILE_SIZE*uint32(int32(-cam[1])+bfs.pVisited.GetY())), float32(engine.TILE_SIZE), float32(engine.TILE_SIZE), engine.COLOR_BLUE)

	if bfs.hasReach {
		for _, pos := range bfs.FinalPath {
			renderer.DebugDraw(shader, float32(engine.TILE_SIZE*uint32(int32(-cam[0])+pos.GetX())), float32(engine.TILE_SIZE*uint32(int32(-cam[1])+pos.GetY())), float32(engine.TILE_SIZE), float32(engine.TILE_SIZE), engine.COLOR_WHITE_FADED)
		}
	}
	// fmt.Println(bfs.q.values)
}

type PathFindingSearchFunc[T Traversable] func(p *PathFinding[T], lvl *Level, start, goal Pos) []Pos
type PathFinding[T Traversable] struct {
	visitedMap  map[Pos]bool
	searchQueue PrioQueue[T]
	searchFunc  PathFindingSearchFunc[T]
}

func CreatePathFinding[T Traversable](searchFunc PathFindingSearchFunc[T]) PathFinding[T] {
	return PathFinding[T]{
		searchFunc:  searchFunc,
		visitedMap:  map[Pos]bool{},
		searchQueue: NewPriorityQueue[T](MinPriority),
	}
}

func (p *PathFinding[T]) reset() {
	p.visitedMap = map[Pos]bool{}
	p.searchQueue = NewPriorityQueue[T](MinPriority)
}

func (p *PathFinding[T]) FindPath(lvl *Level, start Pos, goal Pos) []Pos {
	p.reset()
	return p.searchFunc(p, lvl, start, goal)
}
