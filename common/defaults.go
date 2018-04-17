// Copyright 2016 The go-ethereum Authors
// Copyright 2018 Webchain project
// This file is part of Webchain.
//
// Webchain is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Webchain is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Webchain. If not, see <http://www.gnu.org/licenses/>.

package common

import (
	"path/filepath"
	"runtime"
)

const (
	DefaultIPCSocket = "webchaind.ipc"  // Default (relative) name of the IPC RPC socket
	DefaultHTTPHost  = "localhost" // Default host interface for the HTTP RPC server
	DefaultHTTPPort  = 39573        // Default TCP port for the HTTP RPC server
	DefaultWSHost    = "localhost" // Default host interface for the websocket RPC server
	DefaultWSPort    = 39574        // Default TCP port for the websocket RPC server
)

func defaultDataDirParent() string {
	// Try to place the data folder in the user's home dir
	home := HomeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming")
		} else {
			return filepath.Join(home)
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

func dataDir() string {
	if runtime.GOOS == "darwin" {
		return "Webchain"
	} else if runtime.GOOS == "windows" {
		return "Webchain"
	} else {
		return ".webchain"
	}
}

// DefaultDataDir is the default data directory to use for the databases and other
// persistence requirements.
func DefaultDataDir() string {
	// If the parentDir (os-specific) is available, use that.
	if parentDir := defaultDataDirParent(); parentDir != "" {
		return filepath.Join(parentDir, dataDir())
	} else {
		return parentDir
	}
}
