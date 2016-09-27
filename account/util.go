package account

import (
    "time"
    "strings"
    "strconv"
    "fmt"
    "net/http"
    "encoding/json"
    "D"
)
/**
 * 2016-02-12 の形を崩す
 * @param  {[type]} allTime string)       (y, m, d int [description]
 * @return {[type]}         [description]
 */
func divideTime(allTime string) (y, m, d int, err error) {
	divTime := strings.Split(allTime, "-")
	y, err = strconv.Atoi(divTime[0])
	if err != nil {
		y, m, d = 0, 0, 0
		return
	}
	m, err = strconv.Atoi(divTime[1])
	if err != nil {
		y, m, d = 0, 0, 0
		return
	}
	d, err = strconv.Atoi(divTime[2])
	if err != nil {
		y, m, d = 0, 0, 0
		return 
	}
	return
}

func checkStrTime(strTime string) bool {
	_, err := strTimeToTime(strTime)
	if err != nil {
		return false
	}
	return true
}

//%v-%v-%v　のtimeを　time.Timeに変換する
func strTimeToTime(strTime string) (time.Time, error){
	y, m, d, err := divideTime(strTime)
	if err != nil {
		return time.Now(), fmt.Errorf("間違った数字が呼ばれました")
	}
	date := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
	if date.Year() == y && date.Month() == time.Month(m) && date.Day() == d {
        return date, nil
    }
    return time.Now(), fmt.Errorf("%d-%d-%d is not exist", y, m, d)
}

//入力された期日をこれから使いやすい形に変える
func getInputAjustedTimes(from, to string) (rentalFrom, rentalTo time.Time) {
	rentalFrom, _ = strTimeToTime(from)
	if to == "" {
		//最後の日が指定してなかった場合
		rentalTo = rentalFrom
		to = from
	} else {
		rentalTo, _ = strTimeToTime(to)
	}
	return
}

//jsonを作成して、jsonのエラーが合ったらjsonのエラーが起きたと返す。json_string似できたらそれを返す
func jsonMarshalAndResponseError(js interface{}, w http.ResponseWriter) (jsStr string, isSuccess bool) {
	isSuccess = false
	resjs, err := json.Marshal(js)
	if err != nil {
		responseServerError(w)
		return "", isSuccess
	}
	jsStr = string(resjs)
	isSuccess = true
	return jsStr, isSuccess
}

//サーバーエラーのレスポンスを返す
func responseServerError(w http.ResponseWriter) {
    w = setResponseJsonHeader(500, w)
	fmt.Fprintf(w, "%v", D.ServerErrorMessage)
}

//エラーのレスポンスを返す
func responseError(httpState int, message string, w http.ResponseWriter) {
	resErr := resError{}
	w = setResponseJsonHeader(httpState, w)
	resErr.Error.ErrorMessage = message
	resErrJs, isSuccess := jsonMarshalAndResponseError(resErr, w)
	if !isSuccess {
		return
	}
	fmt.Fprintf(w, "%v", resErrJs)
}

//成功したかどうかのjsonを返す
func getResponseSuccessResult(isSucceses bool) string {
	if !isSucceses {
		return "{\"response\":{\"is_success\":false}}"
	}
	return "{\"response\":{\"is_success\":true}}"
}

//jsonを返す時のheaderを書く
func setResponseJsonHeader(state int, w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(state)
	return w
}