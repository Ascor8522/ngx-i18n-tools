package common

import (
	"slices"
	"testing"
)

func TestTranslationManager_SetSourceLocale(t *testing.T) {
	translationManager := TranslationManager{}

	translationManager.EnsureLocale("en")
	translationManager.SetSourceLocale("en")

	if translationManager.sourceLocale != "en" {
		t.Error("Expected source locale to be 'en'")
	}
}

func TestTranslationManager_SetSourceLocale_Add(t *testing.T) {
	translationManager := TranslationManager{}

	translationManager.SetSourceLocale("en")

	if translationManager.sourceLocale != "en" {
		t.Error("Expected source locale to be 'en'")
	}
}

func TestTranslationManager_EnsureLocale(t *testing.T) {
	translationManager := TranslationManager{}

	translationManager.EnsureLocale("en")

	if !translationManager.HasLocale("en") {
		t.Error("Expected locale 'en' to be added")
	}
}

func TestTranslationManager_EnsureLocale_AlreadyExists(t *testing.T) {
	translationManager := TranslationManager{}

	translationManager.EnsureLocale("en")
	translationManager.EnsureLocale("en")

	if len(translationManager.locales) != 1 {
		t.Error("Expected locale 'en' to be added only once")
	}
}

func TestTranslationManager_HasLocale(t *testing.T) {
	translationManager := TranslationManager{}

	translationManager.EnsureLocale("en")

	if !translationManager.HasLocale("en") {
		t.Error("Expected locale 'en' to be added")
	}
}

func TestTranslationManager_HasLocale_NotExists(t *testing.T) {
	translationManager := TranslationManager{}

	if translationManager.HasLocale("en") {
		t.Error("Expected locale 'en' to not exist")
	}
}

func TestTranslationManager_AddTranslations_InSourceLocale(t *testing.T) {
	translationManager := TranslationManager{}

	translationManager.SetSourceLocale("en")
	err := translationManager.AddTranslations(KeyValueMap{
		"key1": "value1",
		"key2": "value2",
	}, "en")
	if err != nil {
		t.Error("Expected no error")
	}

	if translationManager.translations["key1"]["en"] != "value1" {
		t.Error("Expected translation to be added")
	}
	if translationManager.translations["key2"]["en"] != "value2" {
		t.Error("Expected translation to be added")
	}
}

func TestTranslationManager_AddTranslations_InNonSourceLocale_Addition(t *testing.T) {
	translationManager := TranslationManager{}

	translationManager.SetSourceLocale("en")
	translationManager.EnsureLocale("fr")
	err := translationManager.AddTranslations(KeyValueMap{
		"key1": "value1",
		"key2": "value2",
	}, "en")
	if err != nil {
		t.Error("Expected no error")
	}
	err = translationManager.AddTranslations(KeyValueMap{
		"key1": "value1",
	}, "fr")
	if err != nil {
		t.Error("Expected no error")
	}

	if translationManager.translations["key1"]["fr"] != "value1" {
		t.Error("Expected translation to be added")
	}
	if translationManager.translations["key2"]["fr"] != "" {
		t.Error("Expected new translation to be added")
	}
}

func TestTranslationManager_AddTranslations_InNonSourceLocale_Existing(t *testing.T) {
	translationManager := TranslationManager{}

	translationManager.SetSourceLocale("en")
	translationManager.EnsureLocale("fr")
	err := translationManager.AddTranslations(KeyValueMap{
		"key1": "value1",
		"key2": "value2",
	}, "en")
	if err != nil {
		t.Error("Expected no error")
	}
	err = translationManager.AddTranslations(KeyValueMap{
		"key1": "value1",
	}, "fr")
	if err != nil {
		t.Error("Expected no error")
	}
	err = translationManager.AddTranslations(KeyValueMap{
		"key2": "value2",
	}, "fr")
	if err != nil {
		t.Error("Expected no error")
	}

	if translationManager.translations["key1"]["fr"] != "value1" {
		t.Error("Expected translation to be added")
	}
	if translationManager.translations["key2"]["fr"] != "value2" {
		t.Error("Expected translation to be added")
	}
}

func TestTranslationManager_AddTranslations_InNonSourceLocale_Deletion(t *testing.T) {
	translationManager := TranslationManager{}

	translationManager.SetSourceLocale("en")
	translationManager.EnsureLocale("fr")
	err := translationManager.AddTranslations(KeyValueMap{
		"key1": "value1",
	}, "en")
	if err != nil {
		t.Error("Expected no error")
	}
	err = translationManager.AddTranslations(KeyValueMap{
		"key1": "value1",
		"key2": "value2",
	}, "fr")
	if err != nil {
		t.Error("Expected no error")
	}

	if translationManager.translations["key1"]["fr"] != "value1" {
		t.Error("Expected translation to be added")
	}
	if translationManager.translations["key2"] == nil {
		t.Error("Expected deleted translation to be added")
	}
	if translationManager.translations["key2"]["fr"] != "value2" {
		t.Error("Expected translation to be added")
	}
	if slices.Contains(translationManager.sourceKeys, "key2") {
		t.Error("Expected deleted translation not to be added to sources")
	}
}
