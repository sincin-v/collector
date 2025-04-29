package compress

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
)

func Compress(data bytes.Buffer) (*bytes.Buffer, error) {
	var b bytes.Buffer

	compressWriter := gzip.NewWriter(&b)
	_, err := compressWriter.Write(data.Bytes())
	defer func() {
		if errCompressWrite := compressWriter.Close(); errCompressWrite != nil {
			err = errors.Join(err, fmt.Errorf("close compress write error: %s", errCompressWrite))
		}
	}()
	if err != nil {

		return nil, err
	}
	errCompresWriter := compressWriter.Close()
	if errCompresWriter != nil {
		return nil, errCompresWriter
	}
	return &b, nil

}
