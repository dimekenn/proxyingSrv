package configs

type Configs struct {
	SrvConfig *SrvConfig `json:"srv_config"`
	JsonKeys *JsonKeys `json:"json_keys"`
}

type JsonKeys struct {
	Url string `json:"url"`
	Method string `json:"method"`
	Body string `json:"body"`
	Headers string `json:"headers"`
	Queries string `json:"queries"`
}

type SrvConfig struct {
	Port string `json:"port"`
}

func NewConfig() *Configs {
	return &Configs{
		JsonKeys: &JsonKeys{},
		SrvConfig: &SrvConfig{},
	}
}
