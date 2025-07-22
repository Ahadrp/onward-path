package mm

import (
	"fmt"
	"log"
	"onward-path/internal/ipc"
)

var (
	IPC *ipc.IPC
)

func Load() error {
	IPC = ipc.New()
	if err := IPC.Load(); err != nil {
		log.Panic("IPC has not been initilized")
	}

	fmt.Println("All modules have been loaded")
	return nil
}

func Run() error {
	if IPC == nil {
		log.Panic("IPC has not been initilized")
	}
	if err := IPC.Run(); err != nil {
		log.Panic("Error while running IPC: ", err)
	}

	fmt.Println("All modules have been run")
	return nil
}
