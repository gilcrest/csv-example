package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gilcrest/csv/lib/movie"
	"github.com/gocarina/gocsv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	debug := flag.Bool("debug", false, "sets log level to debug")
	arc := flag.Bool("archive", true, "archives file upon processing")

	flag.Parse()

	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Returns a slice of strings which have the names of all
	// files matching the pattern given (movielist*)
	files, err := filepath.Glob("../in/movielist*")
	log.Debug().Msgf("files=%s", files)

	if err != nil {
		panic(err)
	}

	// for each file in the list
	for _, f := range files {

		// Removes the prefix "movielist" and returns everything after it
		suffix := strings.TrimLeft(f, "../in/movielist")
		log.Debug().Msgf("suffix=%s", suffix)
		// Creates a new prefix with the archive folder in the path
		prefix := "../archive/movielist"
		log.Debug().Msgf("prefix=%s", prefix)
		// concatenate the prefix and suffix strings to create the full
		// path/filename
		cf := prefix + suffix

		// Archive files if flag is set to true
		if *arc {
			err := archive(f, cf)
			if err != nil {
				panic(err)
			}
		}

		moviesFile, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			panic(err)
		}
		defer moviesFile.Close()

		// initialize movie as a slice of nil movie.Raw pointers
		raw := []*movie.Raw{}

		// Unmarshal moviesFile into movies
		if err := gocsv.UnmarshalFile(moviesFile, &raw); err != nil {
			panic(err)
		}

		dfile := movie.File{}
		sp := []*movie.Processed{}

		err = movie.SetRawFile(&dfile, raw, sp)
		if err != nil {
			panic(err)
		}

		err = dfile.Process()
		if err != nil {
			panic(err)
		}

		err = writeGoodFile(&dfile, suffix)
		if err != nil {
			log.Fatal().Err(err).Msg("Error from writeGoodFile")
		}

		err = writeBadFile(&dfile, suffix)
		if err != nil {
			log.Error().Err(err).Msg("Error from writeBadFile")
		}

		if *debug {
			err = printFile(&dfile, "good")
			if err != nil {
				log.Error().Err(err).Msg("Error from printFile")
			}
		}

		// remove file from /in directory
		if *arc {
			err = os.Remove(f)
			if err != nil {
				log.Error().Err(err).Msg("Error from os.Remove")
			}
		}
	}
}

func archive(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func writeGoodFile(d *movie.File, suffix string) error {

	fileName := strings.Join([]string{"../out/movielist", suffix}, "")

	file := filepath.Clean(fileName)

	gf, err := os.Create(file)
	if err != nil {
		return err
	}

	defer gf.Close()

	// Use this to save the CSV back to the file
	err = gocsv.MarshalFile(&d.Proc, gf)
	if err != nil {
		return err
	}

	return nil

}

func writeBadFile(d *movie.File, suffix string) error {

	fileName := strings.Join([]string{"../bad/movielist", suffix}, "")

	file := filepath.Clean(fileName)

	err := d.ProcessBad()

	bf, err := os.Create(file)
	if err != nil {
		return err
	}

	defer bf.Close()

	// Use this to save the CSV back to the file
	err = gocsv.MarshalFile(&d.Unproc, bf)
	if err != nil {
		return err
	}

	return nil

}

func printFile(d *movie.File, file string) error {

	if file == "good" {
		// Get file contents as CSV string
		csvContent, err := gocsv.MarshalString(d.Proc)
		if err != nil {
			return err
		}

		// Display all contents as CSV string
		fmt.Println(csvContent)

	} else {
		// Get file contents as CSV string
		csvContent, err := gocsv.MarshalString(d.Raw)
		if err != nil {
			panic(err)
		}

		// Display all contents as CSV string
		fmt.Println(csvContent)

	}

	return nil
}
