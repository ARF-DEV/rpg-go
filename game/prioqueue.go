package game

import (
	"encoding/json"
	"fmt"
)

type prioQueueLessFunc[T any] func(first, second pqItem[T]) bool

// TODO move to engine
type PrioQueue[T any] struct {
	less prioQueueLessFunc[T]
	q    []pqItem[T]
}

func NewPriorityQueue[T any](less prioQueueLessFunc[T]) PrioQueue[T] {
	return PrioQueue[T]{
		less: less,
	}
}

// Pop
// left, right, isRight, isLeft
// Push
// parent func

func (p *PrioQueue[T]) Push(val T, priorityValue int32) {

	p.q = append(p.q, pqItem[T]{
		Data: val,
		Pv:   priorityValue,
	})
	curIdx := len(p.q) - 1

	// loop until root, if less is true than swap, if not break
	for curIdx != 0 {
		parentIdx, parentItem := p.parent(curIdx)
		if p.less(p.q[curIdx], parentItem) {
			p.q[parentIdx], p.q[curIdx] = p.q[curIdx], p.q[parentIdx]
			curIdx = parentIdx
			continue
		}
		break
	}
}

func (p *PrioQueue[T]) parent(idx int) (int, pqItem[T]) {
	parentIdx := (idx - 1) / 2
	parentItem := p.q[parentIdx]
	return parentIdx, parentItem
}

func (p *PrioQueue[T]) Pop() T {

	if len(p.q) == 0 {
		var result T
		return result
	}

	popVal := p.q[0].Data
	p.q[0] = p.q[len(p.q)-1]
	p.q = p.q[:len(p.q)-1]

	// loop until leaf, if less swap, if not break
	for curIdx := 0; ; {

		if !p.isLeft(curIdx) && !p.isRight(curIdx) {
			break
		}
		fmt.Println(curIdx)
		swapCandidateIdx := curIdx
		// fmt.Println(p.isLeft(curIdx), p.less(p.q[curIdx], p.left(curIdx)), p.q[curIdx].Pv, p.left(curIdx).Pv)
		// fmt.Println(p.isRight(curIdx), p.less(p.q[curIdx], p.right(curIdx), p.q[curIdx].Pv, p.right(curIdx).Pv)
		if p.isLeft(curIdx) && p.less(p.left(curIdx), p.q[curIdx]) {
			if p.isRight(curIdx) && p.less(p.right(curIdx), p.left(curIdx)) {
				swapCandidateIdx = curIdx*2 + 2
			} else {
				swapCandidateIdx = curIdx*2 + 1
			}
		}

		if p.isRight(curIdx) && p.less(p.right(curIdx), p.q[curIdx]) {
			if p.isLeft(curIdx) && p.less(p.left(curIdx), p.right(curIdx)) {
				swapCandidateIdx = curIdx*2 + 1
			} else {
				swapCandidateIdx = curIdx*2 + 2
			}
		}

		if swapCandidateIdx == curIdx {
			break
		}

		p.q[swapCandidateIdx], p.q[curIdx] = p.q[curIdx], p.q[swapCandidateIdx]
		curIdx = swapCandidateIdx
	}

	return popVal
}

func (p *PrioQueue[T]) isRight(idx int) bool {
	rcIdx := idx*2 + 2
	return rcIdx < len(p.q)
}

func (p *PrioQueue[T]) isLeft(idx int) bool {
	lcIdx := idx*2 + 1
	return lcIdx < len(p.q)
}

func (p *PrioQueue[T]) right(idx int) pqItem[T] {
	rcIdx := idx*2 + 2
	return p.q[rcIdx]
}
func (p *PrioQueue[T]) left(idx int) pqItem[T] {
	lcIdx := idx*2 + 1
	return p.q[lcIdx]
}

func (p PrioQueue[T]) String() string {
	b, _ := json.MarshalIndent(p.q, "", "\t")
	return string(b)
}

type pqItem[T any] struct {
	Data T
	Pv   int32 //priority value
}
