package entity

type FileBody struct {
	path string
	body []byte
}

func NewFileBody(path string, body []byte) *FileBody {
	file := &FileBody{}

	file.SetPath(path)
	file.SetBody(body)

	return file
}

func (f *FileBody) GetPath() string {
	return f.path
}

func (f *FileBody) SetPath(path string) {
	f.path = path
}

func (f *FileBody) GetBody() []byte {
	return f.body
}

func (f *FileBody) SetBody(body []byte) {
	f.body = body
}
