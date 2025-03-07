package service

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	"wsync/conf"
	"wsync/global"

	"github.com/fsnotify/fsnotify"
)

// If the accepter is configured, the sender will perform synchronization; otherwise, 
// the puller will perform synchronization
func SenderWatch() {
	if _, err := os.Stat(conf.GlobalConfig.Sender.Dir); os.IsNotExist(err) {
		log.Fatalf("sender dir does not exist: %s", conf.GlobalConfig.Sender.Dir)
	}

	if err := addWatchTarget(); err != nil {
		log.Fatal(err)
	}

	log.Println("Watching:", conf.GlobalConfig.Sender.Dir)
	go func() {
	BreakForLoop:
		for {
			select {
			case event, isopen := <-global.GlobalWsyncer.Watcher.Events:
				if !isopen {
					break BreakForLoop
				}

				fileExt := path.Ext(event.Name)
				if strings.Contains(fileExt, ".sw") || strings.Contains(fileExt, "~") || strings.Contains(event.Name, "~") {
					continue
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Printf("event name: %[1]v, op: %[2]v\n", event.Name, event.Op)
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Printf("event name: %[1]v, op: %[2]v\n", event.Name, event.Op)
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Printf("event name: %[1]v, op: %[2]v\n", event.Name, event.Op)
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					log.Printf("event name: %[1]v, op: %[2]v\n", event.Name, event.Op)
				}

				global.GlobalWsyncer.SenderUpdateTime = time.Now().Unix()

				if conf.GlobalConfig.Accepter.Dir != "" && conf.GlobalConfig.Accepter.User != "" && conf.GlobalConfig.Accepter.IP != "" {
					if err := SyncFile(conf.GlobalConfig.Sender.Dir, conf.GlobalConfig.Accepter.User, conf.GlobalConfig.Accepter.IP, conf.GlobalConfig.Accepter.Dir); err != nil {
						log.Fatal(err.Error())
					}
					continue
				}

				if err := pullerExecuteSync(); err != nil {
					log.Println(err.Error())
				}

			case err, isopen := <-global.GlobalWsyncer.Watcher.Errors:
				if !isopen {
					break BreakForLoop
				}
				log.Printf("watcher errors: %s\n", err.Error())
			}
		}
	}()
}

// Add watching file and directory
func addWatchTarget() error {
	err := global.GlobalWsyncer.Watcher.Add(strings.TrimRight(conf.GlobalConfig.Sender.Dir, "/"))
	if err != nil {
		return err
	}

	err = filepath.WalkDir(strings.TrimRight(conf.GlobalConfig.Sender.Dir, "/"), func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return global.GlobalWsyncer.Watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func pullerExecuteSync() error {
	url := fmt.Sprintf("https://%[1]v/watch", conf.GlobalConfig.Puller.ListenAddr)
	ishttp, err := global.ServerIsHttp(url)
	if err != nil {
		return err
	}
	if ishttp {
		url = fmt.Sprintf("http://%s", strings.Split(url, "://")[1])
	} else {
		url = fmt.Sprintf("https://%s", strings.Split(url, "://")[1])
	}
	body := &struct {
		SenderUpdateTime int64 `json:"senderupdatetime"`
	}{
		SenderUpdateTime: global.GlobalWsyncer.SenderUpdateTime,
	}
	code, bytes, err := global.HttpRequest("POST", url, body)
	if code != http.StatusOK || err != nil {
		return fmt.Errorf("send watch request error: %[1]v, %[2]s, %[3]v", code, bytes, err)
	}

	return nil
}