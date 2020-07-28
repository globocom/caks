package resources

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ConvertToMebiBytes(value string)(int,error){

	if strings.HasSuffix(value,"Gi"){
		valueInGibiByte, err := convertMemoryValueToInteger(value, "Gi")
		if err != nil {
			return 0,err
		}
		valueInMebiByte := valueInGibiByte * 1024
		return valueInMebiByte, nil

	} else if strings.HasSuffix(value, "Mi"){
		valueInMebiByte, err := convertMemoryValueToInteger(value, "Mi")
		if err != nil {
			return 0, err
		}
		return valueInMebiByte, nil
	}

	return 0, errors.New("memory unit not supported")
}

func ConvertIntegerToStringMebiBytes(value int)string{
	return fmt.Sprintf("%dMi",value)
}

func convertMemoryValueToInteger(value string, unit string)(int, error){
	valueArray := strings.Split(value,unit)
	if len(valueArray) > 0{
		valueInteger, err := strconv.Atoi(valueArray[0])
		if err != nil {
			return 0, err
		}
		return valueInteger, nil
	} else{
		return 0,errors.New(fmt.Sprintf("error in format value %s",value))
	}
}