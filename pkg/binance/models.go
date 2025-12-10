package binance

// NotificationPayload represents the outer structure of a Binance WebSocket message.
type NotificationPayload struct {
	Type  string `json:"type"`
	Topic string `json:"topic"`
	Data  string `json:"data"`
}

// AnnouncementData represents the parsed content of the Data field.
type AnnouncementData struct {
	CatalogId   int    `json:"catalogId"`
	CatalogName string `json:"catalogName"`
	PublishDate int64  `json:"publishDate"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	Disclaimer  string `json:"disclaimer"`
}
