package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ka1fe1/crypto-monitoring/pkg/binance"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alter/dingding"
)

type BinanceAnnouncementService struct {
	subscriber *binance.CmsSubscriber
	dingBot    *dingding.DingBot
}

func NewBinanceAnnouncementService(apiKey, secretKey, proxyURL string, dingBot *dingding.DingBot) *BinanceAnnouncementService {
	sub := binance.NewCmsSubscriber(apiKey, secretKey)
	if proxyURL != "" {
		sub.SetProxy(proxyURL)
	}
	return &BinanceAnnouncementService{
		subscriber: sub,
		dingBot:    dingBot,
	}
}

func (s *BinanceAnnouncementService) Start(ctx context.Context) {
	topics := []string{"com_announcement_en", "com_announcement_cn"}

	s.subscriber.Run(ctx, topics, func(msg []byte) {
		log.Printf("Received Binance Message: %s", msg)

		var payload binance.NotificationPayload
		if err := json.Unmarshal(msg, &payload); err == nil && payload.Type == "DATA" {
			var data binance.AnnouncementData
			if err := json.Unmarshal([]byte(payload.Data), &data); err == nil {
				// Format Markdown
				title := fmt.Sprintf("%s", data.Title)
				text := fmt.Sprintf("**Topic**: %s\n\n**Catalog**: %s\n\n**Time**: %s\n\n%s",
					payload.Topic,
					data.CatalogName,
					time.Unix(data.PublishDate/1000, 0).Format(time.RFC3339),
					data.Body)

				if err := s.dingBot.SendMarkdown(title, text, nil, false); err != nil {
					log.Printf("Failed to send DingTalk markdown: %v", err)
				}
				return
			}
		}

		// Fallback to Text
		if err := s.dingBot.SendText(string(msg), nil, false); err != nil {
			log.Printf("Failed to send DingTalk message: %v", err)
		}
	})
}
