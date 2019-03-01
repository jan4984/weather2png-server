#include <signal.h>
#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
#include <time.h>

#include <map>
#include <iostream>
#include <pistache/net.h>
#include <pistache/http.h>
#include <pistache/endpoint.h>
//WTF https://stackoverflow.com/questions/3002101/how-can-i-violate-encapsulation-property
#define private public
#include <pngwriter.h>
#undef private
#include <fstream>

#include <curl/curl.h>
#include <rapidjson/document.h>

using namespace Pistache;
using namespace Pistache::Http;
static char* fontPath = "./msyh.ttf";

static const int width = 600;
static const int height = 800;
static auto mapOfWeekdayName =std::map<int, const char*>{
        {0, "%d %02d-%02d 星期一"},{1, "%d %02d-%02d 星期二"},{2, "%d %02d-%02d 星期三"},
        {3, "%d %02d-%02d 星期四"},{4, "%d %02d-%02d 星期五"},{5, "%d %02d-%02d 星期六"},{6, "%d %02d-%02d 星期天"}
};

static void pngError(ResponseStream& writer, const std::string& error){
    char pngFileNameBuf[512];
    auto pngFileNamePtr = std::tmpnam(pngFileNameBuf);
    time_t rawtime;
    struct tm * timeinfo;
    char nowTimeBuffer [80];

    time (&rawtime);
    rawtime += 8 * 3600;//BeiJing Time
    timeinfo = gmtime(&rawtime);
    {
        pngwriter png(width, height, 1.0, pngFileNamePtr);
        png.bit_depth_ = 8;

        sprintf(nowTimeBuffer, "%02d:%02d", timeinfo->tm_hour, timeinfo->tm_sec);
        png.plot_text_utf8(fontPath, 80, 150, 700, 0., nowTimeBuffer, 0., 0., 0.);
        auto weekdayName = mapOfWeekdayName[timeinfo->tm_wday - 1];
        sprintf(nowTimeBuffer, weekdayName, timeinfo->tm_mon + 1, timeinfo->tm_mday);
        png.plot_text_utf8(fontPath, 40, 120, 620, 0., nowTimeBuffer, 0., 0., 0.);

        png.plot_text_utf8(fontPath, 48, 10, 300, 0., (char*)error.c_str(), 0., 0., 0.);
        png.close();
    }

    {
        char pngFileBuf[4096];
        std::ifstream f(pngFileNamePtr, std::ios_base::binary | std::ios_base::in);
        do {
            f.read(pngFileBuf, sizeof(pngFileBuf));
            if(f.gcount() <= 0){
                break;
            }
            writer.write(pngFileBuf, f.gcount());
        } while (true);
    }
    writer << ends;
    std::remove(pngFileNamePtr);
}

struct Live{
    std::string weather;
    std::string wind;
    std::string temperature;
    std::string aqi;
};

struct Forecast{
    std::string weather;
    std::string temperature;
    std::string date;
};

Live live;
Forecast tomorrow;
Forecast houtian;
Forecast dahoutian;
Forecast dadahoutian;

static void pngWeather(ResponseStream& writer){
    char pngFileNameBuf[512];
    auto pngFileNamePtr = std::tmpnam(pngFileNameBuf);
    time_t rawtime;
    struct tm * timeinfo;
    char nowTimeBuffer [80];

    time (&rawtime);
    rawtime += 8 * 3600;//BeiJing Time
    timeinfo = gmtime(&rawtime);
    {
        pngwriter png(width, height, 1.0, pngFileNamePtr);
        png.bit_depth_ = 8;

        //时间
        sprintf(nowTimeBuffer, "%02d:%02d", timeinfo->tm_hour, timeinfo->tm_min);
        png.plot_text_utf8(fontPath, 80, 50, 700, 0., nowTimeBuffer, 0., 0., 0.);
        //日期
        auto weekdayName = mapOfWeekdayName[timeinfo->tm_wday - 1];
        sprintf(nowTimeBuffer, weekdayName, timeinfo->tm_year + 1900, timeinfo->tm_mon + 1, timeinfo->tm_mday);
        png.plot_text_utf8(fontPath, 40, 90, 620, 0., nowTimeBuffer, 0., 0., 0.);
        //当前天气
        sprintf(nowTimeBuffer, "%s %s", live.weather.c_str(), live.temperature.c_str());
        png.plot_text_utf8(fontPath, 40, 350, 700, 0., nowTimeBuffer, 0., 0., 0.);

        //分割线
        png.line_blend(0, 600, 600, 600, 1.0, 0., 0., 0.);

        int y = 500;
        //明天天气
        png.plot_text_utf8(fontPath, 24, 10, y, 0., "明天", 0., 0., 0.);
        sprintf(nowTimeBuffer, "%s %s", tomorrow.weather.c_str(), tomorrow.temperature.c_str());
        png.plot_text_utf8(fontPath, 40, 50, y-50, 0., nowTimeBuffer, 0., 0., 0.);

        y-=120;
        png.plot_text_utf8(fontPath, 24, 10, y, 0., "后天", 0., 0., 0.);
        sprintf(nowTimeBuffer, "%s %s", houtian.weather.c_str(), houtian.temperature.c_str());
        png.plot_text_utf8(fontPath, 40, 50, y-50, 0., nowTimeBuffer, 0., 0., 0.);

        y-=120;
        png.plot_text_utf8(fontPath, 24, 10, y, 0., (char*)dahoutian.date.c_str(), 0., 0., 0.);
        sprintf(nowTimeBuffer, "%s %s", dahoutian.weather.c_str(), dahoutian.temperature.c_str());
        png.plot_text_utf8(fontPath, 40, 50, y-50, 0., nowTimeBuffer, 0., 0., 0.);

        y-=120;
        png.plot_text_utf8(fontPath, 24, 10, y, 0., (char*)dadahoutian.date.c_str(), 0., 0., 0.);
        sprintf(nowTimeBuffer, "%s %s", dadahoutian.weather.c_str(), dadahoutian.temperature.c_str());
        png.plot_text_utf8(fontPath, 40, 50, y-50, 0., nowTimeBuffer, 0., 0., 0.);

        png.close();
    }

    {
        char pngFileBuf[4096];
        std::ifstream f(pngFileNamePtr, std::ios_base::binary | std::ios_base::in);
        do {
            f.read(pngFileBuf, sizeof(pngFileBuf));
            if(f.gcount() <= 0){
                break;
            }
            writer.write(pngFileBuf, f.gcount());
        } while (true);
    }
    writer << ends;
    std::remove(pngFileNamePtr);
}

static size_t curl_write_cb(void *data, size_t size, size_t nmemb, void* userdata){
    auto byteCount = size * nmemb;
    auto downloadFile = reinterpret_cast<std::ostream*> (userdata);
    downloadFile->write((const char*)data, byteCount);
    if(!downloadFile->good()){
        std::cout << "write update content failed:" << strerror(errno) << std::endl;
        return 0;
    }
    return byteCount;
}


static void getWeather(const std::string &city, ResponseStream &writer) {
    auto apiKey = std::getenv("API_KEY");
    auto curlH = curl_easy_init();
    if(!curlH) {
        pngError(writer, "CURL创建失败");
        return;
    }
    if(city == ""){
        pngWeather(writer);
        return;
    }
    std::stringstream strStream;
    curl_easy_setopt(curlH, CURLOPT_URL,
            (std::string("http://apis.juhe.cn/simpleWeather/query?key=") + apiKey + "&&city=" + city).c_str());
    curl_easy_setopt(curlH, CURLOPT_WRITEFUNCTION, curl_write_cb);
    curl_easy_setopt(curlH, CURLOPT_WRITEDATA, &strStream);
    auto res = curl_easy_perform(curlH);
    if(res == CURLE_OK){
        long rspCode = 0;
        curl_easy_getinfo(curlH, CURLINFO_RESPONSE_CODE, &rspCode);
        if(rspCode<200 || rspCode >= 300){
            pngError(writer, std::string("HTTP返回") + std::to_string(rspCode));
            return;
        }

        {
            using namespace rapidjson;
            Document doc;
            auto jsonStr = strStream.str();
            if(jsonStr == ""){
                pngError(writer, "查询天气数据失败");
                return;
            }
            std::cout << jsonStr << std::endl;
            doc.Parse(jsonStr.c_str());
            if(doc.GetParseError()) {
                pngError(writer, "解析JSON失败");
                return;
            }

            if(std::string("查询成功!").find(doc["reason"].GetString()) == std::string::npos){
                pngError(writer, doc["reason"].GetString());
                return;
            }
            auto liveObj = doc["result"].GetObject()["realtime"].GetObject();
            live.temperature = std::string() + liveObj["temperature"].GetString() + "℃";
            live.weather = std::string(liveObj["info"].GetString());
            if(liveObj.HasMember("aqi")){
                auto aqi = std::atoi(liveObj["aqi"].GetString());
                if(aqi == 0)
                    live.aqi = "无空气数据";
                else{
                    if(aqi <= 50) live.aqi = "空气优";
                    else if(aqi < 100) live.aqi = "空气良";
                    else if(aqi < 200) live.aqi = "轻度污染";
                    else if(aqi < 300) live.aqi = "中度污染";
                    else live.aqi = "重度污染";
                }
            } else {
                live.aqi = "无空气数据";
            }
            live.wind = std::string(liveObj["direct"].GetString()) + " " + liveObj["power"].GetString();

            auto tomorrowObj = doc["result"].GetObject()["future"].GetArray()[0].GetObject();
            tomorrow.temperature = tomorrowObj["temperature"].GetString();
            tomorrow.weather = tomorrowObj["weather"].GetString();

            auto houtianObj = doc["result"].GetObject()["future"].GetArray()[1].GetObject();
            houtian.temperature = houtianObj["temperature"].GetString();
            houtian.weather = houtianObj["weather"].GetString();

            auto dahoutianObj = doc["result"].GetObject()["future"].GetArray()[2].GetObject();
            dahoutian.temperature = dahoutianObj["temperature"].GetString();
            dahoutian.weather = dahoutianObj["weather"].GetString();
            dahoutian.date = dahoutianObj["date"].GetString();

            auto dadahoutianObj = doc["result"].GetObject()["future"].GetArray()[3].GetObject();
            dadahoutian.temperature = dadahoutianObj["temperature"].GetString();
            dadahoutian.weather = dadahoutianObj["weather"].GetString();
            dadahoutian.date = dadahoutianObj["date"].GetString();
        }

        pngWeather(writer);
        return;
    }
    pngError(writer, std::string("CURL返回") + std::to_string(res));
    return;
}

struct HelloHandler : public Http::Handler {
HTTP_PROTOTYPE(HelloHandler)
    void onRequest(const Http::Request& req, Http::ResponseWriter writer) override{
        if(req.resource() == "/hello") {
            writer.send(Http::Code::Ok, "Hello, World!");
            return;
        }

        if(req.resource() == "/time") {
            auto stream = writer.stream(Http::Code::Ok);
            getWeather("", stream);
            return;
        }

        if(req.resource() == "/weather") {
            auto q = req.query();
            if(q.has("city")){
                auto city = q.get("city").get();
                std::vector<char> png;
                writer.headers().add<Header::ContentType>(MIME(Image, Png));
                auto stream = writer.stream(Http::Code::Ok);
                getWeather(city, stream);
                return;
            }
        }

        writer.send(Http::Code::Not_Found);
    }
};

static void my_handler(int s);

Endpoint endpoint("*:9080");

int main() {
    struct sigaction sigIntHandler;

    sigIntHandler.sa_handler = my_handler;
    sigemptyset(&sigIntHandler.sa_mask);
    sigIntHandler.sa_flags = 0;

    sigaction(SIGINT, &sigIntHandler, NULL);
    sigaction(SIGTERM, &sigIntHandler, NULL);
    curl_global_init(0);
    //Http::listenAndServe<HelloHandler>("*:9080");
    endpoint.init(Endpoint::options());
    endpoint.setHandler(make_handler<HelloHandler>());
    endpoint.serve();
}

static void my_handler(int s){
    printf("Caught signal %d\n",s);
    endpoint.shutdown();
    exit(1);
}
