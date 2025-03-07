package global

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func SystemSignal() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for s := range stop {
		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			ReleaseSources()
			log.Fatalf("interrupt signal: %s", s.String())
		default:
			log.Printf("unknown signal: %s", s.String())
		}
	}
}

func ReleaseSources() {
	GlobalWsyncer.Watcher.Close()
	log.Println("close watcher")
	GlobalWsyncer.Cancelfunc()
	log.Println("cancel func")
	GlobalWsyncer.Wg.Wait()
	log.Println("wait group done")
}
