package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/ipfs/kubo/client/rpc"
)

type App struct {
	node    *rpc.HttpApi
	pinList []string
	ipfsCmd *exec.Cmd
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup() {
	var err error

	fmt.Println("Starting IPFS daemon...")
	err = a.startIPFS()
	if err != nil {
		fmt.Println("Error starting IPFS daemon:", err)
		return
	}

	fmt.Println("Waiting for IPFS to be ready...")
	err = a.waitForIPFS()
	if err != nil {
		fmt.Println("Error waiting for IPFS:", err)
		return
	}

	a.node, err = rpc.NewLocalApi()
	if err != nil {
		fmt.Println("Error connecting to local IPFS node:", err)
		return
	}

	err = a.runSetupScript()
	if err != nil {
		fmt.Println("Error during IPFS setup:", err)
		return
	}

	cid := "bafkreidtuosuw37f5xmn65b3ksdiikajy7pwjjslzj2lxxz2vc4wdy3zku"
	err = a.pinCID(cid)
	if err != nil {
		fmt.Println("Error pinning CID:", err)
		return
	}

	fmt.Println("IPFS daemon is ready!")
	fmt.Println("Loading pinned CIDs...")
	a.loadPinList()
	fmt.Println("Initializing system tray...")
	a.createTray()
}

func (a *App) waitForIPFS() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for IPFS to be ready")
		default:
			cmd := exec.Command("ipfs", "id")
			err := cmd.Run()
			if err == nil {
				return nil // IPFS is ready
			}
			time.Sleep(1 * time.Second)
		}
	}
}
