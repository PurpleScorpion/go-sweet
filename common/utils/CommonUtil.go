package utils

var (
	imgPath string = ""
)

func InitUtils() {
	imgPath = ValueString("${sweet.img.path}")
}
