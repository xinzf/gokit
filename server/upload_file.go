package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type UploadFile struct {
	ctx *gin.Context
}

type File struct {
	file        *multipart.FileHeader
	ctx         *gin.Context
	_dataReaded bool
	_data       []byte
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

func (this *File) GetMimeType() string {
	data, err := this.RawBody()
	if err != nil {
		return "application/octet-stream"
	}

	return http.DetectContentType(data)
}

func (this *File) SaveFile(dst string) error {
	return this.ctx.SaveUploadedFile(this.file, dst)
}

func (this *File) RawBody() ([]byte, error) {
	if this._dataReaded {
		return this._data, nil
	}

	this._dataReaded = true

	f, err := this.file.Open()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = f.Close()
	}()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	this._data = data
	return data, nil
}

func (this *File) Size() int64 {
	return this.file.Size
}

func (this *File) FileName() string {
	return this.file.Filename
}
