// Code generated by cmd/cgo -godefs; DO NOT EDIT.
// cgo -godefs defs_aix.go

// Added for go1.11 compatibility
//go:build aix
// +build aix

package socket

type iovec struct {
	Base *byte
	Len  uint64
}

type msghdr struct {
	Name       *byte
	Namelen    uint32
	Iov        *iovec
	Iovlen     int32
	Control    *byte
	Controllen uint32
	Flags      int32
}

type mmsghdr struct {
	Hdr       msghdr
	Len       uint32
	Pad_cgo_0 [4]byte
}

type cmsghdr struct {
	Len   uint32
	Level int32
	Type  int32
}

type sockaddrInet struct {
	Len    uint8
	Family uint8
	Port   uint16
	Addr   [4]byte /* in_addr */
	Zero   [8]uint8
}

type sockaddrInet6 struct {
	Len      uint8
	Family   uint8
	Port     uint16
	Flowinfo uint32
	Addr     [16]byte /* in6_addr */
	Scope_id uint32
}

const (
	sizeofIovec  = 0x10
	sizeofMsghdr = 0x30

	sizeofSockaddrInet  = 0x10
	sizeofSockaddrInet6 = 0x1c
)
