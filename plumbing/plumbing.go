package plumbing

import (
	"log"
	"os"
	"crypto/sha1"
    "encoding/hex"
	"fmt"
	"bytes"
    "compress/zlib"
)

func readFile(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	fileSize := len(content)
	header := fmt.Sprintf("blob %d\x00", fileSize)
	byteStream := append([]byte(header), content...)
	return byteStream, nil
}

func HashFile(filePath string) {
	stream, err := readFile(filePath)
	if err != nil {
		log.Fatalf("Error while getting content of the %s: %s", filePath, err)
	}

	hash := sha1.New()
	hash.Write(stream)
	sha1_hash := hex.EncodeToString(hash.Sum(nil))
	fmt.Println("Hash: ", sha1_hash)

	compressed, err := compress(stream) 
	if err != nil {
		log.Fatalf("Error while compressing %s: %s", filePath, err)
	}
	fmt.Printf("Compressed (%d bytes): %x\n", len(compressed), compressed)
}

func compress(data []byte) ([]byte, error) {
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
