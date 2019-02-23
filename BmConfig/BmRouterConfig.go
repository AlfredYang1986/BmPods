package BmConfig

import "github.com/alfredyang1986/blackmirror/bmconfighandle"

type BmRouterConfig struct {
	Host string
	Port string
	TmpDir string
}

func (br *BmRouterConfig) GenerateConfig() {
	configPath := "Resources/routerconfig.json"
	profileItems := bmconfig.BMGetConfigMap(configPath)

	br.Host = profileItems["Host"].(string)
	br.Port = profileItems["Port"].(string)
	br.TmpDir = profileItems["TmpDir"].(string)
}