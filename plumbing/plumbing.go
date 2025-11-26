package plumbing

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"path/filepath"
	"encoding/binary"
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

func createIndexInstance(path, hash string) (StageEntry, error) {
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

func createHeaderForIndex(count int) []byte {
	header := make([]byte, 12)
	copy(header[0:4], []byte("DIRC"))
	binary.BigEndian.PutUint32(header[4:8], 2)
	binary.BigEndian.PutUint32(header[8:12], uint32(count))
	return header
}

func createStagingEntry(entry StageEntry) []byte {
	var buffer bytes.Buffer

	binary.Write(&buffer, binary.BigEndian, entry.CTimeSec)
	binary.Write(&buffer, binary.BigEndian, entry.CTimeNano)
	binary.Write(&buffer, binary.BigEndian, entry.MTimeSec)
	binary.Write(&buffer, binary.BigEndian, entry.MTimeNano)
	binary.Write(&buffer, binary.BigEndian, entry.Dev)
	binary.Write(&buffer, binary.BigEndian, entry.Ino)
	binary.Write(&buffer, binary.BigEndian, entry.Mode)
	binary.Write(&buffer, binary.BigEndian, entry.Uid)
	binary.Write(&buffer, binary.BigEndian, entry.Gid)
	binary.Write(&buffer, binary.BigEndian, entry.Size)

	hashBytes, _ := hex.DecodeString(entry.Hash)
	buffer.Write(hashBytes)

	pathLen := len(entry.Path)
	if pathLen > 0xFFF {
		pathLen = 0xFFF
	}

	binary.Write(&buffer, binary.BigEndian, uint16(pathLen))
	buffer.WriteString(entry.Path)

	buffer.WriteByte(0x00)
	totalLen := 62 + pathLen + 1 
	padding := (8 - (totalLen % 8)) % 8

	buffer.Write(make([]byte, padding))
	return buffer.Bytes()
}

func UpdateIndex(entries []IndexEntry) error {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})

	var indexBuf bytes.Buffer
	indexBuf.Write(createHeader(len(entries)))
	for _, entry := range entries {
		indexBuf.Write(createEntryBlock(entry))
	}

	indexBuf.Write(digest[:])
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		os.Mkdir(".git", 0755)
	}

	return os.WriteFile(index, indexBuf.Bytes(), 0644)
}

func WalkDir(rootPath string) <-chan []byte {
	entries := make(chan []byte) 
	go func() {
		defer close(entries)

		_ = filepath.WalkDir(rootPath, func(currPath string, d fs.DirEntry, err error) error { 
			if err != nil {
				return err
			}

			if d.IsDir() {
				if d.Name() == ".git" || d.Name() == ".gitgood" {
					return filepath.SkipDir
				}
				return nil
			}

			normPath := filepath.ToSlash(currPath)
			if strings.Contains(normPath, ".git") || strings.Contains(normPath, ".gitgood") {
				return nil
			}

			content, err := ReadFile(currPath)
			if err != nil {
				return err
			}
			entries <- content
			
			return nil 
		})
	}()
	return entries
}
