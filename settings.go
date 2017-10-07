package settings

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

type Settings struct {
	Data      interface{}
	SavePaths []string
	SavePath  string
}

// Returns new settings object. If no file in paths exists, saves
// supplied data in first option. Otherwise, loads from first path
// where the file exists.
func NewSettings(data interface{}, paths []string) (*Settings, error) {
	set := Settings{
		Data:      data,
		SavePaths: paths,
	}

	for _, filePath := range set.SavePaths {
		path := expandTilde(filePath)

		if _, err := os.Stat(path); !os.IsNotExist(err) {
			set.SavePath = path
		}
	}
	if set.SavePath == "" {
		set.SavePath = expandTilde(set.SavePaths[0])
		if err := set.Save(); err != nil {
			return nil, err
		}
	}

	if err := set.Load(); err != nil {
		return nil, err
	}

	if err := set.Save(); err != nil {
		return nil, err
	}

	return &set, nil
}
func expandTilde(path string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir

	if len(path) > 1 && path[:2] == "~/" {
		path = filepath.Join(dir, path[2:])
	}
	return path
}

func (s *Settings) Save() error {
	data, err := json.MarshalIndent(s.Data, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(s.SavePath, data, 0644)
	return err
}

func (s *Settings) Load() error {
	fileData, err := ioutil.ReadFile(s.SavePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(fileData, s.Data)
	return err
}
