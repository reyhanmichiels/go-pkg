package configreader

import (
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/reyhanmichiels/go-pkg/files"
	"github.com/spf13/viper"
)

const (
	JSONType string = "json"
)

type Interface interface {
	ReadConfig(cfg interface{})
}

type Options struct {
	ConfigFile string
}

type configReader struct {
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

	return &configReader{
		option: option,
		vp:     vp,
	}
}

func (c *configReader) ReadConfig(cfg interface{}) {
	if files.GetExtension(filepath.Base(c.option.ConfigFile)) == JSONType {
		c.resolveJSONRef()
	}

	decoderConfig := &mapstructure.DecoderConfig{
		Result:           cfg,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToSliceHookFunc(","),
			stringToTimeDurationHookFunc(),
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

// internally modified string to duration parser hooks function to handle empty string
func stringToTimeDurationHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(time.Duration(5)) {
			return data, nil
		}

		// If data is empty string return zero duration
		if data.(string) == "" {
			return time.Duration(0), nil
		}

		// Convert it by parsing
		return time.ParseDuration(data.(string))
	}
}

func (c *configReader) resolveJSONRef() {
	refmap := make(map[string]interface{})
	refregxp := regexp.MustCompile(`^\\$ref:#\\/(.*)$`)
	for _, k := range c.vp.AllKeys() {
		refpath := c.vp.GetString(k)
		if refregxp.MatchString(refpath) {
			v, ok := refmap[refpath]
			if !ok {
				refkey := refregxp.ReplaceAllString(refpath, "$1")
				refkey = strings.ToLower(strings.ReplaceAll(refkey, "/", "."))
				refmap[refpath] = c.vp.Get(refkey)
				c.vp.Set(k, refmap[refpath])
			} else {
				c.vp.Set(k, v)
			}
		}
	}
}
