package movie

// Constants for Max Length, etc.
const (
	TitleMaxLen           = 10
	LengthInMinutesMaxLen = 1000
)

// Global vars for checking if a particular value exists in a list
var (
	TitleVals = []string{"", "Repo Man", "Goonies"}
)

// File struct has the Raw csv data as well as Good and Bad records
type File struct {
	Raw    []*Raw
	Proc   []*Processed
	Unproc []*Unprocessed
}

// Raw represents data from csv prior to any cleansing
type Raw struct {
	Index           int    `csv:"-"`
	Error           string `csv:"-"`
	Title           string `csv:"Title"`
	LengthInMinutes string `csv:"Length in Minutes"`
	GenreList       string `csv:"Genre List"`
	ReleaseDate     string `csv:"Release Date"`
}

// Processed represents the processed data after validation
type Processed struct {
	Index           int    `csv:"-"`
	Title           string `csv:"Title"`
	LengthInMinutes string `csv:"Length in Minutes"`
	GenreList       string `csv:"Genre List"`
	ReleaseDate     string `csv:"Release Date"`
}

// Unprocessed represents the errant data after validation
type Unprocessed struct {
	Error           string `csv:"Error"`
	Title           string `csv:"Title"`
	LengthInMinutes string `csv:"Length in Minutes"`
	GenreList       string `csv:"Genre List"`
	ReleaseDate     string `csv:"Release Date"`
}
