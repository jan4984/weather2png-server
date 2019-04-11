package weather2png

import (
	"github.com/oliveagle/jsonpath"
	"os"
	"net/http"
	"encoding/json"
	"strings"
	"errors"
	"strconv"
	"io/ioutil"
	"log"
)

type WeatherInfo struct {
	Live, Tomorrow, Houtian, Dahoutian, Dadahoutian struct {
		Aqi         string
		Weather     string
		Temperature string
		Date        string
	}
	Today struct {
		Wind    string
		Weather string
	}
	Err string
}

func FetchWeather(city string) WeatherInfo {
	retErr := func(err error) WeatherInfo {
		panic(err)
		return WeatherInfo{
			Err: err.Error(),
		}
	}
	apikey := os.Getenv("JUHE_APIKEY")
	if apikey == "" {
		return retErr(errors.New("没有定义JUHE API KEY"))
	}
	rsp, err := http.Get("http://apis.juhe.cn/simpleWeather/query?key=" + apikey + "&city=" + city)
	if err != nil {
		return retErr(err)
	}

	defer rsp.Body.Close()
	rspTxt, _ := ioutil.ReadAll(rsp.Body)
	log.Println("rsp:", string(rspTxt))
	var doc interface{}
	err = json.Unmarshal(rspTxt, &doc)
	if err != nil {
		return retErr(err)
	}

	v, err := jsonpath.JsonPathLookup(doc, "$.reason")
	if err != nil {
		return retErr(err)
	}

	if vv, ok := v.(string); ok != true || !strings.Contains(vv, "查询成功") {
		return retErr(errors.New("查询结果：" + vv))
	}

	wi := WeatherInfo{}

	v, _ = jsonpath.JsonPathLookup(doc, "$.result.realtime.temperature")
	wi.Live.Temperature = v.(string)
	v, _ = jsonpath.JsonPathLookup(doc, "$.result.realtime.info")
	wi.Live.Weather = v.(string)
	v, _ = jsonpath.JsonPathLookup(doc, "$.result.realtime.aqi")
	vv, _ := strconv.ParseInt(v.(string), 10, 32)
	if (vv <= 50) {
		wi.Live.Aqi = "空气优"
	} else if (vv < 100) {
		wi.Live.Aqi = "空气良"
	} else if (vv < 200) {
		wi.Live.Aqi = "轻度污染"
	} else if (vv < 300) {
		wi.Live.Aqi = "中度污染"
	} else {
		wi.Live.Aqi = "重度污染"
	}
	v, _ = jsonpath.JsonPathLookup(doc, "$.result.realtime.direct")
	wi.Today.Wind = v.(string);
	v, _ = jsonpath.JsonPathLookup(doc, "$.result.realtime.power")
	wi.Today.Wind += v.(string)

	v, _ = jsonpath.JsonPathLookup(doc, "$.result.future[0].weather")
	wi.Today.Weather = v.(string)
	v, _ = jsonpath.JsonPathLookup(doc, "$.result.future[0].temperature")
	wi.Today.Weather += " " + v.(string)

	v, _ = jsonpath.JsonPathLookup(doc, "$.result.future[1].temperature")
	wi.Tomorrow.Temperature = v.(string)
	v, _ = jsonpath.JsonPathLookup(doc, "$.result.future[1].weather")
	wi.Tomorrow.Weather = v.(string)
	v, _ = jsonpath.JsonPathLookup(doc, "$.result.future[1].date")
	wi.Tomorrow.Date = v.(string)

	v, _ = jsonpath.JsonPathLookup(doc, "$.result.future[2].temperature")
	wi.Houtian.Temperature = v.(string)
	v, _ = jsonpath.JsonPathLookup(doc, "$.result.future[2].weather")
	wi.Houtian.Weather = v.(string)
	v, _ = jsonpath.JsonPathLookup(doc, "$.result.future[2].date")
	wi.Houtian.Date = v.(string)

	v, _ = jsonpath.JsonPathLookup(doc, "$.result.future[3].temperature")
	wi.Dahoutian.Temperature = v.(string)
	v, _ = jsonpath.JsonPathLookup(doc, "$.result.future[3].weather")
	wi.Dahoutian.Weather = v.(string)
	v, _ = jsonpath.JsonPathLookup(doc, "$.result.future[3].date")
	wi.Dahoutian.Date = v.(string)

	v, _ = jsonpath.JsonPathLookup(doc, "$.result.future[4].temperature")
	wi.Dadahoutian.Temperature = v.(string)
	v, _ = jsonpath.JsonPathLookup(doc, "$.result.future[4].weather")
	wi.Dadahoutian.Weather = v.(string)
	v, _ = jsonpath.JsonPathLookup(doc, "$.result.future[4].date")
	wi.Dadahoutian.Date = v.(string)

	return wi
}
