package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"

	. "github.com/Ascor8522/ngx-i18n-tools/common"
	"github.com/fatih/color"
)

const (
	placeholderInValueRegex = `<x[\s\t\n\r]*[\s\S]*?id="([\s\S]*?)"[\s\S]*?(?:\/>|>[\s\t\n\r]*<\/x>)`
	placeholderInTextRegex  = `\$\{\{([\s\S]*?)\}\}`
	placeholderSprintf      = "${{%s}}"
	defaultFilePermissions  = 0600
	unmarshalStringFormat   = "<root>%s</root>"
)

type Xliff struct {
	XMLName struct{} `xml:"xliff"`
	File    File     `xml:"file"`
	Xmlns   string   `xml:"xmlns,attr"`
	Version string   `xml:"version,attr"`
}

type File struct {
	Body struct {
		TransUnits []TransUnit `xml:"trans-unit"`
	} `xml:"body"`
	SourceLanguage Locale `xml:"source-language,attr"`
	DataType       string `xml:"datatype,attr"`
	Original       string `xml:"original,attr"`
}

type TransUnit struct {
	ID           Key            `xml:"id,attr"`
	DataType     string         `xml:"datatype,attr"`
	Source       SourceTarget   `xml:"source"`
	SourceStr    string         `xml:"-"` // Same as Source, but with all the placeholders replaced with their string representation.
	X            []X            `xml:"-"` // Placeholders found in the Source string.
	ContextGroup []ContextGroup `xml:"context-group"`
	Target       SourceTarget   `xml:"target,omitempty"`
}

// Must use a dedicated struct as we cannot use tag `xml:",innerxml"` together with `xml:"source"` or `xml:"target"`.
type SourceTarget struct {
	InnerXML string `xml:",innerxml"`
}

// A placeholder inside a source or target element.
type X struct {
	XMLName   xml.Name `xml:"x"`
	ID        string   `xml:"id,attr"`
	EquivText string   `xml:"equiv-text,attr"`
}

type ContextGroup struct {
	Purpose string `xml:"purpose,attr"`
	Context []struct {
		ContextType string `xml:"context-type,attr"`
		Value       string `xml:",chardata"`
	} `xml:"context"`
}

func getPathXlf(path Path) (Xliff, error) {
	xlf := Xliff{}
	err := xlf.read(path)

	return xlf, err
}

func (x *Xliff) read(path Path) error {
	fileContent, err := os.ReadFile(string(path))
	if err != nil {
		return err
	}

	err = xml.Unmarshal(fileContent, x)
	if err != nil {
		return err
	}

	for i := range x.File.Body.TransUnits {
		// Since we cannot unmarshal mixed content, we cannot unmarshal the placeholders in the source tag.
		// We simply extract it as raw text, then:
		// - extract the placeholders
		// - create a string version of the source string with the placeholders in their string representation.
		// - unescape the source string
		err = x.File.Body.TransUnits[i].fixRead()
		if err != nil {
			return err
		}
	}

	return err
}

func (x *Xliff) getKeyValues() KeyValueMap {
	keyValueMap := KeyValueMap{}

	for _, transUnit := range x.File.Body.TransUnits {
		keyValueMap[transUnit.ID] = Value(transUnit.SourceStr)
	}

	return keyValueMap
}

func (x *Xliff) write(path Path, translations KeyValueMap) error {
	for key, value := range translations {
		index := slices.IndexFunc(x.File.Body.TransUnits, func(transUnit TransUnit) bool { return transUnit.ID == key })
		if index == -1 {
			// If Excel file has additional keys that are not in the source xlf file, we ignore them.
			continue
		}

		err := x.File.Body.TransUnits[index].setTarget(value)
		if err != nil {
			return err
		}
	}

	bytes, err := xml.MarshalIndent(x, "", "  ")
	if err != nil {
		return err
	}

	// The XML header needs to be added manually.
	bytes = append([]byte(xml.Header), bytes...)

	return os.WriteFile(string(path), bytes, defaultFilePermissions)
}

func (tu *TransUnit) fixRead() error {
	// Replace placeholders in the source string with their string representation.
	regex := regexp.MustCompile(placeholderInValueRegex)
	sourceStr := regex.ReplaceAllStringFunc(tu.Source.InnerXML, func(placeholder string) string {
		placeholderMatch := regex.FindStringSubmatch(placeholder)
		if placeholderMatch == nil {
			log.Fatalf("could not find and replace placeholder %s in value %s", strconv.Quote(placeholder), strconv.Quote(tu.Source.InnerXML))
		}

		placeholderId := placeholderMatch[1]

		return fmt.Sprintf(placeholderSprintf, placeholderId)
	})

	// Since the source string was raw XML, we need to unescape it ourself.
	sourceStr, err := unescape(sourceStr)
	if err != nil {
		return err
	}
	tu.SourceStr = sourceStr

	// Extract placeholders from the source string.
	placeholders, err := extractPlaceholdersFromXMLString(tu.Source.InnerXML)
	if err != nil {
		return err
	}
	tu.X = placeholders

	return nil
}

func (tu *TransUnit) setTarget(value Value) error {
	colorGrayString := color.RGB(128, 128, 128).SprintFunc()

	placeholderIDs := extractPlaceholderIDsFromTextString(string(value))
	placeholderCount := len(placeholderIDs)

	// Ensure number of placeholders is the same in source and target.
	if placeholderCount != len(tu.X) {
		log.Printf("%s in %s\n"+
			"\tplaceholder count in translation does not match placeholder count in source string.\n"+
			"\tsource had %s (%v) but translation has %s (%v).\n"+
			"\tsource string was %s",
			color.YellowString("[WARN]"),
			color.CyanString(strconv.Quote(string(tu.ID))),
			color.MagentaString(strconv.Itoa(len(tu.X))), tu.X,
			color.RedString(strconv.Itoa(placeholderCount)), placeholderIDs,
			colorGrayString(strconv.Quote(tu.SourceStr)))
	}

	// Ensure all placeholders in source are present in target.
	for _, x := range tu.X {
		contains := slices.ContainsFunc(placeholderIDs, func(placeHolderId string) bool { return placeHolderId == x.ID })
		if !contains {
			log.Printf("%s in %s\n"+
				"\tplaceholder %s present in source string is missing from translation %s",
				color.YellowString("[WARN]"),
				color.CyanString(strconv.Quote(string(tu.ID))),
				color.MagentaString(strconv.Quote(x.ID)),
				colorGrayString(strconv.Quote(string(value))))
		}
	}

	targetStr := implacePlaceholders(tu, value)
	tu.Target.InnerXML = targetStr

	return nil
}

// Replace, in a string, the string representation of placeholders by a corresponding XML tags.
// Original placeholders are re-used if found, otherwise a new one are created.
func implacePlaceholders(tu *TransUnit, value Value) string {
	regex := regexp.MustCompile(placeholderInTextRegex)
	valueStr := regex.ReplaceAllStringFunc(string(value), func(placeholder string) string {
		result := regex.FindStringSubmatch(placeholder)
		if result == nil {
			log.Fatalf("could not find and replace placeholder %s in value %s", strconv.Quote(placeholder), strconv.Quote(string(value)))
		}

		placeholderId := result[1]

		var placeholderObj X
		index := slices.IndexFunc(tu.X, func(x X) bool { return x.ID == placeholderId })
		if index == -1 {
			log.Printf("%s in %s\n"+
				"\tcould not find corresponding placeholder in source string; created a made up one",
				color.YellowString("[WARN]"),
				color.CyanString(strconv.Quote(string(tu.ID))))
			placeholderObj = X{
				ID:        placeholderId,
				EquivText: placeholderId,
			}
		} else {
			placeholderObj = tu.X[index]
		}

		placeholderStr, err := xml.Marshal(placeholderObj)
		if err != nil {
			log.Fatal(err)
		}

		return string(placeholderStr)
	})

	return valueStr
}

func extractPlaceholdersFromXMLString(str string) ([]X, error) {
	// This struct is only used to unmarshal the placeholders from the raw string in SourceTarget.Value.
	type Tmp struct {
		X []X `xml:"x,omitempty"`
	}

	var tmp Tmp
	err := xml.Unmarshal([]byte(fmt.Sprintf("<root>%s</root>", str)), &tmp)
	if err != nil {
		return nil, err
	}

	return tmp.X, nil
}

// Extract the ID of the placeholders from a string that might contain some.
func extractPlaceholderIDsFromTextString(str string) []string {
	regex := regexp.MustCompile(placeholderInTextRegex)
	matches := regex.FindAllStringSubmatch(str, -1)
	if matches == nil {
		matches = [][]string{}
	}
	var results []string
	for _, placeholderFoundArr := range matches {
		results = append(results, placeholderFoundArr[1])
	}

	return results
}

func unescape(str string) (string, error) {
	type Tmp struct {
		Text string `xml:",chardata"`
	}
	var tmp Tmp
	err := xml.Unmarshal([]byte(fmt.Sprintf(unmarshalStringFormat, str)), &tmp)
	if err != nil {
		return "", err
	}

	return tmp.Text, nil
}

func (x X) String() string {
	return x.ID
}
