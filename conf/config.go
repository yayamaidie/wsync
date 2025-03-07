package conf

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"wsync/global"

	"gopkg.in/yaml.v2"
)

var GlobalConfig *wsyncConfig

const config_type = "config.yaml"

var config_dir string

type wsyncConfig struct {
	Role     string        `yaml:"role"`
	Sender   *senderConf   `yaml:"sender"`
	Accepter *accepterConf `yaml:"accepter"`
	Puller   *pullerConf   `yaml:"puller"`
}

func init() {
	flag.StringVar(&config_dir, "conf", "./", "")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -conf /path/to/config.yaml(default: ./) \n", os.Args[0])
	}
	flag.Parse()

	bytes, err := global.FileReadBytes(configFile())
	if err != nil {
		flag.Usage()
		log.Printf("open file failed: %s, %s", configFile(), err.Error())
		os.Exit(1)
	}

	GlobalConfig = &wsyncConfig{}

	err = yaml.Unmarshal(bytes, GlobalConfig)
	if err != nil {
		log.Printf("yaml unmarshal failed: %s", err.Error())
		os.Exit(1)
	}
}

func configFile() string {
	configfilepath := path.Join(config_dir, config_type)
	return configfilepath
}
