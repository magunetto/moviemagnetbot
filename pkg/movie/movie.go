package movie

// Movie object
type Movie struct {
	imdbID    string
	TMDbID    int
	TMDbURL   string
	Title     string
	Date      string
	mediaType string
}

// New returns a new Movie
func New() Movie {
	return Movie{}
}

// IMDbID returns IMDb ID of a Movie
func (m Movie) IMDbID() string {
	return m.imdbID
}
