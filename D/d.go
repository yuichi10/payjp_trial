package D

const (
    //エラーメッセージ
    ServerErrorMessage = "server_error"
    ProcessingErrorMessage = "processing_error"
    InvalidParamKeyMessage = "invalid_param_key"
    JSONErrorMessage = "json_error"

    //予約時のエラー
    BookingDaysErrorMessage = "book_day_error"

    //ブッキングの決まり日数
    CAN_BOOK_DAY_FROM_RENTAL_FROM = -4
    DAY_LIMIT = 58
    BOOK_MARGIN_DAYS = 3

    //料金の設定
    CancelRatioBefore4Day = 0.3 
    CancelRatioBefore3Day = 0.3
    CancelRatioBefore2Day = 0.4
    CancelRatioBefore1Day = 0.5
    ManegementChargeRatio = 0.1
    InsurancePriceRatio = 0.1
)

const (
	//オーダーのキャンセルのステータス
	STATUS_FAILED_DELAY_CANCEL = -3 //遅延によるキャンセルの失敗
	STATUS_FAILED_CANCEL_PAY   = -2 //有料キャンセルの失敗
	STATUS_FAILED_CANCEL_FREE  = -1 //無料キャンセルの失敗
	ORDER_STATE_CANCEL_NONE    = 0  //キャンセル無し
	ORDER_STATE_CANCEL_FREE    = 1  //無料のキャンセル
	ORDER_STATE_CANCEL_PAID    = 2  //有料のキャンセル
	ORDER_STATE_CANCEL_DELAY   = 3  //遅延キャンセル
)

const (
	//オーダーステータス
	STATUS_CONTINUE_DELAY_FAILED        = -5 //遅延続行の実売上取得を失敗した時
	STATUS_FAILED_REAL_SALE             = -4 //実売上の取得に失敗
	STATUS_FAILED_CONSENT_PAY_BACK      = -3 //オーダーがキャンセルされた時に仮売上をキャンセル失敗した時
	STATUS_FAILED_CONSENT               = -2 //オーダーが同意されなかった時
	STATUS_FAILED_PROVISION_SALE        = -1 //仮売上が取れなかった時
	STATUS_GET_PROVISION_SALE           = 1  //仮売上をとった
	STATUS_GET_CONSENT                  = 2  //同意が取れた時
	STATUS_CANCEL                       = 3  //キャンセルされた時
	STATUS_GET_REAL_SALE                = 4  //実売上をとった
	STATUS_WRITE_ON_CSV                 = 5  //CSVに書き出し
	STATUS_CONTINUE_DELAY               = 6  //遅延キャンセルをキャンセルして続ける場合
	STATUS_CONTINUE_DELAY_GET_REAL_SALE = 7  //遅延キャンセルをして実売上をとった場合
	STATUS_FINISH                       = 99 //すべての工程を終了
)
