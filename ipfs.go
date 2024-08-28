package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func (a *App) startIPFS() error {
	// Check if IPFS daemon is already running
	if a.ipfsCmd != nil && a.ipfsCmd.Process != nil {
		return nil
	}

	// Start IPFS daemon
	cmd := exec.Command("ipfs", "daemon")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start IPFS daemon: %w", err)
	}

	a.ipfsCmd = cmd
	return nil
}

func (a *App) runSetupScript() error {
	ctx := context.Background()

	// Create the multiaddr
	multiAddr, err := multiaddr.NewMultiaddr("/dnsaddr/ipfs.teia.art/p2p/12D3KooWP84PmvN2ncA2vDCzoea2DGgBsEgxRreiMWpvZdpEgtrq")
	if err != nil {
		return fmt.Errorf("error creating multiaddr: %w", err)
	}

	// Convert multiaddr to peer.AddrInfo
	addrInfo, err := peer.AddrInfoFromP2pAddr(multiAddr)
	if err != nil {
		return fmt.Errorf("error converting multiaddr to peer.AddrInfo: %w", err)
	}

	// Connect to the peer
	for i := 0; i < 5; i++ { // Try 5 times
		err = a.node.Swarm().Connect(ctx, *addrInfo) // <<-- fails to connect
		if err == nil {
			fmt.Println("Successfully connected to peer")
			return nil // Successfully connected
		}
		fmt.Printf("Attempt %d: Error connecting to swarm: %v\n", i+1, err)
		time.Sleep(5 * time.Second) // Wait 5 seconds before retrying
	}

	return fmt.Errorf("failed to connect after 5 attempts")
}

func (a *App) stopIPFS() {
	if a.ipfsCmd != nil && a.ipfsCmd.Process != nil {
		a.ipfsCmd.Process.Kill()
		a.ipfsCmd.Wait()
	}
}
