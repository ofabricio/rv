package data

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"path"
)

func downloadGoogleDriveFileWithURL(URL string) (io.Reader, error) {

	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	fileID := path.Base(path.Dir(u.Path))

	return downloadGoogleDriveFile(fileID)
}

func downloadGoogleDriveFile(fileID string) (io.Reader, error) {

	res, err := http.Get("https://drive.google.com/uc?export=download&id=" + fileID)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	d, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(d), nil
}
