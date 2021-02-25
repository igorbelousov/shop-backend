package web

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

func Upload(dir string, file multipart.File, fileName string) {
	os.MkdirAll(dir, 0777)
	filename := fmt.Sprintf("%s/%s", dir, fileName)

	outputfile, err := os.Create(filename)
	if err != nil {
		fmt.Println("Create file error ", err)
	}
	defer outputfile.Close()

	io.Copy(outputfile, file)
}
