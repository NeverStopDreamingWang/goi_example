package build

import (
	"testing"

	"goi_example/backend/server"
)

// 调试  go test -run TestStart
func TestStart(t *testing.T) {
	server.Start()
}
