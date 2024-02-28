package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ExtractBinFiles(path string) ([]string, error) {
	var binFiles []string

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".bin") && !strings.HasSuffix(file.Name(), ".0.bin") && !strings.HasSuffix(file.Name(), ".1.bin") {
			binFiles = append(binFiles, filepath.Join(path, file.Name()))
		}
	}

	return binFiles, nil
}

func PathTabCache() (string, error) {
	path, ok := os.LookupEnv("LOCALAPPDATA")
	if !ok {
		return "", fmt.Errorf("could not resolve LOCALAPPDATA")
	}
	path = filepath.Join(path, "Packages", "Microsoft.WindowsNotepad_8wekyb3d8bbwe", "LocalState", "TabState")
	return path, nil
}

func GetTabCacheFiles() ([]string, error) {
	tabCachePath, err := PathTabCache()
	if err != nil {
		return nil, fmt.Errorf("error retrieving cache dir path: %s", err)
	}

	binFiles, err := ExtractBinFiles(tabCachePath)
	if err != nil {
		return nil, fmt.Errorf("error extracting .bin files: %v", err)
	}

	return binFiles, nil
}
