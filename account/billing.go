package account

import (
	"D"
	"dbase"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"net/http"
	"strconv"
)

type resError struct {
	Error struct {
		ErrorMessage string `json:"error_message"`
	} `json:"error"`
}

type resToken struct {
	Token string `json:"token"`
}

const jsonErrorResponse = "{\"error\":{\"error_message\":\"invalid_param_key\"}}"

func TestDB() {
	db := dbase.OpenDB()
	defer db.Close()
	db.DB()
	orderType := new(Order)
	db.First(orderType)
	fmt.Println(orderType)
}

const (
	Number   = "card_number"
	CVC      = "card_cvc"
	ExpMonth = "card_month"
	ExpYear  = "card_year"
	TokenID  = "ID"
	Card     = "Card"

	//getPayjpToken
	PayjpToken = "token"

)

//Tokenの発行（本番では各アプリが実行、テストでユーザー登録に必要)
func GetToken(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var err error
	//値の取得
	number := r.Form.Get(Number)
	cvcStr := r.Form.Get(CVC)
	var cvc int
	if cvc, err = strconv.Atoi(cvcStr); err != nil {
		responseError(400, D.InvalidParamKeyMessage, w)
		return
	}
	expMonthStr := r.Form.Get(ExpMonth)
	var expMonth int
	if expMonth, err = strconv.Atoi(expMonthStr); err != nil {
		responseError(400, D.InvalidParamKeyMessage, w)
		return
	}
	expYearStr := r.Form.Get(ExpYear)
	var expYear int
	if expYear, err = strconv.Atoi(expYearStr); err != nil {
		responseError(400, D.InvalidParamKeyMessage, w)
		return
	}

	//tokenの取得
	token, err := getPayjpToken(number, cvc, expMonth, expYear)
	if err != nil {
		responseError(200, D.ProcessingErrorMessage, w)
		return
	}
	var jsToken *simplejson.Json
	if jsToken, err = simplejson.NewJson([]byte(token)); err != nil {
		//json化でエラーが起こったとき
		responseServerError(w)
		return
	}
	if isErr, code := isError(jsToken); isErr {
		//payjpで何かしらエラーが起こった場合
		responseError(200, code, w)
		return
	}
	//成功時
	res := resToken{}
	res.Token = jsToken.Get(TokenID).MustString()
	resJs, err := json.Marshal(res)
	if err != nil {
		//jsonによるエラー
		responseServerError(w)
		return
	}
	w = setResponseJsonHeader(200, w)
	fmt.Fprintf(w, "%v", string(resJs))
	return
}

//クレジットカードの登録(新規作成)
func AddCardInfo(w http.ResponseWriter, r *http.Request) {
	db := dbase.OpenDB()
	defer db.Close()
	r.ParseForm()
	token := r.Form.Get(PayjpToken)
	userID := r.Form.Get(UserID)
	cusRes, err := addPayjpCustomer(token)
	if err != nil {
		responseError(400, D.ProcessingErrorMessage, w)
		return
	}
	//ユーザーのクレジットカスタマーIDを追加
	user := GetUserInfo(userID, db)
	db.Model(&user).Update(CustomerID, cusRes.ID)
	resJsStr := getResponseSuccessResult(true)
	fmt.Fprintf(w, "%v", string(resJsStr))
}

func PublishOrder(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	db := dbase.OpenDB()
	defer db.Close()
	//order := Order{}
	rTo := r.Form.Get(RentalTo)
	if !checkStrTime(rTo) && rTo != ""{
		//レンタル終了日のチェック
		responseError(400, D.InvalidParamKeyMessage, w)
		return
	}
	rFrom := r.Form.Get(RentalFrom)
	if !checkStrTime(rFrom) {
		//レンタル開始日のチェック
		responseError(400, D.InvalidParamKeyMessage, w)
		return
	}
	//rentalFrom, rentalTo := getInputAjustedTimes(rFrom, rTo)
}



