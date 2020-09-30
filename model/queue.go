package model

import (
	"sync"
)

type queue struct {
	mutext sync.Mutex
	head   *node
	last   *node
	size   int
}

type Queue interface {
	Enqueue(value interface{})
	Dequeue() interface{}
	Size() int
	IsEmpty() bool
	Front() interface{}
}

func NewQueue() Queue {
	return &queue{sync.Mutex{}, nil, nil, 0}
}

func (q *queue) Enqueue(value interface{}) {
	q.mutext.Lock()
	defer q.mutext.Unlock()
	if q.head == nil {
		q.head = &node{value: value}
		q.last = q.head
	} else {
		q.last.next = &node{value: value}
		q.last.next.prev = q.last
		q.last = q.last.next
	}
	q.size++
}

func (q *queue) Dequeue() interface{} {
	q.mutext.Lock()
	defer q.mutext.Unlock()
	if q.head == nil {
		return nil
	}
	q.size--
	temp := q.head.value
	q.head = q.head.next
	if q.head == nil {
		q.last = q.head
	} else {
		q.head.prev = nil
	}
	return temp
}

func (q *queue) Front() interface{} {
	q.mutext.Lock()
	defer q.mutext.Unlock()
	if q.head == nil {
		return nil
	}
	return q.head.value
}

func (q *queue) IsEmpty() bool { return q.head == nil }

func (q *queue) Size() int { return q.size }
