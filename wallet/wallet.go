package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"os"

	"github.com/guiwoo/gocoin/utils"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
}

const (
	fileName string = "gocoin.wallet"
)

var w *wallet

func hasWalletFile() bool {
	// has a wallet already ?
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func createPrivateKey() *ecdsa.PrivateKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privKey
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)
	err = os.WriteFile(fileName, bytes, 0644)
	utils.HandleErr(err)
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if hasWalletFile() {
			// yes -> restore from file
		} else {
			// no -> 1. create prv key, 2. save to file
			key := createPrivateKey()
			persistKey(key)
			w.privateKey = key
		}
	}
	return w
}
