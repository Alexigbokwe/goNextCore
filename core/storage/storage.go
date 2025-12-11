package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// Disk defines the interface for file storage
type Disk interface {
	// Put saves a file to the storage
	Put(file *multipart.FileHeader, folder string) (string, error)
	// PutFile saves a generic reader to storage
	PutFile(name string, content io.Reader, folder string) (string, error)
	// Get retrieves a file's content
	Get(path string) (io.ReadCloser, error)
	// Delete removes a file
	Delete(path string) error
	// Exists checks if a file exists
	Exists(path string) bool
	// URL returns the accessible URL for the file
	URL(path string) string
}

// LocalDisk implementation
type LocalDisk struct {
	Root    string // e.g. "./public"
	BaseURL string // e.g. "http://localhost:5050/public"
}

func NewLocalDisk(root string, baseURL string) *LocalDisk {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		os.MkdirAll(root, 0755)
	}
	return &LocalDisk{Root: root, BaseURL: baseURL}
}

func (d *LocalDisk) Put(file *multipart.FileHeader, folder string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	return d.PutFile(filename, src, folder)
}

func (d *LocalDisk) PutFile(name string, content io.Reader, folder string) (string, error) {
	fullPath := filepath.Join(d.Root, folder)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		os.MkdirAll(fullPath, 0755)
	}

	dstPath := filepath.Join(fullPath, name)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, content); err != nil {
		return "", err
	}

	// Return relative path
	return filepath.Join(folder, name), nil
}

func (d *LocalDisk) Get(path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(d.Root, path)
	return os.Open(fullPath)
}

func (d *LocalDisk) Delete(path string) error {
	fullPath := filepath.Join(d.Root, path)
	return os.Remove(fullPath)
}

func (d *LocalDisk) Exists(path string) bool {
	fullPath := filepath.Join(d.Root, path)
	_, err := os.Stat(fullPath)
	return !os.IsNotExist(err)
}

func (d *LocalDisk) URL(path string) string {
	// Simple join, might need better URL handling
	return fmt.Sprintf("%s/%s", d.BaseURL, path)
}
