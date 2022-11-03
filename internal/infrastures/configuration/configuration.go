package configuration

import (
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

var GlobalConfig *Configuration

type Configuration struct {
	Env           string `yaml:"env"`
	AdvertiseAddr string `yaml:"advertise_addr"`
	GRPCBind      string `yaml:"grpc_bind"`
	HTTPBind      string `yaml:"http_bind"`
	Auth          struct {
		Key string `yaml:"key"`
		TTL int    `yaml:"ttl"`
	} `yaml:"auth"`
}

func init() {
	Reload()
}

func Reload() {
	data, err := os.ReadFile("configs/config.yml")
	if err != nil {
		log.Panic().Msgf("%v", err)
		return
	}

	if err := yaml.Unmarshal(data, &GlobalConfig); err != nil {
		log.Panic().Msgf("%v", err)
		return
	}

	log.Info().Msgf("%v", GlobalConfig)
}
