package finder

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// FileFinder ...
type FileFinder struct {
	rootDir   string
	filename  string
	minSize   int64
	maxSize   int64
	totalSize int64 // accumulates while finding matches
	totalHits int64
}

// NewFileFinder returns a new FileFinder
func NewFileFinder(inDir string) (*FileFinder, error) {
	if _, err := os.Stat(inDir); err != nil && os.IsNotExist(err) {
		return nil, err
	}
	rootDir, err := filepath.Abs(inDir)
	if err != nil {
		return nil, err
	}
	find := FileFinder{
		rootDir: rootDir,
	}
	return &find, nil
}

// Filename sets the filename to search for, using wildcards
func (find *FileFinder) Filename(s string) {
	if s == "" {
		s = "*"
	}
	find.filename = s
}

// MinSize sets the min size
func (find *FileFinder) MinSize(n int64) {
	find.minSize = n
}

// MaxSize sets the max size
func (find *FileFinder) MaxSize(n int64) {
	find.maxSize = n
}

// SearchAndPrint performs search, printing matches to current pattern
func (find *FileFinder) SearchAndPrint() {
	log.Println("Searching in", find.rootDir, find.renderCriterias())

	filepath.Walk(find.rootDir, func(fp string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Println(err) // can't walk here,
			return nil       // but continue walking elsewhere
		}
		if fi.IsDir() {
			return nil // not a file.  ignore.
		}
		matched, err := filepath.Match(find.filename, fi.Name())
		if err != nil {
			log.Println(err) // malformed pattern
			return err       // this is fatal.
		}
		if find.minSize != 0 {
			if fi.Size() < find.minSize {
				// log.Println("DEBUG: skipping, file too small:", fi.Name(), prettyDataSize(fi.Size()))
				matched = false
			}
		}
		if find.maxSize != 0 {
			if fi.Size() > find.maxSize {
				// log.Println("DEBUG: skipping, file too big:", fi.Name(), prettyDataSize(fi.Size()))
				matched = false
			}
		}
		if matched {
			fmt.Println(fp, prettyDataSize(fi.Size()))
			find.totalSize += fi.Size()
			find.totalHits++
		}
		return nil
	})

	fmt.Println("Found", find.totalHits, "files in", prettyDataSize(find.totalSize))
}

// present data size in proper scale, like "512KiB" or "700GiB"
func prettyDataSize(val int64) string {
	if val < 1024 {
		return fmt.Sprintf("%d", val) + " bytes"
	}
	v := float64(val)
	if v < 1024*1024 {
		return fmt.Sprintf("%.1f", v/1024) + " KiB"
	}
	if v < 1024*1024*1024 {
		return fmt.Sprintf("%.1f", v/(1024*1024)) + " MiB"
	}
	if v < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.1f", v/(1024*1024*1024)) + " GiB"
	}
	return fmt.Sprintf("%.1f", v/(1024*1024*1024*1024)) + " TiB"
}

func (find *FileFinder) renderCriterias() string {
	res := []string{}
	if find.minSize != 0 {
		res = append(res, "at least "+prettyDataSize(find.minSize))
	}
	if find.maxSize != 0 {
		res = append(res, "at max "+prettyDataSize(find.maxSize))
	}
	return strings.Join(res, ", ")
}
