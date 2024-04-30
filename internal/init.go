package internal

import (
	"cats-social/common/configs"
	"fmt"
)

func Run() {
	runtimeConfig, err := configs.NewConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n\n", runtimeConfig)
}
