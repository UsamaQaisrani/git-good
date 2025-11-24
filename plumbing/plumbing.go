package plumbing

import (
	"os"
	"crypto/sha1"
    "encoding/hex"
	"fmt"
	"bytes"
    "compress/zlib"
)

func ReadFile(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	fileSize := len(content)
	header := fmt.Sprintf("blob %d\x00", fileSize)
	byteStream := append([]byte(header), content...)
	return byteStream, nil
}

func HashFile(content []byte) string {
	hash := sha1.New()
	hash.Write(content)
	sha1_hash := hex.EncodeToString(hash.Sum(nil))
	fmt.Println("Hash: ", sha1_hash)
	return sha1_hash
}

func Compress(data []byte) ([]byte, error) {
	var buff bytes.Buffer
	w := zlib.NewWriter(&buff)

	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
