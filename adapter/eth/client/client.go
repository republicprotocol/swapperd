package client

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/beth-go"
	"github.com/republicprotocol/swapperd/foundation"
)

type NetworkConfig struct {
	URL     string          `json:"url"`
	Network string          `json:"network"`
	Tokens  []EthereumToken `json:"tokens"`
}

type EthereumToken struct {
	Name           string `json:"name"`
	TokenAddress   string `json:"tokenAddress"`
	SwapperAddress string `json:"swapperAddress"`
}

type Client struct {
	beth.Account
	SwapperContracts map[string]string
	TokenContracts   map[string]string
	Network          string
}

func New(config NetworkConfig, key *ecdsa.PrivateKey) (*Client, error) {
	acc, err := beth.NewAccount(config.URL, key)
	if err != nil {
		return nil, err
	}
	swapperContracts := map[string]string{}
	tokenContracts := map[string]string{}
	for _, token := range config.Tokens {
		swapperContracts[token.Name] = token.SwapperAddress
		tokenContracts[token.Name] = token.TokenAddress
	}
	return &Client{
		acc,
		swapperContracts,
		tokenContracts,
		config.Network,
	}, nil
}

func (client *Client) GetTokenAddress(token foundation.Token) common.Address {
	return common.HexToAddress(client.TokenContracts[token.Name])
}

func (client *Client) GetSwapperAddress(token foundation.Token) common.Address {
	return common.HexToAddress(client.SwapperContracts[token.Name])
}

func (client *Client) Transact(ctx context.Context, preCon func() bool, tx func(*bind.TransactOpts) (*types.Transaction, error), postCon func() bool, confirmations int64) error {
	if err := client.Account.Transact(ctx, preCon, tx, postCon, confirmations); err != beth.ErrPreConditionCheckFailed {
		return err
	}
	return nil
}

func (client *Client) Conn() *ethclient.Client {
	cl := client.EthClient()
	return cl.EthClient()
}

func (client *Client) FormatTransactionView(msg, txHash string) string {
	switch client.Network {
	case "kovan":
		return fmt.Sprintf("%s, the transaction can be viewed at https://kovan.etherscan.io/tx/%s", msg, txHash)
	case "ropsten":
		return fmt.Sprintf("%s, the transaction can be viewed at https://ropsten.etherscan.io/tx/%s", msg, txHash)
	case "mainnet":
		return fmt.Sprintf("%s, the transaction can be viewed at https://etherscan.io/tx/%s", msg, txHash)
	default:
		panic(fmt.Sprintf("Unknown network :%s", client.Network))
	}
}
