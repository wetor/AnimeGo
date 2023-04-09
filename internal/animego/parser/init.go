package parser

var (
	TMDBFailSkip           bool
	TMDBFailUseTitleSeason bool
	TMDBFailUseFirstSeason bool
)

type Options struct {
	TMDBFailSkip           bool
	TMDBFailUseTitleSeason bool
	TMDBFailUseFirstSeason bool
}

func Init(opts *Options) {
	TMDBFailSkip = opts.TMDBFailSkip
	TMDBFailUseTitleSeason = opts.TMDBFailUseTitleSeason
	TMDBFailUseFirstSeason = opts.TMDBFailUseFirstSeason
}
