package main

import (
	"log"

	finder "github.com/martinlindhe/filefinder/lib"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	inDir    = kingpin.Arg("inDir", "Input directory.").String()
	minSize  = kingpin.Flag("min-size", "Minimum size in bytes.").String()
	maxSize  = kingpin.Flag("max-size", "Maximum size in bytes.").String()
	filename = kingpin.Flag("filename", "Filename wildcard match, eg: *.mp3").String()
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

	finder, err := finder.NewFileFinder(*inDir)
	if err != nil {
		log.Fatal(err)
	}
	finder.Filename(*filename)
	finder.MinSize(*minSize)
	finder.MaxSize(*maxSize)
	finder.SearchAndPrint()
}
