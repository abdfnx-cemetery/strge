// +build !windows

package home

// Copyright 2013-2018 Docker, Inc.
// NOTE: this package has originally been copied from github.com/docker/docker.

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/gepis/strge/pkg/unshare"
)

// Key returns the env var name for the user's home dir based on
// the platform being run on
func Key() string {
	return "HOME"
}

// Get returns the home directory of the current user with the help of
// environment variables depending on the target operating system.
// Returned path should be used with "path/filepath" to form new paths.
//
// If linking statically with cgo enabled against glibc, ensure the
// osusergo build tag is used.
//
// If needing to do nss lookups, do not disable cgo or set osusergo.
func Get() string {
	home, _ := unshare.HomeDir()
	return home
}

// GetShortcutString returns the string that is shortcut to user's home directory
// in the native shell of the platform running on.
func GetShortcutString() string {
	return "~"
}

// GetRuntimeDir returns XDG_RUNTIME_DIR.
// XDG_RUNTIME_DIR is typically configured via pam_systemd.
// GetRuntimeDir returns non-nil error if XDG_RUNTIME_DIR is not set.
//
// See also https://standards.freedesktop.org/basedir-spec/latest/ar01s03.html
func GetRuntimeDir() (string, error) {
	if xdgRuntimeDir := os.Getenv("XDG_RUNTIME_DIR"); xdgRuntimeDir != "" {
		return xdgRuntimeDir, nil
	}
	return "", errors.New("could not get XDG_RUNTIME_DIR")
}

// StickRuntimeDirContents sets the sticky bit on files that are under
// XDG_RUNTIME_DIR, so that the files won't be periodically removed by the system.
//
// StickyRuntimeDir returns slice of sticked files.
// StickyRuntimeDir returns nil error if XDG_RUNTIME_DIR is not set.
//
// See also https://standards.freedesktop.org/basedir-spec/latest/ar01s03.html
func StickRuntimeDirContents(files []string) ([]string, error) {
	runtimeDir, err := GetRuntimeDir()
	if err != nil {
		// ignore error if runtimeDir is empty
		return nil, nil
	}
	runtimeDir, err = filepath.Abs(runtimeDir)
	if err != nil {
		return nil, err
	}
	var sticked []string
	for _, f := range files {
		f, err = filepath.Abs(f)
		if err != nil {
			return sticked, err
		}
		if strings.HasPrefix(f, runtimeDir+"/") {
			if err = stick(f); err != nil {
				return sticked, err
			}
			sticked = append(sticked, f)
		}
	}
	return sticked, nil
}

func stick(f string) error {
	st, err := os.Stat(f)
	if err != nil {
		return err
	}
	m := st.Mode()
	m |= os.ModeSticky
	return os.Chmod(f, m)
}

// GetDataHome returns XDG_DATA_HOME.
// GetDataHome returns $HOME/.local/share and nil error if XDG_DATA_HOME is not set.
//
// See also https://standards.freedesktop.org/basedir-spec/latest/ar01s03.html
func GetDataHome() (string, error) {
	if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
		return xdgDataHome, nil
	}
	home := Get()
	if home == "" {
		return "", errors.New("could not get either XDG_DATA_HOME or HOME")
	}
	return filepath.Join(home, ".local", "share"), nil
}

// GetConfigHome returns XDG_CONFIG_HOME.
// GetConfigHome returns $HOME/.config and nil error if XDG_CONFIG_HOME is not set.
//
// See also https://standards.freedesktop.org/basedir-spec/latest/ar01s03.html
func GetConfigHome() (string, error) {
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		return xdgConfigHome, nil
	}
	home := Get()
	if home == "" {
		return "", errors.New("could not get either XDG_CONFIG_HOME or HOME")
	}
	return filepath.Join(home, ".config"), nil
}

// GetCacheHome returns XDG_CACHE_HOME.
// GetCacheHome returns $HOME/.cache and nil error if XDG_CACHE_HOME is not set.
//
// See also https://standards.freedesktop.org/basedir-spec/latest/ar01s03.html
func GetCacheHome() (string, error) {
	if xdgCacheHome := os.Getenv("XDG_CACHE_HOME"); xdgCacheHome != "" {
		return xdgCacheHome, nil
	}
	home := Get()
	if home == "" {
		return "", errors.New("could not get either XDG_CACHE_HOME or HOME")
	}
	return filepath.Join(home, ".cache"), nil
}
