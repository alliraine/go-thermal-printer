//go:build !windows

package service

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"
)

func TestNewConfigServicePermissionDenied(t *testing.T) {
	if os.Getenv("GO_WANT_PERMISSION_DENIED_HELPER") == "1" {
		runPermissionDeniedHelper(t)
		return
	}

	if runtime.GOOS == "windows" {
		t.Skip("permission handling differs on Windows")
	}

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.toml")
	if err := os.WriteFile(configPath, []byte(""), 0o644); err != nil {
		t.Fatalf("failed to create config file: %v", err)
	}

	if err := os.Chmod(tempDir, 0o000); err != nil {
		t.Fatalf("failed to remove directory permissions: %v", err)
	}
	defer func() {
		if err := os.Chmod(tempDir, 0o755); err != nil {
			t.Fatalf("failed to restore directory permissions: %v", err)
		}
	}()

	t.Setenv("CONFIG_PATH", configPath)

	_, err := NewConfigService()
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			return
		}

		t.Fatalf("expected permission error, got %v", err)
	}

	if os.Getuid() != 0 {
		t.Fatalf("expected permission error when reading config file")
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestNewConfigServicePermissionDenied")
	cmd.Env = append(os.Environ(), "GO_WANT_PERMISSION_DENIED_HELPER=1", "CONFIG_PATH="+configPath)

	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("helper process failed: %v\n%s", err, string(output))
	}
}

func runPermissionDeniedHelper(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("permission handling differs on Windows")
	}

	if os.Getuid() == 0 {
		if err := syscall.Setgid(65534); err != nil {
			t.Fatalf("failed to drop gid: %v", err)
		}
		if err := syscall.Setuid(65534); err != nil {
			t.Fatalf("failed to drop uid: %v", err)
		}
	}

	t.Setenv("CONFIG_PATH", os.Getenv("CONFIG_PATH"))

	_, err := NewConfigService()
	if err == nil {
		t.Fatalf("expected permission error when reading config file")
	}

	if !errors.Is(err, os.ErrPermission) {
		t.Fatalf("expected permission error, got %v", err)
	}
}
