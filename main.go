package main

import (
	"log"

	"wsync/conf"
	"wsync/global"
	"wsync/service"
	"wsync/webserver"
)

func main() {
	global.NewWsyncer()
	global.GlobalWsyncer.Role = conf.GlobalConfig.Role
	log.Printf("role: %s\n", global.GlobalWsyncer.Role)

	switch conf.GlobalConfig.Role {
	case "sender":
		service.SenderWatch()
	case "accepter":

	case "puller":
		switch conf.GlobalConfig.Puller.PullMethod {
		case "peried":
			service.TimerPullerSync()
		case "web":
			webserver.InitWebServer()
		default:
			goto ActiveTerminate
		}
	default:
		goto ActiveTerminate
	}

	global.SystemSignal()

ActiveTerminate:
	global.ReleaseSources()
	log.Fatalf("invalid id: %[1]v\n", conf.GlobalConfig.Role)
}
