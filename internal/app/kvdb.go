package app

import (
	"fmt"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/knadh/koanf/providers/confmap"
)

var (
	// KvDB hold a connection to KV Database
	KvDB *badger.DB
)

// KvDBDefaults set up default configuration for redis client
func KvDBDefaults() {
	Config.Load(confmap.Provider(map[string]interface{}{
		"kvdb.path": "/tmp/badger",
	}, "."), nil)

}

// KvDBInit create the redis client based on koanf configuration
func KvDBInit() {
	db, err := badger.Open(badger.DefaultOptions(Config.String("kvdb.path")))
	if err != nil {
		panic("failed to open KvDB")
	}
	fmt.Println("Connection Opened to KvDB")
	KvDB = db
}

// KvDump dump a map[string]string into db
func KvDump(kv map[string]string) error {
	txn := KvDB.NewTransaction(true)
	for k, v := range kv {
		if err := txn.Set([]byte(k), []byte(v)); err == badger.ErrTxnTooBig {
			_ = txn.Commit()
			txn = KvDB.NewTransaction(true)
			_ = txn.Set([]byte(k), []byte(v))
		}
	}
	return txn.Commit()
}
