package models

type TorrentItem struct {
	AddedOn      int     `json:"added_on"`
	AmountLeft   int     `json:"amount_left"`
	AutoTmm      bool    `json:"auto_tmm"`
	Availability float64 `json:"availability"`
	Category     string  `json:"category"`
	// Completed         int     `json:"completed"`
	// CompletionOn      int     `json:"completion_on"`
	ContentPath       string `json:"content_path"`
	DlLimit           int    `json:"dl_limit"`
	Dlspeed           int    `json:"dlspeed"`
	Downloaded        int    `json:"downloaded"`
	DownloadedSession int    `json:"downloaded_session"`
	Eta               int    `json:"eta"`
	// FLPiecePrio       bool    `json:"f_l_piece_prio"`
	// ForceStart        bool    `json:"force_start"`
	Hash string `json:"hash"`
	// LastActivity     int     `json:"last_activity"`
	MagnetUri string `json:"magnet_uri"`
	// MaxRatio         float64 `json:"max_ratio"`
	MaxSeedingTime int    `json:"max_seeding_time"`
	Name           string `json:"name"`
	// NumComplete      int     `json:"num_complete"`
	// NumIncomplete    int     `json:"num_incomplete"`
	// NumLeechs        int     `json:"num_leechs"`
	// NumSeeds         int     `json:"num_seeds"`
	Priority         int     `json:"priority"`
	Progress         float64 `json:"progress"`
	Ratio            float64 `json:"ratio"`
	RatioLimit       float64 `json:"ratio_limit"`
	SavePath         string  `json:"save_path"`
	SeedingTime      int     `json:"seeding_time"`
	SeedingTimeLimit int     `json:"seeding_time_limit"`
	SeenComplete     int     `json:"seen_complete"`
	// SeqDl            bool    `json:"seq_dl"`
	Size  int    `json:"size"`
	State string `json:"state"`
	// SuperSeeding    bool   `json:"super_seeding"`
	Tags string `json:"tags"`
	// TimeActive      int    `json:"time_active"`
	TotalSize int `json:"total_size"`
	// Tracker         string `json:"tracker"`
	// UpLimit         int    `json:"up_limit"`
	//Uploaded        int    `json:"uploaded"`
	// UploadedSession int    `json:"uploaded_session"`
	// Upspeed         int    `json:"upspeed"`
}

type TorrentContentItem struct {
	Index        int     `json:"index"`        //File index
	Name         string  `json:"name"`         //File name (including relative path)
	Size         int     `json:"size"`         //File size (bytes)
	Progress     float64 `json:"progress"`     //File progress (percentage/100)
	Priority     int     `json:"priority"`     //File priority. See possible values here below
	IsSeed       bool    `json:"is_seed"`      //True if file is seeding/complete
	PieceRange   []int   `json:"piece_range"`  //The first number is the starting piece index and the second number is the ending piece index (inclusive)
	Availability float64 `json:"availability"` //Percentage of file pieces currently available (percentage/100)

}

type SetApplicationPreferencesInnerReq struct {
	Json string `json:"json"`
}
type SetApplicationPreferencesReq struct {
	TorrentContentLayout *string `json:"torrent_content_layout,omitempty"`
}
