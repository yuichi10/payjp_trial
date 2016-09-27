package account

import (
	"encoding/json"
	_ "fmt"
	"github.com/bitly/go-simplejson"
	"github.com/payjp/payjp-go/v1"
	"os"
)

const (
	PayjpError     = "error"
	PayjpStatus    = "status"
	PayjpErrorCode = "code"
)

func openPayjp() *payjp.Service {
	return payjp.New(os.Getenv("PAYJP_PRIVATE_KEY"), nil)
}

//トークンを追加
func getPayjpToken(number string, cvc, month, year int) (string, error) {
	payjpService := openPayjp()
	token, err := payjpService.Token.Create(payjp.Card{
		Number:   number,
		CVC:      cvc,
		ExpMonth: month,
		ExpYear:  year,
	})
	if err != nil {
		return "", err
	}
	js, err := json.Marshal(token)
	return string(js), err
}

//payjpにカスタマーを追加
func addPayjpCustomer(token string) (*payjp.CustomerResponse, error) {
	payjpService := openPayjp()
	res, err := payjpService.Customer.Create(payjp.Customer{
		CardToken: token,
	})
	return res, err
}

func isError(js *simplejson.Json) (bool, string) {
	if val := js.Get(PayjpError).Interface(); val != nil {
		code := js.Get(PayjpError).Get(PayjpErrorCode).MustString()
		return true, code
	}
	return false, ""
}
