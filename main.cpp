#include <time.h>

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

using namespace Pistache;
using namespace Pistache::Http;
static char* fontPath = "./msyh.ttf";

static const int width = 600;
static const int height = 800;
static void pngError(ResponseStream& writer, const std::string& error){
    char pngFileNameBuf[512];
    auto pngFileNamePtr = std::tmpnam(pngFileNameBuf);
    time_t rawtime;
    struct tm * timeinfo;
    char nowTimeBuffer [80];

    time (&rawtime);
    timeinfo = localtime (&rawtime);
    strftime (nowTimeBuffer, sizeof(nowTimeBuffer),"%m-%e %H:%M 星期%u",timeinfo);

    {
        pngwriter png(width, height, 1.0, pngFileNamePtr);
        png.bit_depth_ = 8;
        //if (waitStatus == std::cv_status::timeout) {

        //}
        png.plot_text_utf8(fontPath, 48, 10, 200, 0., nowTimeBuffer, 0., 0., 0.);

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


static void myGet(const std::string& city, ResponseStream& writer) {
    auto curlH = curl_easy_init();
    if(!curlH) {
        pngError(writer, "CURL创建失败");
        return;
    }
    std::stringstream strStream;
    curl_easy_setopt(curlH, CURLOPT_URL, (std::string("http://restapi.amap.com/v3/weather/weatherInfo?key=a4d4aa0d0eac2cbbc74ec84471dac239&extensions=all&city=") + city).c_str());
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
        pngError(writer, "OK");
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

        if(req.resource() == "/") {
            auto q = req.query();
            if(q.has("city")){
                auto city = q.get("city").get();
                std::vector<char> png;
                writer.headers().add<Header::ContentType>(MIME(Image, Png));
                auto stream = writer.stream(Http::Code::Ok);
                myGet(city, stream);
                return;
            }
        }

        writer.send(Http::Code::Not_Found);
    }
};


int main() {
    curl_global_init(0);
    Http::listenAndServe<HelloHandler>("*:9080");
}

