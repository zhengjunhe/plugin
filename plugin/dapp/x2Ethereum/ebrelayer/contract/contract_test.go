package contract

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestLoadABI : test that contract containing named event is successfully loaded
func TestLoadABI(t *testing.T) {
	//Get the ABI ready
	abi := LoadABI(true)

    event := "LogNewProphecyClaim"
	//fmt.Fprintln(os.Stdout, "events is", abi.Events)
	if _, ok := abi.Events["hh"]; !ok {
		t.Fatalf("event:%s doesn't existed", event)
		//panic("not exist")
	}

	require.NotNil(t, abi.Events["hh"])
}
