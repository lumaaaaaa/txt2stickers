package main

type XFBSticker struct {
	Typename string `json:"__typename"`
	StrongID string `json:"strong_id__"`
	ID       string `json:"id"`
	URL      string `json:"url"`
}

type StickersResponse struct {
	Data struct {
		Typename string       `json:"__typename"`
		Stickers []XFBSticker `json:"stickers"`
		Error    interface{}  `json:"error"`
	} `json:"xfb_pair_generate_text2stickers"`
}

type Response struct {
	Data       StickersResponse `json:"data"`
	Extensions struct {
		IsFinal bool `json:"is_final"`
	} `json:"extensions"`
}

type StickerData struct {
	ImageData []byte
	Filename  string
}
