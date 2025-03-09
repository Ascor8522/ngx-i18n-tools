package main

import (
	"log"
	"strconv"

	. "common"
)

func main() {
	log.SetFlags(0)

	log.Println("================================")
	log.Println("ngx-xlf-xlsx")
	log.Println("================================")
	log.Println("")

	err := steps()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("")
	log.Println("================================")
	log.Println("Done!")
	log.Println("Excel file is at:")
	log.Printf("%s\n", XlsxPath)
	log.Println("================================")
}

// Put the logic in a separate function to simply return in case of an error,
// instead of using `log.Fatal` everywhere.
func steps() error {
	translationManager := TranslationManager{}

	log.Println("[1/6]\tReading Angular project configuration")
	mainProject, err := getMainAngularProject()
	if err != nil {
		return err
	}

	sourceLocale := mainProject.getSourceLocale()
	translationManager.SetSourceLocale(sourceLocale)
	nonSourceLocales := mainProject.getNonSourceLocales()
	for _, locale := range nonSourceLocales {
		translationManager.EnsureLocale(locale)
	}

	log.Println("[2/6]\tReading source xlf file")
	sourceXlfPath := mainProject.getLocalesMap()[sourceLocale]
	sourceXlf, err := getPathXlf(sourceXlfPath)
	if err != nil {
		return err
	}
	sourceStringsMap := sourceXlf.getKeyValues()

	err = translationManager.AddTranslations(sourceStringsMap, sourceLocale)
	if err != nil {
		return err
	}

	log.Println("[3/6]\tEnsuring xlsx file exists")
	xlsxFile := Xlsx{}
	err = xlsxFile.EnsureExists(sourceLocale, nonSourceLocales)
	if err != nil {
		return err
	}

	log.Println("[4/6]\tReading xlsx file")
	xlsxData, err := xlsxFile.GetData()
	if err != nil {
		return err
	}
	xlsxDataGrouped := xlsxData.GroupByLocale()

	for locale, keyValueMap := range xlsxDataGrouped {
		if locale == sourceLocale {
			continue
		}

		if locale != sourceLocale && !translationManager.HasLocale(locale) {
			continue
		}

		log.Printf("\tAdding translations for locale %s\n", strconv.Quote(string(locale)))
		err = translationManager.AddTranslations(keyValueMap, locale)
		if err != nil {
			return err
		}
	}

	log.Println("[5/6]\tWriting to xlsx file")
	err = xlsxFile.Write(translationManager.GetExportableTranslations(), sourceLocale, nonSourceLocales)
	if err != nil {
		return err
	}

	log.Println("[6/6]\tWriting xlf files")
	translationsByLocale := translationManager.GetTranslationsByLocale()
	for _, locale := range translationManager.GetNonSourceLocales() {
		log.Printf("\tWriting xlf file for locale %s\n", strconv.Quote(string(locale)))
		localeXlfPath := mainProject.getLocalesMap()[locale]
		translations := translationsByLocale[locale]
		// Make a copy of the source xlf file.
		// This is now the xlf file for the current locale.
		localeXlf := sourceXlf
		err = localeXlf.write(localeXlfPath, translations)
		if err != nil {
			return err
		}
	}

	return nil
}
