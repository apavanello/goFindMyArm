package client

import (
	"context"
	"fmt"
	"time"

	"github.com/apavanello/goFindMyArm/internal/client/remote"
	"github.com/apavanello/goFindMyArm/internal/client/scanner"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// ScanNetwork triggers the scanner
func (a *App) ScanNetwork(password string) (scanner.ScanResults, error) {
	// Scan on port 32000 for 2 seconds
	return scanner.ScanNetwork(32000, 2*time.Second, password)
}

// RebootDevice triggers a remote reboot
func (a *App) RebootDevice(ip string, password string) error {
	// Assuming default port 32001 for now
	return remote.Reboot(ip, 32001, password)
}
