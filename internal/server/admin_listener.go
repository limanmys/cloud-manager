package server

import (
	"context"
	"net"
	"os"
	"os/user"
	"strconv"
	"syscall"

	"github.com/limanmys/cloud-manager/internal/constants"
	"golang.org/x/sys/unix"
)

func Listener() (net.Listener, error) {

	var uid int
	var gid int

	lc := net.ListenConfig{
		Control: func(network, address string, conn syscall.RawConn) error {
			var opErr error
			if err := conn.Control(func(fd uintptr) {
				opErr = syscall.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
			}); err != nil {
				return err
			}
			return opErr
		},
	}
	if _, err := os.Stat(constants.SOCKET_PATH); !os.IsNotExist(err) {
		if err = syscall.Unlink(constants.SOCKET_PATH); err != nil {
			return nil, err
		}
	}
	listener, err := lc.Listen(context.Background(), "unix", constants.SOCKET_PATH)
	if err != nil {
		return nil, err
	}
	if err := os.Chmod(constants.SOCKET_PATH, 0760); err != nil {
		return nil, err
	}
	user_t, err := user.Lookup("cloud-manager")
	if err == nil {
		uid, _ = strconv.Atoi(user_t.Uid)
	}
	group_t, err := user.LookupGroup("cloud-manager")
	if err == nil {
		gid, _ = strconv.Atoi(group_t.Gid)

	}
	if uid != 0 && gid != 0 {
		err = os.Chown(constants.SOCKET_PATH, uid, gid)
		if err != nil {
			return nil, err
		}
	}
	return listener, nil
}
