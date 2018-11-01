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
	s := "0xa328891040853771d5e57d1706cfe5b5d5ae47a504ca23438c69629f5effb823b1fdbc020daf489b0d0017cb972a03131e1c3426cff3dcb9e779c79b7231cc14c8870b245f"
	out := PubAddrRetrievalStr(s)
	require.NotNil(t, out)
}
