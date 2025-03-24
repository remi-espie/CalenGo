package main

import (
	"bytes"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/goki/freetype"
	"github.com/goki/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	fontFace font.Face
	once     sync.Once
)

func main() {

	// Open the template image
	templateImage, err := os.Open("./calendar_template.png")
	if err != nil {
		log.Fatal("Failed to open template image:", err)
	}
	defer func(templateImage *os.File) {
		err := templateImage.Close()
		if err != nil {
			log.Fatal("Failed closing template image")
		}
	}(templateImage)

	// Decode the template image
	template, err := png.Decode(templateImage)
	if err != nil {
		log.Fatal("Failed to decode template image:", err)
	}

	fontFace = loadFont()

	r := gin.Default()

	// Serve the favicon.ico file
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.File("./calendar_template.png")
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World !",
		})
	})

	r.GET("/calendar", func(c *gin.Context) {
		// Parse the request payload
		var req CalendarRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if req.Timestamp == 0 {
			req.Timestamp = time.Now().Unix()
		}

		// Convert the timestamp to time.Time
		calTime := time.Unix(req.Timestamp, 0)

		if req.Locale == "" {
			req.Locale = "UTC"
		}

		if req.Size == 0 || req.Size > 1000 || req.Size < 50 {
			req.Size = 1000
		}

		// Set the desired time locale
		loc, err := time.LoadLocation(req.Locale)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid locale"})
			return
		}
		calTime = calTime.In(loc)

		// Generate the calendar image using calTime
		img := generateCalendarImage(calTime, template, fontFace, req.Size)

		// Create a buffer to hold the image data
		buf := new(bytes.Buffer)

		// Encode the image to PNG and write it to the buffer
		err = png.Encode(buf, img)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode image"})
			return
		}

		// Set the appropriate headers
		c.Header("Content-Type", "image/png")
		c.Header("Content-Disposition", "inline")

		// Write the image data from the buffer to the response
		_, err = c.Writer.Write(buf.Bytes())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write image"})
			return
		}
	})

	err = r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
}

type CalendarRequest struct {
	Timestamp int64  `form:"timestamp"`
	Locale    string `form:"locale"`
	Size      int    `form:"size"`
}

func generateCalendarImage(date time.Time, template image.Image, fontFace font.Face, size int) *image.NRGBA {
	// Set the dimensions of the calendar image
	width := template.Bounds().Dx()

	// Create a new RGBA image
	img := image.NewNRGBA(template.Bounds())

	// Set the font face and size
	fontSize := 128

	draw.Draw(img, img.Bounds(), template, img.Bounds().Min, draw.Src)

	// Create a drawerMonth to write month on the image
	drawerMonth := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.White),
		Face: fontFace,
	}

	textWidthMonth := drawerMonth.MeasureString(date.Format("January")).Round()

	// Calculate the text position
	textXMonth := (width - textWidthMonth) / 2

	drawerMonth.Dot = fixed.P(textXMonth, 100+fontSize)

	drawerMonth.DrawString(date.Format("January"))

	// Create a drawerDate to write day number on the image
	drawerDate := &font.Drawer{
		Dst: img,
		Src: image.NewUniform(color.RGBA{
			R: 0x5B,
			G: 0x51,
			B: 0x44,
			A: 0xFF,
		}),
		Face: fontFace,
	}

	textWidthDate := drawerDate.MeasureString(date.Format("02")).Round()

	// Calculate the text position
	textXDate := (width - textWidthDate) / 2

	drawerDate.Dot = fixed.P(textXDate, 400+fontSize)

	drawerDate.DrawString(date.Format("02"))

	// Create a drawerDay to write day on the image
	drawerDay := &font.Drawer{
		Dst: img,
		Src: image.NewUniform(color.RGBA{
			R: 0x5B,
			G: 0x51,
			B: 0x44,
			A: 0xFF,
		}),
		Face: fontFace,
	}

	textWidthDay := drawerDay.MeasureString(date.Format("Monday")).Round()

	// Calculate the text position
	textXDay := (width - textWidthDay) / 2

	drawerDay.Dot = fixed.P(textXDay, 600+fontSize)

	drawerDay.DrawString(date.Format("Monday"))

	if size != 1000 {
		img = imaging.Resize(img, size, size, imaging.Linear)
	}

	return img
}

func loadFont() font.Face {
	once.Do(func() {
		fontBytes, err := os.ReadFile("./Roboto-Bold.ttf")
		if err != nil {
			log.Fatal("Error reading font file")
		}

		robotoFont, err := freetype.ParseFont(fontBytes)
		if err != nil {
			log.Fatal("Error parsing the font")
		}

		fontFace = truetype.NewFace(robotoFont, &truetype.Options{
			Size:    128,
			DPI:     72,
			Hinting: 0,
		})
	})
	return fontFace
}
