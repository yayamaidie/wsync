package conf

type senderConf struct {
	Dir        string `yaml:"dir"`
	User       string `yaml:"user"`
	IP         string `yaml:"ip"`
}

type accepterConf struct {
	Dir  string `yaml:"dir"`
	User string `yaml:"user"`
	IP   string `yaml:"ip"`
}

type pullerConf struct {
	HTTPS      bool   `yaml:"https"`
	PullMethod string `yaml:"pullmethod"`
	PullPeriod int    `yaml:"pullperiod"`
	Dir        string `yaml:"dir"`
	User       string `yaml:"user"`
	IP         string `yaml:"ip"`
	ListenAddr string `yaml:"addr"`
}
