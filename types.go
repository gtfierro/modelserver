package main

import (
	"encoding/json"
	"errors"
	"github.com/pborman/uuid"
	"time"
)

type Model struct {
	Name                 string   `json:"name"`
	Version              uint64   `json:"version"`
	PredictiveRange      Range    `json:"predictive_range"`
	PredictiveResolution Duration `json:"predictive_resolution"`
	CanPredict           func(Model, Range, time.Duration) bool
	TrainingData         *TrainingData `json:"training_data"` // can be nil
}

var DefaultCanPredict = func(model Model, predictive_range Range, resolution time.Duration) bool {
	if model.TrainingData == nil {
		// TODO: would want custom logic here
		return true
	}
	// predicts data between training data [start, end] + 1 more day
	if model.TrainingData.Range.Start.After(predictive_range.Start) ||
		model.TrainingData.Range.End.Add(time.Duration(24*60*time.Minute)).Before(predictive_range.End) {
		return false
	}

	if model.PredictiveResolution.Duration != time.Duration(60*time.Minute) {
		return false
	}
	return true
}

type Duration struct {
	Duration time.Duration
}

func (dur *Duration) UnmarshalJSON(b []byte) error {
	duration, err := time.ParseDuration("1h")
	if err != nil {
		return err
	}
	dur.Duration = duration
	return nil
}

type Range struct {
	Start, End time.Time
}

func (rng *Range) UnmarshalJSON(b []byte) error {
	var err error
	var m = make(map[string]string)
	if err = json.Unmarshal(b, &m); err != nil {
		return err
	}
	if start, found := m["start"]; !found {
		return errors.New("Range needs 'start'")
	} else if rng.Start, err = time.Parse("2006-01-02 15:04:05 MST", start); err != nil {
		return err
	}

	if end, found := m["end"]; !found {
		return errors.New("Range needs 'end'")
	} else if rng.End, err = time.Parse("2006-01-02 15:04:05 MST", end); err != nil {
		return err
	}
	return nil
}

type TrainingData struct {
	Streams []uuid.UUID
	Range   Range
}

func (data *TrainingData) UnmarshalJSON(b []byte) error {
	var err error
	println(string(b))
	t := &struct {
		Streams []string
		Range   Range
	}{}
	if err = json.Unmarshal(b, t); err != nil {
		return err
	}
	data.Range = t.Range
	data.Streams = make([]uuid.UUID, len(t.Streams))
	for idx, u := range t.Streams {
		data.Streams[idx] = uuid.Parse(u)
	}
	return nil
}
