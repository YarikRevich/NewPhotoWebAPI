package utils

import (
	. "NewPhotoWeb/config"
	"strconv"
	"strings"
)

func replaceAllSpaces(o []string) []string {

	var result []string
	for _, value := range o {
		if len(value) != 0 {
			result = append(result, value)
		}
	}
	return result
}

func GetFileExtension(filename string) string {

	splittedName := strings.Split(filename, ".")
	extension := splittedName[len(splittedName)-1]
	return extension
}

func IsAllowed(url string) bool {
	allowed := []string{
		"static",
		"check_auth",
		"sign_up",
		"sign_in",
	}

	su := strings.Split(strings.TrimSuffix(strings.TrimPrefix(url, "/"), "/"), "/")
	if len(su) == 0 {
		return false
	}
	for _, value := range allowed {
		if value == su[0] {
			return true
		}
	}
	return false
}

func GetWishedSize(s string) (uint, uint) {
	splitdogsymbol := strings.Split(s, "@")[1]
	splitxsymbol := strings.Split(splitdogsymbol, "x")

	height, err := strconv.Atoi(splitxsymbol[0])

	if err != nil {
		Logger.Errorln(err.Error())
	}

	width, err := strconv.Atoi(splitxsymbol[0])

	if err != nil {
		Logger.Errorln(err.Error())
	}

	return uint(height), uint(width)
}

func GetCleanTags(dirty string) []string {
	var response []string
	split := strings.Split(dirty, ";")
	for _, value := range split {
		if value != "" {
			response = append(response, value)
		}
	}
	return response
}
