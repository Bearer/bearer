package util

import (
	"regexp"
	"strings"

	"github.com/bearer/bearer/pkg/report/schema"
)

func StripQuotes(value string) string {
	return strings.Trim(value, `"'`+"`")
}

func ConvertToSimpleType(value string) string {
	simplified := strings.ToLower(value)
	simplified = strings.Split(simplified, " ")[0]

	numberMap := []string{"bit", `tinyint(\(\d?\))?`, "smallint", "mediumint", "int", "integer", "bigint", `float(\(\d?,\d?\))`, `double(\(\d?,\d?\))?`, `decimal(\(\d?,\d?\))`, `dec`}
	for _, typeValue := range numberMap {
		reg := regexp.MustCompile(typeValue)
		if reg.Match([]byte(simplified)) {
			return schema.SimpleTypeNumber
		}
	}

	dateMap := []string{"date", "datetime", "timestamp", "time", "year"}
	for _, typeValue := range dateMap {
		if strings.Contains(simplified, typeValue) {
			return schema.SimpleTypeDate
		}
	}

	stringMap := []string{`char(\(\d?\))?`, `varchar(\(\d?\))?`, `character(\(\d?\))?`, "tinytext", "mediumtext", "longtext"}
	for _, typeValue := range stringMap {
		reg := regexp.MustCompile(typeValue)
		if reg.Match([]byte(simplified)) {
			return schema.SimpleTypeString
		}
	}

	binaryMap := []string{"binary", "varbinary", "tinyblob", "mediumblob", "longblob"}
	for _, typeValue := range binaryMap {
		if strings.Contains(simplified, typeValue) {
			return schema.SimpleTypeBinary
		}
	}

	booleanMap := []string{"bool", "boolean"}
	for _, typeValue := range booleanMap {
		if strings.Contains(simplified, typeValue) {
			return schema.SimpleTypeBool
		}
	}

	return schema.SimpleTypeObject
}
