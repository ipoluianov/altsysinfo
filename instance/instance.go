// Package instance keeps one copy of the application running per user: the
// settings and the window layout are written by one process only. A second
// start asks the running copy to show its window and quits.
package instance

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	lockFileName = "instance.lock"

	// showCommand is what a second start sends to the running copy
	showCommand = "show"
)

// dialTimeout is how long a second start tries to reach the running copy,
// which may be starting up itself
const dialTimeout = 3 * time.Second

// Instance is the running copy: it holds the lock while it works with the files
type Instance struct {
	lock     *os.File
	listener net.Listener
}

// Acquire makes this process the running copy. When another copy is
// running, it is asked to show its window and ok is false: this one must quit.
// When the lock cannot be taken at all (e.g. no access to the directory),
// the application runs unprotected rather than not at all.
func Acquire(dir string) (c *Instance, ok bool) {
	c = &Instance{}
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Println("Instance lock error:", err)
		return c, true
	}
	lock, err := os.OpenFile(filepath.Join(dir, lockFileName), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Instance lock error:", err)
		return c, true
	}
	if !tryLock(lock) {
		lock.Close()
		if err := askToShow(socketPath(dir)); err != nil {
			fmt.Println("AltSysInfo is already running, but does not answer:", err)
		}
		return nil, false
	}
	c.lock = lock

	// A socket left by a copy that crashed: the lock shows no one uses it
	socket := socketPath(dir)
	os.Remove(socket)
	c.listener, err = net.Listen("unix", socket)
	if err != nil {
		// Still the only copy, just a second start cannot show it
		fmt.Println("Instance socket error:", err)
		c.listener = nil
	}
	return c, true
}

// Serve calls onShow from its own goroutine each time the application is started again
func (c *Instance) Serve(onShow func()) {
	if c == nil || c.listener == nil {
		return
	}
	go func() {
		for {
			conn, err := c.listener.Accept()
			if err != nil {
				return // closed
			}
			conn.SetReadDeadline(time.Now().Add(time.Second))
			line, _ := bufio.NewReader(conn).ReadString('\n')
			conn.Close()
			if strings.TrimSpace(line) == showCommand {
				onShow()
			}
		}
	}()
}

// Close releases the lock; call it when the application has stopped writing its files
func (c *Instance) Close() {
	if c == nil {
		return
	}
	if c.listener != nil {
		c.listener.Close() // removes the socket file
		c.listener = nil
	}
	if c.lock != nil {
		// The lock file stays: removing it could let two copies lock different files
		c.lock.Close()
		c.lock = nil
	}
}

// socketPath returns where the running copy of the directory listens. A unix
// socket path is limited to about 100 bytes, which a long home path can exceed,
// so the short runtime directory is used when there is one (Linux). The name
// tells the directories apart: each has its own running copy.
func socketPath(dir string) string {
	if runtimeDir := os.Getenv("XDG_RUNTIME_DIR"); runtimeDir != "" {
		sum := sha256.Sum256([]byte(dir))
		return filepath.Join(runtimeDir, "altsysinfo-"+hex.EncodeToString(sum[:8])+".sock")
	}
	return filepath.Join(dir, "instance.sock")
}

// askToShow asks the running copy to show its window
func askToShow(socketPath string) error {
	// Let the running copy take the focus, which Windows gives only to the foreground process
	allowForeground()
	deadline := time.Now().Add(dialTimeout)
	for {
		conn, err := net.DialTimeout("unix", socketPath, time.Second)
		if err == nil {
			defer conn.Close()
			conn.SetWriteDeadline(time.Now().Add(time.Second))
			_, err = conn.Write([]byte(showCommand + "\n"))
			return err
		}
		if time.Now().After(deadline) {
			return err
		}
		time.Sleep(100 * time.Millisecond)
	}
}
