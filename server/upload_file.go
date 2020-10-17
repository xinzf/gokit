package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
)

type UploadFile struct {
	ctx *gin.Context
}

type File struct {
	file *multipart.FileHeader
	ctx  *gin.Context
}

func (this *UploadFile) Get(name string) (*File, error) {
	f, err := this.ctx.FormFile(name)
	if err != nil {
		return nil, err
	}
	return &File{file: f, ctx: this.ctx}, nil
}

func (this *UploadFile) GetFiles(name string) []*File {
	form, _ := this.ctx.MultipartForm()
	files := form.File[fmt.Sprintf("%s[]", name)]

	f := make([]*File, 0)
	for _, file := range files {
		f = append(f, &File{
			file: file,
			ctx:  this.ctx,
		})
	}
	return f
}

func (this *File) SaveFile(dst string) error {
	return this.ctx.SaveUploadedFile(this.file, dst)
}

func (this *File) RawBody() ([]byte, error) {
	f, err := this.file.Open()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = f.Close()
	}()

	return ioutil.ReadAll(f)
}

func (this *File) Size() int64 {
	return this.file.Size
}

func (this *File) FileName() string {
	return this.file.Filename
}
