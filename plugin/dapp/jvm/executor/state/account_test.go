package state

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewContractAccount(t *testing.T) {
	contractAccount := NewContractAccount("", nil)
	if nil != contractAccount {
		assert.Equal(t, nil, contractAccount)
	}
}
