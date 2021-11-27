package wallet

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/guiwoo/gocoin/utils"
)

const (
	signature     string = "4ae631c0fb49e6292f16016f64a67e19b8a1fe41d9f181e812b3b4068b1290819b2f91d9c48bc110d551910fd8a73ef1633eaa4195ef25c77b67aaf2da9e78"
	privateKey    string = "30770201010420c79486bba6abe92873742d86f81dbf7cb9b5b90a950c91110bebd31899b29c9aa00a06082a8648ce3d030107a144034200044a0056608dc665806ad4d35e590b256557b0a58264a5e6d9ddd82c29add6d2458cffbf0835e2b957cdb308aa729877988ebb35c4be035d86af498cd72549ab21"
	hashedMessage string = "1c4863cd55b5a4413fd59f054af57ba3c75c0698b3851d70f99b8de2d5c7338f"
)

func Start() {

	privateBytes, err := hex.DecodeString(privateKey)
	utils.HandleErr(err)

	private, err := x509.ParseECPrivateKey(privateBytes)
	utils.HandleErr(err)

	sigBytes, err := hex.DecodeString(signature)
	rBtyes := sigBytes[:len(sigBytes)/2]
	sBytes := sigBytes[len(sigBytes)/2:]

	var bigR, bigS = big.Int{}, big.Int{}

	bigR.SetBytes(rBtyes)
	bigS.SetBytes(sBytes)

	hashBytes, err := hex.DecodeString(hashedMessage)
	utils.HandleErr(err)

	ok := ecdsa.Verify(&private.PublicKey, hashBytes, &bigR, &bigS)
	fmt.Println(ok)
}
