package routing

import (
	"fmt"
	"testing"
)

func TestNecastPeerInfo(t *testing.T) {
	info := NewNePeerInfo(1)
	info.NewMsg()
	info.Confirmation()
	info.ReceivedConfirmation()
	fmt.Println(info.Score())
}
