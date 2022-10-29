package blacklist

import (
	"sync"
	"time"

	"github.com/elmasy-com/columbus-server/config"
)

var list sync.Map

func Init() {
	list = sync.Map{}
}

func IsBlocked(ip string) bool {

	notbefore, exist := list.Load(ip)
	if !exist {
		return false
	}

	if nb, ok := notbefore.(time.Time); !ok {
		panic("Invalid type in Blacklist")
	} else if time.Now().Before(nb) {
		return true
	}

	// delete expired entries
	list.Delete(ip)

	return false
}

// Block adds ip to the blacklist.
func Block(ip string) {
	list.Store(ip, time.Now().Add(config.BlacklistTime))
}
