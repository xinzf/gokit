package registry

import (
	"testing"
)

func Test_consulWatcher_getRecurse(t *testing.T) {
	new(consulWatcher).getRecurse("gateway/ding.api.litudai.com", false, 10)

	select {}
}
