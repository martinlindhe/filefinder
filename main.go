package main

import (
	"fmt"

	"gopkg.in/alecthomas/kingpin.v2"

	"log"
	"os"
	"path/filepath"
)

var (
	inDir   = kingpin.Arg("inDir", "Input directory.").String()
	minSize = kingpin.Flag("min-size", "Minimum size in bytes.").Int64()
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
	finder.minSize = *minSize
	if err != nil {
		log.Fatal(err)
	}

	finder.SearchAndPrint()
}

type FileFinder struct {
	rootDir string
	minSize int64
}

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

func (find *FileFinder) SearchAndPrint() {

	log.Println("Searching in", find.rootDir)

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

		// check additional criteria
		if find.minSize != 0 {
			if fi.Size() < find.minSize {
				// log.Println("DEBUG: file too small so hiding", fi.Name())
				matched = false
			}
		}

		if matched {
			fmt.Println(fp, prettyDataSize(fi.Size()))
		}
		return nil
	})
}

// present data size in proper scale, like "512KiB" or "700GiB"
func prettyDataSize(val int64) string {
	if val < 1024 {
		return fmt.Sprintf("%d", val) + "b"
	}
	v := float64(val)
	if v < 1024*1024 {
		return fmt.Sprintf("%.1f", v/1024) + "KiB"
	}
	if v < 1024*1024*1024 {
		return fmt.Sprintf("%.1f", v/(1024*1024)) + "MiB"
	}
	if v < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.1f", v/(1024*1024*1024)) + "GiB"
	}
	return fmt.Sprintf("%.1f", v/(1024*1024*1024*1024)) + "TiB"
}

// exists reports whether the named file or directory exists.
func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
