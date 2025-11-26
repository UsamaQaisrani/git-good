package plumbing

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
)

type StageEntry struct {
	CTimeSec  uint32
	CTimeNano uint32
	MTimeSec  uint32
	MTimeNano uint32
	Dev       uint32
	Ino       uint32
	Mode      uint32
	Uid       uint32
	Gid       uint32
	Size      uint32
	Hash      string
	Flags     uint16
	Path      string
}

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

func createNewStagingEntry(path, hash string) (StageEntry, error) {
	info, err := os.Stat(path)
	if err != nil {
		return StageEntry(), err
	}

	stats := info.sys().(*syscall.Stat_t)
	return StageEntry{
		CTimeSec:  uint32(stat.Ctim.Sec),
		CTimeNano: uint32(stat.Ctim.Nsec),
		MtimeSec:  uint32(stat.Mtim.Sec),
		MTimeNano: uint32(stat.Mtim.Nsec),
		Dev:       uint32(stat.Dev),
		Ino:       uint32(stat.Ino),
		Uid:       uint32(stat.Uid),
		Gid:       uint32(stat.Gid),
		Size:      uint32(stat.Size),
		Hash:      hash,
		Path:      path,
	}, nil
}
