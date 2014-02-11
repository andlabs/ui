// 11 february 2014
package main

import (
	"fmt"
	"sync"
)

// I need a way to get a sysData for a given HWND or a given HWND/control ID. So, this.

type sdcEntry struct {
	s			*sysData
	members		map[_HMENU]*sysData
}

var (
	sysDatas = map[_HWND]*sdcEntry{}
	sysDatasLock sys.Mutex
)

func addSysData(hwnd _HWND, s *sysData) {
	sysDatasLock.Lock()
	defer sysDatasLock.Unlock()
	sysDatas[hwnd] = &sdcEntry{
		s:			s,
		members:		map[_HMENU]*sysData{},
	}
}

func addIDSysData(hwnd _HWND, id _HMENU, s *sysData) {
	sysDatasLock.Lock()
	defer sysDatasLock.Unlock()
	if ss, ok := sysDatas[hwnd]; ok {
		ss.members[id] = s
	}
	panic(fmt.Sprintf("adding ID %d to nonexistent HWND %d\n", id, hwnd))
}

func getSysData(hwnd _HWND) *sysData {
	sysDatasLock.Lock()
	defer sysDatasLock.Unlock()
	if ss, ok := sysDatas[hwnd]; ok {
		return ss.s
	}
	panic(fmt.Sprintf("getting nonexistent HWND %d\n", hwnd))
}

func getIDSysData(hwnd _HWND, id _HMENU) *sysData {
	sysDatasLock.Lock()
	defer sysDatasLock.Unlock()
	if ss, ok := sysDatas[hwnd]; ok {
		if xx, ok := ss.members[id]; ok {
			return xx
		}
		panic(fmt.Sprintf("getting nonexistent ID %d for HWND %d\n", id, hwnd))
	}
	panic(fmt.Sprintf("getting nonexistent HWND %d\n", hwnd))
}
