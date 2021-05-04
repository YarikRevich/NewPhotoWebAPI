package utils

import (
	. "NewPhotoWeb/config"
	"strconv"
	"strings"
)

var (
	stuffurls = []string{
		"static",
		"account",
	}
)

func replaceAllSpaces(o []string)[]string{
	
	var result []string
	for _, value := range o{
		if len(value) != 0{
			result = append(result, value)
		}
	}
	return result
}

func GetFileExtension(filename string)(string){

	splittedName := strings.Split(filename, ".")
	extension := splittedName[len(splittedName)-1]
	return extension
}


func ContainsStuffUrls(url string)bool{
	//Checks whether current url is stuff url ...
	
	su := replaceAllSpaces(strings.Split(url, "/"))
	if len(su) == 0{
		return true
	}
	for _, value := range stuffurls{
		if value == su[0]{
			return true
		}
	}	
	return false
}

func GetWishedSize(s string)(uint, uint){
	splitdogsymbol := strings.Split(s, "@")[1]
	splitxsymbol := strings.Split(splitdogsymbol, "x")

	height, err := strconv.Atoi(splitxsymbol[0])

	if err != nil{
		Logger.Errorln(err.Error())
	}

	width, err := strconv.Atoi(splitxsymbol[0])

	if err != nil{
		Logger.Errorln(err.Error())
	}

	return uint(height), uint(width)
}

func GetCleanTags(dirty string)[]string{
	var response []string
	split := strings.Split(dirty, ";")
	for _, value := range split{
		if value != ""{
			response = append(response, value)
		}
	}
	return response
}
