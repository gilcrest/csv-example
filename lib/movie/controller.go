package movie

import (
	"errors"
	"fmt"
	"time"
)

// ErrMissingField used for Required fields
type ErrMissingField string

func (e ErrMissingField) Error() string {
	return string(e) + " is required"
}

// ErrInvalidDate used for Date fields
type ErrInvalidDate struct {
	fieldName string
	err       error
}

func (e ErrInvalidDate) Error() string {
	errStr := fmt.Sprintf("%s is invalid: %s", e.fieldName, e.err.Error())
	return errStr
}

// Checks if a string value exists in a given slice
func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

// validates that a date can be parsed in the given format MM/DD/YYYY
func validateDate(date string) error {
	_, err := time.Parse("01/02/2006", date)
	if err != nil {
		return err
	}

	return nil

}

// validates field max length
func validateMaxLen(s string, fieldName string, maxLen int) (int, error) {
	l := len([]rune(s))
	if l > maxLen {
		errStr := fmt.Sprintf("%s string length > max, Max: %d, Actual: %d", fieldName, maxLen, l)
		return 0, errors.New(errStr)
	}

	return l, nil
}

// SetRawFile performs input checks on raw data from csv input
func SetRawFile(d *File, r []*Raw, p []*Processed) error {

	// Add raw data from file to struct
	d.Raw = r

	// Add index as a record #
	for i, r := range d.Raw {
		r.Index = i
	}

	return nil
}

// Process contents of movielist file
func (d *File) Process() error {

	var (
		err error
	)

	for _, r := range d.Raw {
		p := &Processed{}

		err = p.setTitle(r.Title)
		if err != nil {
			r.Error = err.Error()
			continue
		}

		err = p.setLengthInMinutes(r.LengthInMinutes)
		if err != nil {
			r.Error = err.Error()
			continue
		}

		err = p.setGenreList(r.GenreList)
		if err != nil {
			r.Error = err.Error()
			continue
		}

		err = p.setReleaseDate(r.ReleaseDate)
		if err != nil {
			r.Error = err.Error()
			continue
		}

		d.Proc = append(d.Proc, p)
	}

	return nil
}

// ProcessBad loops through raw csv records with an error and
// adds them to the Unprocessed slice in Dfile
func (d *File) ProcessBad() error {

	for _, r := range d.Raw {
		b := &Unprocessed{}

		if r.Error != "" {
			b.Error = r.Error
			b.Title = r.Title
			b.LengthInMinutes = r.LengthInMinutes
			b.GenreList = r.GenreList
			b.ReleaseDate = r.ReleaseDate

			d.Unproc = append(d.Unproc, b)
		}
	}

	return nil
}
