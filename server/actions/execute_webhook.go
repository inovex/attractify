package actions

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

const timeout = time.Second * 30

type executeWebhook struct {
	URL          string          `json:"url"`
	ReturnResult bool            `json:"returnResult"`
	Properties   json.RawMessage `json:"params"`
}

type payload struct {
	Event            string          `json:"event"`
	ActionID         uuid.UUID       `json:"actionId"`
	UserID           string          `json:"userId"`
	Channel          string          `json:"channel"`
	UserProperties   json.RawMessage `json:"userProperties,omitempty"`
	CustomProperties json.RawMessage `json:"customProperties,omitempty"`
	Timestamp        time.Time       `json:"timestamp"`
}

type result struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
	Signature  string `json:"signature"`
	Error      string `json:"error"`
}

func (h Hook) ExecuteWebhook() (*result, error) {
	config := executeWebhook{}
	if err := json.Unmarshal(h.Config, &config); err != nil {
		return nil, err
	}

	userProps := json.RawMessage("null")
	if h.Properties != nil {
		userProps = *h.Properties
	}

	props, _ := strconv.Unquote(string(config.Properties))
	pl := payload{
		Event:            h.Event,
		ActionID:         h.Action.ID,
		UserID:           h.UserID,
		Channel:          h.Channel,
		UserProperties:   userProps,
		CustomProperties: json.RawMessage(props),
		Timestamp:        time.Now().UTC(),
	}

	buf, err := json.Marshal(pl)
	if err != nil {
		return nil, err
	}

	key, err := h.App.DB.GetKeyForOrganization(context.Background(), h.OrganizationID)
	if err != nil {
		return nil, err
	}
	signature := h.calcSignature(buf, key)

	req, err := http.NewRequest("POST", config.URL, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Attractify-Signature", signature)
	req.Header.Set("User-Agent", "Attractify Webhook Client")

	res := result{Signature: signature}

	client := &http.Client{Timeout: timeout}
	webhookRes, err := client.Do(req)
	if err != nil {
		res.StatusCode = -1
		res.Error = err.Error()
		return nil, err
	}

	body, err := ioutil.ReadAll(webhookRes.Body)
	res.StatusCode = webhookRes.StatusCode
	res.Body = string(body)

	h.App.Logger.Debug("actions.webhook", zap.Any("result", res))

	if config.ReturnResult {
		return &res, err
	}

	return nil, err
}

func (h Hook) calcSignature(digest, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(digest)
	return hex.EncodeToString(mac.Sum(nil))
}
