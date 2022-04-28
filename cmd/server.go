package main

import (
	"io"
	"net/http"
	"os"
	"time"
	"unicode/utf8"

	weather2png_server "github.com/jan4984/weather2png-server"
)

var weakDays = [7]string{"天", "一", "二", "三", "四", "五", "六"}

func draw(png *weather2png_server.PngWriter, wi *weather2png_server.WeatherInfo, writer io.Writer) {
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
	if wi.Err != nil {
		png.Text("错误!", 10, 200, 40)
		png.Text(wi.Err.Error(), 10, 300, 40)
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
	weather2png_server.Start()
	// time.Sleep(time.Second * 5)
	// f, _ := os.Create("out.png")
	// png := weather2png_server.NewPngWriter(600, 800, os.Getenv("TTF_PATH"))
	// wi := weather2png_server.Get("101020500")
	// draw(png, &wi, f)
	http.HandleFunc("/update", func(writer http.ResponseWriter, request *http.Request) {
		params := request.URL.Query()
		writer.Header().Set("Content-Type", "image/png")
		city := os.Getenv("LOCATION_ID")
		if city == "" {
			city = params.Get("city")
		}
		wi := weather2png_server.Get(city)
		png := weather2png_server.NewPngWriter(600, 800, os.Getenv("TTF_PATH"))
		draw(png, &wi, writer)
	})
	http.ListenAndServe(":10008", nil)
}
