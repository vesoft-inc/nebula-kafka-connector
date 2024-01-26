package utils

import (
	"io"
	"mime/multipart"
	"os"
)

func SaveFormFile(fh *multipart.FileHeader, dest string) (int64, error) {
	src, err := fh.Open()
	if err != nil {
		return 0, err
	}
	defer src.Close()

	out, err := os.Create(dest)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	return io.Copy(out, src)
}
