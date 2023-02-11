package tiktok

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/janitorjeff/jeff-bot/core"
	"github.com/janitorjeff/jeff-bot/frontends"

	"github.com/janitorjeff/gosafe"
)

var Hooks = gosafe.Map[string, int]{}

var ErrHookNotFound = errors.New("Wasn't monitoring, what are you even trynna do??")

type TTSResp struct {
	Data struct {
		SKey     string `json:"s_key"`
		VStr     string `json:"v_str"`
		Duration string `json:"duration"`
		Speaker  string `json:"speaker"`
	} `json:"data"`
	Extra struct {
		LogID string `json:"log_id"`
	} `json:"extra"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

// TTS will return a slice of bytes containing the audio generated by the TikTok
// TTS. You need to have the TikTokSessionID global set.
func TTS(text string) ([]byte, error) {
	reqURL := "https://api16-normal-useast5.us.tiktokv.com/media/api/text/speech/invoke/?"
	reqURL += "text_speaker=" + "en_us_002"
	reqURL += "&req_text=" + url.QueryEscape(text)
	reqURL += "&speaker_map_type=0&aid=1233"

	client := &http.Client{}
	req, err := http.NewRequest("POST", reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Cookie": {"sessionid=" + core.TikTokSessionID},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data TTSResp
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(data.Data.VStr)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

// Play will, if necessary join the appropriate voice channel, and start playing
// the TTS specified by text.
func Play(sp core.Speaker, text string) error {
	audio, err := TTS(text)
	if err != nil {
		return err
	}

	err = sp.Join()
	if err != nil {
		return err
	}

	state := &core.State{}
	state.Set(core.Play)

	buf := ioutil.NopCloser(bytes.NewReader(audio))
	core.FFmpegBufferPipe(sp, buf, state)

	return nil
}

// Start will create a hook and will monitor all incoming messages, if they
// are from twitch and match the specified username then the the TTS audio will
// be sent to the specified speaker.
func Start(sp core.Speaker, twitchUsername string) {
	id := core.Hooks.Register(func(m *core.Message) {
		if m.Frontend != frontends.Twitch || m.Here.Name() != twitchUsername {
			return
		}
		Play(sp, m.Raw)
	})
	Hooks.Set(twitchUsername, id)
}

// Stop will delete the hook created by Start. Returns ErrHookNotFound if the
// hook doesn't exist.
func Stop(twitchUsername string) error {
	id, ok := Hooks.Get(twitchUsername)
	if !ok {
		return ErrHookNotFound
	}
	core.Hooks.Delete(id)
	Hooks.Delete(twitchUsername)
	return nil
}
