// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (darwin || dragonfly || freebsd || netbsd || openbsd) && !readline_notty
// +build darwin dragonfly freebsd netbsd openbsd
// +build !readline_notty

package readline

import "golang.org/x/sys/unix"

const ioctlReadTermios = unix.TIOCGETA
const ioctlWriteTermios = unix.TIOCSETA
