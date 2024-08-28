package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ipfs/boxo/path"
)

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

func (a *App) pinCID(cid string) error {
	p, err := path.NewPath("/ipfs/" + cid)
	if err != nil {
		return fmt.Errorf("error creating path: %w", err)
	}

	err = a.node.Pin().Add(context.Background(), p)
	if err != nil {
		return fmt.Errorf("error pinning CID: %w", err)
	}

	return nil
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
