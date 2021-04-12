package io

import (
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

// CreateDirIfNotExist : create directory if directory does not exist
func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Debug("Creating directory for ", dir)

		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Error("Failed to create directory ", dir)
			return err
		}
	}

	return nil
}

// RemoveFileIfExist: remove the given file if exists and is a regular file
func RemoveFileIfExist(file string) error {
	info, err := os.Stat(file)

	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("%s is a directory", file)
	}

	return os.Remove(file)
}

// RemoveDirIfExist: remove the given directoryc if exists and is a directory
func RemoveDirIfExist(dir string) error {
	info, err := os.Stat(dir)

	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", dir)
	}

	return os.RemoveAll(dir)
}

// Move: move the given source file to the given destination, this function
// deals with copy between partitions
func Move(source string, destination string) error {
	// Try default rename, works on same partition
	err := os.Rename(source, destination)
	if err == nil {
		return nil
	}

	// If default rename fails, perform move copy
	inputFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %s", err)
	}

	outputFile, err := os.Create(destination)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}

	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}

	// The copy was successful, so now delete the original file
	err = os.Remove(source)
	if err != nil {
		return fmt.Errorf("failed removing original file: %s", err)
	}

	return nil
}

// SymLink : create a symbolic link from sourcePath to targetPath
func SymLink(targetPath string, sourcePath string) {
	log.Debug("Creating simbolic link ", sourcePath, " -> ", targetPath)
	err := os.Symlink(sourcePath, targetPath)

	if err != nil {
		log.Error(
			"Unable to create new symlink ", sourcePath, " -> ", targetPath,
			", check if you have the right permissions and the symbolic link does not exists")
		panic(err)
	}
}

// ForceSymLink : create a symbolic link from sourcePath to targetPath, delete if already exists
func ForceSymLink(targetPath string, sourcePath string) {
	if IsSymlink(targetPath) {
		log.Debug("Removing existing simbolic link ", targetPath)
		err := os.Remove(targetPath)

		if err != nil {
			log.Error(
				"Unable to remove symlink ", targetPath,
				", check if you have the right permissions")
			panic(err)
		}
	}

	SymLink(targetPath, sourcePath)
}

// IsSymlink : check file exists and is symlink
func IsSymlink(path string) bool {
	file, err := os.Lstat(path)
	if os.IsNotExist(err) || err != nil {
		log.Debug(err)
		return false
	}

	// Not a symlink
	if err != nil {
		return false
	}

	return file.Mode()&os.ModeSymlink != 0
}

// SetExecutable: set a file as executable, with 0755 permissions
func SetExecutable(file string) error {
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		log.Error("File ", file, " does not exists")
		return err
	}

	if info.IsDir() {
		log.Error("File ", file, " is a directory")
		return fmt.Errorf("file %s is a directory", file)
	}

	return os.Chmod(file, 0755)
}
