package resources

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

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