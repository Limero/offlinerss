package main

import (
	"fmt"
	"os"

	"github.com/limero/offlinerss/log"
)

func run(args []string) error {
	syncOnlyTo := false
	for _, arg := range args {
		switch arg {
		case "to":
			syncOnlyTo = true
		case "-v":
			log.DebugEnabled = true
		default:
			log.Warn("Unknown argument %q", arg)
		}
	}

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

	totalActions := syncToActions.Total()
	if totalActions > 0 {
		if err := AuthServer(server); err != nil {
			return err
		}
		if err := SyncToServer(server, syncToActions); err != nil {
			return err
		}
	}

	if syncOnlyTo {
		if totalActions > 0 {
			if err := ReplaceReferenceDBsWithUserDBs(clients); err != nil {
				return err
			}
		}
	} else {
		if totalActions == 0 {
			// No actions were synced to server, so we haven't authenticated yet
			if err := AuthServer(server); err != nil {
				return err
			}
		}
		folders, err := GetNewFromServer(server)
		if err != nil {
			return err
		}

		TransformFolders(folders)

		if err := SyncClients(clients, folders); err != nil {
			return err
		}
	}

	if err := symlinkClientPaths(clients); err != nil {
		return err
	}

	log.Info("Everything synced!")

	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Println("Error:", err)
	}
}
