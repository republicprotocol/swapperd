package watch_test

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/republicprotocol/atom-go/adapters/atoms/btc"
	"github.com/republicprotocol/atom-go/adapters/atoms/eth"
	btcclient "github.com/republicprotocol/atom-go/adapters/clients/btc"
	ethclient "github.com/republicprotocol/atom-go/adapters/clients/eth"
	"github.com/republicprotocol/atom-go/adapters/config"
	"github.com/republicprotocol/atom-go/adapters/keystore"
	"github.com/republicprotocol/atom-go/adapters/store/leveldb"
	"github.com/republicprotocol/atom-go/drivers/btc/regtest"
	"github.com/republicprotocol/atom-go/services/swap"
	. "github.com/republicprotocol/atom-go/services/watch"

	ax "github.com/republicprotocol/atom-go/adapters/info/eth"
	net "github.com/republicprotocol/atom-go/adapters/networks/eth"
	wal "github.com/republicprotocol/atom-go/adapters/wallet/mock"
	"github.com/republicprotocol/atom-go/domains/match"
)

var _ = Describe("Ethereum - Bitcoin Atomic Swap using Watch", func() {

	var aliceWatch, bobWatch Watch
	var aliceOrderID, bobOrderID [32]byte

	rand.Read(aliceOrderID[:])
	rand.Read(bobOrderID[:])

	BeforeSuite(func() {
		var aliceInfo, bobInfo swap.Info
		var aliceNet, bobNet swap.Network
		var aliceSendValue, bobSendValue *big.Int
		var aliceCurrency, bobCurrency uint32
		var alice, bob *bind.TransactOpts
		var swapID [32]byte

		rand.Read(swapID[:])

		aliceCurrency = 1
		bobCurrency = 0

		var confPath = "/Users/susruth/go/src/github.com/republicprotocol/atom-go/secrets/config.json"
		var ksPathA = "/Users/susruth/go/src/github.com/republicprotocol/atom-go/secrets/keystoreA.json"
		var ksPathB = "/Users/susruth/go/src/github.com/republicprotocol/atom-go/secrets/keystoreB.json"

		config, err := config.LoadConfig(confPath)
		Expect(err).ShouldNot(HaveOccurred())
		keystoreA := keystore.NewKeystore(ksPathA)
		keystoreB := keystore.NewKeystore(ksPathB)

		ganache, err := ethclient.Connect(config)
		Expect(err).ShouldNot(HaveOccurred())

		keysA, err := keystoreA.LoadKeys()
		Expect(err).ShouldNot(HaveOccurred())

		keysB, err := keystoreB.LoadKeys()
		Expect(err).ShouldNot(HaveOccurred())

		aliceEthKey := keysA[0]
		bobEthKey := keysB[0]
		aliceBtcKey := keysA[1]
		bobBtcKey := keysB[1]

		aliceAddrBytes, err := aliceBtcKey.GetAddress()
		Expect(err).ShouldNot(HaveOccurred())
		bobAddrBytes, err := bobBtcKey.GetAddress()
		Expect(err).ShouldNot(HaveOccurred())

		aliceAddr := string(aliceAddrBytes)
		bobAddr := string(bobAddrBytes)

		ownerECDSA, err := crypto.HexToECDSA("2aba04ee8a322b8648af2a784144181a0c793f1a2e80519418f3d20bbfb22249")
		Expect(err).ShouldNot(HaveOccurred())
		owner := bind.NewKeyedTransactor(ownerECDSA)

		aliceEthAddr, err := aliceEthKey.GetAddress()
		Expect(err).ShouldNot(HaveOccurred())
		ganache.Transfer(common.BytesToAddress(aliceEthAddr), owner, 1000000000000000000)

		bobEthAddr, err := bobEthKey.GetAddress()
		Expect(err).ShouldNot(HaveOccurred())
		ganache.Transfer(common.BytesToAddress(bobEthAddr), owner, 1000000000000000000)

		time.Sleep(5 * time.Second)
		connection, err := btcclient.Connect(config)
		Expect(err).ShouldNot(HaveOccurred())

		aliceSendValue = big.NewInt(10000000)
		bobSendValue = big.NewInt(10000000)

		go func() {
			err = regtest.Mine(connection)
			Expect(err).ShouldNot(HaveOccurred())
		}()
		time.Sleep(5 * time.Second)

		_aliceAddr, err := btcutil.DecodeAddress(aliceAddr, connection.ChainParams)
		Expect(err).ShouldNot(HaveOccurred())
		_bobAddr, err := btcutil.DecodeAddress(bobAddr, connection.ChainParams)
		Expect(err).ShouldNot(HaveOccurred())

		btcvalue, err := btcutil.NewAmount(5.0)
		Expect(err).ShouldNot(HaveOccurred())

		connection.Client.SendToAddress(_aliceAddr, btcvalue)
		connection.Client.SendToAddress(_bobAddr, btcvalue)

		_aliceWIF, err := aliceBtcKey.GetKeyString()
		Expect(err).ShouldNot(HaveOccurred())

		aliceWIF, err := btcutil.DecodeWIF(_aliceWIF)
		Expect(err).ShouldNot(HaveOccurred())

		err = connection.Client.ImportPrivKey(aliceWIF)
		Expect(err).ShouldNot(HaveOccurred())

		_bobWIF, err := bobBtcKey.GetKeyString()
		Expect(err).ShouldNot(HaveOccurred())

		bobWIF, err := btcutil.DecodeWIF(_bobWIF)
		Expect(err).ShouldNot(HaveOccurred())

		err = connection.Client.ImportPrivKey(bobWIF)
		Expect(err).ShouldNot(HaveOccurred())

		alice = bind.NewKeyedTransactor(aliceEthKey.GetKey())
		bob = bind.NewKeyedTransactor(bobEthKey.GetKey())

		aliceNet, err = net.NewEthereumNetwork(ganache, alice)
		Expect(err).Should(BeNil())

		bobNet, err = net.NewEthereumNetwork(ganache, bob)
		Expect(err).Should(BeNil())

		aliceInfo, err = ax.NewEtereumAtomInfo(ganache, alice)
		Expect(err).Should(BeNil())

		bobInfo, err = ax.NewEtereumAtomInfo(ganache, bob)
		Expect(err).Should(BeNil())

		atomMatch := match.NewMatch(aliceOrderID, bobOrderID, aliceSendValue, bobSendValue, aliceCurrency, bobCurrency)
		mockWallet := wal.NewMockWallet()

		err = mockWallet.SetMatch(atomMatch)
		Expect(err).Should(BeNil())

		reqAlice, err := eth.NewEthereumAtom(ganache, aliceEthKey)
		Expect(err).Should(BeNil())

		reqBob := btc.NewBitcoinAtom(connection, bobBtcKey)
		resAlice := btc.NewBitcoinAtom(connection, aliceBtcKey)

		resBob, err := eth.NewEthereumAtom(ganache, bobEthKey)
		Expect(err).Should(BeNil())

		aliceStr := swap.NewSwapStore(leveldb.NewLDBStore("/db"))
		bobStr := swap.NewSwapStore(leveldb.NewLDBStore("/db"))

		aliceWatch = NewWatch(aliceNet, aliceInfo, mockWallet, reqAlice, resAlice, aliceStr)
		bobWatch = NewWatch(bobNet, bobInfo, mockWallet, reqBob, resBob, bobStr)
	})

	It("can do an eth - btc atomic swap", func() {
		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer GinkgoRecover()

			err := aliceWatch.Run(aliceOrderID)
			fmt.Println(err)
			Expect(err).ShouldNot(HaveOccurred())

			fmt.Println("Done 1")
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer GinkgoRecover()

			err := bobWatch.Run(bobOrderID)
			fmt.Println(err)
			Expect(err).ShouldNot(HaveOccurred())

			fmt.Println("Done 2")
		}()

		wg.Wait()
	})
})
