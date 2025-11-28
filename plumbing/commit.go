package plumbing

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type Node struct {
	Name     string
	Mode     uint32
	Hash     string
	Children []*Node
}

func ReadIndex() error {
	content, err := ReadFile(index)
	if err != nil {
		return err
	}

	header := string(content[0:4])
	if string(header) != "DIRC" {
		return errors.New("Invalid header, not an index file.")
	}

	signature := binary.BigEndian.Uint32(content[4:8])
	entryCount := binary.BigEndian.Uint32(content[8:12])
	fmt.Println(header)
	fmt.Println("Version:", signature)
	fmt.Println("Entires count:", entryCount)

	var entries []StageEntry
	i := 12

	for j := 0; j < 7; j++ {
		cTimeSec := binary.BigEndian.Uint32(content[i : i+4])
		i += 4
		cTimeNano := binary.BigEndian.Uint32(content[i : i+4])
		i += 4
		mTimeSec := binary.BigEndian.Uint32(content[i : i+4])
		i += 4
		mTimeNano := binary.BigEndian.Uint32(content[i : i+4])
		i += 4
		dev := binary.BigEndian.Uint32(content[i : i+4])
		i += 4
		ino := binary.BigEndian.Uint32(content[i : i+4])
		i += 4
		mode := binary.BigEndian.Uint32(content[i : i+4])
		i += 4
		uid := binary.BigEndian.Uint32(content[i : i+4])
		i += 4
		gid := binary.BigEndian.Uint32(content[i : i+4])
		i += 4
		size := binary.BigEndian.Uint32(content[i : i+4])
		i += 4
		hash := hex.EncodeToString(content[i : i+20])
		i += 20
		pathLen := binary.BigEndian.Uint16(content[i : i+2])
		i += 2
		path := string(content[i : i+int(pathLen)])

		i += int(pathLen)

		// Skipping the null byte (0x00)
		i += 1

		// Reading the padding
		entrySize := 62 + int(pathLen) + 1
		padding := (8 - (entrySize % 8)) % 8
		i += padding

		fmt.Printf("%d %s %d %s\n", mode, hash, 0, path)

		entry := StageEntry{
			CTimeSec:  cTimeSec,
			CTimeNano: cTimeNano,
			MTimeSec:  mTimeSec,
			MTimeNano: mTimeNano,
			Dev:       dev,
			Ino:       ino,
			Mode:      mode,
			Uid:       uid,
			Gid:       gid,
			Size:      size,
			Hash:      hash,
			Path:      path,
		}

		entries = append(entries, entry)
	}

	return nil
}

func GenerateTree() (*Node, error) {
	root := "."
	tree := &Node {
		Name: filepath.Base(root),
		Mode: 0x81A4,
	}

	nodeMap := map[string]*Node{
		root:tree,
	}

	err := filepath.WalkDir(root, func(currPath string, d fs.DirEntry, err error) error {
			normPath := filepath.ToSlash(currPath)
			if strings.Contains(normPath, ".git") || strings.Contains(normPath, git) {
				return nil
			}

			parent := nodeMap[filepath.Dir(currPath)]

			node := &Node{
				Name: d.Name(),
				Mode: 0x81A4,
			}

			if !d.IsDir() {
				content, err := ReadFile(currPath)
				if err != nil {
					return err
				}

				hash := HashFile(content)
				node.Hash = hash
			}
			
			parent.Children = append(parent.Children, node)

			if d.IsDir() {
				nodeMap[currPath] = node
			}

		return nil
	})


	if err != nil {
		return tree, err
	}

	fmt.Println("Map:", nodeMap)
	return tree, nil
}
