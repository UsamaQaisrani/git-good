# git-good

This is a minimal implementation of Git (Version Control System), built only to get a better understanding of git and how it works under the hood, also to learn Go as well.

**Current Status:** Minimal functionality of VCS is complete, reads files, adds to index and creates tree and stages the files. Might implement further git commands in future as well for fun.

## What it does
It implements core Git functionality but stores everything in a **`.gitgood`** folder.

## Supported Commands

* **`init`**: Reads the staging area files.
* **`read-index`**: Reads the staging area files.
* **`add <path>`** Adds the directory to the staging file.
* **`write-tree`**: Scans the directory and creates a Tree Object (a snapshot of your files).
