package constants

import "syscall"

// fromStatT converts a syscall.Stat_t type to a constants.Stat_t type
func fromStatT(s *syscall.Stat_t) (*StatT, error) {
	return &StatT{size: s.Size,
		mode: s.Mode,
		uid:  s.Uid,
		gid:  s.Gid,
		rdev: uint64(s.Rdev),
		mtim: s.Mtim}, nil
}

// FromStatT converts a syscall.Stat_t type to a constants.Stat_t type
// This is exposed on Linux as pkg/archive/changes uses it.
func FromStatT(s *syscall.Stat_t) (*StatT, error) {
	return fromStatT(s)
}
