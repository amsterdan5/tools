package logger

import (
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/petermattis/goid"
)

var traceMap = map[int64]string{}
var lock = sync.RWMutex{}

func GetTraceId() (int64, string) {
	gid := goid.Get()

	traceId, ok := traceMap[gid]
	if !ok {
		traceId = genTraceId()
		SetTraceId(gid, traceId)
	}
	return gid, traceId
}

func SetTraceId(gid int64, traceId string) {
	lock.Lock()
	defer lock.Unlock()

	if gid == 0 {
		gid = goid.Get()
	}

	if traceId == "" {
		delete(traceMap, gid)
		return
	}

	traceMap[gid] = traceId
}

func genTraceId() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}

func DelTraceId() {
	lock.Lock()
	defer lock.Unlock()
	delete(traceMap, goid.Get())
}
