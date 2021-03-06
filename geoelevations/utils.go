package geoelevations

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"log"

	"io/ioutil"
)

func gzipBytes(b *[]byte) (*[]byte, error) {
	buf := new(bytes.Buffer)
	gz := gzip.NewWriter(buf)

	_, err := gz.Write(*b)
	if err != nil {
		return nil, err
	}

	err = gz.Close()
	if err != nil {
		return nil, err
	}

	out := buf.Bytes()

	return &out, nil
}

func ungzipBytes(b *[]byte) (*[]byte, error) {
	r, err := gzip.NewReader(ioutil.NopCloser(bytes.NewBuffer(*b)))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	bb, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return &bb, nil
}

func unzipFile(fileName string) ([]byte, error) {
	r, err := zip.OpenReader(fileName)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {
		fmt.Printf("Contents of %s:\n", f.Name)
		rc, err := f.Open()
		if err != nil {
			log.Printf("Error reading %s from %s: %s", f.Name, fileName, err.Error())
			return nil, err
		}
		defer rc.Close()

		bytes, err := ioutil.ReadAll(rc)
		if err != nil {
			log.Printf("Error reading %s from %s: %s", f.Name, fileName, err.Error())
			return nil, err
		}

		return bytes, nil
	}

	return nil, errors.New(fmt.Sprintf("No file in .zip %s", fileName))
}
