package movie

import (
	"errors"
)

func (p *Processed) setTitle(s string) error {
	const fieldName = "Title"

	var (
		err error
	)

	if s == "" {
		return ErrMissingField(fieldName)
	}

	_, err = validateMaxLen(s, fieldName, TitleMaxLen)
	if err != nil {
		return err
	}

	if !stringInSlice(s, TitleVals) {
		return errors.New("Title not in Title List")
	}

	p.Title = s

	return nil
}

func (p *Processed) setLengthInMinutes(s string) error {
	const fieldName = "LengthInMinutes"

	var (
		err error
	)

	if s == "" {
		return ErrMissingField(fieldName)
	}

	_, err = validateMaxLen(s, fieldName, LengthInMinutesMaxLen)
	if err != nil {
		return err
	}

	p.LengthInMinutes = s

	return nil
}

func (p *Processed) setGenreList(s string) error {
	const fieldName = "GenreList"

	if s == "" {
		return ErrMissingField(fieldName)
	}

	p.GenreList = s

	return nil
}

func (p *Processed) setReleaseDate(s string) error {
	const fieldName = "ReleaseDate"

	var (
		err error
	)

	if s == "" {
		return ErrMissingField(fieldName)
	}

	err = validateDate(s)
	if err != nil {
		return ErrInvalidDate{fieldName, err}
	}

	p.ReleaseDate = s

	return nil
}
