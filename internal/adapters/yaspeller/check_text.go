package yaspeller

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Kugeki/kode_test_task/internal/domain"
	"io"
	"net/url"
	"strconv"
)

type CheckTextReq struct {
	Text    string `json:"text"`
	Lang    string `json:"lang"`
	Options int    `json:"options"`
	Format  string `json:"format"`
}

func (r *CheckTextReq) Encode() io.Reader {
	values := url.Values{}
	values.Add("text", r.Text)
	values.Add("lang", r.Lang)
	values.Add("options", strconv.Itoa(r.Options))
	values.Add("format", r.Format)

	var buf bytes.Buffer
	buf.WriteString(values.Encode())

	return &buf
}

var CheckTextEndpoint = "checkText"

func (c *Client) CheckText(ctx context.Context, text, lang, format string, opts ...CheckOpt) (domain.SpellResults, error) {
	options := 0
	for _, op := range opts {
		options = op(options)
	}

	req := &CheckTextReq{
		Text:    text,
		Lang:    lang,
		Options: options,
		Format:  format,
	}

	resp, err := c.post(ctx, CheckTextEndpoint, "application/x-www-form-urlencoded", req.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var results domain.SpellResults
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}
