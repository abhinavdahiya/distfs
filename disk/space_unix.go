// +build darwin dragonfly freebsd linux netbsd openbsd

package disk

import (
	"syscall"

	"github.com/pkg/errors"
)

const (
	RootUnix = "/"
)

var (
	ErrCalcSpace = "Error calculating disk space"
	ErrCalcFree  = "Error calculating free disk space"
)

// Space returns total and free bytes available in `/`.
// Think of it as "df" UNIX command.
func Space() (int64, int64, error) {
	s := syscall.Statfs_t{}
	err = syscall.Statfs(RootUnix, &s)
	if err != nil {
		return -1, -1, errors.Wrap(err, ErrCalcSpace)
	}
	total := s.Bsize * s.Blocks
	free := s.Bsize * s.Bfree
	return total, free, nil
}
