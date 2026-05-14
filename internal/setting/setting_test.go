package setting

import (
	"encoding/json"
	"os"
	"testing"
)

func TestNewSettingCreatesDefaultSettingFile(t *testing.T) {
	appDataDir := t.TempDir()

	st, err := newSettingWithAppDataDir(appDataDir)
	if err != nil {
		t.Fatalf("newSettingWithAppDataDir() error = %v", err)
	}
	if st.AppDataDir() != appDataDir {
		t.Fatalf("AppDataDir = %q", st.AppDataDir())
	}
	if st.Get(KeyRepositoryName) != RepositoryLocal {
		t.Fatalf("Repository = %q", st.Get(KeyRepositoryName))
	}
	if st.Get(KeyRepositoryLocalTable) != DefaultTable {
		t.Fatalf("Local table = %q", st.Get(KeyRepositoryLocalTable))
	}
	if _, err := os.Stat(settingFilePath(appDataDir)); err != nil {
		t.Fatalf("setting file was not created: %v", err)
	}
}

func TestNewSettingLoadsPersistedSetting(t *testing.T) {
	appDataDir := t.TempDir()
	st := New(appDataDir)
	if err := st.Set(KeyRepositoryLocalTable, "work"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	if err := st.Set("CUSTOM_THEME", "light"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	loaded, err := newSettingWithAppDataDir(appDataDir)
	if err != nil {
		t.Fatalf("newSettingWithAppDataDir() error = %v", err)
	}
	if loaded.Get(KeyRepositoryLocalTable) != "work" {
		t.Fatalf("Local table = %q", loaded.Get(KeyRepositoryLocalTable))
	}
	if loaded.Get("CUSTOM_THEME") != "light" {
		t.Fatalf("CUSTOM_THEME = %q", loaded.Get("CUSTOM_THEME"))
	}
}

func TestNewSettingKeepsDefaultsWhenPersistedSettingIsPartial(t *testing.T) {
	appDataDir := t.TempDir()
	if err := os.MkdirAll(appDataDir, 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(settingFilePath(appDataDir), []byte(`{"CUSTOM":"value"}`), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	st, err := newSettingWithAppDataDir(appDataDir)
	if err != nil {
		t.Fatalf("newSettingWithAppDataDir() error = %v", err)
	}
	if st.Get(KeyRepositoryName) != RepositoryLocal || st.Get(KeyRepositoryLocalTable) != DefaultTable || st.Get("CUSTOM") != "value" {
		t.Fatalf("defaults were not kept: %+v", st.Values())
	}
}

func TestNewSettingReturnsInvalidJSONError(t *testing.T) {
	appDataDir := t.TempDir()
	if err := os.MkdirAll(appDataDir, 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(settingFilePath(appDataDir), []byte("not json"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if _, err := newSettingWithAppDataDir(appDataDir); err == nil {
		t.Fatal("newSettingWithAppDataDir() expected invalid JSON error")
	}
}

func TestSetGetAndValues(t *testing.T) {
	st := New(t.TempDir())
	if err := st.Set("FOO", "bar"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	if st.Get("FOO") != "bar" {
		t.Fatalf("FOO = %q", st.Get("FOO"))
	}
	values := st.Values()
	values["FOO"] = "changed"
	if st.Get("FOO") != "bar" {
		t.Fatal("Values() should return a copy")
	}
}

func TestSetWritesSettingFile(t *testing.T) {
	appDataDir := t.TempDir()
	st := New(appDataDir)
	if err := st.Set("CUSTOM_THEME", "dark"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	data, err := os.ReadFile(settingFilePath(appDataDir))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	var got map[string]string
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if got["CUSTOM_THEME"] != "dark" {
		t.Fatalf("CUSTOM_THEME = %q", got["CUSTOM_THEME"])
	}
}

func TestGetReturnsSingletonSetting(t *testing.T) {
	st, err := Get()
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	again, err := Get()
	if err != nil {
		t.Fatalf("Get() second call error = %v", err)
	}
	if st.AppDataDir() != again.AppDataDir() || st.Get(KeyRepositoryName) != again.Get(KeyRepositoryName) {
		t.Fatalf("Get() returned different settings: first=%+v second=%+v", st.Values(), again.Values())
	}
}
