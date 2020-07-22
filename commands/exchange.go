package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/webdonalds/discord-bot/background"
)

type ExchangeCommand struct{}

const exchageHelpMsg = "사용법: !환율 <통화쌍> <액수>\n예시: !환율 USDKRW 4\n사용 가능한 통화쌍 리스트는 깃헙 페이지 참조"

const availableList = `
AUDCNY AUDEUR AUDGBP AUDJPY AUDKRW AUDUSD
BRLCNY BRLEUR BRLGBP BRLJPY BRLKRW BRLUSD
CADCNY CADEUR CADGBP CADJPY CADKRW CADUSD
CHFCNY CHFEUR CHFGBP CHFJPY CHFKRW CHFUSD
CNYAUD CNYBRL CNYCAD CNYCHF CNYEUR CNYGBP CNYHKD CNYINR CNYJPY CNYKRW CNYMXN CNYRUB CNYTHB CNYTWD CNYUSD CNYVND
EURAUD EURBRL EURCAD EURCHF EURCNY EURGBP EURHKD EURINR EURJPY EURKRW EURMXN EURRUB EURTHB EURTWD EURUSD EURVND
GBPAUD GBPBRL GBPCAD GBPCHF GBPCNY GBPEUR GBPHKD GBPINR GBPJPY GBPKRW GBPMXN GBPRUB GBPTHB GBPTWD GBPUSD GBPVND
HKDCNY HKDEUR HKDGBP HKDJPY HKDKRW HKDUSD INRCNY
INREUR INRGBP INRJPY INRKRW INRUSD
JPYAUD JPYBRL JPYCAD JPYCHF JPYCNY JPYEUR JPYGBP JPYHKD JPYINR JPYKRW JPYMXN JPYRUB JPYTHB JPYTWD JPYUSD JPYVND
KRWAUD KRWBRL KRWCAD KRWCHF KRWCNY KRWEUR KRWGBP KRWHKD KRWINR KRWJPY KRWMXN KRWRUB KRWTHB KRWTWD KRWUSD KRWVND
MXNCNY MXNEUR MXNGBP MXNJPY MXNKRW MXNUSD
RUBCNY RUBEUR RUBGBP RUBJPY RUBKRW RUBUSD
THBCNY THBEUR THBGBP THBJPY THBKRW THBUSD
TWDCNY TWDEUR TWDGBP TWDJPY TWDKRW TWDUSD
USDAUD USDBRL USDCAD USDCHF USDCNY USDEUR USDGBP USDHKD USDINR USDJPY USDKRW USDMXN USDRUB USDTHB USDTWD USDVND
VNDCNY VNDEUR VNDGBP VNDJPY VNDKRW VNDUSD
`

const api = "https://earthquake.kr:23490/query/" // 참고: https://jaeheon.kr/12

func NewExchangeCommand() Command {
	return &ExchangeCommand{}
}

func (*ExchangeCommand) CommandTexts() []string {
	return []string{"환율"}
}

func (*ExchangeCommand) Execute(args []string, _ *discordgo.MessageCreate) (string, background.Watcher, error) {
	// 예외 처리
	if len(args) != 2 {
		return exchageHelpMsg, nil, nil
	}

	//통화쌍
	currencyPair := args[0]
	orig_price := args[1]

	if strings.Contains(availableList, currencyPair) == false {
		return "지원하지 않는 국가 코드입니다.", nil, nil
	}

	targetPrice, parseErr := strconv.ParseFloat(args[1], 64) // 환율 적용 하고 싶은 액수

	if parseErr != nil {
		return "금액은 숫자로만 적어주세요.", nil, nil
	}

	// API 읽어오기
	url := api + currencyPair

	req, err := http.Get(url)

	if err != nil {
		return "환율 API 에러: ", nil, err
	}

	defer req.Body.Close()

	data, err := ioutil.ReadAll(req.Body)

	if err != nil {
		return "환율 API Read 에러", nil, err
	}

	// Json 파싱
	var res map[string]interface{}
	jsonErr := json.Unmarshal([]byte(data), &res)

	if jsonErr != nil {
		return "환율 json 처리 에러", nil, jsonErr
	}
	exchangeData := res[currencyPair].([]interface{})

	// 통화 단위
	orig := currencyPair[0:3]
	target := currencyPair[3:6]

	// json 데이터 캐싱
	price := exchangeData[0]
	change := exchangeData[1]
	changePercent := exchangeData[2]

	actualPrice := price.(float64) * targetPrice // 환율 적용한 액수

	// 결과물 출력
	msg := fmt.Sprintf("%s: %s %s 는 %s %s 입니다. (1 %s = %.2f %s)\n전일 대비 변화량: %.1f %s.\n전일 대비 변화량(%%): %.1f%%.",
		currencyPair, orig_price, orig, strconv.FormatFloat(actualPrice, 'f', -1, 32), target,
		orig, price, target,
		change, target, changePercent)

	return msg, nil, nil
}
