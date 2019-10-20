package language

import (
	"errors"
	"github.com/abadojack/whatlanggo"
)

var (
	ErrTextEmpty = errors.New("text is empty")
)

func detectLanguage(text string) (string, float64, error) {
	if len(text) < 1 {
		return "", 0, ErrTextEmpty
	}

	info := whatlanggo.DetectWithOptions(text, whatlanggo.Options{
		Whitelist: map[whatlanggo.Lang]bool{
			whatlanggo.Eng: true,
			whatlanggo.Tur: true,
			whatlanggo.Spa: true,
			whatlanggo.Deu: true,
			whatlanggo.Fra: true,
			whatlanggo.Arb: true,
		},
		Blacklist: nil,
	})

	return info.Lang.Iso6391(), info.Confidence, nil
}
