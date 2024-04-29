package char

import (
	"github.com/liuzl/gocc"
	"github.com/rs/zerolog/log"
)

type STConvertCharFilter struct {
	cc *gocc.OpenCC
}

func NewSTConvertCharFilter(conversion string) (*STConvertCharFilter, error) {
	cc, err := gocc.New(conversion)
	if err != nil {
		return nil, err
	}
	return &STConvertCharFilter{cc: cc}, nil
}

func (t STConvertCharFilter) Filter(input []byte) []byte {
	out, err := t.cc.Convert(string(input))
	if err != nil {
		log.Error().Err(err).Msg("stconvert error")
	}
	return []byte(out)
}
