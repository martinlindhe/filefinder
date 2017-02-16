package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	inDir   = kingpin.Arg("inDir", "Input directory.").String()
	minSize = kingpin.Flag("min-size", "Minimum size in bytes.").Int64()
	maxSize = kingpin.Flag("max-size", "Maximum size in bytes.").Int64()
)

func init() {
	log.SetFlags(log.Lshortfile)
	kingpin.Parse()
}

func main() {
	if *inDir == "" {
		// log.Println("DEBUG: No indir provided, using current")
		*inDir = "./"
	}

	finder, err := NewFileFinder(*inDir)
	if err != nil {
		log.Fatal(err)
	}

	finder.minSize = *minSize
	finder.maxSize = *maxSize
	finder.SearchAndPrint()
}

// FileFinder ...
type FileFinder struct {
	rootDir   string
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
		matched, err := filepath.Match("*", fi.Name()) // XXX "*.mp3" to match only this extension
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
