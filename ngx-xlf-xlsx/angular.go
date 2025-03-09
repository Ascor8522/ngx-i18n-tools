package main

import (
	"encoding/json"
	"errors"
	"maps"
	"os"

	. "common"
)

const (
	AngularConfigPath      = "./angular.json"
	projectTypeApplication = "application"
	SourceXlfPath          = "./src/locale/messages.xlf"
)

type ConfigFile struct {
	Projects map[string]Project `json:"projects"`
}

type Project struct {
	ProjectType string `json:"projectType"`
	I18n        struct {
		SourceLocale Locale        `json:"sourceLocale"`
		Locales      LocalePathMap `json:"locales"` // This does not include the source locale.
	} `json:"i18n"`
}

func getMainAngularProject() (Project, error) {
	var angularConfig ConfigFile
	err := angularConfig.read()
	if err != nil {
		return Project{}, err
	}

	return angularConfig.getMainProject()
}

func (a *ConfigFile) read() error {
	fileContent, err := os.ReadFile(AngularConfigPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(fileContent, a)
}

// Angular config files can contain multiple projects.
// Those projects can be applications or libraries.
// We are only interested in the main application project.
func (a *ConfigFile) getMainProject() (Project, error) {
	for _, project := range a.Projects {
		if project.ProjectType == projectTypeApplication {
			return project, nil
		}
	}

	return Project{}, errors.New("main angular project not found")
}

func (p *Project) getSourceLocale() Locale {
	return p.I18n.SourceLocale
}

func (p *Project) getNonSourceLocales() []Locale {
	var locales []Locale

	for locale := range p.I18n.Locales {
		if locale == p.I18n.SourceLocale {
			continue
		}
		locales = append(locales, locale)
	}

	return locales
}

// A map with the locales and the path to the corresponding xlf file.
// This also includes the source locale, as opposed to `p.I18n.Locales`.
func (p *Project) getLocalesMap() LocalePathMap {
	m := LocalePathMap{}

	maps.Copy(m, p.I18n.Locales)
	m[p.getSourceLocale()] = SourceXlfPath

	return m
}
