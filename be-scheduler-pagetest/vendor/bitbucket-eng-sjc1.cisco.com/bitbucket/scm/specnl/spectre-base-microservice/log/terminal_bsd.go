// +build darwin freebsd openbsd netbsd dragonfly

package log // spectre

import "syscall"

const ioctlReadTermios = syscall.TIOCGETA

type Termios syscall.Termios
