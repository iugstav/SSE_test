package queue

// type PriorityComparer interface {
// 	constraints.Signed | ~string
// }

// type PriorityQueueElement[T comparable, P PriorityComparer] struct {
// 	Val      T
// 	Priority P
// }

// type PriorityQueue[T comparable, P PriorityComparer] [10]*PriorityQueueElement[T, P]

type PriorityQueueElement struct {
	Val      string `json:"value"`
	Priority int    `json:"priority"`
}

type PriorityQueue []*PriorityQueueElement

func (p PriorityQueue) Len() int {
	return len(p)
}

func (p PriorityQueue) Less(i int, j int) bool {
	return p[i].Priority > p[j].Priority
}

func (p PriorityQueue) IsFull() bool {
	return len(p) == 10
}

func (p PriorityQueue) IsEmpty() bool {
	return len(p) == 0
}

func (p *PriorityQueue) Swap(i int, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

func (p *PriorityQueue) Push(x any) {
	if p.IsFull() {
		return
	}

	element := x.(PriorityQueueElement)
	*p = append(*p, &element)
}

func (p *PriorityQueue) Pop() any {
	if p.IsEmpty() {
		var empty any
		return empty
	}

	oldQueue := *p
	oldQueueLen := oldQueue.Len()

	returnVal := oldQueue[oldQueueLen-1]

	*p = oldQueue[0 : oldQueueLen-1]

	return returnVal
}
