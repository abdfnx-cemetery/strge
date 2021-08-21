package constants

import "runtime"

// dupe = Default Unix Path Env

const dupe = "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"

func DefaultPathEnv(platform string) string {
	if runtime.GOOS == "windows" {
		if platform != runtime.GOOS && LCOWSupported() {
			return dupe
		}

		return ""
	}

	return dupe
}
