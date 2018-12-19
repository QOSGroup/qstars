package slim

import "testing"

func TestEncrypt(t *testing.T) {
	key := "mattwanggoaescbc"
	plaintext := "v/04ti1qAfUnzVfbftGO6sI18WM1Ny1QfyZKA/f2tx62bOWwbasg/wF4OLXP973NEmkapSgODJoFaT5JTmA9lQ=="
	out := AesEncrypt(key, plaintext)
	t.Log(out)
}

func TestDecrypt(t *testing.T) {
	key := "mattwanggoaescbc"
	cipertext := "aNxF2-jrFaW_Tz2dsszvDgHD4vwv0V_YKn83MQIiueZgnrC5EV00f-8PekQjHQn5D0ohfP-UL2wuKuFcVeCLX_VzElcoO8jfJgyXoQr7pMbSPrw6yhA2rTAYMIzoyE1YYP9R1yezUn8="
	out := AesDecrypt(key, cipertext)
	t.Log(out)
}
