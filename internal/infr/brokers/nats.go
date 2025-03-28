package brokers

import (
	"context"
	"errors"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog/log"
)

var _conn *nats.Conn

func MustConnect() {
	nc, err := nats.Connect("nats://192.168.0.147:4222")
	if err != nil {
		log.Fatal().Err(err).Msg("Не удалось подключиться к NATS")
	}
	log.Info().Msg("Подключено к NATS на 192.168.0.147:4222")

	_conn = nc
}

type NatsBroker struct {
	*nats.Conn
}

func NewNatsBroker() *NatsBroker {
	return &NatsBroker{_conn}
}

func (b *NatsBroker) GetStream(url string) (<-chan []byte, error) {
	return exe(b.Conn, url)
}

func exe(nc *nats.Conn, url string) (<-chan []byte, error) {
	return nil, nil
}

func ps(nc *nats.Conn, url string) <-chan []byte {
	out := make(chan []byte)

	ctx := context.Background()

	// Create a JetStream management interface
	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка получения JetStream контекста")
	}

	s, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "DOWNLOAD",
		Subjects: []string{"DOWNLOADER.*"},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка получения JetStream контекста")
	}

	c, err := s.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "CONS",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка получения JetStream контекста")
	}

	nc.Publish("download.request", []byte(url))

	it, _ := c.Messages()

	// Подписываемся на `DOWNLOADER.complete`
	doneSub, err := nc.Subscribe("DOWNLOADER.complete", func(msg *nats.Msg) {
		it.Stop()
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка подписки на завершение")
	}

	go func() {
		defer doneSub.Unsubscribe()
		defer close(out)

		for {
			msg, err := it.Next()
			if err != nil {
				if errors.Is(err, jetstream.ErrMsgIteratorClosed) {
					log.Info().Msg("Чанки закончились")
					break
				}
				log.Error().Err(err).Msg("Ошибка при получении чанка")
				continue
			}

			out <- msg.Data()

			msg.Ack()
		}
	}()

	return out
}
