package common

import (
	"maps"
	"os"
	"slices"
	"sort"
	"unsafe"

	"github.com/xuri/excelize/v2"
)

const (
	basePath         = "translations.xlsx"
	defaultSheetName = "Sheet1"
	keyColumnLabel   = "key"
	columnWidth      = 50
)

type Xlsx struct {
}

func (x *Xlsx) GetData() (KeyLocaleValueMap, error) {
	keyLocaleValueMap := KeyLocaleValueMap{}

	workbook, err := excelize.OpenFile(basePath)
	if err != nil {
		return nil, err
	}
	defer workbook.Close()

	worksheetName := workbook.GetSheetName(0)

	rows, err := workbook.GetRows(worksheetName)
	if err != nil {
		return nil, err
	}

	var locales []string
	for i, row := range rows {
		if i == 0 {
			row = row[1:]
			locales = row

			continue
		}

		key := row[0]
		values := row[1:]

		keyLocaleValueMap[Key(key)] = LocaleValueMap{}

		for j, locale := range locales {
			if j < len(values) {
				keyLocaleValueMap[Key(key)][Locale(locale)] = Value(values[j])
			} else {
				keyLocaleValueMap[Key(key)][Locale(locale)] = defaultTranslationValue
			}
		}
	}

	return keyLocaleValueMap, nil
}

func (x *Xlsx) EnsureExists(sourceLocale Locale, nonSourceLocales []Locale) error {
	_, err := os.Stat(basePath)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}

	workbook := excelize.NewFile()
	defer workbook.Close()

	return x.Write(KeyLocaleValueMap{}, sourceLocale, nonSourceLocales)
}

func (x *Xlsx) Write(translations KeyLocaleValueMap, sourceLocale Locale, nonSourceLocales []Locale) error {
	workbook := excelize.NewFile()
	defer workbook.Close()

	sheetIndex, err := workbook.NewSheet(defaultSheetName)
	if err != nil {
		return err
	}

	worksheetName := workbook.GetSheetName(sheetIndex)

	var locales []Locale
	locales = append(locales, sourceLocale)
	locales = append(locales, nonSourceLocales...)

	// Write header.
	headerLocales := locales
	headerLocales = append([]Locale{keyColumnLabel}, headerLocales...)
	for i, locale := range headerLocales {
		cellAddress, err := excelize.CoordinatesToCellName(i+1, 0+1)
		if err != nil {
			return err
		}
		err = workbook.SetCellValue(worksheetName, cellAddress, string(locale))
		if err != nil {
			return err
		}
	}

	// Write translations.
	translationKeys := slices.Collect(maps.Keys(translations))
	translationKeysStr := *((*[]string)(unsafe.Pointer(&translationKeys)))
	sort.Strings(translationKeysStr)
	translationKeys = *((*[]Key)(unsafe.Pointer(&translationKeysStr)))

	for i, key := range translationKeys {
		cellAddress, err := excelize.CoordinatesToCellName(0+1, i+1+1)
		if err != nil {
			return err
		}
		err = workbook.SetCellValue(worksheetName, cellAddress, string(key))
		if err != nil {
			return err
		}

		for j, locale := range locales {
			cellAddress, err := excelize.CoordinatesToCellName(j+1+1, i+1+1)
			if err != nil {
				return err
			}
			err = workbook.SetCellValue(worksheetName, cellAddress, string(translations[key][locale]))
			if err != nil {
				return err
			}
		}
	}

	// Set column width.
	startCol, err := excelize.ColumnNumberToName(1)
	if err != nil {
		return err
	}
	endCol, err := excelize.ColumnNumberToName(len(headerLocales))
	if err != nil {
		return err
	}
	err = workbook.SetColWidth(worksheetName, startCol, endCol, columnWidth)
	if err != nil {
		return err
	}

	return workbook.SaveAs(basePath)
}
