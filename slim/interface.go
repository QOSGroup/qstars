package slim

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/QOSGroup/qstars/slim/funcInlocal/bech32local"
	"github.com/QOSGroup/qstars/slim/funcInlocal/ed25519local"
	"github.com/QOSGroup/qstars/slim/funcInlocal/respwrap"
	"github.com/pkg/errors"
	"io/ioutil"
	"math/big"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

//genStdSendTx for the Tx send operation
// NewInt constructs BigInt from int64
func NewInt(n int64) Int {
	return Int{big.NewInt(n)}
}

func (i BigInt) Int64() int64 {
	if !i.i.IsInt64() {
		panic("Int64() out of bound")
	}
	return i.i.Int64()
}

type BigInt struct {
	i *big.Int
}

func add(i *big.Int, i2 *big.Int) *big.Int { return new(big.Int).Add(i, i2) }

// Add adds BigInt from another
func (i BigInt) Add(i2 BigInt) (res BigInt) {
	res = BigInt{add(i.i, i2.i)}
	// Check overflow
	if res.i.BitLen() > 255 {
		panic("BigInt overflow")
	}
	return
}

func (bi BigInt) IsNil() bool {
	return bi.i == nil
}

func (i BigInt) NilToZero() BigInt {
	if i.IsNil() {
		return ZeroInt()
	}
	return i
}

// ZeroInt returns BigInt value with zero
func ZeroInt() BigInt { return BigInt{big.NewInt(0)} }

func (i BigInt) String() string {
	return i.i.String()
}

// MarshalAmino defines custom encoding scheme
func (i BigInt) MarshalAmino() (string, error) {
	if i.i == nil { // Necessary since default Uint initialization has i.i as nil
		i.i = new(big.Int)
	}
	return marshalAmino(i.i)
}

// UnmarshalAmino defines custom decoding scheme
func (i *BigInt) UnmarshalAmino(text string) error {
	if i.i == nil { // Necessary since default BigInt initialization has i.i as nil
		i.i = new(big.Int)
	}
	return unmarshalAmino(i.i, text)
}

// MarshalJSON defines custom encoding scheme
func (i BigInt) MarshalJSON() ([]byte, error) {
	if i.i == nil { // Necessary since default Uint initialization has i.i as nil
		i.i = new(big.Int)
	}
	return marshalJSON(i.i)
}

// UnmarshalJSON defines custom decoding scheme
func (i *BigInt) UnmarshalJSON(bz []byte) error {
	if i.i == nil { // Necessary since default BigInt initialization has i.i as nil
		i.i = new(big.Int)
	}
	return unmarshalJSON(i.i, bz)
}

// MarshalAmino for custom encoding scheme
func marshalAmino(i *big.Int) (string, error) {
	bz, err := i.MarshalText()
	return string(bz), err
}

// UnmarshalAmino for custom decoding scheme
func unmarshalAmino(i *big.Int, text string) (err error) {
	return i.UnmarshalText([]byte(text))
}

// MarshalJSON for custom encoding scheme
// Must be encoded as a string for JSON precision
func marshalJSON(i *big.Int) ([]byte, error) {
	text, err := i.MarshalText()
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(text))
}

// UnmarshalJSON for custom decoding scheme
// Must be encoded as a string for JSON precision
func unmarshalJSON(i *big.Int, bz []byte) error {
	var text string
	err := json.Unmarshal(bz, &text)
	if err != nil {
		return err
	}
	return i.UnmarshalText([]byte(text))
}

// 函数：int64 转化为 []byte
func Int2Byte(in int64) []byte {
	var ret = bytes.NewBuffer([]byte{})
	err := binary.Write(ret, binary.BigEndian, in)
	if err != nil {
		fmt.Printf("Int2Byte error:%s", err.Error())
		return nil
	}

	return ret.Bytes()
}

type BaseCoin struct {
	Name   string `json:"coin_name"`
	Amount BigInt `json:"amount"`
}

type TxStd struct {
	ITx       ITx         `json:"itx"`      //ITx接口，将被具体Tx结构实例化
	Signature []Signature `json:"sigature"` //签名数组
	ChainID   string      `json:"chainid"`  //ChainID: 执行ITx.exec方法的链ID
	MaxGas    BigInt      `json:"maxgas"`   //Gas消耗的最大值
}

func (tx *TxStd) GetSignData() []byte {
	if tx.ITx == nil {
		panic("ITx shouldn't be nil in TxStd.GetSignData()")
		return nil
	}

	ret := tx.ITx.GetSignData()
	ret = append(ret, []byte(tx.ChainID)...)
	ret = append(ret, Int2Byte(tx.MaxGas.Int64())...)

	return ret
}

// 签名：每个签名者外部调用此方法
func (tx *TxStd) SignTx(privkey ed25519local.PrivKey, nonce int64) (signedbyte []byte, err error) {
	if tx.ITx == nil {
		return nil, errors.New("Signature txstd err(itx is nil)")
	}

	sigdata := append(tx.GetSignData(), Int2Byte(nonce)...)
	signedbyte, err = privkey.Sign(sigdata)
	if err != nil {
		return nil, err
	}

	return
}

type ITx interface {
	GetSignData() []byte //获取签名字段
}

//var _ txs.ITx = (*TransferTx)(nil)

type Signature struct {
	Pubkey    ed25519local.PubKey `json:"pubkey"`    //可选
	Signature []byte              `json:"signature"` //签名内容
	Nonce     int64               `json:"nonce"`     //nonce的值
}

// 调用 NewTxStd后，需调用TxStd.SignTx填充TxStd.Signature(每个TxStd.Signer())
func NewTxStd(itx ITx, cid string, mgas BigInt) (rTx *TxStd) {
	rTx = &TxStd{
		itx,
		[]Signature{},
		cid,
		mgas,
	}

	return
}

func genStdSendTx(sendTx ITx, priKey ed25519local.PrivKeyEd25519, chainid string, nonce int64) *TxStd {
	gas := NewBigInt(int64(0))
	stx := NewTxStd(sendTx, chainid, gas)
	signature, _ := stx.SignTx(priKey, nonce)
	stx.Signature = []Signature{Signature{
		Pubkey:    priKey.PubKey(),
		Signature: signature,
		Nonce:     nonce,
	}}

	return stx
}

func getAddrFromBech32(bech32Addr string) (address []byte) {
	//prefix, bz, err := bech32local.DecodeAndConvert(bech32Addr)
	_, bz, _ := bech32local.DecodeAndConvert(bech32Addr)
	//fmt.Printf("the prefix is %s\n", prefix)
	address = bz
	//if prefix != "address" {
	//	return nil, errors.Wrap(err, "Valid Address string should begin with")
	//}
	return
}

type Address []byte

func (add Address) Bytes() []byte {
	return add[:]
}

func (add Address) String() string {
	bech32Addr, err := bech32local.ConvertAndEncode(PREF_ADD, add.Bytes())
	if err != nil {
		panic(err)
	}
	return bech32Addr
}

func (add Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(add.String())
}

// 将Bech32编码的地址Json进行UnMarshal
func (add *Address) UnmarshalJSON(bech32Addr []byte) error {
	var s string
	err := json.Unmarshal(bech32Addr, &s)
	if err != nil {
		return err
	}
	add2 := getAddrFromBech32(s)
	//if err != nil {
	//	return err
	//}
	*add = add2
	return nil
}

type BaseCoins []*BaseCoin
type QSCs = BaseCoins

func (coins BaseCoins) String() string {
	if len(coins) == 0 {
		return ""
	}

	out := ""
	for _, coin := range coins {
		out += fmt.Sprintf("%v,", coin.String())
	}
	return out[:len(out)-1]
}

func (coin *BaseCoin) String() string {
	return fmt.Sprintf("%v%v", coin.Amount, coin.Name)
}

type TransItem struct {
	Address Address `json:"addr"` // 账户地址
	QOS     BigInt  `json:"qos"`  // QOS
	QSCs    QSCs    `json:"qscs"` // QSCs
}

type TransferTx struct {
	Senders   []TransItem `json:"senders"`   // 发送集合
	Receivers []TransItem `json:"receivers"` // 接收集合
}

// 签名字节
func (tx TransferTx) GetSignData() (ret []byte) {
	for _, sender := range tx.Senders {
		ret = append(ret, sender.Address...)
		ret = append(ret, (sender.QOS.NilToZero()).String()...)
		ret = append(ret, sender.QSCs.String()...)
	}
	for _, receiver := range tx.Receivers {
		ret = append(ret, receiver.Address...)
		ret = append(ret, (receiver.QOS.NilToZero()).String()...)
		ret = append(ret, receiver.QSCs.String()...)
	}

	return ret
}

func warpperTransItem(addr Address, coins []BaseCoin) TransItem {
	var ti TransItem
	ti.Address = addr
	ti.QOS = NewBigInt(0)

	for _, coin := range coins {
		if coin.Name == "qos" {
			ti.QOS = ti.QOS.Add(coin.Amount)
		} else {
			ti.QSCs = append(ti.QSCs, &coin)
		}
	}

	return ti
}

// NewTransfer ...
func NewTransfer(sender Address, receiver Address, coin []BaseCoin) ITx {
	var sendTx TransferTx

	sendTx.Senders = append(sendTx.Senders, warpperTransItem(sender, coin))
	sendTx.Receivers = append(sendTx.Receivers, warpperTransItem(receiver, coin))

	return sendTx
}

func (coins Coins) Len() int           { return len(coins) }
func (coins Coins) Less(i, j int) bool { return coins[i].Denom < coins[j].Denom }
func (coins Coins) Swap(i, j int)      { coins[i], coins[j] = coins[j], coins[i] }

var _ sort.Interface = Coins{}

type Coins []Coin

func (coins Coins) Sort() Coins {
	sort.Sort(coins)
	return coins
}

func (coins Coins) IsZero() bool {
	for _, coin := range coins {
		if !coin.IsZero() {
			return false
		}
	}
	return true
}

func (coins Coins) IsValid() bool {
	switch len(coins) {
	case 0:
		return true
	case 1:
		return !coins[0].IsZero()
	default:
		lowDenom := coins[0].Denom
		for _, coin := range coins[1:] {
			if coin.Denom <= lowDenom {
				return false
			}
			if coin.IsZero() {
				return false
			}
			// we compare each coin against the last denom
			lowDenom = coin.Denom
		}
		return true
	}
}

func ParseCoins(coinsStr string) (coins Coins, err error) {
	coinsStr = strings.TrimSpace(coinsStr)
	if len(coinsStr) == 0 {
		return nil, nil
	}

	coinStrs := strings.Split(coinsStr, ",")
	for _, coinStr := range coinStrs {
		coin, err := ParseCoin(coinStr)
		if err != nil {
			return nil, err
		}
		coins = append(coins, coin)
	}

	// Sort coins for determinism.
	coins.Sort()

	// Validate coins before returning.
	if !coins.IsValid() {
		return nil, fmt.Errorf("parseCoins invalid: %#v", coins)
	}

	return coins, nil
}

type Int struct {
	i *big.Int
}

func (i Int) IsZero() bool {
	return i.i.Sign() == 0
}

func (i Int) Int64() int64 {
	if !i.i.IsInt64() {
		panic("Int64() out of bound")
	}
	return i.i.Int64()
}

type Coin struct {
	Denom  string `json:"denom"`
	Amount Int    `json:"amount"`
}

func (coin Coin) IsZero() bool {
	return coin.Amount.IsZero()
}

var (
	// Denominations can be 3 ~ 16 characters long.
	reDnm  = `[[:alpha:]][[:alnum:]]{2,15}`
	reAmt  = `[[:digit:]]+`
	reSpc  = `[[:space:]]*`
	reCoin = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, reAmt, reSpc, reDnm))
)

func ParseCoin(coinStr string) (coin Coin, err error) {
	coinStr = strings.TrimSpace(coinStr)

	matches := reCoin.FindStringSubmatch(coinStr)
	if matches == nil {
		err = fmt.Errorf("invalid coin expression: %s", coinStr)
		return
	}
	denomStr, amountStr := matches[2], matches[1]

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return
	}

	return Coin{denomStr, NewInt(int64(amount))}, nil
}

func NewBigInt(n int64) BigInt {
	return BigInt{big.NewInt(n)}
}

type BaseAccount struct {
	AccountAddress Address             `json:"account_address"` // account address
	Publickey      ed25519local.PubKey `json:"public_key"`      // public key
	Nonce          int64               `json:"nonce"`           // identifies tx_status of an account
}

type QOSAccount struct {
	BaseAccount `json:"base_account"`
	QOS         BigInt `json:"qos"`  // coins in public chain
	QSCs        QSCs   `json:"qscs"` // varied QSCs
}

//only need the following arguments, it`s enough!
func QSCtransferSendStr(addrto, coinstr, privkey, chainid string) string {
	//generate the receiver address, i.e. "addrto" with the following format
	to := getAddrFromBech32(addrto)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//generate the sender address, i.e. the "from" part as the input with privkey in hex string format
	//_, addrben32, priv := utility.PubAddrRetrievalFromAmino(privkey, cmCdc)
	var key ed25519local.PrivKeyEd25519
	ts := "{\"type\": \"tendermint/PrivKeyEd25519\",\"value\": \"" + privkey + "\"}"
	//bz, _ := base64.StdEncoding.DecodeString(privkey)
	//Cdc.MustUnmarshalBinaryBare(bz, &key)
	err := Cdc.UnmarshalJSON([]byte(ts), &key)
	if err != nil {
		fmt.Println(err)
	}
	priv := key
	addrben32, _ := bech32local.ConvertAndEncode(PREF_ADD, key.PubKey().Address().Bytes())
	from := getAddrFromBech32(addrben32)
	//coins generate from input
	var ccs []BaseCoin
	coins, err := ParseCoins(coinstr)
	if err != nil {
		fmt.Println(err)
	}
	for _, coin := range coins {
		ccs = append(ccs, BaseCoin{
			Name:   coin.Denom,
			Amount: NewBigInt(coin.Amount.Int64()),
		})
	}

	//Get "nonce" from the func QSCQueryAccountGet
	AccountStr := QSCQueryAccountGet(addrben32)
	accb := []byte(AccountStr)
	data := respwrap.RPCResponse{}
	err = Cdc.UnmarshalJSON(accb, &data)
	rawresp := data.Result
	acc := QOSAccount{}
	Cdc.UnmarshalJSON(rawresp, &acc)

	//coins check to further improvement
	/*	var qcoins types.Coins
		for _, qsc := range acc.QSCs {
			amount := qsc.Amount
			qcoins = append(qcoins, types.NewCoin(qsc.Name, types.NewInt(amount.Int64())))
		}
		qcoins = append(qcoins, types.NewCoin("qos", types.NewInt(acc.QOS.Int64())))

		if !qcoins.IsGTE(coins) {
			fmt.Println("Address %s doesn't have enough coins to pay for this transaction.", from)
		}
	*/
	var nn int64
	nn = int64(acc.Nonce)
	nn++

	//New transfer for QOS transaction
	t := NewTransfer(from, to, ccs)
	msg := genStdSendTx(t, priv, chainid, nn)
	jasonpayload, err := Cdc.MarshalJSON(msg)
	if err != nil {
		fmt.Println(err)
	}
	datas := bytes.NewBuffer(jasonpayload)
	aurl := Accounturl + "txSend"
	req, _ := http.NewRequest("POST", aurl, datas)
	req.Header.Set("Content-Type", "application/json")
	clt := http.Client{}
	resp, _ := clt.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	output := string(body)
	return output
}

//const QSCResultMapperName = "qstarsResult"
//
//func QOSCommitResultCheck(txhash, height string) string {
//	qstarskey := "heigth:" + height + ",hash:" + txhash
//	d, err := config.GetCLIContext().QSCCliContext.QueryStore([]byte(qstarskey), QSCResultMapperName)
//
//	log.Fatalf("QueryStore: %+v, %+v\n", d, err)
//	if err != nil {
//		return "null"
//	}
//	if d == nil {
//		return "null"
//	}
//	var res []byte
//	err = Cdc.UnmarshalBinaryBare(d, &res)
//	if err != nil {
//		return "null"
//	}
//
//	return string(res)
//}
