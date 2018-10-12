package stub

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestAccountCreateStr(t *testing.T) {
	out:=AccountCreateStr()
	require.NotNil(t,out)
}

func TestQSCQueryAccountGet(t *testing.T) {
	url := "http://localhost:1317/accounts/cosmosaccaddr1nskqcg35k8du3ydhntkcqjxtk254qv8me943mv"
	out := QSCQueryAccountGet(url)
	require.NotNil(t,out)
}