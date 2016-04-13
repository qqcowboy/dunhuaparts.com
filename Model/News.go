package Model

import (
	"fmt"
)

type News struct {
	UserName string
	Password string
}

func init() {
	fmt.Println("modules News init")
}
