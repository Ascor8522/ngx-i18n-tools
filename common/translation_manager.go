package common

import (
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
)

const (
	defaultTranslationValue = ""
	multipleSpacesRegex     = `\s{2,}`
)

type TranslationManager struct {
	locales      []Locale
	sourceLocale Locale
	translations KeyLocaleValueMap
	sourceKeys   []Key
}

func (tm *TranslationManager) SetSourceLocale(locale Locale) {
	tm.EnsureLocale(locale)
	tm.sourceLocale = locale
}

func (tm *TranslationManager) EnsureLocale(locale Locale) {
	if tm.HasLocale(locale) {
		return
	}
	tm.locales = append(tm.locales, locale)
	tm.ensureTranslationsInLocale(locale)
}
func (tm *TranslationManager) HasLocale(locale Locale) bool {
	return slices.Contains(tm.locales, locale)
}

func (tm *TranslationManager) ensureTranslationsInLocale(locale Locale) {
	for _, value := range tm.translations {
		if _, ok := value[locale]; !ok {
			value[locale] = defaultTranslationValue
		}
	}
}

func (tm *TranslationManager) AddTranslations(valueMap KeyValueMap, locale Locale) error {
	if tm.sourceLocale == "" {
		// We don't know if it's a source translation or not, since there is no source locale yet.
		// So it could be a deleted translation.
		// We'll return an error just in case.
		return fmt.Errorf("trying to add translations in locale %s but source locale not set yet", strconv.Quote(string(locale)))
	}

	tm.EnsureLocale(locale)

	multipleSpaces := regexp.MustCompile(multipleSpacesRegex)
	for key, value := range valueMap {
		if locale == tm.sourceLocale {
			tm.ensureSourceKey(key)
		}
		tm.ensureTranslationsForKey(key)
		value = Value(strings.Trim(string(value), " "))
		value = Value(multipleSpaces.ReplaceAllString(string(value), " "))
		tm.translations[key][locale] = value
	}

	return nil
}

func (tm *TranslationManager) ensureSourceKey(key Key) {
	if tm.sourceKeys == nil {
		tm.sourceKeys = []Key{}
	}
	if slices.Contains(tm.sourceKeys, key) {
		return
	}
	tm.sourceKeys = append(tm.sourceKeys, key)
}

func (tm *TranslationManager) ensureTranslationsForKey(key Key) {
	if tm.translations == nil {
		tm.translations = make(KeyLocaleValueMap)
	}
	if tm.translations[key] == nil {
		tm.translations[key] = make(LocaleValueMap)
	}
	for _, locale := range tm.locales {
		if _, ok := tm.translations[key][locale]; !ok {
			tm.translations[key][locale] = defaultTranslationValue
		}
	}
}
func (tm *TranslationManager) GetNonSourceLocales() []Locale {
	var locales []Locale

	for _, locale := range tm.locales {
		if locale == tm.sourceLocale {
			continue
		}
		locales = append(locales, locale)
	}

	sort.Slice(locales, func(i, j int) bool {
		return locales[i] < locales[j]
	})

	return locales
}

func (tm *TranslationManager) GetExportableTranslations() KeyLocaleValueMap {
	keyLocaleValueMap := KeyLocaleValueMap{}

	for _, key := range tm.sourceKeys {
		keyLocaleValueMap[key] = tm.translations[key]
	}

	return keyLocaleValueMap
}

func (tm *TranslationManager) GetTranslationsByLocale() LocaleKeyValueMap {
	return tm.translations.GroupByLocale()
}
