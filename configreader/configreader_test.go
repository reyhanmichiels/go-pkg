package configreader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_configreader_ReadConfig(t *testing.T) {
	type MockConfig struct {
		ValueString string
		ValueInt    int
		ValueBool   bool
	}

	mockRes := MockConfig{
		ValueString: "string",
		ValueInt:    1,
		ValueBool:   true,
	}

	tests := []struct {
		name    string
		args    string
		wantErr bool
		want    MockConfig
	}{
		{
			name:    "success",
			args:    "./files/test.json",
			wantErr: false,
			want:    mockRes,
		},
		{
			name:    "failed decode to struct",
			args:    "./files/failed-test.json",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configReader := Init(Options{
				ConfigFile: tt.args,
			})

			cfg := MockConfig{}
			if tt.wantErr {
				assert.Panics(t, func() {
					configReader.ReadConfig(&cfg)
				})
			} else {
				configReader.ReadConfig(&cfg)
				assert.Equal(t, tt.want, cfg)
			}
		})
	}
}

func Test_configreader_Init(t *testing.T) {
	tests := []struct {
		name    string
		args    Options
		wantErr bool
	}{
		{
			name: "wrong file extension",
			args: Options{
				ConfigFile: "anything.yaml",
			},
			wantErr: true,
		},
		{
			name: "file not exist",
			args: Options{
				ConfigFile: "anything.json",
			},
			wantErr: true,
		},
		{
			name: "success",
			args: Options{
				ConfigFile: "./files/test.json",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				assert.Panics(t, func() {
					Init(tt.args)
				})
			} else {
				config := Init(tt.args)
				assert.NotNil(t, config)
			}
		})
	}
}
