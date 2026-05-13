package app

import "testing"

func TestAppVersionDefault(t *testing.T) {
	if APP_VERSION == "" {
		t.Fatal("APP_VERSION should not be empty")
	}
}
