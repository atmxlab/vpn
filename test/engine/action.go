package engine

import "github.com/atmxlab/vpn/test"

type simpleAction struct {
	callback func(a test.App)
}

func (s *simpleAction) Handle(a test.App) {
	s.callback(a)
}

func newSimpleAction(callback func(a test.App)) test.Action {
	return &simpleAction{callback: callback}
}
