package master

import "fmt"

func printFlag (flag, value interface{})string{
	return fmt.Sprintf("--%s=%v",flag, value)
}

