package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/ceriath/rsa-shamir-secret-sharing/s4"
	"github.com/spf13/cobra"
)

var gn, gk int
var keysize int
var saveKey, gSaveShares bool

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().IntVarP(&gn, "sharecount", "n", 5, "number of shares to generate")
	generateCmd.Flags().IntVarP(&gk, "threshold", "k", 3, "number of required shares to reconstruct the secret")
	generateCmd.Flags().IntVarP(&keysize, "bitsize", "b", 2048, "bitsize of the generated key")
	generateCmd.Flags().BoolVarP(&saveKey, "createKeyfile", "c", false, "set to true if you want to save the private key to a file")
	generateCmd.Flags().BoolVarP(&gSaveShares, "saveShares", "f", false, "set to true if you want to save the shares to a file")
}

var generateCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generates a shared RSA KeyPair",
	Long:  `Generates a public RSA key and n shares for the private key of which k are required for reconstruction`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := rand.Reader

		key, err := rsa.GenerateKey(reader, keysize)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		publicKey := key.PublicKey
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if saveKey {
			savePEMKey("private.pem", key)
		}
		savePublicPEMKey("public.pem", publicKey)

		shares, err := s4.Split(x509.MarshalPKCS1PrivateKey(key), byte(gk), byte(gn))
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		for _, share := range shares {
			fmt.Printf("Share %d:\n%s\n\n", share.Index, share.GetBase64())
		}

		if gSaveShares {
			storeShares(shares)
		}
	},
}

func savePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}
	defer func() {
		err = outFile.Close()
		if err != nil {
			panic("can't close file")
		}
	}()

	b := x509.MarshalPKCS1PrivateKey(key)

	var privateKey = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: b,
	}

	err = pem.Encode(outFile, privateKey)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}
}

func savePublicPEMKey(fileName string, pubkey rsa.PublicKey) {
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&pubkey)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemfile, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}
	defer func() {
		err = pemfile.Close()
		if err != nil {
			panic("can't close file")
		}
	}()

	err = pem.Encode(pemfile, pemkey)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}
}
