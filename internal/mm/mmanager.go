package mm

import (
	"fmt"
	"log"
	"onward-path/api"
	"onward-path/internal/ipc"
	"onward-path/internal/usr"
	"onward-path/internal/xui"
)

var (
	IPC *ipc.IPC
	XUI *xui.XUI
	USR *usr.USR
	API *api.API
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

	USR = usr.New()
	if err := USR.Load(); err != nil {
		log.Panic("USR has not been initilized")
	}

	API = api.New()
	if err := API.Load(); err != nil {
		log.Panic("API has not been initilized")
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

	if USR == nil {
		log.Panic("USR has not been initilized")
	}
	if err := USR.Run(); err != nil {
		log.Panic("Error while running USR: ", err)
	}

	if API == nil {
		log.Panic("API has not been initilized")
	}
	if err := API.Run(); err != nil {
		log.Panic("Error while running API: ", err)
	}

	fmt.Println("All modules have been run")
	return nil
}
