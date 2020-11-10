package main

import (
	"fmt"
	"math/big"

	"github.com/fentec-project/gofe/data"
	"github.com/fentec-project/gofe/innerprod/simple"
)

func main() {
	l := 3                  // length of input vectors
	bound := big.NewInt(10) // upper bound for input vector coordinates
	modulusLength := 2048   // bit length of prime modulus p
	trustedEnt, _ := simple.NewDDHPrecomp(l, modulusLength, bound)
	msk, mpk, _ := trustedEnt.GenerateMasterKeys()
	y := data.NewVector([]*big.Int{big.NewInt(1), big.NewInt(1), big.NewInt(4)})
	fmt.Println("y: ", y)
	fsKey, _ := trustedEnt.DeriveKey(msk, y)

	enc := simple.NewDDHFromParams(trustedEnt.Params)
	x := data.NewVector([]*big.Int{big.NewInt(3), big.NewInt(4), big.NewInt(5)})
	fmt.Println("x: ", x)
	cipher, _ := enc.Encrypt(x, mpk)

	dec := simple.NewDDHFromParams(trustedEnt.Params)
	xy, _ := dec.Decrypt(cipher, fsKey, y)
	fmt.Println("x.y:", xy)
}


