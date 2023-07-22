package qrcode

import (
	"bufio"
	"fmt"
	"image"
	"net/http"
	"os"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

func qr() error {
	encOpts := []qrcode.EncodeOption{
		qrcode.WithEncodingMode(qrcode.EncModeByte),

		// 纠错程度
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionHighest),

		// 版本
		qrcode.WithVersion(9),
	}

	qrc, err := qrcode.NewWith("https://github.com/bougou", encOpts...)
	if err != nil {
		return fmt.Errorf("could not generate QRCode: %v", err)
	}

	fileToBeUploaded := "./header.jpg"
	file, err := os.Open(fileToBeUploaded)

	if err != nil {
		return fmt.Errorf("open: %s", err)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	bytes := make([]byte, size)

	// read file into bytes
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)
	if err != nil {
		return fmt.Errorf("read bytes: %s", err)
	}

	filetype := http.DetectContentType(bytes)
	fmt.Println("filetype: ", filetype)

	file.Seek(0, 0)
	i, x, err := image.Decode(buffer)
	if err != nil {
		return fmt.Errorf("decode: %s", err)
	}
	fmt.Println("s, ", x)

	opts := []standard.ImageOption{
		// 控制像素的形状，默认是实心正方形
		// standard.WithCustomShape(NewSmallerCircle(1)),
		standard.WithLogoImage(i),
		// standard.WithQRWidth(20),
		standard.WithHalftone("./header.jpg"),
	}

	w, err := standard.New("./qrcode.jpeg", opts...)
	if err != nil {
		return fmt.Errorf("standard.New failed: %v", err)
	}

	// save file
	if err = qrc.Save(w); err != nil {
		return fmt.Errorf("could not save image: %v", err)
	}

	return nil
}
