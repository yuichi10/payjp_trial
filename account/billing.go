package account

import (
	"D"
	"dbase"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"net/http"
	"strconv"
	"time"
	"github.com/jinzhu/gorm"
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

func TestDB(w http.ResponseWriter, r *http.Request) {
	db := dbase.OpenDB()
	defer db.Close()
	db.DB()
	rF,_ := strTimeToTime("2018-8-14")
	rT,_ := strTimeToTime("2018-8-16")
	c := countOtherOverlapBook(&rF, &rT, 5, db)
	fmt.Println(c)
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
	uID := r.Form.Get(UserID)
	var userID uint64
	var err error
	if userID, err = strconv.ParseUint(uID, 10, 64); err != nil {
		responseError(400, D.InvalidParamKeyMessage, w)
		return
	}
	iID := r.Form.Get(ItemID)
	var itemID uint64
	if itemID, err = strconv.ParseUint(iID, 10, 64); err != nil {
		responseError(400, D.InvalidParamKeyMessage, w)
		return
	}
	rentalFrom, rentalTo := getInputAjustedTimes(rFrom, rTo)
	order := Order{}
	item := getItemInfo(itemID, db)
	if item == nil {
		fmt.Println("itemが見つからない")
		responseError(400, D.InvalidParamKeyMessage, w)
		return
	}
	order.UserID = uint(userID)
	order.RentalFrom = &rentalFrom
	order.RentalTo = &rentalTo
	order.ItemID = item.ID
	order.BasePrice = item.BasePrice
	order.DailyCharge = item.DailyCharge
	order.DepositFee = item.DepositFee
	order.InsurancePrice = int (float64(order.BasePrice) * D.InsurancePriceRatio)
	purePrice := order.calcPureRentalPrice()
	order.ManagementCharge = int (float64 (purePrice) * D.ManegementChargeRatio)
	order.Amount = purePrice + order.InsurancePrice + order.ManagementCharge
	db.Create(&order)
	orderJs, isSuccess := jsonMarshalAndResponseError(order, w)
	if !isSuccess {
		return
	} 
	fmt.Fprintf(w, "%v", orderJs)
	return
}


/**
 * かせるかどうかの日にち判定
 */
func checkRentalDay(from, to *time.Time, itemID string, db *gorm.DB) bool {
	var able bool
	//仮売上のリミットから借りれるかどうかの判断
	if able = checkRentalProvisonLimit(from); !able {
		return able
	}
	//もうすでにその期間借りられてないかどうかのチェック
	//if able = checkDoubleBooking(from, to, itemID, db); !able {
	//	return able
	//}
	//利用日から考えて利用できないかどうかのチェック
	if able = checkRentalDayStart(from); !able {
		return able
	}
	//前後のレンタルの日程を調べてマージンが足りなかった場合予約できないようにする
	return true
}

//レンタルがスタートする人予約できる日の制限をチェックする
func checkRentalDayStart(from *time.Time) bool {
	nowDay := time.Now()
	nowDay = timeToTimeYMD(nowDay)
	canRentalDay := from.AddDate(0, 0, D.CAN_BOOK_DAY_FROM_RENTAL_FROM)
	subTime := canRentalDay.Sub(nowDay)
	fmt.Printf(" today: %v\nrental from : %v\ncanRental: %v\nsubMinutes: %v\n\n", nowDay, from, canRentalDay, subTime.Hours())
	if subTime.Minutes() < 0 {
		return false
	}
	return true
}

//仮売上の日程からチェック
func checkRentalProvisonLimit(from *time.Time) bool {
	nowTime := time.Now()
	nowTime = time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), nowTime.Hour(), nowTime.Minute(), 0, 0, time.UTC)
	//契約できる日かどうか(今はとりあえず仮売上の日にちを超えないようになってるかどうか)
	subDays := calcSubDate(&nowTime, from)
	fmt.Printf("%v : %v 時間差: %v \n", from, nowTime, subDays)
	if subDays > D.DAY_LIMIT {
		return false
	}
	return true
}
/*
//その日にもう借りられてないかどうか
func checkDoubleBooking(tFrom, tTo time.Time, itemID string, db *gorm.DB) bool {
	//始まりか終わりどちらかが利用期間にかかってる
	//SELECT count(*) FROM orders WHERE (item_id=4 AND (status=1 OR status=2)) AND ('2016-8-22' BETWEEN rental_from AND rental_to OR '2016-8-22' BETWEEN rental_from AND rental_to);
	var count int = 0
	marginFrom := tFrom.AddDate(0,0,-D.BOOK_MARGIN_DAYS)
	marginTo := tTo.AddDate(0,0,D.BOOK_MARGIN_DAYS)
	from := timeToStrYMD(marginFrom)
	to := timeToStrYMD(marginTo)
	fmt.Printf("from -> to : %v -> %v \n", from, to)
	//ステータスのsql
	dbState := fmt.Sprintf("(%v=%v)", ORDER_STATUS, STATUS_GET_CONSENT)
	dbWhereTime := fmt.Sprintf("('%v' BETWEEN %v AND %v OR '%v' BETWEEN %v AND %v)", from, RENTAL_FROM, RENTAL_TO, to, RENTAL_FROM, RENTAL_TO)
	dbSql := fmt.Sprintf("SELECT count(*) FROM %v WHERE (%v=%v AND %v) AND %v", ORDER, ITEM_ID, itemID, dbState, dbWhereTime)
	fmt.Printf("sql: %v \n", dbSql)
	res, err := db.Query(dbSql)
	var count1 int
	if err != nil {
		return false
	}
	for res.Next() {
		if err := res.Scan(&count1); err != nil {
			return false
		}
	}
	count += count1
	//レンタルする間に他のレンタルがある場合
	dbSql = fmt.Sprintf("SELECT count(*) FROM %v WHERE (%v=%v AND %v) AND ('%v'<%v AND '%v'>%v)", ORDER, ITEM_ID, itemID, dbState, from, RENTAL_FROM, to, RENTAL_TO)
	fmt.Printf("sql: %v \n", dbSql)
	res, err = db.Query(dbSql)
	var count2 int
	if err != nil {
		return false
	}
	for res.Next() {
		if err := res.Scan(&count2); err != nil {
			return false
		}
	}
	count += count2
	if count == 0 {
		return true
	}
	return false
}
*/


