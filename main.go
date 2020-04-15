package main

import (
	"context"
	"log"

	"github.com/uenoryo/gcp-env/gcpenv"
)

func main() {
	env := gcpenv.New(&gcpenv.Config{})
	if err := env.Fetch(context.Background()); err != nil {
		log.Println(err.Error())
	}
	log.Println(env.Map())
}
