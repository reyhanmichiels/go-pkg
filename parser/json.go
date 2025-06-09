package parser

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/reyhanmichiels/go-pkg/v2/codes"
	"github.com/reyhanmichiels/go-pkg/v2/errors"
	"github.com/reyhanmichiels/go-pkg/v2/log"
)

type jsonConfig string

const (
	defaultConfig jsonConfig = "default"
	fastestConfig jsonConfig = "fastest"
	customConfig  jsonConfig = "custom"
)

type JSONOptions struct {
	Config                        jsonConfig
	IndentionStep                 int
	MarshalFloatWith6Digits       bool
	EscapeHTML                    bool
	SortMapKeys                   bool
	UseNumber                     bool
	DisallowUnknownFields         bool
	TagKey                        string
	OnlyTaggedField               bool
	ValidateJSONRawMessage        bool
	ObjectFieldMustBeSimpleString bool
	CaseSensitive                 bool
}

type JSONInterface interface {
	// Marshal go structs into bytes
	Marshal(orig interface{}) ([]byte, error)

	// Unmarshal bytes into go structs
	Unmarshal(blob []byte, dest interface{}) error
}

type jsonParser struct {
	API    jsoniter.API
	logger log.Interface
}

func initJSON(opt JSONOptions, log log.Interface) JSONInterface {
	var jsonAPI jsoniter.API
	switch opt.Config {
	case defaultConfig:
		jsonAPI = jsoniter.ConfigDefault
	case fastestConfig:
		jsonAPI = jsoniter.ConfigFastest
	case customConfig:
		jsonAPI = jsoniter.Config{
			IndentionStep:                 opt.IndentionStep,
			MarshalFloatWith6Digits:       opt.MarshalFloatWith6Digits,
			EscapeHTML:                    opt.EscapeHTML,
			SortMapKeys:                   opt.SortMapKeys,
			UseNumber:                     opt.UseNumber,
			DisallowUnknownFields:         opt.DisallowUnknownFields,
			TagKey:                        opt.TagKey,
			OnlyTaggedField:               opt.OnlyTaggedField,
			ValidateJsonRawMessage:        opt.ValidateJSONRawMessage,
			ObjectFieldMustBeSimpleString: opt.ObjectFieldMustBeSimpleString,
			CaseSensitive:                 opt.CaseSensitive,
		}.Froze()
	default:
		jsonAPI = jsoniter.ConfigCompatibleWithStandardLibrary
	}

	p := &jsonParser{
		API:    jsonAPI,
		logger: log,
	}

	return p
}

// Marshal input blobs into go structs.
func (p *jsonParser) Marshal(orig interface{}) ([]byte, error) {
	stream := p.API.BorrowStream(nil)
	defer p.API.ReturnStream(stream)
	stream.WriteVal(orig)
	result := make([]byte, stream.Buffered())
	if stream.Error != nil {
		return nil, errors.NewWithCode(codes.CodeJSONMarshalError, stream.Error.Error())
	}
	copy(result, stream.Buffer())
	return result, nil
}

// Unmarshal input blobs into go structs.
func (p *jsonParser) Unmarshal(blob []byte, dest interface{}) error {
	iter := p.API.BorrowIterator(blob)
	defer p.API.ReturnIterator(iter)
	iter.ReadVal(dest)
	if iter.Error != nil {
		return errors.NewWithCode(codes.CodeJSONUnmarshalError, iter.Error.Error())
	}
	return nil
}
