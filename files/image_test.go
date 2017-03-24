package files

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/BurntSushi/graphics-go/graphics"
)

func TestSeparationAlpha(t *testing.T) {
	src, err := LoadImage("./test/texture1_0.png")
	if err != nil {
		log.Fatal(err)
	}
	rgbImage, aImage := SeparationAlpha(src)
	err = saveImage("./test/rgbImage.jpg", rgbImage)
	if err != nil {
		log.Fatal(err)
	}

	err = saveImage("./test/aImage.png", aImage)
	if err != nil {
		log.Fatal(err)
	}

}

func TestQuarteredImage(t *testing.T) {
	src, err := LoadImage("./test/200.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	dsts := QuarteredImage(src)
	for i, dst := range dsts {
		err = saveImage(fmt.Sprintf("./test/%d.jpg", i), dst)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func TestOversamplingThumbnail(t *testing.T) {
	src, err := LoadImage("./test/200.jpg")
	if err != nil {
		log.Fatal(err)
	}
	dst, err := OversamplingThumbnail(src, float64(2))
	if err != nil {
		log.Fatal(err)
	}
	err = saveImage("./test/100.jpg", dst)
	if err != nil {
		log.Fatal(err)
	}

	dst, err = OversamplingThumbnail(src, float64(4))
	if err != nil {
		log.Fatal(err)
	}
	err = saveImage("./test/50.jpg", dst)
	if err != nil {
		log.Fatal(err)
	}
}

func TestOversamplingEnlarge(t *testing.T) {
	src, err := LoadImage("./test/50.jpg")
	if err != nil {
		log.Fatal(err)
	}
	dst, err := OversamplingEnlarge(src)
	if err != nil {
		log.Fatal(err)
	}
	err = saveImage("./test/101.jpg", dst)
	if err != nil {
		log.Fatal(err)
	}

	// 缩略图的大小

	ddst := image.NewRGBA(image.Rect(0, 0, 960, 600))
	// 产生缩略图,等比例缩放
	err = graphics.Thumbnail(ddst, src)
	if err != nil {
		log.Fatal(err)
	}

	err = saveImage("./test/102.png", ddst)
	if err != nil {
		log.Fatal(err)
	}

	//dst, err = OversamplingThumbnail(src, float64(4))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = saveImage("./test/50.jpg", dst)
	//if err != nil {
	//	log.Fatal(err)
	//}
}

// LoadImage decodes an image from a file.
func LoadImage(path string) (img image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	if strings.HasSuffix(file.Name(), ".jpg") || strings.HasSuffix(file.Name(), ".jpeg") {
		return jpeg.Decode(file)
	}
	img, _, err = image.Decode(file)
	return
}

// 保存Png图片
func saveImage(path string, img image.Image) (err error) {
	// 需要保存的文件
	imgfile, err := os.Create(path)
	defer imgfile.Close()
	// 以PNG格式保存文件
	if strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg") {
		return jpeg.Encode(imgfile, img, &jpeg.Options{90})
	}
	return png.Encode(imgfile, img)
}
