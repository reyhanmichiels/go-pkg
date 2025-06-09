package parser

import (
	"github.com/reyhanmichiels/go-pkg/v2/log"
)

type Parser interface {
	JSONParser() JSONInterface
}

type Options struct {
	JSONOptions JSONOptions
}

type parser struct {
	json JSONInterface
}

func InitParser(log log.Interface, opt Options) Parser {
	return &parser{
		json: initJSON(opt.JSONOptions, log),
	}
}

func (p *parser) JSONParser() JSONInterface {
	return p.json
}
