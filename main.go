package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"os"
	"time"
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

	fontFace := loadFont()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
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

		// Set the desired time locale
		loc, err := time.LoadLocation(req.Locale)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid locale"})
			return
		}
		calTime = calTime.In(loc)

		// Generate the calendar image using calTime
		img := generateCalendarImage(calTime, template, fontFace)

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
}

func generateCalendarImage(date time.Time, template image.Image, fontFace font.Face) *image.RGBA {
	// Set the dimensions of the calendar image
	width := template.Bounds().Dx()
	//height := template.Bounds().Dy()

	// Create a new RGBA image
	img := image.NewRGBA(template.Bounds())

	// Set the font face and size
	fontSize := 52

	draw.Draw(img, img.Bounds(), template, img.Bounds().Min, draw.Src)

	// Create a drawer to write month on the image
	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.White),
		Face: fontFace,
	}

	textWidth := drawer.MeasureString(date.Format("January")).Round()

	// Calculate the text position
	textX := (width - textWidth) / 2

	drawer.Dot = fixed.P(textX, 30+fontSize)

	drawer.DrawString(date.Format("January"))

	// Create a drawer to write day number on the image
	drawer = &font.Drawer{
		Dst: img,
		Src: image.NewUniform(color.RGBA{
			R: 0x5B,
			G: 0x51,
			B: 0x44,
			A: 0xFF,
		}),
		Face: fontFace,
	}

	textWidth = drawer.MeasureString(date.Format("02")).Round()

	// Calculate the text position
	textX = (width - textWidth) / 2

	drawer.Dot = fixed.P(textX, 210+fontSize)

	drawer.DrawString(date.Format("02"))

	// Create a drawer to write day on the image
	drawer = &font.Drawer{
		Dst: img,
		Src: image.NewUniform(color.RGBA{
			R: 0x5B,
			G: 0x51,
			B: 0x44,
			A: 0xFF,
		}),
		Face: fontFace,
	}

	textWidth = drawer.MeasureString(date.Format("Monday")).Round()

	// Calculate the text position
	textX = (width - textWidth) / 2

	drawer.Dot = fixed.P(textX, 140+fontSize)

	drawer.DrawString(date.Format("Monday"))

	return img
}

func loadFont() font.Face {
	fontBytes, err := os.ReadFile("./Roboto-Bold.ttf")
	if err != nil {
		log.Fatal("Error reading font file")
	}

	robotoFont, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatal("Error parsing the font")
	}

	return truetype.NewFace(robotoFont, &truetype.Options{
		Size:    52,
		DPI:     72,
		Hinting: 0,
	})
}
