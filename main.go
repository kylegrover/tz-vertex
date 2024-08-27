/* tz-vertez-systray - Tezos NFT IPFS Companion App by Teia

manages an ipfs node and automatically pins NFTs owned or minted by the user

to do:
[ ] add wallet entry ui and store in config
[ ] fetch minted/collected nfts based on wallet entries
[ ] add tray menu to show wallet entries
[ ] add tray menu to pin/unpin wallet entries
[ ] add tray menu to show pinned entries
[ ] add tray menu to open ipfs webui
[ ] move away from exec.Command and use go-ipfs-api

--- ai generated goals: ---
*/

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/systray"
	"fyne.io/systray/example/icon"
	"github.com/ipfs/boxo/path"
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

	// Start IPFS daemon
	err = a.startIPFS()
	if err != nil {
		fmt.Println("Error starting IPFS daemon:", err)
		return
	}

	a.node, err = rpc.NewLocalApi()
	if err != nil {
		fmt.Println("Error connecting to local IPFS node:", err)
		return
	}

	// Wait for IPFS to be ready and run setup
	err = a.waitForIPFSAndSetup()
	if err != nil {
		fmt.Println("Error during IPFS setup:", err)
		return
	}

	a.loadPinList()
	a.createTray()
}

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

func (a *App) waitForIPFSAndSetup() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Wait for IPFS to be ready
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for IPFS to be ready")
		default:
			_, err := a.node.Key().Self(ctx)
			if err == nil {
				// IPFS is ready, proceed with setup
				return a.runSetupScript()
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func (a *App) runSetupScript() error {
	setupCommands := []string{
		"ipfs bootstrap add /dnsaddr/ipfs.teia.art/p2p/12D3KooWP84PmvN2ncA2vDCzoea2DGgBsEgxRreiMWpvZdpEgtrq",
		"ipfs swarm peering add /p2p/12D3KooWP84PmvN2ncA2vDCzoea2DGgBsEgxRreiMWpvZdpEgtrq",
		"ipfs swarm connect /p2p/12D3KooWP84PmvN2ncA2vDCzoea2DGgBsEgxRreiMWpvZdpEgtrq",
	}

	for _, cmd := range setupCommands {
		parts := strings.Fields(cmd)
		command := exec.Command(parts[0], parts[1:]...)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		err := command.Run()
		if err != nil {
			return fmt.Errorf("error running command '%s': %w", cmd, err)
		}
	}

	return nil
}

func (a *App) stopIPFS() {
	if a.ipfsCmd != nil && a.ipfsCmd.Process != nil {
		a.ipfsCmd.Process.Kill()
		a.ipfsCmd.Wait()
	}
}

func (a *App) loadPinList() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	pinListPath := filepath.Join(homeDir, ".tz-vertex_pinlist.json")
	data, err := ioutil.ReadFile(pinListPath)
	if err != nil {
		if os.IsNotExist(err) {
			a.pinList = []string{}
			return
		}
		fmt.Println("Error reading pin list:", err)
		return
	}

	err = json.Unmarshal(data, &a.pinList)
	if err != nil {
		fmt.Println("Error parsing pin list:", err)
		return
	}
}

func (a *App) savePinList() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	pinListPath := filepath.Join(homeDir, ".tz-vertex_pinlist.json")
	data, err := json.Marshal(a.pinList)
	if err != nil {
		fmt.Println("Error encoding pin list:", err)
		return
	}

	err = ioutil.WriteFile(pinListPath, data, 0644)
	if err != nil {
		fmt.Println("Error saving pin list:", err)
		return
	}
}

func (a *App) AddCID(cid string) {
	a.pinList = append(a.pinList, cid)
	a.savePinList()
	go a.pinCID(cid)
}

func (a *App) RemoveCID(cid string) {
	for i, pinnedCID := range a.pinList {
		if pinnedCID == cid {
			a.pinList = append(a.pinList[:i], a.pinList[i+1:]...)
			break
		}
	}
	a.savePinList()
	go a.unpinCID(cid)
}

func (a *App) GetPinList() []string {
	return a.pinList
}

func (a *App) pinCID(cid string) {
	p, err := path.NewPath("/ipfs/" + cid)
	if err != nil {
		fmt.Println("Error creating path:", err)
		return
	}
	err = a.node.Pin().Add(context.Background(), p)
	if err != nil {
		fmt.Println("Error pinning CID:", err)
	}
}

func (a *App) unpinCID(cid string) {
	p, err := path.NewPath("/ipfs/" + cid)
	if err != nil {
		fmt.Println("Error creating path:", err)
		return
	}
	err = a.node.Pin().Rm(context.Background(), p)
	if err != nil {
		fmt.Println("Error unpinning CID:", err)
	}
}

func (a *App) createTray() {
	systray.Run(a.onReady, a.onExit)
}

func (a *App) onReady() {
	systray.SetIcon(icon.Data) // Set your icon data
	systray.SetTitle("TZ Vertex")
	systray.SetTooltip("Tezos NFT IPFS Companion App by Teia")

	mShowPinList := systray.AddMenuItem("Show Pin List", "Display the list of pinned CIDs")
	mAddCID := systray.AddMenuItem("Add CID", "Add a new CID to pin")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	go func() {
		for {
			select {
			case <-mShowPinList.ClickedCh:
				a.showPinList()
			case <-mAddCID.ClickedCh:
				a.showAddCIDDialog()
			case <-mQuit.ClickedCh:
				systray.Quit()
				a.stopIPFS()
				return
			}
		}
	}()
}

func (a *App) onExit() {
	// Perform any cleanup here if needed
	a.stopIPFS()
}

func (a *App) showPinList() {
	pinListStr := "Pinned CIDs:\n" + strings.Join(a.pinList, "\n")
	fmt.Println(pinListStr)
}

func (a *App) showAddCIDDialog() {
	var cid string
	fmt.Print("Enter the CID to pin: ")
	fmt.Scanln(&cid)

	if cid != "" {
		a.AddCID(cid)
	}
}

func main() {
	app := NewApp()
	app.startup()
}
