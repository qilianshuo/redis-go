package logger

import (
	"os"
	"testing"
)

func TestCheckNotExist(t *testing.T) {
	tmpFile := "test_check_not_exist.tmp"
	// Ensure file does not exist
	_ = os.Remove(tmpFile)
	if !checkNotExist(tmpFile) {
		t.Errorf("Expected file to not exist")
	}
	// Create file
	f, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	f.Close()
	if checkNotExist(tmpFile) {
		t.Errorf("Expected file to exist")
	}
	_ = os.Remove(tmpFile)
}

func TestCheckPermission(t *testing.T) {
	dir := "test_permission_dir"
	_ = os.RemoveAll(dir)
	if checkPermission(dir) {
		t.Errorf("Expected no permission error for non-existent dir")
	}
	// Create dir
	if err := os.Mkdir(dir, 0755); err != nil {
		t.Fatalf("Failed to create dir: %v", err)
	}
	if checkPermission(dir) {
		t.Errorf("Expected permission for created dir")
	}
	_ = os.RemoveAll(dir)
}

func TestIsNotExistMkDir(t *testing.T) {
	dir := "test_mkdir_dir"
	_ = os.RemoveAll(dir)
	if err := isNotExistMkDir(dir); err != nil {
		t.Errorf("Failed to create dir: %v", err)
	}
	if checkNotExist(dir) {
		t.Errorf("Dir was not created")
	}
	_ = os.RemoveAll(dir)
}

func TestMkDir(t *testing.T) {
	dir := "test_mkdir"
	_ = os.RemoveAll(dir)
	if err := mkDir(dir); err != nil {
		t.Errorf("mkDir failed: %v", err)
	}
	if checkNotExist(dir) {
		t.Errorf("Dir was not created")
	}
	_ = os.RemoveAll(dir)
}

func TestMustOpen(t *testing.T) {
	dir := "test_mustopen_dir"
	file := "test.log"
	_ = os.RemoveAll(dir)
	f, err := mustOpen(file, dir)
	if err != nil {
		t.Fatalf("mustOpen failed: %v", err)
	}
	if f == nil {
		t.Fatalf("mustOpen returned nil file")
	}
	f.Close()
	// Clean up
	_ = os.RemoveAll(dir)
}
