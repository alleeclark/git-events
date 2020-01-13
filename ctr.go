package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"git-events/git"
	"io/ioutil"
	"log"
	"os"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh"
)

var ctrCommand = cli.Command{
	Name:        "ctr",
	Usage:       "ctr a gitevent",
	ArgsUsage:   "[flags] <ref>",
	Description: "ctr",
	Subcommands: cli.Commands{
		cli.Command{
			Name:    "subscribe",
			Aliases: []string{"sub"},
			Flags: []cli.Flag{
				cli.StringSliceFlag{Name: "topics", Value: &cli.StringSlice{"Added", "Deleted", "Modified", "Renamed", "Copied"}, Hidden: true, Usage: "Topics you would like to filter by"},
			},
			Action: func(c *cli.Context) error {
				client := git.NewEventsServiceClient(nil)
				eventClient, err := client.Event(context.Background(), &git.EventRequest{
					Topics: c.GlobalStringSlice("topics"),
				})
				if err != nil {
					return err
				}
				resp, err := eventClient.Recv()
				if err != nil {
					return err
				}
				m := &jsonpb.Marshaler{}
				if err := m.Marshal(os.Stdout, resp); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to marshal message %v", err)
				}
				return nil
			},
		},
	},
}

func generateKeys(privateKeyPath, publicKeyPath string) string {
	bitSize := 4096

	privateKey, err := generatePrivateKey(bitSize)
	if err != nil {
		log.Fatal(err.Error())
	}

	publicKeyBytes, err := generatePublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Fatal(err.Error())
	}

	privateKeyBytes := encodePrivateKeyToPEM(privateKey)

	err = writeKeyToFile(privateKeyBytes, privateKeyPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = writeKeyToFile([]byte(publicKeyBytes), publicKeyPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(publicKeyBytes)
}

// generatePrivateKey creates a RSA Private Key of specified byte size
func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	log.Println("Private Key generated")
	return privateKey, nil
}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

// generatePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func generatePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	log.Println("Public key generated")
	return pubKeyBytes, nil
}

// writePemToFile writes keys to a file
func writeKeyToFile(keyBytes []byte, saveFileTo string) error {
	err := ioutil.WriteFile(saveFileTo, keyBytes, 0600)
	if err != nil {
		return err
	}

	log.Printf("Key saved to: %s", saveFileTo)
	return nil
}
