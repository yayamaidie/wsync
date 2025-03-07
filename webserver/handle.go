package webserver

import (
	"log"
	"net/http"
	"os"
	"wsync/conf"
	"wsync/service"

	"github.com/gin-gonic/gin"
)

func watchHandle(ctx *gin.Context) {
	body := &struct {
		SenderUpdateTime int64 `json:"senderupdatetime"`
	}{}
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if _, err := os.Stat(conf.GlobalConfig.Puller.Dir); os.IsNotExist(err) {
		if err := os.MkdirAll(conf.GlobalConfig.Puller.Dir, 0644); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	if err := service.SyncFile(conf.GlobalConfig.Puller.Dir, conf.GlobalConfig.Sender.User, conf.GlobalConfig.Sender.IP, conf.GlobalConfig.Sender.Dir); err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, nil)
}
