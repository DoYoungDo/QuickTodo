package setting

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"todo_list/internal/app"
)

const (
	KeyRepositoryName       = "REPOSITORY_NAME"
	KeyRepositoryLocalTable = "REPOSITORY_LOCAL_TABLE"

	RepositoryLocal = "local"
	DefaultTable    = "default"
	SettingFileName = "setting.json"
	HistoryFileName = "history.json"
	HistoryLimit    = 20
)

type Setting struct {
	appDataDir string
	values     map[string]string
	history    map[string][]string
}

var Get = sync.OnceValues(newSetting)

func New(appDataDir string) Setting {
	return Setting{
		appDataDir: appDataDir,
		values:     defaultValues(),
		history: map[string][]string{
			KeyRepositoryName:       {RepositoryLocal},
			KeyRepositoryLocalTable: {DefaultTable},
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
		values := map[string]string{}
		if err := json.Unmarshal(data, &values); err != nil {
			return Setting{}, err
		}
		for key, value := range values {
			st.set(key, value)
		}
	} else if !os.IsNotExist(err) {
		return Setting{}, err
	}
	historyFile := historyFilePath(appDataDir)
	if data, err := os.ReadFile(historyFile); err == nil {
		history := map[string][]string{}
		if err := json.Unmarshal(data, &history); err != nil {
			return Setting{}, err
		}
		for key, values := range history {
			st.history[key] = slices.Clone(values)
		}
	} else if !os.IsNotExist(err) {
		return Setting{}, err
	}
	return st, st.save()
}

func (s Setting) AppDataDir() string {
	return s.appDataDir
}

func (s Setting) Get(key string) string {
	return s.values[key]
}

func (s Setting) Set(key string, value string) error {
	s.set(key, value)
	s.addHistory(key, value)
	return s.save()
}

func (s Setting) Delete(keys ...string) error {
	for _, key := range keys {
		if value, ok := defaultValues()[key]; ok {
			s.set(key, value)
			s.clearHistory(key)
			continue
		}
		delete(s.values, key)
		delete(s.history, key)
	}
	return s.save()
}

func (s Setting) set(key string, value string) {
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

func (s Setting) History(key string) []string {
	return slices.Clone(s.history[key])
}

func (s Setting) save() error {
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
	if err := os.WriteFile(settingFilePath(s.appDataDir), append(data, '\n'), 0644); err != nil {
		return err
	}
	return s.saveHistory()
}

func (s Setting) addHistory(key string, value string) {
	if s.history == nil {
		s.history = map[string][]string{}
	}
	if slices.Contains(s.history[key], value) {
		return
	}
	s.history[key] = append(s.history[key], value)
	if len(s.history[key]) > HistoryLimit {
		s.history[key] = s.history[key][len(s.history[key])-HistoryLimit:]
	}
}

func (s Setting) clearHistory(key string) {
	if s.history == nil {
		s.history = map[string][]string{}
	}
	s.history[key] = []string{}
}

func (s Setting) saveHistory() error {
	data, err := json.MarshalIndent(s.history, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(historyFilePath(s.appDataDir), append(data, '\n'), 0644)
}

func settingFilePath(appDataDir string) string {
	return filepath.Join(appDataDir, SettingFileName)
}

func historyFilePath(appDataDir string) string {
	return filepath.Join(appDataDir, HistoryFileName)
}

func defaultValues() map[string]string {
	return map[string]string{
		KeyRepositoryName:       RepositoryLocal,
		KeyRepositoryLocalTable: DefaultTable,
	}
}
