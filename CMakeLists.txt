cmake_minimum_required(VERSION 3.10)
project(weather2png_server)

set(CMAKE_CXX_STANDARD 14)

add_executable(weather2png_server main.cpp)

find_package(Threads)
find_package(CURL REQUIRED)
target_link_libraries(weather2png_server
        "/usr/local/lib/libpistache.a"
        "/home/ubuntu/sources/pngwriter/out/libPNGwriter.a"
        "-lpng -lfreetype"
        "${CMAKE_THREAD_LIBS_INIT}"
        "${CURL_LIBRARIES}")
target_include_directories(weather2png_server PRIVATE
        "/home/ubuntu/sources/pistache/third-party/rapidjson/include"
        "/usr/include/freetype2"
        "${CURL_INCLUDE_DIR}"
        )
