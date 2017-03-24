package files

import (
	"errors"
	"image"
	"image/color"
)

func QuarteredImage(src image.Image) []*image.RGBA {
	sw := src.Bounds().Dx()
	sh := src.Bounds().Dy()
	images := []*image.RGBA{
		image.NewRGBA(image.Rect(0, 0, sw>>1, sh>>1)),
		image.NewRGBA(image.Rect(0, 0, sw>>1, sh)),
		image.NewRGBA(image.Rect(0, 0, sw>>1, sh>>1)),
	}
	for h := 0; h < sh; h++ {
		remainderH := h % 2
		half := h >> 1
		for w := 0; w < sw; w++ {
			remainderW := w % 2
			setH := half
			if remainderW != remainderH {
				setH = h
			}
			images[remainderH+remainderW].Set(w>>1, setH, src.At(w, h))
		}
	}
	return images
}

func SeparationAlpha(src image.Image) (*image.RGBA, *image.RGBA) {
	sw := src.Bounds().Dx()
	sh := src.Bounds().Dy()
	rgbImage := image.NewRGBA(image.Rect(0, 0, sw, sh))
	aImage := image.NewRGBA(image.Rect(0, 0, sw, sh))
	for h := 0; h < sh; h++ {
		for w := 0; w < sw; w++ {
			rc, gc, bc, ac := src.At(w, h).RGBA()
			rgbImage.Set(w, h, color.RGBA{
				R: uint8(rc >> 8),
				G: uint8(gc >> 8),
				B: uint8(bc >> 8),
			})
			aImage.Set(w, h, color.RGBA{
				A: uint8(ac >> 8),
			})
		}
	}
	return rgbImage, aImage
}

func OversamplingThumbnail(src image.Image, scale float64) (image.Image, error) {
	if scale <= float64(1) {
		return nil, errors.New("Thumbnail scale must greater than 1")
	}
	sw := float64(src.Bounds().Dx())
	sh := float64(src.Bounds().Dy())
	if sw < scale || sh < scale {
		return nil, errors.New("Image`s with and height must greater than thumbnail scale")
	}
	dst := image.NewRGBA(image.Rect(0, 0, int(sw/scale), int(sh/scale)))
	db := dst.Bounds()
	absScale := int(scale)
	sLen := uint32(absScale * absScale)
	var rc, gc, bc, ac, rt, gt, bt, at uint32
	for h := db.Min.Y; h < db.Max.Y; h++ {
		for w := db.Min.X; w < db.Max.X; w++ {
			rt, gt, bt, at = 0, 0, 0, 0
			for a := 0; a < absScale; a++ {
				for b := 0; b < absScale; b++ {
					rc, gc, bc, ac = src.At(int(float64(w)*scale)+a, int(float64(h)*scale)+b).RGBA()
					rt += rc
					gt += gc
					bt += bc
					at += ac
				}
			}
			dst.Set(w, h, color.RGBA{
				R: uint8((rt / sLen) >> 8),
				G: uint8((gt / sLen) >> 8),
				B: uint8((bt / sLen) >> 8),
				A: uint8((at / sLen) >> 8),
			})
		}
	}
	return dst, nil
}

func OversamplingEnlarge(src image.Image) (image.Image, error) {
	scale := 2
	dst := image.NewRGBA(image.Rect(0, 0, int(src.Bounds().Dx()*scale), int(src.Bounds().Dy()*scale)))
	db := dst.Bounds()
	cols := make([]color.Color, scale*scale)
	var rc, gc, bc, ac, rt, gt, bt, at uint32
	for h := db.Min.Y; h < db.Max.Y; h++ {
		for w := db.Min.X; w < db.Max.X; w++ {
			for a := 0; a < scale; a++ {
				for b := 0; b < scale; b++ {
					i := a*scale + b
					cols[i] = src.At(w+a, h+b)
				}
			}
			dst.Set(w*scale, h*scale, cols[0])
			alternate := append(make([]color.Color, 0, 3), cols[0])

			rt, gt, bt, at = 0, 0, 0, 0
			nLen := len(cols)
			for i := nLen - 1; i >= 0; i-- {
				rc, gc, bc, ac = cols[i].RGBA()
				rt += rc
				gt += gc
				bt += bc
				at += ac
			}
			nCol := color.RGBA{
				R: uint8((rt / uint32(nLen)) >> 8),
				G: uint8((gt / uint32(nLen)) >> 8),
				B: uint8((bt / uint32(nLen)) >> 8),
				A: uint8((at / uint32(nLen)) >> 8),
			}
			dst.Set(w*scale+1, h*scale+1, nCol)
			alternate = append(alternate, nCol, cols[2])

			rt, gt, bt, at = 0, 0, 0, 0
			nLen = len(alternate)
			for i := nLen - 1; i >= 0; i-- {
				rc, gc, bc, ac = alternate[i].RGBA()
				rt += rc
				gt += gc
				bt += bc
				at += ac
			}
			nCol = color.RGBA{
				R: uint8((rt / uint32(nLen)) >> 8),
				G: uint8((gt / uint32(nLen)) >> 8),
				B: uint8((bt / uint32(nLen)) >> 8),
				A: uint8((at / uint32(nLen)) >> 8),
			}
			dst.Set(w*scale, h*scale+1, nCol)
			alternate[nLen-1] = cols[1]

			rt, gt, bt, at = 0, 0, 0, 0
			alternate[len(alternate)-1] = cols[1]
			nLen = len(alternate)
			for i := nLen - 1; i >= 0; i-- {
				rc, gc, bc, ac = alternate[i].RGBA()
				rt += rc
				gt += gc
				bt += bc
				at += ac
			}
			nCol = color.RGBA{
				R: uint8((rt / uint32(nLen)) >> 8),
				G: uint8((gt / uint32(nLen)) >> 8),
				B: uint8((bt / uint32(nLen)) >> 8),
				A: uint8((at / uint32(nLen)) >> 8),
			}
			dst.Set(w*scale+1, h*scale, nCol)
		}
	}
	return dst, nil
}
