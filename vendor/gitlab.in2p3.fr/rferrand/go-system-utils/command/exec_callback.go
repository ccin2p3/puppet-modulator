package command

import (
	"os/exec"
)

type ExecCallback func(*exec.Cmd)
type ExecCallbacks struct {
	cbs []ExecCallback
}

func NewExecCallbacks() ExecCallbacks {
	return ExecCallbacks{}
}

func (cbs *ExecCallbacks) AddCallback(cb ExecCallback) {
	(*cbs).cbs = append((*cbs).cbs, cb)
}

func (cbs *ExecCallbacks) RunCallbacks(command *exec.Cmd) {
	for _, cb := range (*cbs).cbs {
		go func(ccb ExecCallback) {
			ccb(command)
		}(cb)
	}
}

func (cbs *ExecCallbacks) Empty() bool {
	if len(cbs.cbs) > 0 {
		return false
	}
	return true
}
