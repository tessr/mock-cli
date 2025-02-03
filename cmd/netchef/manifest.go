package main

// This file is basically cribbed from the real netchef

// Manifest represents the top-level configuration for a network deployment
type Manifest struct {
	Name string   `yaml:"name"`
	Type string   `yaml:"type"`
	L1   L1Config `yaml:"l1"`
	L2   L2Config `yaml:"l2"`
}

// L1Config represents the L1 chain configuration
type L1Config struct {
	Name    string `yaml:"name"`
	ChainID uint64 `yaml:"chain_id"`
}

// L2Config represents the L2 chain configuration
type L2Config struct {
	Chains []Chain `yaml:"chains"`
}

// Chain represents a chain configuration
type Chain struct {
	Name    string `yaml:"name"`
	ChainID uint64 `yaml:"chain_id"`
}
