package main

import (
	"fmt"
	"strings"

	"fyne.io/systray"
	"fyne.io/systray/example/icon"
)

func (a *App) createTray() {
	fmt.Println("Creating tray...")
	systray.Run(a.onReady, a.onExit)
}

func (a *App) onReady() {
	systray.SetIcon(icon.Data)
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
