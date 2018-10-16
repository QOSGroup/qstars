package bank

import (
	sdk "github.com/QOSGroup/qstars/types"
	"encoding/json"
)
//----------------------------------------
// Input

// Transaction Input
type Input struct {
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

// Return bytes to sign for Input
func (in Input) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(in)
	if err != nil {
		panic(err)
	}

	return MustSortJSON(bin)
}

// ValidateBasic - validate transaction input
func (in Input) ValidateBasic() sdk.Error {
	if len(in.Address) == 0 {
		return sdk.ErrInvalidAddress(in.Address.String())
	}
	if !in.Coins.IsValid() {
		return sdk.ErrInvalidCoins(in.Coins.String())
	}
	if !in.Coins.IsPositive() {
		return sdk.ErrInvalidCoins(in.Coins.String())
	}
	return nil
}

// NewInput - create a transaction input, used with MsgSend
func NewInput(addr sdk.AccAddress, coins sdk.Coins) Input {
	input := Input{
		Address: addr,
		Coins:   coins,
	}
	return input
}

//----------------------------------------
// Output

// Transaction Output
type Output struct {
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

// Return bytes to sign for Output
func (out Output) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(out)
	if err != nil {
		panic(err)
	}
	return MustSortJSON(bin)
}

// ValidateBasic - validate transaction output
func (out Output) ValidateBasic() sdk.Error {
	if len(out.Address) == 0 {
		return sdk.ErrInvalidAddress(out.Address.String())
	}
	if !out.Coins.IsValid() {
		return sdk.ErrInvalidCoins(out.Coins.String())
	}
	if !out.Coins.IsPositive() {
		return sdk.ErrInvalidCoins(out.Coins.String())
	}
	return nil
}

// NewOutput - create a transaction output, used with MsgSend
func NewOutput(addr sdk.AccAddress, coins sdk.Coins) Output {
	output := Output{
		Address: addr,
		Coins:   coins,
	}
	return output
}

func SortJSON(toSortJSON []byte) ([]byte, error) {
	var c interface{}
	err := json.Unmarshal(toSortJSON, &c)
	if err != nil {
		return nil, err
	}
	js, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return js, nil
}

// MustSortJSON is like SortJSON but panic if an error occurs, e.g., if
// the passed JSON isn't valid.
func MustSortJSON(toSortJSON []byte) []byte {
	js, err := SortJSON(toSortJSON)
	if err != nil {
		panic(err)
	}
	return js
}