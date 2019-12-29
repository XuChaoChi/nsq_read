package pqueue

import (
	"container/heap"
)

//队列的item定义
type Item struct {
	//具体的内容
	Value    interface{}
	//优先级的值，nsq中一般是时间戳
	Priority int64
	//索引
	Index    int
}

// this is a priority queue as implemented by a min heap
// ie. the 0th element is the *lowest* value
//这个队列基于最小堆实现
//根元素是优先级最高的元素(创建时间早)
type PriorityQueue []*Item

func New(capacity int) PriorityQueue {
	return make(PriorityQueue, 0, capacity)
}

func (pq PriorityQueue) Len() int {
	return len(pq)
}

//排序时按照优先级降序
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	c := cap(*pq)
	//超过容量时，扩容要原来的2倍
	if n+1 > c {
		npq := make(PriorityQueue, n, c*2)
		copy(npq, *pq)
		*pq = npq
	}
	*pq = (*pq)[0 : n+1]
	item := x.(*Item)
	item.Index = n
	//添加到队列末尾
	(*pq)[n] = item
}

func (pq *PriorityQueue) Pop() interface{} {
	n := len(*pq)
	c := cap(*pq)
	//当实际使用只有现有容量一半并且使用容量大于25的时候缩减到原来的1/2
	if n < (c/2) && c > 25 {
		npq := make(PriorityQueue, n, c/2)
		copy(npq, *pq)
		*pq = npq
	}
	item := (*pq)[n-1]
	item.Index = -1
	//也是从末尾弹出一个元素
	*pq = (*pq)[0 : n-1]
	return item
}

func (pq *PriorityQueue) PeekAndShift(max int64) (*Item, int64) {
	//队列空直接返回
	if pq.Len() == 0 {
		return nil, 0
	}

	item := (*pq)[0]
	//根元素优先值大于传入的时候返回空
	if item.Priority > max {
		return nil, item.Priority - max
	}
	//从堆上移除根元素
	heap.Remove(pq, 0)

	return item, 0
}
