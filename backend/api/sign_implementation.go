package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"math/big"
)

func parsePrivateKey(pemKey string) (crypto.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemKey))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return key, nil
}

type ecdsaSignature struct {
    R *big.Int
    S *big.Int
}

func signMessage(privateKey crypto.PrivateKey, digest []byte) ([]byte, error) {
    switch key := privateKey.(type) {
    case *ecdsa.PrivateKey:
        r, s, err := ecdsa.Sign(rand.Reader, key, digest)
        if err != nil {
            return nil, err
        }

        // Ensure S is less than half the order of the curve
        curveOrder := key.Curve.Params().N
        halfOrder := new(big.Int).Rsh(curveOrder, 1) // Divide curve order by 2

        if s.Cmp(halfOrder) > 0 {
            // If S is larger than half the order, take the complement of S
            s.Sub(curveOrder, s)
        }

        // Create the ASN.1 DER encoded signature
        sig := ecdsaSignature{R: r, S: s}
        signature, err := asn1.Marshal(sig)
        if err != nil {
            return nil, fmt.Errorf("failed to encode signature: %v", err)
        }

        return signature, nil
    default:
        return nil, fmt.Errorf("unsupported private key type")
    }
}