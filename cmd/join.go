package cmd

import (
	"bufio"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ceriath/rsa-shamir-secret-sharing/s4"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(joinCmd)
}

var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "Joins shares to reconstruct the secret",
	Long:  `Joins multiple shares to reconstruct a secret`,
	Run: func(cmd *cobra.Command, args []string) {

		var err error

		shares := make([]s4.Share, 1)
		b64Shares := make([]string, 1)

		b64Shares[0] = readShare()

		shares[0], err = s4.DecodeBase64Share(b64Shares[0])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		for i := 1; i < int(shares[0].RequiredShares); i++ {
			b64Shares = append(b64Shares, readShare())
			var dec s4.Share
			dec, err = s4.DecodeBase64Share(b64Shares[i])
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			shares = append(shares, dec)
		}

		start := time.Now()

		result, err := s4.Reconstruct(shares)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		elapsed := time.Since(start)
		fmt.Printf("Reconstruction took %s\n", elapsed)

		f, err := os.Create("secret.txt")
		if err != nil {
			fmt.Println("Can't write secret to file.")
			fmt.Println(err.Error())
			return
		}

		defer func() {
			err = f.Close()
			if err != nil {
				panic("can't close file")
			}
		}()

		key, err := x509.ParsePKCS1PrivateKey(result)

		if err == nil {
			keybytes := x509.MarshalPKCS1PrivateKey(key)
			var privateKey = &pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: keybytes,
			}

			err = pem.Encode(f, privateKey)
			if err != nil {
				fmt.Println("Can't write secret to file.")
				fmt.Println(err.Error())
				return
			}
		} else {
			_, err = f.Write([]byte(base64.StdEncoding.EncodeToString(result)))
			if err != nil {
				fmt.Println("Can't write secret to file.")
				fmt.Println(err.Error())
				return
			}
		}
	},
}

func readShare() string {
	fmt.Println("Please enter a share: ")
	// byteshare, _ := terminal.ReadPassword(int(syscall.Stdin))
	reader := bufio.NewReader(os.Stdin)
	str, err := reader.ReadString('\n')
	if err != nil {
		panic("Error reading input")
	}

	//fmt.Printf("Read %s\n", str)
	return strings.TrimSpace(str)
}
