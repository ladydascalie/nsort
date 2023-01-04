package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func setupStorage() (storePath string) {
	u, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}
	home := u.HomeDir

	// Ensure the store folder is on disk
	storePath = filepath.Join(home, ".config/nsort")
	if _, err := os.Stat(storePath); err != nil {
		if err := os.MkdirAll(storePath, 0o755); err != nil {
			log.Fatalf("could not create path `%s`: %v", storePath, err)
		}
	}

	return storePath
}

type Storage interface {
	IsMapped(m mapping) (bool, string)
	AddMapping(flag string) error
	DeleteMapping(flag string) error
	GetMapping(ext string) string
}

func main() {
	var (
		targetDirectory string
		mapping         string
		del             string
		update          string
	)

	flag.StringVar(&targetDirectory, "t", ".", "nsort -t dir/")
	flag.StringVar(&mapping, "map", "", "nsort -map go:Source")
	flag.StringVar(&del, "del", "", "nsort -del go:Source")
	flag.StringVar(&update, "upd", "", "nsort -upd go:Source")

	var bykind bool
	flag.BoolVar(&bykind, "by-kind", false, "-by-kind")
	flag.Parse()

	storePath := setupStorage()
	safeguard(targetDirectory)

	store := NewJSONStorage(storePath)

	switch {
	case mapping != "":
		if err := store.AddMapping(mapping); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Successfully added mapping")
		return
	case del != "":
		if err := store.DeleteMapping(del); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Deleted mapping")
		return
	case update != "":
		if err := store.DeleteMapping(update); err != nil {
			fmt.Println(err)
			return
		}
		if err := store.AddMapping(update); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Successfully updated mapping")
		return
	}

	switch bykind {
	case true:
		sortByKind(targetDirectory)
	default:
		sort(targetDirectory, store)
	}
}

type mapping struct {
	Ext    string
	Folder string
}

func (m *mapping) Unpack(mapFlag string) {
	sl := strings.Split(mapFlag, ":")
	m.Ext = sl[0]
	m.Folder = sl[1]
}

// safeguard provides basic security by preventing sortdir from running
// on the user's home directory.
func safeguard(dir string) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error getting pwd: %v", err)
	}
	u, err := user.Current()
	if err != nil {
		log.Fatalf("error getting user: %v", err)
	}

	if dir == "." {
		if pwd == u.HomeDir {
			log.Fatalln("nsort must not be used directly on the home directory")
		}
	} else if dir == u.HomeDir {
		log.Fatalln("nsort must not be used directly on the home directory")
	}
}

func sortByKind(target string) {
	files, err := os.ReadDir(target)
	if err != nil {
		log.Fatalf("could not read directory %q: %v", target, err)
	}

	for _, file := range files {
		ext := filepath.Ext(file.Name())
		// ignore: files without extensions & directories
		if file.Name() == ext || file.IsDir() {
			continue
		}
		// remove dot from extension & make a title
		ext = strings.TrimPrefix(ext, ".")

		// Find or create a directory for this extension mapping
		path := filepath.Join(target, ext)
		if _, err := os.Stat(path); err != nil {
			if err := os.MkdirAll(path, 0o755); err != nil {
				log.Fatalf("could not store path `%s`: %v", path, err)
			}
		}

		// Move the file into the corresponding directory
		oldPath := filepath.Join(target, file.Name())
		newPath := filepath.Join(path, file.Name())

		if err := os.Rename(oldPath, newPath); err != nil {
			log.Printf("could not move file %s into directory %s", file.Name(), ext)
		}
	}
}

// sort borrows heavily from it's forefather sortdir
// but in many ways is cleaner and less confusing
func sort(target string, store Storage) {
	// list files within the target directory
	files, err := os.ReadDir(target)
	if err != nil {
		log.Fatalln(err)
	}

	for _, file := range files {
		ext := filepath.Ext(file.Name())
		// ignore files without extensions
		if file.Name() == ext {
			continue
		}
		// ignore directories
		if file.IsDir() {
			continue
		}
		// remove the dot from the extension
		ext = strings.TrimPrefix(ext, ".")
		// retrieve mapping for this extension, if it exists
		kind := store.GetMapping(ext)
		if kind == "" {
			kind = "Other"
		}

		// Find or create a directory for this extension mapping
		path := filepath.Join(target, kind)
		if _, err := os.Stat(path); err != nil {
			if err := os.MkdirAll(path, 0o755); err != nil {
				log.Fatalf("could not store path `%s`: %v", path, err)
			}
		}

		// Move the file into the corresponding directory
		oldPath := filepath.Join(target, file.Name())
		newPath := filepath.Join(path, file.Name())

		if err := os.Rename(oldPath, newPath); err != nil {
			log.Printf("could not move file %s into directory %s", file.Name(), kind)
		}
	}
}
