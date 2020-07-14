package i18n

import (
	"github.com/BurntSushi/toml"
	i18nLang "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

type TemplateData map[string]interface{}

// Bundle 语言包
type Bundle struct {
	bundle  *i18nLang.Bundle
	loader  *i18nLang.Bundle
	aliases map[string]language.Tag
}

func NewBundle(defaultLanguage language.Tag) *Bundle {
	loader := i18nLang.NewBundle(defaultLanguage)
	loader.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	return &Bundle{
		bundle:  i18nLang.NewBundle(defaultLanguage),
		loader:  loader,
		aliases: make(map[string]language.Tag),
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

// AddAlias 增加语言别名
// 例如
//		应用中有一个zh-Hans-CN语言包,我们想对它增加一个叫"cn"的别名:
//			b.AddAlias("cn", language.Make("zh-Hans-CN"))
// 		那么在调用Localizer传递"cn"即可,不必再提供全称.
func (b *Bundle) AddAlias(alias string, tag language.Tag) {
	b.aliases[alias] = tag
}

func (b *Bundle) RemoveAlias(alias string) {
	delete(b.aliases, alias)
}

func (b *Bundle) Localizer(langTags ...string) *Localizer {
	var langs []string
	for _, eachTag := range langTags {
		if t, ok := b.aliases[eachTag]; ok {
			langs = append(langs, t.String())
		} else {
			langs = append(langs, eachTag)
		}
	}

	return &Localizer{i18nLang.NewLocalizer(b.bundle, langs...)}
}

type Localizer struct {
	*i18nLang.Localizer
}

func (t *Localizer) SimplyLocalize(messageID string, templateData TemplateData, pluralCount interface{}) (string, error) {
	s, err := t.Localizer.Localize(&i18nLang.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
		PluralCount:  pluralCount,
	})

	return s, errors.WithStack(err)
}
