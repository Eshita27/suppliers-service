package config

import (
	"context"
	"fmt"
	"os"
	"time"

	vault "github.com/hashicorp/vault/api"
)

type AppConfig struct {
	MongoURI string
	DBName   string
	Port     string
}

// LoadConfig fetches application configurations securely out of HashiCorp Vault
func LoadConfig() (*AppConfig, error) {
	vaultAddr := os.Getenv("VAULT_ADDR")
	vaultToken := os.Getenv("VAULT_TOKEN")

	// Set fallbacks for native non-docker runs
	if vaultAddr == "" {
		vaultAddr = "http://127.0.0.1:8200"
	}
	if vaultToken == "" {
		vaultToken = "root-token-secret"
	}

	// 1. Initialize core SDK client configs
	vaultCfg := vault.DefaultConfig()
	vaultCfg.Address = vaultAddr

	client, err := vault.NewClient(vaultCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize vault SDK adapter: %w", err)
	}
	client.SetToken(vaultToken)

	// 2. Fetch the secrets engine matrix under KV version 2 path
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// For local testing, we look under the default dev kv mount: secret/data/suppliers
	secret, err := client.KVv2("secret").Get(ctx, "suppliers")
	if err != nil {
		return nil, fmt.Errorf("unable to read configurations from vault engine path: %w", err)
	}

	// 3. Unpack and map values into structural configs
	mongoURI, _ := secret.Data["MONGO_URI"].(string)
	dbName, _ := secret.Data["MONGO_DB_NAME"].(string)

	if mongoURI == "" {
		mongoURI = "mongodb://suppliers-db:27017" // Reliable multi-container fallback fallback
	}
	if dbName == "" {
		dbName = "polyglot_inventory"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &AppConfig{
		MongoURI: mongoURI,
		DBName:   dbName,
		Port:     port,
	}, nil
}
