package main

import (
	"fmt"

	"github.com/fentec-project/gofe/abe"
)

func main() {
	msg := "Attack at dawn!"
	policy := "((0 AND 1) OR (2 AND 3)) AND 5"

	//gamma := []int{0, 2, 3, 5} // owned attributes
	gamma := []string{"0", "2", "3", "5"}

	a := abe.NewFAME()                             // Create the scheme instance
	pubKey, secKey, _ := a.GenerateMasterKeys()    // Create a public key and a master secret key
	msp, _ := abe.BooleanToMSP(policy, false)      // The MSP structure defining the policy
	cipher, _ := a.Encrypt(msg, msp, pubKey)       // Encrypt msg with policy msp under public key pubKey
	keys, _ := a.GenerateAttribKeys(gamma, secKey) // Generate keys for the entity with attributes gamma
	dec, _ := a.Decrypt(cipher, keys, pubKey)      // Decrypt the message
	fmt.Println("prod :", dec)
}
