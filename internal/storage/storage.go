package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func SaveReposPath(path string) error {
	data, err := json.MarshalIndent([]string{path}, "", "  ")
	if err != nil {
		return err
	}
	vinuxDir := path + "/vinux/repos.json"
	return os.WriteFile(
		vinuxDir,
		data,
		0644,
	)
}

func GetRepo(name string) (string, error) {
	reposDir := GetReposDir()
	if !VerifyPathExists(reposDir) {
		return "", os.ErrNotExist
	}

	return filepath.Join(reposDir, "/"+name), nil
}

func ListRepos() ([]string, error) {
	return loadAll()
}

func GetReposDir() string {
	data, err := os.ReadFile(getPathToStorage() + "/repos.json")
	if err != nil {
		if os.IsNotExist(err) {
			return ""
		}
		return ""
	}

	var repos []string
	if len(data) > 0 {
		if err := json.Unmarshal(data, &repos); err != nil {
			return ""
		}
	}

	return repos[0]
}

func VerifyPathExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func getPathToStorage() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	dir := filepath.Dir(exePath)
	return dir
}

func loadAll() ([]string, error) {
	entries, err := os.ReadDir(GetReposDir())
	if err != nil {
		return nil, err
	}

	var folders []string
	for _, entry := range entries {
		if entry.IsDir() {
			folders = append(folders, entry.Name())
		}
	}

	return folders, nil
}
