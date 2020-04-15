package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/uenoryo/gcp-env/gcpenv"
)

func main() {
	if err := _main(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func _main() error {
	var (
		ctx     = context.Background()
		conf    = &gcpenv.Config{}
		flagSet = flag.NewFlagSet("gcpenv", flag.ExitOnError)
	)
	flagSet.StringVar(&conf.ProjectName, "project", "", "プロジェクト名を指定してください")
	flagSet.StringVar(&conf.Version, "version", "latest", "バージョンを指定してください")
	flagSet.StringVar(&conf.Prefix, "prefix", "", "プレフィックスを指定してください")
	flagSet.Parse(os.Args[1:])

	if conf.ProjectName == "" {
		return errors.New("プロジェクト名を指定してください")
	}

	env := gcpenv.New(conf)
	if err := env.Fetch(ctx); err != nil {
		return err
	}
	if err := env.Write(os.Stdout); err != nil {
		return err
	}
	return nil
}
