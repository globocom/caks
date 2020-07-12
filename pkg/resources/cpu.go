package resources

import (
	"fmt"
	"strconv"
	"strings"
)

func ConvertIntegerToStringMilicores(value int)string {
	return fmt.Sprintf("%vm",value)
}

func ConvertToIntegerMiliCores(value string)(int, error) {

	miliCores, err := strconv.Atoi(value)

	if err != nil {
		if strings.HasSuffix(value,"m"){
			miliCores, err = strconv.Atoi(strings.Split(value,"m")[0])
			if err != nil {
				return 0, err
			}
		}
	}else{
		miliCores = miliCores * 1000
	}

	return miliCores, nil
}
