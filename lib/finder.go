package finder

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var (
	info = color.New(color.FgYellow).SprintFunc()
	warn = color.New(color.FgRed).SprintFunc()
)

// FileFinder ...
type FileFinder struct {
	rootDir         string
	filename        string
	dirname         string
	minSize         int64
	maxSize         int64
	totalFileSize   int64 // accumulates while finding matches
	totalFilesFound int64
	totalDirsFound  int64
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

// Dirname sets the dirname to search for, using wildcards
func (find *FileFinder) Dirname(s string) {
	find.dirname = s
}

// MinSize sets the min size
func (find *FileFinder) MinSize(s string) {
	find.minSize = parseDataSize(s)
}

// MaxSize sets the max size
func (find *FileFinder) MaxSize(s string) {
	find.maxSize = parseDataSize(s)
}

// SearchAndPrint performs search, printing matches to current pattern
func (find *FileFinder) SearchAndPrint() {
	fmt.Println(info("Searching in ", find.rootDir, " for ", find.renderCriterias()))

	filepath.Walk(find.rootDir, func(fp string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Println(warn(err)) // can't walk here,
			return nil             // but continue walking elsewhere
		}
		if fi.IsDir() && find.dirname == "" {
			return nil // not a file, and we're not looking for dirnames
		}
		// cannot search for both
		q := find.filename
		if find.dirname != "" {
			q = find.dirname
		}
		matched, err := filepath.Match(q, fi.Name())
		if err != nil {
			log.Println(err) // malformed pattern
			return err       // this is fatal.
		}
		if find.minSize != 0 {
			if fi.Size() < find.minSize {
				matched = false
			}
		}
		if find.maxSize != 0 {
			if fi.Size() > find.maxSize {
				matched = false
			}
		}
		if matched {
			if fi.IsDir() {
				find.totalDirsFound++
			} else {
				fmt.Println(fp, prettyDataSize(fi.Size()))
				find.totalFileSize += fi.Size()
				find.totalFilesFound++
			}
		}
		return nil
	})

	if find.totalFilesFound > 0 {
		fmt.Println(info("Found ", find.totalFilesFound, " files in ", prettyDataSize(find.totalFileSize)))
	}
	if find.totalDirsFound > 0 {
		fmt.Println(info("Found ", find.totalDirsFound, " directories"))
	}
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
	if find.filename != "*" {
		res = append(res, "filename matching "+find.filename)
	}
	return strings.Join(res, ", ")
}

// parseDataSize converts human readable string into a int
func parseDataSize(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	num := ""
	scale := ""

	for _, c := range s {
		switch {
		case c >= '0' && c <= '9':
			num += string(c)
		default:
			scale += string(c)
		}
	}

	scale = strings.ToLower(strings.TrimSpace(scale))

	val, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	switch scale {
	case "":
		return val
	case "k", "kb", "kib":
		return val * 1024
	case "m", "mb", "mib":
		return val * 1024 * 1024
	case "g", "gb", "gib":
		return val * 1024 * 1024 * 1024
	case "t", "tb", "tib":
		return val * 1024 * 1024 * 1024 * 1024
	}

	log.Fatal("Unknown scale", scale)
	return val
}
