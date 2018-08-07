package client

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/republicprotocol/atom-go/adapters/configs/general"
	"github.com/republicprotocol/atom-go/services/renguardClient"
)

type renguardHTTPClient struct {
	ipAddress string
}

// NewrenguardHTTPClient creates a new RenguardClient interface, that interacts
// with RenGuard over http.
func NewrenguardHTTPClient(config config.Config) renguardClient.RenguardClient {
	return &renguardHTTPClient{
		ipAddress: config.RenGuardAddress(),
	}
}

func (client *renguardHTTPClient) ComplainDelayedAddressSubmission(orderID [32]byte) error {
	return client.watch(orderID)
}

func (client *renguardHTTPClient) ComplainDelayedRequestorInitiation(orderID [32]byte) error {
	return client.watch(orderID)
}

func (client *renguardHTTPClient) ComplainWrongRequestorInitiation(orderID [32]byte) error {
	return client.watch(orderID)
}

func (client *renguardHTTPClient) ComplainDelayedResponderInitiation(orderID [32]byte) error {
	return client.watch(orderID)
}

func (client *renguardHTTPClient) ComplainWrongResponderInitiation(orderID [32]byte) error {
	return client.watch(orderID)
}

func (client *renguardHTTPClient) ComplainDelayedRequestorRedemption(orderID [32]byte) error {
	return client.watch(orderID)
}

func (client *renguardHTTPClient) watch(orderID [32]byte) error {
	resp, err := http.Post(fmt.Sprintf("http://"+client.ipAddress+"/watch?orderID="+hex.EncodeToString(orderID[:])), "text", nil)
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		return nil
	}
	return fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
}
