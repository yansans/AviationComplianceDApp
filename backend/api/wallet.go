package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type Wallet struct {
	store WalletStore
}

type WalletStore interface {
	Put(label string, content []byte) error
	Get(label string) ([]byte, error)
	Remove(label string) error
	Exists(label string) bool
	List() ([]string, error)
}

type Identity interface {
	toJSON() ([]byte, error)
	fromJSON(data []byte) (Identity, error)
	Credentials() ([]byte)
}

type X509Identity struct {
	MSP  string `json:"msp"`
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

func NewX509Identity(msp, cert, key string) *X509Identity {
	return &X509Identity{
		MSP:  msp,
		Cert: cert,
		Key:  key,
	}
}

func (i *X509Identity) toJSON() ([]byte, error) {
	return json.Marshal(i)
}

func (i *X509Identity) fromJSON(data []byte) (Identity, error) {
	var identity X509Identity
	err := json.Unmarshal(data, &identity)
	if err != nil {
		return nil, err
	}
	return &identity, nil
}

func (w *Wallet) Put(label string, identity Identity) error {
	data, err := identity.toJSON()
	if err != nil {
		return err
	}
	return w.store.Put(label, data)
}

func (w *Wallet) Get(label string) (X509Identity, error) {
	data, err := w.store.Get(label)
	if err != nil {
		return X509Identity{}, err
	}

	var identity X509Identity
	err = json.Unmarshal(data, &identity)
	if err != nil {
		return X509Identity{}, err
	}

	return identity, nil
}

func (w *Wallet) Remove(label string) error {
	return w.store.Remove(label)
}

func (w *Wallet) Exists(label string) bool {
	return w.store.Exists(label)
}

func (w *Wallet) List() ([]string, error) {
	return w.store.List()
}

func (i *X509Identity) Credentials() ([]byte) {
	credentials := append([]byte(i.Cert), []byte(i.Key)...)
	return credentials
}

func (i *X509Identity) MspID() string {
	return i.MSP
}

func (i *X509Identity) Signer() (func(digest []byte) ([]byte, error), error) {
	privateKey, err := parsePrivateKey(i.Key)
	if err != nil {
		return nil, err
	}

	return func(digest []byte) ([]byte, error) {
		return signMessage(privateKey, digest)
	}, nil
}

func LoadIdentityFromFiles(msp, certPath, keyPath string) (*X509Identity, error) {
	cert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate: %v", err)
	}

	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %v", err)
	}

	return &X509Identity{
		MSP:  msp,
		Cert: string(cert),
		Key:  string(key),
	}, nil
}

func NewWallet(identity Identity, store WalletStore) (*Wallet, error) {

	wallet := &Wallet{store: store}
	err := wallet.Put("user_identity", identity)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

type FileWalletStore struct {
	identities map[string][]byte
}

func (s *FileWalletStore) Put(label string, content []byte) error {
	if s.identities == nil {
		s.identities = make(map[string][]byte)
	}
	s.identities[label] = content
	return nil
}

func (s *FileWalletStore) Get(label string) ([]byte, error) {
	content, exists := s.identities[label]
	if !exists {
		return nil, errors.New("identity not found")
	}
	return content, nil
}

func (s *FileWalletStore) Remove(label string) error {
	delete(s.identities, label)
	return nil
}

func (s *FileWalletStore) Exists(label string) bool {
	_, exists := s.identities[label]
	return exists
}

func (s *FileWalletStore) List() ([]string, error) {
	var labels []string
	for label := range s.identities {
		labels = append(labels, label)
	}
	return labels, nil
}
