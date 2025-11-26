package porcelain

import (
	"fmt"
	"log"
	"usamaqaisrani/git-good/plumbing"
)

// Paths
const git = ".gitgood"
const config = git + "/config"
const refs = git + "/refs"
const refHeads = refs + "/heads"
const head = git + "/HEAD"
const objects = git + "/objects"
const index = git + "/index"

func Init() {
	plumbing.CreateDir(git)
	plumbing.CreateDir(refs)
	plumbing.CreateDir(refHeads)
	plumbing.CreateDir(objects)
	headContent := "ref: refs/heads/master\n"
	configContent := "[core]\n\trepositoryformatversion = 0\n\tfilemode = true\n\tbare = false\n\tlogallrefupdates = true"
	plumbing.WriteFile(head, headContent)
	plumbing.WriteFile(config, configContent)
}


// Stage the file/dir at given path
func Stage(path string) {
	files := plumbing.WalkDir(path)
	var entries []plumbing.StageEntry
	
	for file := range files {
		if file.Err != nil {
			fmt.Printf("Warning: Skipping file due to error: %v\n", file.Err)
			continue
		}

		hash := plumbing.HashFile(file.Content)
		err := plumbing.WriteBlob(file.Content, hash)
		if err != nil {
			log.Fatal("Error while writing blob:", err)
			return
		}

		entry, err := plumbing.CreateIndexInstance(file.Path, hash)
		if err != nil {
			log.Fatal("Error while creating instance for index file:", err)
			return	
		}

		entries = append(entries, entry)
	}

	if len(entries) > 0 {
		err := plumbing.UpdateIndex(entries)
		if err != nil {
			fmt.Println("Failed to write the index file.")
			return
		}
	} else {
		fmt.Println("No files to stage.")
	}
}
