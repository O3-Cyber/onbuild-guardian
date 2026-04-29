package main

import (
	"context"
	"errors"
	"log"
	"os"
	"slices"

	"github.com/goccy/go-yaml"
	"github.com/moby/moby/client"
)

type Config struct {
	Allowed  map[string][]string `yaml:"allowed"`
	AllowAll map[string]bool     `yaml:"allowAll"`
}

var defaultConfig Config

func main() {
	ctx := context.Background()

	args := os.Args
	if len(args) == 1 {
		log.Println("onbuild-guardian")
		log.Print("\n\n")
		log.Println("USAGE: onbuild-guardian <image>")
		os.Exit(0)
	}

	if len(args) < 2 {
		log.Fatal("must provide exactly one input")
	}

	imageName := args[1]

	config, err := os.Open(".onbuild-guardian.yaml")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}

	defer func() {
		_ = config.Close()
	}()

	c := defaultConfig

	if config != nil {
		dec := yaml.NewDecoder(config)
		if err := dec.Decode(&c); err != nil {
			log.Fatal(err)
		}
	}

	if c.AllowAll[imageName] {
		log.Printf("image %s is allow-listed for all onbuild instructions", imageName)
		os.Exit(0)
	}

	cli, err := client.New(client.FromEnv)
	if err != nil {
		log.Fatalf("Failed to create Docker client: %v", err)
	}
	defer cli.Close()

	res, err := cli.ImageInspect(ctx, imageName)
	if err != nil {
		log.Fatalf("Failed to inspect image '%s': %v", imageName, err)
	}

	for _, ob := range res.Config.OnBuild {
		if slices.Contains(c.Allowed[imageName], ob) {
			continue
		}

		log.Fatalf("image %s has disallowed ONBUILD directive: %s", imageName, ob)
	}

	log.Print("no ONBUILD violations found")
	os.Exit(0)
}
