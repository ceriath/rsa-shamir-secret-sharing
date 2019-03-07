package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/ceriath/rsa-shamir-secret-sharing/s4"
	"github.com/spf13/cobra"
)

var sk, sn int
var secret string
var secretpath string
var sSaveShares, deleteSecretfile bool

func init() {
	rootCmd.AddCommand(splitCmd)
	splitCmd.Flags().IntVarP(&sn, "sharecount", "n", 5, "number of shares to generate")
	splitCmd.Flags().IntVarP(&sk, "threshold", "k", 3, "number of required shares to reconstruct the secret")
	splitCmd.Flags().StringVarP(&secret, "secret", "s", "", "the secret to split")
	splitCmd.Flags().StringVarP(&secretpath, "secretpath", "p", "", "filepath to the secret")
	splitCmd.Flags().BoolVarP(&sSaveShares, "saveShares", "f", false, "set to true if you want to save the shares to a file")
	splitCmd.Flags().BoolVarP(&deleteSecretfile, "deleteSecretfile", "d", false, "set to true if you want to delete the file containing the secret")
}

var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "Splits a given secret into shares",
	Long:  `Splits the given secret into n shares of which k are required to reconstruct the secret`,
	Run: func(cmd *cobra.Command, args []string) {

		if secret == "" && secretpath == "" {
			fmt.Printf("No secret provided. See rsasss help for usage.\n")
			return
		}
		if secret != "" && secretpath != "" {
			fmt.Printf("A secret and a secretpath are defined. Please only use one.\n")
			return
		}

		var s []byte

		if secret != "" {
			s = []byte(secret)

		} else {
			var err error
			s, err = ioutil.ReadFile(secretpath)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}

		start := time.Now()

		shares, err := s4.Split(s, byte(sk), byte(sn))
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		elapsed := time.Since(start)

		for _, share := range shares {
			fmt.Printf("Share %d:\n%s\n\n", share.Index, share.GetBase64())
		}

		fmt.Printf("Splitting took %s\n", elapsed)

		if sSaveShares {
			storeShares(shares)
		}

		if deleteSecretfile && secretpath != "" {
			err = os.Remove(secretpath)
			if err != nil {
				fmt.Printf("!!!!! THE SECRETFILE COULD NOT BE DELETED !!!!!\nMake sure to delete it manually!\nError: %s", err.Error())
			}
		}
	},
}

func storeShares(shares []s4.Share) {

	outFile, err := os.Create("shares")
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

	for i, share := range shares {
		fmt.Fprintf(outFile, "Share %d:\n%s\n\n", i, share.GetBase64())
	}
}
