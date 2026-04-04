package logging_test

import (
	"log/slog"
	"testing"

	"github.com/tfcp-site/httpx/logging"
)

func TestNew_notNil(t *testing.T) {
	if logging.New("mentor", slog.LevelInfo) == nil {
		t.Fatal("New returned nil")
	}
}
