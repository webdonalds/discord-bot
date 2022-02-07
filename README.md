# WebdoBot on Discord

## 기능

- 타이머
    - !타이머 <시간> <메시지> (예> !타이머 9h 퇴근시간) 
- 택배 조회 / 실시간 위치 알림
    - !택배 <배송사> <운송장번호> (예> !택배 CJ 1234567890123)
    - !택배 <배송사> <운송장번호> <물품이름> (예> !택배 CJ 1234567890123 한우세트)
- 미세먼지 조회
    - !미세먼지
- 환율 조회
    - !환율 <통화쌍> <액수> (예> !환율 USDKRW 4)
    - [사용 가능한 통화쌍](#환율-조회-사용-가능한-통화쌍)
- 블루 아카이브
    - !몰루 카페
    - !몰루 카페 출석
    - !몰루 알림
    - !몰루 알림 <on/off> (예> !몰루 알림 on)
       

## 개발

### 필요사항

- Go (1.14 혹은 그 이상)

### 환경변수

- `DISCORD_BOT_TOKEN` (required)
- `ERROR_LOG_WEBHOOK_URL`
- `BOT_CHANNEL_ID`
- `BREAKING_NEWS_CHANNEL_ID`
- `MOLLU_CHANNEL_ID`
- `DEV_ARTICLE_CHANNEL_ID`
- `TWITTER_CONSUMER_KEY`
- `TWITTER_CONSUMER_SECRET`
- `TWITTER_ACCESS_TOKEN`
- `TWITTER_ACCESS_SECRET`

### 환율 조회 사용 가능한 통화쌍
* 호주 : `AUDCNY, AUDEUR, AUDGBP, AUDJPY, AUDKRW, AUDUSD`
* 브라질 : `BRLCNY, BRLEUR, BRLGBP, BRLJPY, BRLKRW, BRLUSD`
* 캐나다 : `CADCNY, CADEUR, CADGBP, CADJPY, CADKRW, CADUSD`
* 스위스 : `CHFCNY, CHFEUR, CHFGBP, CHFJPY, CHFKRW, CHFUSD`
* 중국 : `CNYAUD, CNYBRL, CNYCAD, CNYCHF, CNYEUR, CNYGBP, CNYHKD, CNYINR, CNYJPY, CNYKRW, CNYMXN, CNYRUB, CNYTHB, CNYTWD, CNYUSD, CNYVND`
* 유로 : `EURAUD, EURBRL, EURCAD, EURCHF, EURCNY, EURGBP, EURHKD, EURINR, EURJPY, EURKRW, EURMXN, EURRUB, EURTHB, EURTWD, EURUSD, EURVND`
* 영국 : `GBPAUD, GBPBRL, GBPCAD, GBPCHF, GBPCNY, GBPEUR, GBPHKD, GBPINR, GBPJPY, GBPKRW, GBPMXN, GBPRUB, GBPTHB, GBPTWD, GBPUSD, GBPVND`
* 홍콩 : `HKDCNY, HKDEUR, HKDGBP, HKDJPY, HKDKRW, HKDUSD, INRCNY`
* 인도 : `INREUR, INRGBP, INRJPY, INRKRW, INRUSD`
* 일본 : `JPYAUD, JPYBRL, JPYCAD, JPYCHF, JPYCNY, JPYEUR, JPYGBP, JPYHKD, JPYINR, JPYKRW, JPYMXN, JPYRUB, JPYTHB, JPYTWD, JPYUSD, JPYVND`
* 한국 : `KRWAUD, KRWBRL, KRWCAD, KRWCHF, KRWCNY, KRWEUR, KRWGBP, KRWHKD, KRWINR, KRWJPY, KRWMXN, KRWRUB, KRWTHB, KRWTWD, KRWUSD, KRWVND`
* 맥시코 : `MXNCNY, MXNEUR, MXNGBP, MXNJPY, MXNKRW, MXNUSD`
* 러시아 : `RUBCNY, RUBEUR, RUBGBP, RUBJPY, RUBKRW, RUBUSD`
* 태국 : `THBCNY, THBEUR, THBGBP, THBJPY, THBKRW, THBUSD`
* 대만 : `TWDCNY, TWDEUR, TWDGBP, TWDJPY, TWDKRW, TWDUSD`
* 미국 : `USDAUD, USDBRL, USDCAD, USDCHF, USDCNY, USDEUR, USDGBP, USDHKD, USDINR, USDJPY, USDKRW, USDMXN, USDRUB, USDTHB, USDTWD, USDVND`
* 베트남 : `VNDCNY, VNDEUR, VNDGBP, VNDJPY, VNDKRW, VNDUSD`
