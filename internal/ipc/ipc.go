package ipc

import (
	"fmt"
)

type IPC struct {
}

func New() *IPC {
    return &IPC{}
}

func (i IPC) Load() error {
	fmt.Println("IPC module has been loaded")
	return nil
}

func (i IPC) Run() error {
	fmt.Println("IPC module has been run")
	return nil
}
