package weather2png_server

import (
	"fmt"
	"os"
	"testing"
)

func TestRealtime(t *testing.T) {
	os.Setenv("HEFENG_APIKEY", "fdf2a9607de34f4bb815571a86adb33b")
	wi := &WeatherInfo{}
	err := realtime("101010100", wi)
	fmt.Println(err)
	fmt.Printf("%+v", wi)
}

func Test3d(t *testing.T) {
	os.Setenv("HEFENG_APIKEY", "fdf2a9607de34f4bb815571a86adb33b")
	wi := &WeatherInfo{}
	err := threeDays("101010100", wi)
	fmt.Println(err)
	fmt.Printf("%+v", wi)
}
