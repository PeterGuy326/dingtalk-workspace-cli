// Copyright 2026 Alibaba Group
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build linux

package helpers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// linuxClockTicksPerSec is USER_HZ, fixed at 100 by the kernel ABI for the
// starttime field of /proc/<pid>/stat regardless of the scheduler tick rate.
const linuxClockTicksPerSec = 100

// processStartUnix returns the OS-reported start time (unix seconds) of pid.
// ok is false when the process does not exist or /proc is unreadable; callers
// must then skip start-time validation rather than infer death.
func processStartUnix(pid int) (int64, bool) {
	stat, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return 0, false
	}
	// comm (field 2) may itself contain spaces and parens; the fixed-format
	// fields resume after the last ')'.
	s := string(stat)
	i := strings.LastIndexByte(s, ')')
	if i < 0 {
		return 0, false
	}
	fields := strings.Fields(s[i+1:])
	// starttime is overall field 22 (ticks since boot); fields[0] here is
	// field 3 (state), so it sits at index 19.
	if len(fields) < 20 {
		return 0, false
	}
	ticks, err := strconv.ParseInt(fields[19], 10, 64)
	if err != nil {
		return 0, false
	}
	btime, ok := linuxBootTimeUnix()
	if !ok {
		return 0, false
	}
	return btime + ticks/linuxClockTicksPerSec, true
}

// linuxBootTimeUnix reads the boot timestamp (btime) from /proc/stat.
func linuxBootTimeUnix() (int64, bool) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, false
	}
	for _, line := range strings.Split(string(data), "\n") {
		if v, ok := strings.CutPrefix(line, "btime "); ok {
			n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err != nil {
				return 0, false
			}
			return n, true
		}
	}
	return 0, false
}
