package main

import (
	"fmt"
	"os"

	"github.com/limero/offlinerss/log"
)

func run(args []string) error {
	syncOnlyTo := false
	rollback := false
	for _, arg := range args {
		switch arg {
		case "to":
			syncOnlyTo = true
		case "rollback":
			rollback = true
		case "help", "-h", "--help":
			help()
			return nil
		case "-v":
			log.DebugEnabled = true
		default:
			return fmt.Errorf("unknown argument %q", arg)
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
	totalActions := syncToActions.Total()

	if rollback {
		log.Info("Rolling back %d changes since last sync", totalActions)
		return ReplaceUserDBsWithReferenceDBs(clients)
	}

	server := getServer(config.Server)

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

func help() {
	lines := []string{
		"Usage:",
		"  offlinerss [command] [options]",
		"",
		"Running OfflineRSS without a command will perform a 'to' sync and then fetch new items",
		"",
		"Commands:",
		"  to                  Sync only to the server without fetching new items",
		"  rollback            Discard any changes done to the clients since the last sync",
		"  help, -h, --help    Print this help message",
		"Options:",
		"  -v                  Enable debug logs",
	}

	for _, line := range lines {
		fmt.Println(line)
	}
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Println("Error:", err)
	}
}
