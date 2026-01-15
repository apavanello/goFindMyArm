package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/apavanello/goFindMyArm/internal/agent/discovery"
	"github.com/apavanello/goFindMyArm/internal/agent/server"
)

// Version injected at build time
var Version = "dev"

func main() {
	// Flags
	portUDP := flag.Int("udp", 32000, "UDP Port for Discovery")
	portTCP := flag.Int("tcp", 32001, "TCP Port for gRPC Commands")
	disableRemote := flag.Bool("disable-remote", false, "Disable Remote Commands (Reboot/Update)")
	password := flag.String("password", "", "Password for Remote Commands (Required if remote enabled)")
	showVersion := flag.Bool("version", false, "Show Version")

	flag.Parse()

	if *showVersion {
		fmt.Printf("goFindMyArm Agent %s\n", Version)
		return
	}

	if !*disableRemote && *password == "" {
		// In production, maybe fail. For dev, warn.
		fmt.Println("[WARN] Remote commands enabled but no password provided! Auth will fail.")
	}

	fmt.Printf("[INFO] Starting Agent %s\n", Version)
	fmt.Printf("[INFO] UDP Discovery: :%d\n", *portUDP)
	if *disableRemote {
		fmt.Println("[INFO] Remote Commands: DISABLED")
	} else {
		fmt.Printf("[INFO] Remote Commands: :%d\n", *portTCP)
	}

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS Signals (SIGINT, SIGTERM)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Block main thread until signal
	go func() {
		sig := <-sigs
		fmt.Printf("\n[INFO] Received signal: %v. Shutting down...\n", sig)
		cancel()
	}()

	// Start Discovery Loop
	go func() {
		if err := discovery.Serve(*portUDP, *password, *disableRemote, Version); err != nil {
			fmt.Printf("[ERROR] Discovery Server failed: %v\n", err)
			cancel()
		}
	}()

	// Start gRPC Server if enabled
	if !*disableRemote {
		go func() {
			if err := server.Serve(*portTCP, *password); err != nil {
				fmt.Printf("[ERROR] gRPC Server failed: %v\n", err)
				cancel()
			}
		}()
	}

	<-ctx.Done()
	fmt.Println("[INFO] Shutdown complete.")
	// Small delay to ensure logs flush if needed
	time.Sleep(500 * time.Millisecond)
}
