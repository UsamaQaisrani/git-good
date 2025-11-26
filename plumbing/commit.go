package plumbing

import (
	"fmt"
	"errors"
	"encoding/binary"
)

func ReadIndex() error {
	content, err := ReadFile(".gitgood/index")
	if err != nil {
		return err
	}

	header := content[0:4]
	if string(header) != "DIRC" {
		return errors.New("Invalid header, not an index file.")
	}

	signature := binary.BigEndian.Uint32(content[4:8])
	fmt.Println("Version:", signature)
	entryCount := binary.BigEndian.Uint32(content[8:12])
	fmt.Println("Entires count:", entryCount)

	return nil
}
