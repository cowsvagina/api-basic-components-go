package i18n

import (
	"github.com/BurntSushi/toml"
	i18nLang "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

// Bundle 语言包
type Bundle struct {
	bundle *i18nLang.Bundle
	loader *i18nLang.Bundle
	tagMap map[string]language.Tag
}

func NewBundle(defaultLanguage language.Tag) *Bundle {
	loader := i18nLang.NewBundle(defaultLanguage)
	loader.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	return &Bundle{
		bundle: i18nLang.NewBundle(defaultLanguage),
		loader: loader,
		tagMap: make(map[string]language.Tag),
	}
}

func (b *Bundle) MustLoadFiles(f map[string][]language.Tag) {
	if err := b.LoadFiles(f); err != nil {
		panic(err)
	}
}

func (b *Bundle) LoadFiles(files map[string][]language.Tag) error {
	for filePath, tags := range files {
		f, err := b.loader.LoadMessageFile(filePath)
		if err != nil {
			return errors.WithStack(err)
		}

		for _, t := range tags {
			err := b.bundle.AddMessages(t, f.Messages...)
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}

	return nil
}
