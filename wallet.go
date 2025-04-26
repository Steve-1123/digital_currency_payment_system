package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
	Address    string
}

func NewWallet() (*Wallet, error) {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}
	pubKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	hash := sha256.Sum256(pubKey)
	address := hex.EncodeToString(hash[:])

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  pubKey,
		Address:    address,
	}, nil
}

func (w *Wallet) Sign(data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	r, s, err := ecdsa.Sign(rand.Reader, w.PrivateKey, hash[:])
	if err != nil {
		return nil, err
	}
	signature := append(r.Bytes(), s.Bytes()...)
	return signature, nil
}

func (w *Wallet) VerifySignature(data, signature []byte, address string) bool {
	hash := sha256.Sum256(data)
	r := big.Int{}
	s := big.Int{}
	sigLen := len(signature)
	r.SetBytes(signature[:sigLen/2])
	s.SetBytes(signature[sigLen/2:])

	curve := elliptic.P256()
	x, y := elliptic.Unmarshal(curve, w.PublicKey)
	pubKey := ecdsa.PublicKey{Curve: curve, X: x, Y: y}
	return ecdsa.Verify(&pubKey, hash[:], &r, &s)
}
