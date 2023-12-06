package main

import (
	"fmt"

	"github.com/limero/offlinerss/log"
)

func run() error {
	config, err := getConfig()
	if err != nil {
		return err
	}

	clients := getClients(config.Clients)

	syncToActions, err := GetSyncToActions(clients)
	if err != nil {
		return err
	}

	server := getServer(config.Server)

	folders, err := SyncServer(server, syncToActions)
	if err != nil {
		return err
	}

	if err := SyncClients(clients, folders); err != nil {
		return err
	}

	log.Info("Everything synced!")

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println("Error:", err)
	}
}
