package broadcaster

import (
	"sync"

	"github.com/TF2Stadium/wsevent"
)

var steamIdSocketMap = make(map[string]*wsevent.Client)
var steamIdSocketMapLock = new(sync.RWMutex)

func SetSocket(steamid string, so *wsevent.Client) {
	steamIdSocketMapLock.Lock()
	defer steamIdSocketMapLock.Unlock()

	steamIdSocketMap[steamid] = so
}

func RemoveSocket(steamid string) {
	steamIdSocketMapLock.Lock()
	defer steamIdSocketMapLock.Unlock()

	delete(steamIdSocketMap, steamid)
}

func GetSocket(steamid string) (so *wsevent.Client, success bool) {
	steamIdSocketMapLock.RLock()
	defer steamIdSocketMapLock.RUnlock()

	so, success = steamIdSocketMap[steamid]
	return
}

func IsConnected(steamid string) bool {
	_, ok := GetSocket(steamid)
	return ok
}
