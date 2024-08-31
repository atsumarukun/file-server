package entity

type FileBody struct {
	path string
	body []byte
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
