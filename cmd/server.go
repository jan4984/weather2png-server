package main

import (
	"github.com/jan4984/weather2png"
	"net/http"
	"sync"
	"time"
	"io"
	"os"
	"unicode/utf8"
)

type wiTTL struct {
	wi weather2png.WeatherInfo
	at time.Time
}

var weakDays = [7]string{"天", "一", "二", "三", "四", "五", "六"}

func draw(png *weather2png.PngWriter, wi *weather2png.WeatherInfo, writer io.Writer) {
	defer func() {
		png.Reset(writer)
	}()
	now := time.Now().UTC().Add(time.Hour * 8)

	png.Text(wi.Today.Wind, 10, 40, 30)
	png.Text(wi.Today.Weather, 600-utf8.RuneCount([]byte(wi.Today.Weather))*30, 40, 30)
	png.VerticalLine(0, 50, 600)
	startY := 150
	y := startY
	png.Text(now.Format("15:04"), 20, y, 100)
	y += 60
	png.Text(now.Format("2006-1-2 星期"+weakDays[now.Weekday()]), 220, y, 40)
	y += 20
	png.VerticalLine(0, y, 600)
	if wi.Err != "" {
		png.Text("错误!", 10, 200, 40)
		png.Text(wi.Err, 10, 300, 40)
		return
	}
	png.Text(wi.Live.Weather+" "+wi.Live.Temperature+"℃", 300, startY, 50)
	png.Text(wi.Live.Aqi, 20, y-20, 40)

	y += 40
	png.Text("明天 "+wi.Tomorrow.Date, 10, y, 35)
	y += 65
	png.Text(wi.Tomorrow.Weather+" "+wi.Tomorrow.Temperature, 66, y, 50)
	y += 80

	png.Text("后天 "+wi.Houtian.Date, 10, y, 35)
	y += 65
	png.Text(wi.Houtian.Weather+" "+wi.Houtian.Temperature, 66, y, 50)
	y += 80

	png.Text(wi.Dahoutian.Date, 10, y, 35)
	y += 65
	png.Text(wi.Dahoutian.Weather+" "+wi.Dahoutian.Temperature, 66, y, 50)
	y += 80

	png.Text(wi.Dadahoutian.Date, 10, y, 35)
	y += 65
	png.Text(wi.Dadahoutian.Weather+" "+wi.Dadahoutian.Temperature, 66, y, 50)
	y += 80
}

func main() {
	if os.Getenv("TTF_PATH") == "" {
		panic("env TTF_PATH not defined")
	}
	var png = weather2png.NewPngWriter(600, 800, os.Getenv("TTF_PATH"))
	wis := sync.Map{}
	http.HandleFunc("/update", func(writer http.ResponseWriter, request *http.Request) {
		params := request.URL.Query()
		writer.Header().Set("Content-Type", "image/png")
		if v, has := params["city"]; has {
			city := v[0]
			if v, has := wis.Load(city);
				has && time.Now().Unix() < v.(*wiTTL).at.Unix()+1800 && v.(*wiTTL).wi.Err == "" {
				//has info and live and not error
				wi := v.(*wiTTL).wi
				draw(png, &wi, writer)
				return
			}
			wi := weather2png.FetchWeather(city)
			wis.Store(city, &wiTTL{wi, time.Now()})
			draw(png, &wi, writer)
			return
		}

		draw(png, &weather2png.WeatherInfo{
			Err: "未指定城市",
		}, writer)
	})

	http.ListenAndServe(":10008", nil)

}
