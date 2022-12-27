package version

import (
	"fmt"
	"strconv"
	"strings"
)

type Version uint32

func newVer(major, minor, patch uint8) Version {
	return Version(uint32(major)<<16 | uint32(minor)<<8 | uint32(patch))
}

func newVerFromStr(s string) Version {
	vs := strings.Split(strings.TrimPrefix(s, "v"), ".")
	if len(vs) < 3 {
		return newVer(0, 1, 0)
	}

	major, _ := strconv.Atoi(vs[0])
	minor, _ := strconv.Atoi(vs[1])
	patch, _ := strconv.Atoi(vs[2])

	return newVer(uint8(major), uint8(minor), uint8(patch))
}

// Ints returns (major, minor, patch) versions
func (ve Version) Ints() (uint32, uint32, uint32) {
	v := uint32(ve)
	return (v & majorOnlyMask) >> 16, (v & minorOnlyMask) >> 8, v & patchOnlyMask
}

func (ve Version) String() string {
	vmj, vmi, vp := ve.Ints()
	return fmt.Sprintf("%d.%d.%d", vmj, vmi, vp)
}

// semver versions of the rpc api exposed
var (
	GamblerVersion = newVerFromStr(PROJECT_VERSION)
)

//nolint:varcheck,deadcode
const (
	majorMask = 0xff0000
	minorMask = 0xffff00
	patchMask = 0xffffff

	majorOnlyMask = 0xff0000
	minorOnlyMask = 0x00ff00
	patchOnlyMask = 0x0000ff
)
