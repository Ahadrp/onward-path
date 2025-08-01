package mm

import (
	"fmt"
	"log"
	"onward-path/internal/ipc"
	"onward-path/internal/xui"
)

var (
	IPC *ipc.IPC
	XUI *xui.XUI
)

func Load() error {
	IPC = ipc.New()
	if err := IPC.Load(); err != nil {
		log.Panic("IPC has not been initilized")
	}

	XUI = xui.New()
	if err := XUI.Load(); err != nil {
		log.Panic("XUI has not been initilized")
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

	if XUI == nil {
		log.Panic("XUI has not been initilized")
	}
	if err := XUI.Run(); err != nil {
		log.Panic("Error while running XUI: ", err)
	}

	fmt.Println("All modules have been run")
	return nil
}
