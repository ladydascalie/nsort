package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/rapidloop/skv"
)

var (
	targetDirectory string
	storePath       string
	store           *skv.KVStore
	home            string
)

func init() {
	u, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}
	home = u.HomeDir

	// Ensure the store folder is on disk
	storePath = filepath.Join(home, ".config/nsort")
	if _, err := os.Stat(storePath); err != nil {
		os.MkdirAll(storePath, 0755)

	}

	// Open the store
	store, err = skv.Open(storePath + "/mappings.db")
	if err != nil {
		log.Fatalln(err)
	}

	// if we just created the store, throw in the default mappings
	for ext, dir := range defaultMappings {
		store.Put(ext, dir)
	}
}

func main() {
	flag.StringVar(&targetDirectory, "t", ".", "nsort -t dir/")
	mapping := flag.String("map", "", "nsort -map go:Source")
	del := flag.String("del", "", "nsort -del go:Source")
	update := flag.String("upd", "", "nsort -upd go:Source")
	flag.Parse()

	Safeguard(targetDirectory)

	switch {
	case *mapping != "":
		if err := addMapping(*mapping); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Successfully added mapping")
		return
	case *del != "":
		if err := delMapping(*del); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Deleted mapping")
		return
	case *update != "":
		if err := delMapping(*update); err != nil {
			fmt.Println(err)
			return
		}
		if err := addMapping(*update); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Successfully updated mapping")
		return
	}

	sort(targetDirectory)
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

// Safeguard provides basic security by preventing sortdir from running
// on the user's home directory.
func Safeguard(dir string) {
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

func isMapped(mapping mapping) (bool, string) {
	var kind string
	if err := store.Get(mapping.Ext, &kind); err == nil {
		return true, kind
	}
	return false, kind
}

func delMapping(mapFlag string) error {
	var m mapping
	m.Unpack(mapFlag)
	if err := store.Delete(m.Ext); err != nil {
		return err
	}
	return nil
}

func addMapping(mapFlag string) error {
	var m mapping
	m.Unpack(mapFlag)
	ok, kind := isMapped(m)
	if ok {
		return fmt.Errorf("[ %s ] already mapped to [ %s ]", m.Ext, kind)
	}

	if err := store.Put(m.Ext, m.Folder); err != nil {
		return err
	}

	return nil
}

// sort borrows heavily from it's forefather sortdir
// but in many ways is cleaner and less confusing
func sort(target string) {
	// list files within the target directory
	files, err := ioutil.ReadDir(target)
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
		var kind string
		// remove the dot from the extension
		ext = strings.TrimPrefix(ext, ".")
		// retrieve mapping for this extension, if it exists
		if err := store.Get(ext, &kind); err != nil {
			fmt.Printf("could not find mapping for: %s\n", file.Name())
			continue
		}

		// Find or create a directory for this extension mapping
		path := filepath.Join(targetDirectory, kind)
		if _, err := os.Stat(path); err != nil {
			os.MkdirAll(path, 0755)
		}

		// Move the file into the corresponding directory
		oldPath := filepath.Join(targetDirectory, file.Name())
		newPath := filepath.Join(path, file.Name())

		if err := os.Rename(oldPath, newPath); err != nil {
			log.Printf("could not move file %s into directory %s", file.Name(), kind)
		}
	}
}

var defaultMappings = map[string]string{
	"mp3":  "Music",
	"aac":  "Music",
	"flac": "Music",
	"ogg":  "Music",
	"wma":  "Music",
	"m4a":  "Music",
	"aiff": "Music",
	"wav":  "Music",
	"amr":  "Music",
	// Videos
	"flv":  "Videos",
	"ogv":  "Videos",
	"avi":  "Videos",
	"mp4":  "Videos",
	"mpg":  "Videos",
	"mpeg": "Videos",
	"3gp":  "Videos",
	"mkv":  "Videos",
	"ts":   "Videos",
	"webm": "Videos",
	"vob":  "Videos",
	"wmv":  "Videos",
	// Pictures
	"png":  "Pictures",
	"jpeg": "Pictures",
	"gif":  "Pictures",
	"jpg":  "Pictures",
	"bmp":  "Pictures",
	"svg":  "Pictures",
	"webp": "Pictures",
	"psd":  "Pictures",
	"tiff": "Pictures",
	// Archives
	"rar":  "Archives",
	"zip":  "Archives",
	"7z":   "Archives",
	"gz":   "Archives",
	"bz2":  "Archives",
	"tar":  "Archives",
	"dmg":  "Archives",
	"tgz":  "Archives",
	"xz":   "Archives",
	"iso":  "Archives",
	"cpio": "Archives",
	// Documents
	"txt":  "Documents",
	"pdf":  "Documents",
	"doc":  "Documents",
	"docx": "Documents",
	"odf":  "Documents",
	"xls":  "Documents",
	"xlsv": "Documents",
	"xlsx": "Documents",
	"ppt":  "Documents",
	"pptx": "Documents",
	"ppsx": "Documents",
	"odp":  "Documents",
	"odt":  "Documents",
	"ods":  "Documents",
	"md":   "Documents",
	"json": "Documents",
	"csv":  "Documents",
	// Books
	"mobi": "Books",
	"epub": "Books",
	"chm":  "Books",
	// deb Packages
	"deb": "DEBPackages",
	// Programs
	"exe": "Programs",
	"msi": "Programs",
	// RPM Packages
	"rpm": "RPMPackages",
}
