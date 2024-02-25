package configs

type WebCfg struct {
	NavCfg *NavConfig
}

func NewWebCfg() *WebCfg {
	cfg := &WebCfg{NavCfg: NewNavConfig()}
	return cfg
}
