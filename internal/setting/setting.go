package setting

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"todo_list/internal/app"
)

const (
	KeyRepositoryName       = "REPOSITORY_NAME"
	KeyRepositoryLocalTable = "REPOSITORY_LOCAL_TABLE"

	RepositoryLocal = "local"
	DefaultTable    = "default"
	SettingFileName = "setting.json"
)

type Setting struct {
	appDataDir string
	values     map[string]string
}

var Get = sync.OnceValues(newSetting)

func New(appDataDir string) Setting {
	return Setting{
		appDataDir: appDataDir,
		values: map[string]string{
			KeyRepositoryName:       RepositoryLocal,
			KeyRepositoryLocalTable: DefaultTable,
		},
	}
}

func newSetting() (Setting, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return Setting{}, err
	}
	return newSettingWithAppDataDir(filepath.Join(userConfigDir, app.APP_NAME))
}

func newSettingWithAppDataDir(appDataDir string) (Setting, error) {
	st := New(appDataDir)
	settingFile := settingFilePath(appDataDir)
	if data, err := os.ReadFile(settingFile); err == nil {
		values, err := parseValues(data)
		if err != nil {
			return Setting{}, err
		}
		st.values = values
	} else if !os.IsNotExist(err) {
		return Setting{}, err
	}
	st.fillDefaults()
	return st, st.Save()
}

func (s Setting) AppDataDir() string {
	return s.appDataDir
}

func (s Setting) Get(key string) string {
	return s.values[key]
}

func (s Setting) Set(key string, value string) {
	if s.values == nil {
		s.values = map[string]string{}
	}
	s.values[key] = value
}

func (s Setting) Values() map[string]string {
	values := make(map[string]string, len(s.values))
	for key, value := range s.values {
		values[key] = value
	}
	return values
}

func (s Setting) Save() error {
	if s.appDataDir == "" {
		return errors.New("app data dir is empty")
	}
	if err := os.MkdirAll(s.appDataDir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s.values, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(settingFilePath(s.appDataDir), append(data, '\n'), 0644)
}

func (s Setting) fillDefaults() {
	if s.values == nil {
		s.values = map[string]string{}
	}
	if s.values[KeyRepositoryName] == "" {
		s.values[KeyRepositoryName] = RepositoryLocal
	}
	if s.values[KeyRepositoryLocalTable] == "" {
		s.values[KeyRepositoryLocalTable] = DefaultTable
	}
}

func settingFilePath(appDataDir string) string {
	return filepath.Join(appDataDir, SettingFileName)
}

func parseValues(data []byte) (map[string]string, error) {
	values := map[string]string{}
	if err := json.Unmarshal(data, &values); err == nil {
		return values, nil
	}

	var legacy map[string]any
	if err := json.Unmarshal(data, &legacy); err != nil {
		return nil, err
	}
	if repository, ok := legacy["repository"].(string); ok {
		values[KeyRepositoryName] = repository
	} else if repository, ok := legacy["repository"].(map[string]any); ok {
		if mode, ok := repository["repository_mode"].(string); ok {
			values[KeyRepositoryName] = mode
		}
		if local, ok := repository["local"].(map[string]any); ok {
			if dataFile, ok := local["dataFile"].(string); ok && dataFile != "" {
				values[KeyRepositoryLocalTable] = listNameFromDataFile(dataFile)
			}
		}
	}
	if repositories, ok := legacy["repositories"].(map[string]any); ok {
		if local, ok := repositories[RepositoryLocal].(map[string]any); ok {
			if listName, ok := local["listName"].(string); ok {
				values[KeyRepositoryLocalTable] = listName
			}
		}
	}
	for key, value := range legacy {
		if text, ok := value.(string); ok {
			values[key] = text
		}
	}
	return values, nil
}

func listNameFromDataFile(dataFile string) string {
	base := filepath.Base(dataFile)
	if ext := filepath.Ext(base); ext != "" {
		return base[:len(base)-len(ext)]
	}
	return base
}
