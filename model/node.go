package model

type node struct {
	value interface{}
	next  *node
	prev  *node
}
