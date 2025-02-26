package common

import (
	"maps"
	"slices"
)

// A locale, with or without country; e.g. "en-US", "fr", "nl-BE", etc.
type Locale string

// A path to a file; used for xlf files; e.g. "src/locale/messages.de.xlf".
type Path string

// A translation key.
type Key string

// A translation value.
type Value string

type LocalePathMap map[Locale]Path

type LocaleValueMap map[Locale]Value

type KeyValueMap map[Key]Value

type KeyLocaleValueMap map[Key]LocaleValueMap

type LocaleKeyValueMap map[Locale]KeyValueMap

func (l LocalePathMap) GetLocales() []Locale {
	return slices.Collect(maps.Keys(l))
}

func (l LocaleValueMap) GetLocales() []Locale {
	return slices.Collect(maps.Keys(l))
}

// Groups the translations by locale instead of by key.
// Excel provides translations in a key-locale format, but we need them in a locale-key format for the xlf files.
func (k KeyLocaleValueMap) GroupByLocale() LocaleKeyValueMap {
	localeKeyValueMap := LocaleKeyValueMap{}

	for key, localeValueMap := range k {
		for locale, value := range localeValueMap {
			if _, ok := localeKeyValueMap[locale]; !ok {
				localeKeyValueMap[locale] = KeyValueMap{}
			}

			localeKeyValueMap[locale][key] = value
		}
	}

	return localeKeyValueMap
}
