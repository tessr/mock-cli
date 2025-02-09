package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	cReset  = "\033[0m"
	cRed    = "\033[31m"
	cGreen  = "\033[32m"
	cYellow = "\033[33m"
	cBlue   = "\033[34m"
	cPurple = "\033[35m"
	cCyan   = "\033[36m"
	cGray   = "\033[37m"
	cWhite  = "\033[97m"
)

const usageString = `
Usage: netchef <command> [options]

Commands:
  generate           Generate the default configuration or output.
  deploy             Deploy using optional input files. If these files are not provided, 
                     deployer will use the default files inventory.yaml, manifest.yaml 
                     and state.json, in your working directory.
					 
    --inventory <file>       (Optional) Path to the inventory file, which describes the services
                             that netchef will deploy.
    --manifest <file>        (Optional) Path to the manifest file, which describes the overall devnet.
    --state <file>           (Optional) Path to the state file generated by OP Deployer.

Examples:
  netchef generate
  netchef deploy --inventory inventory.yaml --manifest manifest.yaml
  netchef deploy --state state.json
`

const (
	defaultInventory = "inventory.yaml"
	defaultManifest  = "manifest.yaml"
	defaultState     = "state.json"
)

func main() {
	if len(os.Args) < 2 {
		log.Print(usageString)
		os.Exit(1)
	}

	command := os.Args[1]
	if command == "generate" {
		// generate files for a default network, with sensible defaults
		err := generate()
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}

	if command == "deploy" {
		// Define flags
		deployCmd := flag.NewFlagSet("deploy", flag.ExitOnError)
		inventory := deployCmd.String("inventory", "", "Path to the inventory file, which describes the services that netchef will deploy.")
		manifest := deployCmd.String("manifest", "", "Path to the manifest file, which describes the overall devnet, including chain ID and chain name.")
		deployerState := deployCmd.String("state", "", "Path to the state.json file generated by OP Deployer.")

		// Parse the command-line arguments, skipping the first argument (the subcommand)
		deployCmd.Parse(os.Args[2:])

		// Read files or default values
		_, err := readFileOrDefault(*inventory, defaultInventory)
		if err != nil {
			os.Exit(1)
		}

		manifestContent, err := readFileOrDefault(*manifest, defaultManifest)
		if err != nil {
			os.Exit(1)
		}

		_, err = readFileOrDefault(*deployerState, defaultState)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Parse manifest.yaml to get the chainID
		parsedManifest, err := parseManifest(manifestContent)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var chainID = 0
		if len(parsedManifest.L2.Chains) > 0 {
			chainID = int(parsedManifest.L2.Chains[0].ChainID)
		} else {
			log.Println("chain ID unset")
		}

		deploy(parsedManifest.Name, chainID)
		os.Exit(0)
	}

	log.Print(usageString)
	os.Exit(1)
}

func generate() error {
	name := randomChainName()
	chainID := rand.Intn(10000000000)
	log.Printf("🧑🏻‍🍳  Cooking up inventory, manifest, and state files with default state and randomly generated chain name and ID... ")
	time.Sleep(time.Second)
	log.Printf("🍤  Generated %snetwork_name=%s%s and %schain_id=%d%s...", cBlue, name, cReset, cBlue, chainID, cReset)
	time.Sleep(time.Second)

	// write manifest file
	log.Printf("🥄  Cooking up a manifest file at %s%s%s...", cBlue, defaultManifest, cReset)
	b := []byte(fmt.Sprintf(manifest, name, time.Now().Format("2006-01-02"), name, chainID))
	err := os.WriteFile(defaultManifest, b, 0644)
	if err != nil {
		return err
	}

	time.Sleep(time.Second)

	// write inventory file
	log.Printf("🥄  Cooking up an inventory file at %s%s%s with all default services... ", cBlue, defaultInventory, cReset)
	b = []byte(fmt.Sprintf(inventory, name))
	err = os.WriteFile(defaultInventory, b, 0644)
	if err != nil {
		return err
	}
	time.Sleep(time.Second)

	// write state.json
	log.Println("🚀  Calling into op-deployer to generate L1 contracts and genesis state...")
	time.Sleep(2 * time.Second)
	log.Printf("✅  %sL1 contracts successfully deployed to Sepolia through OPCM!%s", cGreen, cReset)
	time.Sleep(time.Second)
	log.Println("🚀  Genesis state generated by op-deployer!")
	time.Sleep(time.Second)
	log.Printf("🥄  Writing the genesis created by op-deployer to %s%s%s...", cBlue, defaultState, cReset)
	b = []byte(state)
	err = os.WriteFile(defaultState, b, 0644)
	if err != nil {
		return err
	}
	time.Sleep(2 * time.Second)

	log.Printf("✅  %sRequired files generated successfully for chain_id=%d.%s", cGreen, chainID, cReset)

	return nil
}

func deploy(chainName string, chainID int) {
	log.Printf("🧑🏻‍🍳  Cooking up a new network named %s%s%s with %schainID=%d%s", cBlue, chainName, cReset, cGreen, chainID, cReset)

	// launch nodes
	log.Printf("🥄  Cooking up an op-geth sequencer node at version op-geth/v1.101411.4-rc.4...")
	time.Sleep(time.Second)
	log.Printf("🥄  Cooking up an op-reth sequencer node at version op-reth/v1.1.5...")
	time.Sleep(time.Second)
	log.Printf("🥄  Cooking up an op-geth rpc node at version op-geth/v1.101411.4-rc.4...")
	time.Sleep(time.Second)
	log.Printf("🥄  Cooking up an op-reth rpc node at version op-reth/v1.1.5...")
	time.Sleep(2 * time.Second)

	// launch services
	log.Printf("🥄  Cooking up a proposer service at version op-proposer/v1.10.0-rc.2...")
	time.Sleep(time.Second)
	log.Printf("🥄  Cooking up a batcher service at version op-batcher/v1.10.0...")
	time.Sleep(time.Second)
	log.Printf("🥄  Cooking up a challenger service at version op-challenger/v1.3.1-rc.4...")
	time.Sleep(time.Second)
	log.Printf("🥄  Cooking up a proxyd service...")
	time.Sleep(time.Second)
	log.Printf("🥄  Cooking up an op-conductor service...")
	time.Sleep(2 * time.Second)

	// buen apetito!
	log.Printf("✅  %sYour network is now available at %shttps://%s-0.optimism.io%s", cGreen, cBlue, chainName, cReset)
	log.Printf("🧑🏻‍🍳  Buen apetito!")
}

func readFileOrDefault(filename, defaultFilename string) (string, error) {
	if filename == "" {
		filename = defaultFilename
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("File %s not found, using default: %s\n", filename, defaultFilename)
		filename = defaultFilename
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Failed to read file %s. %sDid you mean to run netchef generate?%s\n", filename, cRed, cReset)
		return "", err
	}
	return string(content), nil
}

func parseManifest(content string) (*Manifest, error) {
	var manifest Manifest
	err := yaml.Unmarshal([]byte(content), &manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to parse manifest.yaml: %v", err)
	}
	return &manifest, nil
}
