package event

import (
	"bytes"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/traPtitech/traQ/model"
	"github.com/traPtitech/traQ/repository"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	headerTRAQBotEvent             = "X-TRAQ-BOT-EVENT"
	headerTRAQBotRequestID         = "X-TRAQ-BOT-REQUEST-ID"
	headerTRAQBotVerificationToken = "X-TRAQ-BOT-TOKEN"
	headerUserAgent                = "User-Agent"
	ua                             = "traQ_Bot_Processor/1.0"
)

var eventSendCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "traq",
	Name:      "bot_event_send_count_total",
}, []string{"bot_id", "status"})

type dispatcherImpl struct {
	client http.Client
	l      *zap.Logger
	repo   repository.BotRepository
}

func NewDispatcher(logger *zap.Logger, repo repository.BotRepository) Dispatcher {
	return &dispatcherImpl{
		client: http.Client{
			Jar:     nil,
			Timeout: 5 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		l:    logger.Named("bot.dispatcher"),
		repo: repo,
	}
}

func (d *dispatcherImpl) Send(b *model.Bot, event model.BotEventType, body []byte) (ok bool) {
	reqID := uuid.Must(uuid.NewV4())

	req, _ := http.NewRequest(http.MethodPost, b.PostURL, bytes.NewReader(body))
	req.Header.Set(headerUserAgent, ua)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	req.Header.Set(headerTRAQBotEvent, event.String())
	req.Header.Set(headerTRAQBotRequestID, reqID.String())
	req.Header.Set(headerTRAQBotVerificationToken, b.VerificationToken)

	start := time.Now()
	res, err := d.client.Do(req)
	stop := time.Now()

	if err != nil {
		eventSendCounter.WithLabelValues(b.ID.String(), "ne").Inc()
		d.writeLog(&model.BotEventLog{
			RequestID: reqID,
			BotID:     b.ID,
			Event:     event,
			Body:      string(body),
			Error:     err.Error(),
			Code:      -1,
			Latency:   stop.Sub(start).Nanoseconds(),
			DateTime:  time.Now(),
		})
		return false
	}
	_ = res.Body.Close()

	if res.StatusCode == http.StatusNoContent {
		eventSendCounter.WithLabelValues(b.ID.String(), "ok").Inc()
	} else {
		eventSendCounter.WithLabelValues(b.ID.String(), "ng").Inc()
	}

	d.writeLog(&model.BotEventLog{
		RequestID: reqID,
		BotID:     b.ID,
		Event:     event,
		Body:      string(body),
		Code:      res.StatusCode,
		Latency:   stop.Sub(start).Nanoseconds(),
		DateTime:  time.Now(),
	})
	return res.StatusCode == http.StatusNoContent
}

func (d *dispatcherImpl) writeLog(log *model.BotEventLog) {
	if err := d.repo.WriteBotEventLog(log); err != nil {
		d.l.Warn("failed to write log", zap.Error(err), zap.Any("eventLog", log))
	}
}
