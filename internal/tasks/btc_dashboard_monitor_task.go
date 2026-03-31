package tasks

import (
	"fmt"
	"time"

	"github.com/ka1fe1/crypto-monitoring/internal/service"
	"github.com/ka1fe1/crypto-monitoring/pkg/logger"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils"
	"github.com/ka1fe1/crypto-monitoring/pkg/utils/alter/dingding"
)

type BtcDashboardMonitorTask struct {
	svc              service.BtcDashboardService
	dingBot          *dingding.DingBot
	ticker           *time.Ticker
	stop             chan bool
	interval         time.Duration
	quietHoursParams utils.QuietHoursParams
	lastRunTime      time.Time
}

func NewBtcDashboardMonitorTask(svc service.BtcDashboardService, dingBot *dingding.DingBot, intervalSeconds int, quietHoursParams utils.QuietHoursParams) *BtcDashboardMonitorTask {
	interval := time.Duration(intervalSeconds) * time.Second
	if interval <= 0 {
		interval = 43200 * time.Second // default to 12 hours
	}

	return &BtcDashboardMonitorTask{
		svc:              svc,
		dingBot:          dingBot,
		stop:             make(chan bool),
		interval:         interval,
		quietHoursParams: quietHoursParams,
	}
}

func (t *BtcDashboardMonitorTask) Start() {
	t.ticker = time.NewTicker(t.interval)
	go func() {
		for {
			select {
			case <-t.ticker.C:
				t.run()
			case <-t.stop:
				t.ticker.Stop()
				return
			}
		}
	}()
}

func (t *BtcDashboardMonitorTask) Stop() {
	t.stop <- true
}

func (t *BtcDashboardMonitorTask) run() {
	if !utils.ShouldExecTask(t.quietHoursParams, t.lastRunTime, t.interval) {
		return
	}
	t.lastRunTime = time.Now()

	metrics, err := t.svc.FetchAndCalculateMetrics()
	if err != nil {
		logger.Error("BtcDashboardMonitorTask fetch metrics failed: %v", err)
		return
	}

	markdownReport := t.svc.GenerateMarkdownReport(metrics)
	var title string
	if t.dingBot.Keyword != "" {
		title = fmt.Sprintf("%s BTC 宏观周期指标", t.dingBot.Keyword)
	} else {
		title = "BTC 宏观周期指标"
	}

	err = t.dingBot.SendMarkdown(title, markdownReport, nil, false)
	if err != nil {
		logger.Error("BtcDashboardMonitorTask failed sending dingtalk message: %v", err)
	} else {
		logger.Info("BtcDashboardMonitorTask sent markdown report successfully")
	}
}
