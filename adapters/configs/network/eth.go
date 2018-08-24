package network

// EthereumNetwork are the parameters required to create an ethereum client
type EthereumNetwork struct {
	Chain              string `json:"chain"`
	URL                string `json:"url"`
	Network            string `json:"network"`
	RenExAtomicSwapper string `json:"renExAtomicSwapper"`
	RenExAtomicInfo    string `json:"renExAtomicInfo"`
	RenExSettlement    string `json:"renExSettlement"`
	Orderbook          string `json:"orderbook"`
}

func (network *Config) GetEthereumNetwork() EthereumNetwork {
	network.mu.RLock()
	defer network.mu.RUnlock()
	return network.Ethereum
}

func (config *Config) SetEthereumNetwork(ethereumConfig EthereumNetwork) {
	config.mu.Lock()
	defer config.mu.Unlock()
	config.Ethereum = ethereumConfig
	config.Update()
}
