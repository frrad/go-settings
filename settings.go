package settings

import (
	"encoding/json"
	"io/ioutil"
	"os"
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

	for _, filePath := range paths {
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			set.SavePath = filePath
		}
	}
	if set.SavePath == "" {
		set.SavePath = set.SavePaths[0]
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
