// Copyright 2018 Lucas Gameiro
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package procfs

import (
	"os"
	"bufio"
	"strings"
	"strconv"
)

// ProcStat provides memory information about the process,
// read from /proc/[pid]/smaps.
type ProcSmap struct {
	// The process ID.
	PID int
	// The command name of a process.
	Comm string
	// The size of the mapping
	Size int
	// The amount of the mapping that is currently resident in RAM (RSS)
	Rss int
	// The process' proportional share of this mapping (PSS)
	Pss int
	// The number of clean pages in the mapping
	Shared_Clean int
	// The number of dirty pages in the mapping
	Shared_Dirty int
	// The number of clean private pages in the mapping
	Private_Clean int
	// The number of dirty private pages in the mapping
	Private_Dirty int
	// The amount of memory currently marked as referenced or accessed
	Referenced int
	// The amount of memory that does not belong to any file
	Anonymous int
	// How much would-be-anonymous memory is also used, but out on swap.
	Swap int

	fs FS
}

func (p Proc) NewSmap() (ProcSmap, error) {
	file, err := os.Open(p.path("smaps"))
	if err != nil {
		return ProcSmap{}, err
	}
	defer file.Close()

	comm, err := p.Comm()
	if err != nil {
		return ProcSmap{}, err
	}

	var ps = ProcSmap{
		PID: p.PID,
		Comm: comm,
		Size: 0,
		Rss: 0,
		Pss: 0,
		Shared_Clean: 0,
		Shared_Dirty: 0,
		Private_Clean: 0,
		Private_Dirty: 0,
		Referenced: 0,
		Anonymous: 0,
		Swap: 0,
		fs: p.fs,
	}
	var st []string
	var value int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		st = strings.Split(scanner.Text(), " ")
		_, st = st[len(st)-1], st[:len(st)-1]
		value, _ = strconv.Atoi(st[len(st)-1])
		switch st[0] {
		case "Size:":
			ps.Size += value
		case "Rss:":
			ps.Rss += value
		case "Pss:":
			ps.Pss += value
		case "Shared_Clean:":
			ps.Shared_Clean += value
		case "Shared_Dirty:":
			ps.Shared_Dirty += value
		case "Private_Clean:":
			ps.Private_Clean += value
		case "Private_Dirty:":
			ps.Private_Dirty += value
		case "Referenced:":
			ps.Referenced += value
		case "Anonymous:":
			ps.Anonymous += value
		case "Swap:":
			ps.Swap += value
		}
	}

	if err := scanner.Err(); err != nil {
		return ProcSmap{}, err
	}

	return ps, nil
}

// ProportionalMemory returns the virtual memory size in bytes.
func (s ProcSmap) ProportionalMemoryBytes() int {
	return s.Pss * 1024
}