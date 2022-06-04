package weather2png_server

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
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
	Err error
}

var client = http.Client{
	Timeout: time.Second * 5,
	//my router outof update for the ca-cert
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func doReq(path, loc string, params url.Values) (map[string]interface{}, error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("WRONG URL %s : %w", path, err)
	}
	params.Set("key", os.Getenv("HEFENG_APIKEY"))
	params.Set("location", loc)
	u.RawQuery = params.Encode()

	resp, err := client.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("GET %s failed:%w", path, err)
	}
	defer resp.Body.Close()
	//reader, err := gzip.NewReader(resp.Body)
	//defer reader.Close()
	//respBody, err := ioutil.ReadAll(reader)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("READ %s failed:%w", path, err)
	}
	var doc map[string]interface{}
	err = json.Unmarshal(respBody, &doc)
	if err != nil {
		return nil, fmt.Errorf("UNMARSHAL %s %s failed:%w", path, string(respBody), err)
	}
	if doc["code"] != "200" {
		return nil, fmt.Errorf("WRONG CODE %s %s failed", path, string(respBody))
	}
	return doc, nil
}

var wis = sync.Map{}

func Start() {
	defaultLoc := os.Getenv("LOCATION_ID")
	if defaultLoc != "" {
		wis.Store(defaultLoc, &WeatherInfo{})
	}
	go func() {
		n := int64(0)
		for {
			locs := make(map[string]*WeatherInfo)
			wis.Range(func(key, value any) bool {
				locs[key.(string)] = value.(*WeatherInfo)
				return true
			})
			for loc, wi := range locs {
				wi.Err = realtime(loc, wi)
				if wi.Err != nil {
					fmt.Println(loc, wi.Err)
				}
			}
			for loc, wi := range locs {
				if wi.Today.Weather == "" {
					wi.Err = threeDays(loc, wi)
				} else if n%(10*60) == 0 {
					wi.Err = threeDays(loc, wi)
				}
				if wi.Err != nil {
					fmt.Println(loc, wi.Err)
				}
			}
			time.Sleep(time.Minute)
			n = n + 60
		}
	}()
}

func Get(loc string) WeatherInfo {
	wi := WeatherInfo{}
	wiAny, has := wis.LoadOrStore(loc, &wi)
	if !has {
		time.Sleep(time.Second * 6)
	}

	return *wiAny.(*WeatherInfo)
}

func realtime(loc string, wi *WeatherInfo) error {
	doc, err := doReq("https://devapi.qweather.com/v7/weather/now", loc, url.Values{})
	if err != nil {
		return err
	}
	doc = doc["now"].(map[string]interface{})
	wi.Live.Temperature = doc["temp"].(string)
	wi.Live.Weather = doc["text"].(string)
	//wi.Live.Wind = doc["windDir"].(string) + doc["windScale"].(string) + "级"

	doc, err = doReq("https://devapi.qweather.com/v7/air/now", loc, url.Values{})
	if err != nil {
		return err
	}
	doc = doc["now"].(map[string]interface{})
	wi.Live.Aqi = doc["category"].(string)
	return nil
}

func threeDays(loc string, wi *WeatherInfo) error {
	doc, err := doReq("https://devapi.qweather.com/v7/weather/7d", loc, url.Values{})
	if err != nil {
		return err
	}
	doc0 := doc["daily"].([]interface{})[0].(map[string]interface{})
	doc1 := doc["daily"].([]interface{})[1].(map[string]interface{})
	doc2 := doc["daily"].([]interface{})[2].(map[string]interface{})
	doc3 := doc["daily"].([]interface{})[3].(map[string]interface{})
	doc4 := doc["daily"].([]interface{})[4].(map[string]interface{})
	for _, doc = range []map[string]interface{}{doc0, doc1, doc2, doc3, doc4} {
		if doc["textDay"] == doc["textNight"] {
			doc["textNight"] = ""
		} else {
			doc["textNight"] = fmt.Sprintf("-%v", doc["textNight"])
		}
	}
	wi.Today.Weather = fmt.Sprintf("%v%v %v/%v℃", doc0["textDay"], doc0["textNight"], doc0["tempMin"], doc0["tempMax"])
	wi.Today.Wind = fmt.Sprintf("%v%v级", doc0["windDirDay"], doc0["windScaleDay"])
	wi.Tomorrow.Weather = fmt.Sprintf("%v%v %v/%v℃", doc1["textDay"], doc1["textNight"], doc1["tempMin"], doc1["tempMax"])
	wi.Tomorrow.Date = doc1["fxDate"].(string)
	wi.Houtian.Weather = fmt.Sprintf("%v%v %v/%v℃", doc2["textDay"], doc2["textNight"], doc2["tempMin"], doc2["tempMax"])
	wi.Houtian.Date = doc2["fxDate"].(string)
	wi.Dahoutian.Weather = fmt.Sprintf("%v%v %v/%v℃", doc3["textDay"], doc3["textNight"], doc3["tempMin"], doc3["tempMax"])
	wi.Dahoutian.Date = doc3["fxDate"].(string)
	wi.Dadahoutian.Weather = fmt.Sprintf("%v%v %v/%v℃", doc4["textDay"], doc4["textNight"], doc4["tempMin"], doc4["tempMax"])
	wi.Dadahoutian.Date = doc4["fxDate"].(string)
	return nil
}
