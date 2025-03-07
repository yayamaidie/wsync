package global

import (
	"context"
	"log"
	"os"
	"os/user"
	"sync"

	"github.com/fsnotify/fsnotify"
)

var GlobalWsyncer *wsyncer

type wsyncer struct {
	Role string

	Watcher *fsnotify.Watcher

	Wg *sync.WaitGroup

	Cancelctx  context.Context
	Cancelfunc context.CancelFunc

	SenderUpdateTime int64

	WorkDir  string
	WorkUser string
}

func NewWsyncer() {
	cancelctx, cancelfunc := context.WithCancel(context.Background())
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	workdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	userinfo, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	GlobalWsyncer = &wsyncer{
		Watcher:    watcher,
		Wg:         &sync.WaitGroup{},
		Cancelctx:  cancelctx,
		Cancelfunc: cancelfunc,
		WorkDir:    workdir,
		WorkUser:   userinfo.Username,
	}
}
