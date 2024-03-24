package define

type Configure struct {
	ServerUrl string `mapstructure:"ServerUrl" json:"ServerUrl" yaml:"ServerUrl"`
	PortList  string `mapstructure:"PortList" json:"PortList" yaml:"PortList"`
}

var (
	File    string
	Url     string
	Port    string
	TimeOut int
	OutPUT  string
)
