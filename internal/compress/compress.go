package compress

import (
	"bytes"
	"compress/gzip"
)

func Compress(data bytes.Buffer) (*bytes.Buffer, error) {
	var b bytes.Buffer

	compressWriter := gzip.NewWriter(&b)
	defer compressWriter.Close()
	_, err := compressWriter.Write(data.Bytes())
	if err != nil {

		return nil, err
	}
	return &b, nil

}
