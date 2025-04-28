package scrape

import "math/rand"

var IMAGE_LIST = []string{
	"thumbnails/blue2.jpg",
	"thumbnails/blue3.jpg",
	"thumbnails/blue4.jpg",
	"thumbnails/blue_green1.jpg",
	"thumbnails/blue_green2.jpg",
	"thumbnails/blue_green3.jpg",
	"thumbnails/blue_green4.jpg",
	"thumbnails/blue_green5.jpg",
	"thumbnails/blue.jpg",
	"thumbnails/color1.jpg",
	"thumbnails/color2.jpg",
	"thumbnails/color3.jpg",
	"thumbnails/color4.jpg",
	"thumbnails/cyan1.jpg",
	"thumbnails/cyan2.jpg",
	"thumbnails/cyan3.jpg",
	"thumbnails/cyan4.jpg",
	"thumbnails/cyan_orange1.jpg",
	"thumbnails/cyan_orange2.jpg",
	"thumbnails/cyan_orange3.jpg",
	"thumbnails/cyan_orange4.jpg",
	"thumbnails/green1.jpg",
	"thumbnails/green2.jpg",
	"thumbnails/green3.jpg",
	"thumbnails/green4.jpg",
	"thumbnails/light_colour1.jpg",
	"thumbnails/light_colour2.jpg",
	"thumbnails/light_colour3.jpg",
	"thumbnails/light_colour4.jpg",
	"thumbnails/light-orange1.jpg",
	"thumbnails/light-orange2.jpg",
	"thumbnails/light-orange3.jpg",
	"thumbnails/light-orange4.jpg",
	"thumbnails/orange1.jpg",
	"thumbnails/orange2.jpg",
	"thumbnails/orange3.jpg",
	"thumbnails/orange4.jpg",
	"thumbnails/orange_pink1.jpg",
	"thumbnails/orange_pink2.jpg",
	"thumbnails/orange_pink3.jpg",
	"thumbnails/orange_pink4.jpg",
	"thumbnails/purple1.jpg",
	"thumbnails/purple2.jpg",
	"thumbnails/purple3.jpg",
	"thumbnails/purple4.jpg",
	"thumbnails/purple5.jpg",
	"thumbnails/purple6.jpg",
	"thumbnails/purple7.jpg",
	"thumbnails/purple8.jpg",
	"thumbnails/red1.jpg",
	"thumbnails/red2.jpg",
	"thumbnails/red3.jpg",
}

func defaultGetImage() string {
	// get random image from IMAGE_LIST
	image := IMAGE_LIST[rand.Intn(len(IMAGE_LIST))]

	return image
}
