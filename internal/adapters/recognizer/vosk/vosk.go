package vosk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Config struct {
	Address    string `env:"VOSK_ADDRESS"`
	SampleRate int    `env:"VOSK_SAMPLE_RATE"`
}

type Client struct {
	Address    string
	buffsize   int
	SampleRate int
	socket     *websocket.Conn
	log        *zap.Logger
}

type Option func(*Client)

func SetLogger(l *zap.Logger) Option {
	return func(c *Client) {
		c.log = l
	}
}

func SetSampleRate(rate int) Option {
	return func(c *Client) {
		if rate > 0 {
			c.SampleRate = rate
		}
	}
}

func New(cfg Config, options ...Option) *Client {
	clt := &Client{
		Address:    cfg.Address,
		buffsize:   8000,
		SampleRate: 16000,
		log:        zap.NewNop(),
	}

	for _, opt := range options {
		opt(clt)
	}

	return clt
}

type Message struct {
	Result []struct {
		Conf  float64
		End   float64
		Start float64
		Word  string
	}
	Text string
}

var m Message

func (c *Client) PostConfigure() error {
	u := url.URL{Scheme: "ws", Host: c.Address, Path: ""}
	c.log.Debug("connecting", zap.String("host", u.String()))

	// Opening websocket connection
	soc, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		c.log.Error("failed connect", zap.Error(err))
		return err
	}
	c.socket = soc
	defer c.socket.Close()

	err = soc.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{\"config\" : { \"sample_rate\" : %v } }", c.SampleRate)))
	if err != nil {
		c.log.Error("failed write message with configure", zap.Error(err))
		return err
	}

	err = c.socket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		c.log.Error("failed write message", zap.Error(err))
		return err
	}

	return nil
}

func (c *Client) Recognize(bufWav []byte, language string) (string, error) {
	u := url.URL{Scheme: "ws", Host: c.Address, Path: ""}
	c.log.Debug("connecting", zap.String("host", u.String()))

	// Opening websocket connection
	soc, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		c.log.Error("failed openning connect", zap.Error(err))
		return "", err
	}
	c.socket = soc
	defer c.socket.Close()

	err = soc.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{\"config\" : { \"sample_rate\" : %v } }", c.SampleRate)))
	if err != nil {
		c.log.Error("failed write message with configure", zap.Error(err))
		return "", err
	}

	f := bytes.NewReader(bufWav)
	for {
		buf := make([]byte, c.buffsize)

		dat, err := f.Read(buf)

		if dat == 0 && err == io.EOF {
			err = c.socket.WriteMessage(websocket.TextMessage, []byte("{\"eof\" : 1}"))
			if err != nil {
				c.log.Error("failed write message", zap.Error(err))
				return "", err
			}
			break
		}
		if err != nil {
			c.log.Error("failed read", zap.Error(err))
			return "", err
		}

		err = c.socket.WriteMessage(websocket.BinaryMessage, buf)
		if err != nil {
			c.log.Error("failed write message", zap.Error(err))
			return "", err
		}

		// Read message from server
		_, _, err = c.socket.ReadMessage()
		if err != nil {
			c.log.Error("failed read message", zap.Error(err))
			return "", err
		}
	}

	// Read final message from server
	_, msg, err := c.socket.ReadMessage()
	if err != nil {
		c.log.Error("failed final read message", zap.Error(err))
		return "", err
	}

	// Closing websocket connection
	c.socket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	// Unmarshalling received message
	err = json.Unmarshal(msg, &m)
	if err != nil {
		c.log.Error("failed unmarshal", zap.Error(err))
		return "", err
	}
	return m.Text, nil
}
