package main

import (
	"encoding/json"
	"os"

	"golang.org/x/crypto/bcrypt"
)

var sd *StreamData

type App struct {
	Name               string `json:"name"`
	EncryptedStreamKey []byte `json:"enc_stream_key"`
}

type StreamData struct {
	Apps []App `json:"apps"`
}

func (sd *StreamData) App(sk string) (*App, bool) {
	if sd == nil {
		return nil, false
	}

	var err error
	for _, app := range sd.Apps {
		if err = bcrypt.CompareHashAndPassword(app.EncryptedStreamKey, []byte(sk)); err == nil {
			return &app, true
		}
	}
	return nil, false
}

func loadStreamData(configPath string) (*StreamData, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	var sd StreamData
	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()
	err = dec.Decode(&sd)
	if err != nil {
		return nil, err
	}

	return &sd, f.Close()
}
