package main

import (
	"fmt"
	"math/big"

	"github.com/fentec-project/gofe/data"
	"github.com/fentec-project/gofe/quadratic"
	"github.com/fentec-project/gofe/sample"
)

func main() {
	l := 2                  // length of input vectors
	bound := big.NewInt(10) // Upper bound for coordinates of vectors x, y, and matrix F

	// Here we fill our vectors and the matrix F (that represents the
	// quadratic function) with random data from [0, bound).
	sampler := sample.NewUniform(bound)
	F, _ := data.NewRandomMatrix(l, l, sampler)
	x, _ := data.NewRandomVector(l, sampler)
	y, _ := data.NewRandomVector(l, sampler)

	sgp := quadratic.NewSGP(l, bound)     // Create scheme instance
	msk, _ := sgp.GenerateMasterKey()     // Create master secret key
	cipher, _ := sgp.Encrypt(x, y, msk)   // Encrypt input vectors x, y with secret key
	key, _ := sgp.DeriveKey(msk, F)       // Derive FE key for decryption
	dec, _ := sgp.Decrypt(cipher, key, F) // Decrypt the result to obtain x^T * F * y
	fmt.Println("prod :", dec)
}
# Functional Encryption Examples
It is based on [gofe](https://github.com/fentec-project/gofe) library.