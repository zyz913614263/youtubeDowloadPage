package server

type UserProfile struct {
	UserName   string
	Email      string
	Total      int
	Day        int
	TotalParse int
	DayParse   int
	// 其他用户信息字段...
}

type Link struct {
	Quality  string
	URL      string
	MimeType string
	Width    int
	Height   int
}
type IndexInfo struct {
	VideoURL     string
	VideoLinks   []*Link
	AudioLinks   []*Link
	UserName     string
	RequestCount int
	ParseCount   int
}

type Format struct {
	FormatID       string  `json:"format_id"`
	URL            string  `json:"url"`
	AudioExt       string  `json:"audio_ext"`
	VideoExt       string  `json:"video_ext"`
	Resolution     string  `json:"resolution"`
	Height         int     `json:"height"`
	ABR            float32 `json:"abr"`
	VBR            float32 `json:"vbr"`
	FormatNote     string  `json:"format_note"`
	FileSize       float32 `json:"filesize"`
	FileSizeApprox float32 `json:"filesize_approx"`
}

type Result struct {
	Title     string   `json:"title"`
	Thumbnail string   `json:"thumbnail"`
	Formats   []Format `json:"formats"`
}

type Audio struct {
	Id     string  `json:"id"`
	Format string  `json:"format"`
	Rate   float32 `json:"rate"`
	Info   string  `json:"info"`
	Size   string  `json:"size"`
}

type Video struct {
	Id     string  `json:"id"`
	Format string  `json:"format"`
	Scale  int     `json:"scale"`
	Rate   float32 `json:"rate"`
	Info   string  `json:"info"`
	Size   string  `json:"size"`
}

type Message struct {
	Website string
	URL     string
	VideoID string
	P       string
}

type HResult struct {
	Website   string    `json:"website"`
	V         string    `json:"v"`
	P         string    `json:"p"`
	Title     string    `json:"title"`
	Thumbnail string    `json:"thumbnail"`
	Best      Best      `json:"best"`
	Available Available `json:"available"`
}

// 使用结构体包装 best 字段
type Best struct {
	Audio Audio `json:"audio"`
	Video Video `json:"video"`
}

// 使用结构体包装 available 字段
type Available struct {
	Audios []Audio `json:"audios"`
	Videos []Video `json:"videos"`
}

type Messages struct {
	Name    string `json:"name" form:"name" binding:"required"`
	Message string `json:"message" form:"message" binding:"required"`
	Time    string `json:"time"`
}
