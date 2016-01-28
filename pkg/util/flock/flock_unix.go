// +build linux darwin freebsd openbsd netbsd dragonfly

/*
Copyright 2016 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package flock

import (
	"fmt"
	"os"
	"strconv"

	"golang.org/x/sys/unix"
)

// os.File has a runtime.Finalizer so the fd will be closed if the struct is
// garbage collected. Let's hold onto a reference so that doesn't happen.
var lockfile *os.File

// Acquire aquires a lock on a file for the duration of the process and writes its
// pid to the file. This method is reentrant but not threadsafe.
func Acquire(path string) error {
	var err error
	if lockfile, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666); err != nil {
		return err
	}

	var lock unix.Flock_t
	lock.Type = unix.F_WRLCK

	if err := unix.FcntlFlock(lockfile.Fd(), unix.F_SETLKW, &lock); err != nil {
		return err
	}
	if _, err := lockfile.Write([]byte(strconv.Itoa(os.Getpid()))); err != nil {
		return fmt.Errorf("failed to write pid to lock file: %v", err)
	}
	return nil
}
