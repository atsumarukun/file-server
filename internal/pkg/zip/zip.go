package zip

import (
	"archive/zip"
	"bytes"
	"file-server/internal/pkg/config"
	"os"
)

func Compress(path string) ([]byte, error) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	if err := compress(w, config.STORAGE_PATH+path, "/"); err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func compress(w *zip.Writer, filePath string, innerPath string) error {
	info, err := os.Lstat(filePath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		innerPath += info.Name() + "/"

		files, err := os.ReadDir(filePath)
		if err != nil {
			return err
		}

		for _, v := range files {
			name := v.Name()
			if v.IsDir() {
				name += "/"
			}
			if err := compress(w, filePath+name, innerPath); err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = innerPath + info.Name()
		header.Method = zip.Deflate

		file, err := w.CreateHeader(header)
		if err != nil {
			return err
		}

		body, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		if _, err := file.Write(body); err != nil {
			return err
		}
	}

	return nil
}
