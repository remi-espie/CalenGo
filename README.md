# 📅 CalenGo - A Go Calendar Icon API

## What is this ❓

CalenGo is a simple Go program exposing an API which generate on-the-fly an image of a calendar displaying Month and Date (number and weekday).  
Add a Unix timestamp and a timezone to get a calendar icon depicting the day of the timestamp in your timezone !

The application answer on the url `/ping` with a beautiful `pong` and most importantly on `/calendar` with the arguments `timestamp` and `locale` to get your image !

On `/calendar`, you may add `timestamp` parameter, `locale` parameter, both or none ! Without `timestamp`, it will output today's date. Without `locale`, it will use `UTC` timezone.

## 🚀 Getting started

Just run `go run main.go` from this repository to run the app on port `8080` !

You can also get the latest update from [releases](https://github.com/remi-espie/calengo/releases), however **the only release available was built *with and only* for Linux x86_64**.

You may change your `calendar_template.png` as well as your font (here `Roboto-Bold.ttf`) file, but you'll probably need to update the code.

## 📝 License

This project is under the MIT License - see the [LICENSE](LICENSE) file for details.

## 💻 Dependencies for nerds

Developed with Go 19

And the following dependencies:
- [gin-gonic](https://github.com/gin-gonic/gin)
- [freetype](https://github.com/golang/freetype)


- [Roboto Font](https://fonts.google.com/specimen/Roboto)