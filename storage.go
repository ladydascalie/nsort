package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type JSONStorage struct {
	file    io.ReadWriteCloser
	records map[string]string
}

func NewJSONStorage(storePath string) *JSONStorage {
	store, err := os.OpenFile(storePath+"/mappings.json", os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		log.Fatalf("could not open store: %v", err)
	}

	c := JSONStorage{
		file:    store,
		records: make(map[string]string),
	}

	c.LoadRecords()
	if err := c.WriteRecords(); err != nil {
		log.Fatalf("could not write to store: %v", err)
	}

	return &c
}

func (c *JSONStorage) LoadRecords() {
	records := make(map[string]string)
	if err := json.NewDecoder(c.file).Decode(&records); err != nil {
		log.Fatalf("could not read store: %v", err)
	}

	// merge with the default mappings
	for k, v := range defaultMappings() {
		records[k] = v
	}

	c.records = records
}

func (c *JSONStorage) WriteRecords() error {
	if err := json.NewEncoder(c.file).Encode(c.records); err != nil {
		return fmt.Errorf("could not write to store: %v", err)
	}

	return nil
}

func (c *JSONStorage) IsMapped(m mapping) (bool, string) {
	c.LoadRecords()
	for k, v := range c.records {
		if k == m.Ext && v == m.Folder {
			return true, v
		}
	}

	return false, ""
}

func (c *JSONStorage) DeleteMapping(mapFlag string) error {
	var m mapping
	m.Unpack(mapFlag)

	c.LoadRecords()

	for k, v := range c.records {
		if k == m.Ext && v == m.Folder {
			delete(c.records, k)
		}
	}

	if err := c.WriteRecords(); err != nil {
		return fmt.Errorf("could not write to store: %v", err)
	}

	return nil
}

func (c *JSONStorage) AddMapping(mapFlag string) error {
	var m mapping
	m.Unpack(mapFlag)

	c.LoadRecords()

	for k, v := range c.records {
		if k == m.Ext && v == m.Folder {
			return fmt.Errorf("mapping already exists")
		}
	}

	c.records[m.Ext] = m.Folder

	if err := c.WriteRecords(); err != nil {
		return fmt.Errorf("could not write to store: %v", err)
	}

	return nil
}

func (c *JSONStorage) GetMapping(ext string) string {
	c.LoadRecords()
	return c.records[ext]
}

func defaultMappings() map[string]string {
	return map[string]string{
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
}
