package renderer

import "github.com/brunofjesus/md2pdf/v3/internal/renderer/node"

// states wraps a stack of node.ContainerState pointers.
type states struct {
	stack []*node.ContainerState
}

func (s *states) push(c *node.ContainerState) {
	s.stack = append(s.stack, c)
}

func (s *states) pop() *node.ContainerState {
	var x *node.ContainerState
	x, s.stack = s.stack[len(s.stack)-1], s.stack[:len(s.stack)-1]

	return x
}

func (s *states) peek() *node.ContainerState {
	return s.stack[len(s.stack)-1]
}

func (s *states) parent() *node.ContainerState {
	return s.stack[len(s.stack)-2]
}
