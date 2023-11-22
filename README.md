# 📅 CalenGo - A Go Calendar Icon API

## What is this ❓

CalenGo is a simple Go program exposing an API which generate on-the-fly an image of a calendar displaying Month and Date (number and weekday).  
Add a Unix timestamp and a timezone to get a calendar icon depicting the day of the timestamp in your timezone !

<div align="center">
<img src="https://calengo.remi-espie.me/calendar?size=500" alt="Today's date !"/>
</div>

The application answer on the url `/ping` with a beautiful `pong` and most importantly on `/calendar` with the arguments `timestamp`, `locale` and `size` to get your image !

On `/calendar`, you may add `timestamp` parameter, `locale` parameter, `size` parameter, the three of them or none at all ! 
- `timestamp` parameter allows you to send a custom timestamp a gate an image with its date.  
Without this parameter, it will output today's date. 
- `locale` parameter allows you to give it a custom timezone for your timestamp. It accepts timezones from [IANA Timezone Database](https://www.iana.org/time-zones)'s TZ identifier. You can also get them from [Wikipedia tz database list](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones).  
Without this parameter, it will use `UTC` timezone.
- `size` parameter allows you to output an image with a custom size, between 50x50px and 1000x1000px. Lower or higher values will return an image of 1000x1000px.   
Without this parameter, it will output an image of 1000x1000px.

## 🚀 Getting started

Just run `go run main.go` from this repository to run the app on port `8080` !

You can also get the latest update from [releases](https://github.com/remi-espie/calengo/releases), however **the only release available was built *with and only* for Linux x86_64**.

You can also build and run a docker image from the project using `docker compose up`.

You may change your `calendar_template.png` as well as your font (here `Roboto-Bold.ttf`) file, but you'll probably need to update the code.

## 📝 License

This project is under the MIT License - see the [LICENSE](LICENSE) file for details.

## 💻 Dependencies for nerds

Developed with Go 19

And the following dependencies:
- [gin-gonic](https://github.com/gin-gonic/gin)
- [freetype](https://github.com/golang/freetype)


- [Roboto Font](https://fonts.google.com/specimen/Roboto)
