package service

import (
	"os"
	"time"
	"wsync/conf"
	"wsync/global"

	"log"
)

func TimerPullerSync() {
	if _, err := os.Stat(conf.GlobalConfig.Puller.Dir); os.IsNotExist(err) {
		log.Fatalf("puller dir does not exist: %s", conf.GlobalConfig.Sender.Dir)
	}

	timer := time.NewTicker(time.Duration(conf.GlobalConfig.Puller.PullPeriod) * time.Second)
	global.GlobalWsyncer.Wg.Add(1)
	go func() {
		defer global.GlobalWsyncer.Wg.Done()
	BreakForLoop:
		for {
			select {
			case <-timer.C:
				if err := SyncFile(conf.GlobalConfig.Puller.Dir, conf.GlobalConfig.Sender.User, conf.GlobalConfig.Sender.IP, conf.GlobalConfig.Sender.Dir); err != nil {
					log.Printf("puller sync file error: %v", err)
				}
			case <-global.GlobalWsyncer.Cancelctx.Done():
				timer.Stop()
				log.Println("stop puller timer")
				break BreakForLoop
			}
		}
	}()
}
