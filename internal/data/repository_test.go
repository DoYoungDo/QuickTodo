package data

import (
	"testing"
	"todo_list/internal/setting"
)

func TestNewRepositoryWithSettingUsesLocalRepository(t *testing.T) {
	st := setting.New(t.TempDir())
	repository, err := NewRepositoryWithSetting(st)
	if err != nil {
		t.Fatalf("NewRepositoryWithSetting() error = %v", err)
	}
	if _, ok := repository.(*LocalRepository); !ok {
		t.Fatalf("repository type = %T, want *LocalRepository", repository)
	}
}

func TestNewRepositoryWithSettingFallsBackToLocalRepository(t *testing.T) {
	st := setting.New(t.TempDir())
	st.Set(setting.KeyRepositoryName, "unknown")
	repository, err := NewRepositoryWithSetting(st)
	if err != nil {
		t.Fatalf("NewRepositoryWithSetting() error = %v", err)
	}
	if _, ok := repository.(*LocalRepository); !ok {
		t.Fatalf("repository type = %T, want *LocalRepository", repository)
	}
}
