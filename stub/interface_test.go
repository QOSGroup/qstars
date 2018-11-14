package stub

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAccountCreateStr(t *testing.T) {
	out := AccountCreateStr()
	require.NotNil(t, out)
}

func TestQSCQueryAccountGet(t *testing.T) {
	//url := "http://localhost:1317/accounts/cosmosaccaddr1nskqcg35k8du3ydhntkcqjxtk254qv8me943mv"
	//out := QSCQueryAccountGet(url)
	//require.NotNil(t, out)
}

func TestAccountRecoverStr(t *testing.T) {
	mncode := "jar dutch hair pluck street legal battle chuckle over hammer fossil material mystery electric during explain spawn aerobic seminar door park artefact resemble recycle"
	out := AccountRecoverStr(mncode)
	require.NotNil(t, out)
}

func TestPubAddrRetrievalStr(t *testing.T) {
	s := "oyiJEECFN3HV5X0XBs/ltdWuR6UEyiNDjGlin17/uCOx/bwCDa9Imw0AF8uXKgMTHhw0Js/z3LnnecebcjHMFMiHCyRf"
	out := PubAddrRetrievalStr(s)
	require.NotNil(t, out)
}
