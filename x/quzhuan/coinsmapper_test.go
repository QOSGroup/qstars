package quzhuan

import (
	"encoding/json"
	"fmt"
	"github.com/QOSGroup/qbase/types"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
}

type ScenesTx struct {
	UserId string //场景唯一标识
	Amount types.BigInt
}

type Test struct {
	ScenesId string
	Rewards  []ScenesTx
}

func TestInitCmd(t *testing.T) {
	aaa := ScenesTx{"13388888888", types.NewInt(888)}
	bbb := ScenesTx{"13399999999", types.NewInt(999)}

	ccc := []ScenesTx{aaa, bbb}

	ddd := Test{"001", ccc}

	result, _ := json.Marshal(ddd)

	strs:=[]string{string(result)}

	fmt.Println(strs)

	eeee,_:=json.Marshal(strs)


	//["{\"ScenesId\":\"001\",\"Rewards\":[{\"UserId\":\"13388888888\",\"Amount\":\"888\"},{\"UserId\":\"13399999999\",\"Amount\":\"999\"}]}"]

	//[{"ScenesId":"001","Rewards":[{"UserId":"13388888888","Amount":"888"},{"UserId":"13399999999","Amount":"999"}]}]

	fmt.Println(string(eeee))

}
