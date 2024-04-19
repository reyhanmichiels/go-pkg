package configreader

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Interface interface {
	ReadConfig(cfg interface{})
}

type Options struct {
	ConfigFile string
}

type configBuilder struct {
	option Options
	vp     *viper.Viper
}

func Init(option Options) Interface {
	vp := viper.New()

	// check file type
	splitedFileName := strings.Split(option.ConfigFile, ".")
	if splitedFileName[len(splitedFileName)-1] != "json" {
		panic(fmt.Errorf("wrong file type, only support json"))
	}

	vp.SetConfigFile(option.ConfigFile)
	if err := vp.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed read in config: %w", err))
	}

	return &configBuilder{
		option: option,
		vp:     vp,
	}
}

func (c *configBuilder) ReadConfig(cfg interface{}) {
	decoderConfig := &mapstructure.DecoderConfig{
		Result: cfg,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToSliceHookFunc(","),
			mapstructure.StringToTimeDurationHookFunc(),
		),
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		panic(fmt.Errorf("failed create decoder: %w", err))
	}

	err = decoder.Decode(c.vp.AllSettings())
	if err != nil {
		panic(fmt.Errorf("failed decode config to struct: %w", err))
	}
}
