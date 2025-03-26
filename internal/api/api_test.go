package api

import (
	"testing"

	"github.com/mehmetalisavas/message-sender/config"
)

func TestNew(t *testing.T) {
	cfg := &config.Config{}

	apiInstance := New(cfg, nil)

	if apiInstance == nil {
		t.Errorf("expected apiInstance to be non-nil")
	}

}
