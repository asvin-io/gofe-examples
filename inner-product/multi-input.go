package main

import (
	"fmt"
	"math/big"

	"github.com/fentec-project/gofe/data"
	"github.com/fentec-project/gofe/innerprod/simple"
	"github.com/fentec-project/gofe/sample"
)

func main() {
	numClients := 2        // number of encryptors
	l := 3                 // length of input vectors
	bound := big.NewInt(5) // upper bound for input vectors

	// Simulate collection of input data.
	// X and Y represent matrices of input vectors, where X are collected
	// from numClients encryptors (omitted), and Y is only known by a single decryptor.
	// Encryptor i only knows its own input vector X[i].
	sampler := sample.NewUniform(bound)
	X, _ := data.NewRandomMatrix(numClients, l, sampler)
	// vecs := make([]data.Vector, 2) // a slice of 2 vectors
	// // fill vecs
	// vecs[0] = data.NewVector([]*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)})
	// vecs[1] = data.NewVector([]*big.Int{big.NewInt(4), big.NewInt(5), big.NewInt(6)})
	// X, _ := data.NewMatrix(vecs)
	Y, _ := data.NewRandomMatrix(numClients, l, sampler)

	fmt.Println("X : ", X)
	fmt.Println("Y : ", Y)
	xy, _ := X.Dot(Y)
	fmt.Println("x.y: ", xy)
	// Trusted entity instantiates scheme instance and generates
	// master keys for all the encryptors. It also derives the FE
	// key derivedKey for the decryptor.
	modulusLength := 2048
	multiDDH, _ := simple.NewDDHMultiPrecomp(numClients, l, modulusLength, bound)
	pubKey, secKey, _ := multiDDH.GenerateMasterKeys()
	derivedKey, _ := multiDDH.DeriveKey(secKey, Y)

	// Different encryptors may reside on different machines.
	// We simulate this with the for loop below, where numClients
	// encryptors are generated.
	encryptors := make([]*simple.DDHMultiClient, numClients)
	for i := 0; i < numClients; i++ {
		encryptors[i] = simple.NewDDHMultiClient(multiDDH.Params)
	}

	// Each encryptor encrypts its own input vector X[i] with the
	// keys given to it by the trusted entity.
	ciphers := make([]data.Vector, numClients)
	for i := 0; i < numClients; i++ {
		cipher, _ := encryptors[i].Encrypt(X[i], pubKey[i], secKey.OtpKey[i])
		ciphers[i] = cipher
	}
	fmt.Println("opt key: ", secKey.OtpKey)
	// Ciphers are collected by decryptor, who then computes
	// inner product over vectors from all encryptors.
	decryptor := simple.NewDDHMultiFromParams(numClients, multiDDH.Params)
	prod, _ := decryptor.Decrypt(ciphers, derivedKey, Y)
	sum := big.NewInt(0)
	for i := 0; i < decryptor.Slots; i++ {
		c, err := decryptor.DDH.Decrypt(ciphers[i], derivedKey.Keys[i], Y[i])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("c", c)
		sum.Add(sum, c)
		fmt.Println("sum ", sum)

	}
	fmt.Println("OTPKey ", derivedKey.OTPKey)
	res := new(big.Int).Sub(sum, derivedKey.OTPKey)
	fmt.Println("res1 ", res)
	res.Mod(res, decryptor.Params.Bound)
	fmt.Println("mod ", res)
	fmt.Println("prod :", prod)
}
