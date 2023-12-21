package main

import (
	"github.com/AnthonyNikitin/go-musthave-shortener-tpl/internal/app/runner"
)

func main() {

	err := runner.RunApplication()
	if err != nil {
		panic(err)
	}
}
