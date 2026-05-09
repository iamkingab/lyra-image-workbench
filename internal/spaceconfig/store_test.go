package spaceconfig

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/y08lin4/image-Workbench-Localhost-Version/internal/spaces"
)

func TestStoreMasksAPIKeyInPublicConfig(t *testing.T) {
	spaceStore, err := spaces.NewFileStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewFileStore() error = %v", err)
	}
	session, err := spaceStore.CreateOrOpenByPassword("R7!Blue#Vault$2026")
	if err != nil {
		t.Fatalf("CreateOrOpenByPassword() error = %v", err)
	}
	store := NewStore(spaceStore)

	rawKey := "  sk-test-secret-1234567890  "
	public, err := store.Update(session.Token, Update{APIKey: &rawKey})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if !public.APIKeySet {
		t.Fatal("APIKeySet should be true")
	}
	if strings.Contains(public.APIKeyPreview, "secret-123456") {
		t.Fatalf("API key preview leaked too much secret: %q", public.APIKeyPreview)
	}
	encoded, _ := json.Marshal(public)
	if strings.Contains(string(encoded), "sk-test-secret-1234567890") {
		t.Fatalf("public config leaked raw API key: %s", encoded)
	}

	private, err := store.Get(session.Token)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if private.APIKey != "sk-test-secret-1234567890" {
		t.Fatalf("private API key was not trimmed/persisted: %q", private.APIKey)
	}
}

func TestStorePersistsDefaultConcurrency(t *testing.T) {
	spaceStore, err := spaces.NewFileStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewFileStore() error = %v", err)
	}
	session, err := spaceStore.CreateOrOpenByPassword("R7!Green#Vault$2026")
	if err != nil {
		t.Fatalf("CreateOrOpenByPassword() error = %v", err)
	}
	store := NewStore(spaceStore)

	defaultPublic, err := store.Public(session.Token)
	if err != nil {
		t.Fatalf("Public() error = %v", err)
	}
	if defaultPublic.DefaultConcurrency != 1 {
		t.Fatalf("default concurrency = %d", defaultPublic.DefaultConcurrency)
	}

	value := 6
	public, err := store.Update(session.Token, Update{DefaultConcurrency: &value})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if public.DefaultConcurrency != 6 {
		t.Fatalf("DefaultConcurrency = %d", public.DefaultConcurrency)
	}

	larger := 99
	public, err = store.Update(session.Token, Update{DefaultConcurrency: &larger})
	if err != nil {
		t.Fatalf("Update(larger) error = %v", err)
	}
	if public.DefaultConcurrency != 99 {
		t.Fatalf("DefaultConcurrency should keep values above 4, got %d", public.DefaultConcurrency)
	}
}
