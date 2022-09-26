package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"task_axxon/internal"
	"task_axxon/internal/configs"
)

//go:embed configs.json
var fs embed.FS

const configName = "configs.json"
func main()  {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//reading json file for configs
	data, readErr := fs.ReadFile(configName)
	if readErr != nil {
		log.Fatal(readErr)
	}

	//creating config entity to deserialize configs.json
	cfg := configs.NewConfig()
	if unmErr := json.Unmarshal(data, &cfg); unmErr != nil {
		log.Fatal(unmErr)
	}

	//channel for error
	errCh := make(chan error, 1)

	//handling signals
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		errCh <- fmt.Errorf("%v", <-sigCh)
		}()

	//new goroutine for REST api server
	go internal.StartHTTPServer(ctx,errCh, cfg)

	log.Fatalf("%v", <-errCh)
}
