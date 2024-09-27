# tz-vertez-systray - Tezos NFT IPFS Companion App by Teia

manages an ipfs node and automatically pins NFTs owned or minted by the user

### dev

build and run the app
```bash
go run main.go
```

build the app
```bash
go build
```


**to do:**
- [ ] add wallet entry ui and store in config
- [ ] fetch minted/collected nfts based on wallet entries
- [ ] add tray menu to show wallet entries
- [ ] add tray menu to pin/unpin wallet entries
- [ ] add tray menu to show pinned entries
- [ ] add tray menu to open ipfs webui
- [ ] move away from exec.Command and use go-ipfs-api



reference:
https://gist.github.com/mattdesl/47f4ea12ea131eed8401bdacf95a1f47
https://nftbiker.xyz/
https://hashquine.github.io/hicetnunc/artists-by-income-3/index.html
https://github.com/hicetnunc2000/hicetnunc/wiki/Tools-made-by-the-community

https://github.com/zir0h/teia-backup - backup tool specifically for teia

---

#### *ai generated overview and todos:*

# TZ Vertex Systray Project Analysis

## Current Functionality

1. Starts and manages an IPFS daemon
2. Connects to a local IPFS node
3. Loads and saves a list of pinned CIDs
4. Provides a system tray interface for basic interactions
5. Allows adding and removing CIDs from the pin list
6. Performs IPFS setup operations (bootstrap, peering, swarm connect)

## Code Structure and Design

The code is structured around a main `App` struct that encapsulates the core functionality. It uses the `fyne.io/systray` library for the system tray interface and `github.com/ipfs/kubo/client/rpc` for IPFS interactions.

### Strengths:
- Clear separation of concerns with methods for different functionalities
- Use of Go's concurrency features (goroutines) for background tasks
- Persistent storage of pinned CIDs

### Areas for Improvement:
1. Error handling: Many errors are just printed to console; consider a more robust error handling strategy
2. Configuration: Hard-coded values (e.g., IPFS peer address) should be moved to a configuration file
3. Use of `exec.Command`: As noted in the TODO, moving to `go-ipfs-api` would be more idiomatic
4. User Interface: Current CLI prompts for adding CIDs could be replaced with GUI dialogs
5. Concurrency: Some operations could benefit from better concurrency control
6. Testing: No tests are present in the current code

## Project Goals (based on TODO list and code analysis)

1. Manage an IPFS node
2. Automatically pin NFTs owned or minted by the user
3. Provide a system tray interface for easy management
4. Integrate with Tezos blockchain for NFT data

## Expanded Project Overview and Steps

1. IPFS Node Management
   - [x] Start and stop IPFS daemon
   - [x] Connect to local IPFS node
   - [x] Perform initial IPFS setup (bootstrap, peering, swarm connect)
   - [ ] Implement proper error handling and recovery for IPFS operations
   - [ ] Add configuration options for IPFS settings

2. Tezos Integration
   - [ ] Implement Tezos API client
   - [ ] Add wallet management functionality
     - [ ] Create UI for adding/removing wallet addresses
     - [ ] Store wallet information securely
   - [ ] Fetch NFT data from Tezos blockchain
     - [ ] Retrieve minted NFTs for each wallet
     - [ ] Retrieve collected NFTs for each wallet
   - [ ] Implement periodic scanning for new NFTs

3. NFT Management
   - [x] Implement basic CID pinning and unpinning
   - [ ] Automate pinning of owned and minted NFTs
   - [ ] Implement smart pinning strategy (e.g., based on storage limits, NFT age)
   - [ ] Add metadata storage for pinned NFTs (e.g., title, artist, collection)

4. User Interface Improvements
   - [x] Basic system tray functionality
   - [ ] Add tray menu to show wallet entries
   - [ ] Add tray menu to pin/unpin wallet entries
   - [ ] Add tray menu to show pinned entries
   - [ ] Add tray menu to open IPFS WebUI
   - [ ] Implement GUI dialogs for user interactions (e.g., adding CIDs, managing wallets)
   - [ ] Add notifications for important events (e.g., new NFT pinned, IPFS node status)

5. Data Management and Storage
   - [ ] Implement a proper database for storing app data (e.g., SQLite)
   - [ ] Create data models for wallets, NFTs, and pin list
   - [ ] Implement data migration strategy for updates

6. Performance and Stability
   - [ ] Replace `exec.Command` usage with `go-ipfs-api`
   - [ ] Implement proper concurrency control for IPFS operations
   - [ ] Add retry mechanisms for failed operations
   - [ ] Implement logging system for better debugging

7. Testing and Quality Assurance
   - [ ] Write unit tests for core functionalities
   - [ ] Implement integration tests for IPFS and Tezos interactions
   - [ ] Set up CI/CD pipeline for automated testing and building

8. Documentation and User Guide
   - [ ] Create detailed README with setup instructions
   - [ ] Write user documentation explaining app features and usage
   - [ ] Document API and core functions for future development

9. Packaging and Distribution
   - [ ] Set up build process for multiple platforms (Windows, macOS, Linux)
   - [ ] Create installers for easy user setup
   - [ ] Implement auto-update mechanism

10. Security Considerations
    - [ ] Implement secure storage for wallet information
    - [ ] Add option for encrypted storage of pinned NFT data
    - [ ] Perform security audit of IPFS and Tezos interactions

This expanded overview provides a more comprehensive roadmap for the project, addressing both the immediate TODO items and long-term considerations for a robust, user-friendly application.