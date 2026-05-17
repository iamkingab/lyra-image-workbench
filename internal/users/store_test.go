package users

import (
	"path/filepath"
	"testing"
	"time"
)

func TestLoginRequiresTOTPWhenEnabled(t *testing.T) {
	store, err := NewStore(filepath.Join(t.TempDir(), "users.json"))
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}
	if _, err := store.Register("alice01", "R7!Blue#Vault$2026", ""); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	setup, err := store.BeginTOTPSetup("alice01")
	if err != nil {
		t.Fatalf("BeginTOTPSetup() error = %v", err)
	}
	code, err := totpCode(setup.Secret, time.Now().Unix()/totpPeriod)
	if err != nil {
		t.Fatalf("totpCode() error = %v", err)
	}
	if err := store.EnableTOTP("alice01", code); err != nil {
		t.Fatalf("EnableTOTP() error = %v", err)
	}
	if _, err := store.Login("alice01", "R7!Blue#Vault$2026", ""); err == nil {
		t.Fatal("Login without TOTP code should fail")
	}
	if _, err := store.Login("alice01", "R7!Blue#Vault$2026", code); err != nil {
		t.Fatalf("Login with TOTP code error = %v", err)
	}
}
