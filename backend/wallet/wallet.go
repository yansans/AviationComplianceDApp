package wallet

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
}

type X509Identity struct {
	MSP  string `json:"msp"`
	Cert string `json:"cert"`
	Key  string `json:"key"`
}
