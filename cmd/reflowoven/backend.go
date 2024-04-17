package main

import (
	"fmt"
	"github.com/minor-industries/rtgraph/schema"
	"github.com/minor-industries/rtgraph/storage"
	"time"
)

type backend struct {
	t0            time.Time
	normalBackend storage.StorageBackend
	profile       Schedule
}

func (b *backend) getProfile(profileName string) (schema.Series, error) {
	profile := map[string]Schedule{
		"reflowoven_profile":   b.profile,
		"reflowoven_profile_1": profile1,
		"reflowoven_profile_2": profile2,
	}[profileName]

	var values []schema.Value

	for _, point := range profile {
		t := point.T()
		ts := b.t0.Add(t)
		values = append(values, schema.Value{
			Timestamp: ts,
			Value:     point.Val,
		})
	}

	for _, value := range values {
		fmt.Println(value.Timestamp, value.Value)
	}

	return schema.Series{
		SeriesName: profileName,
		Values:     values,
	}, nil
}

func (b *backend) LoadDataWindow(seriesName string, start time.Time) (schema.Series, error) {
	switch seriesName {
	case "reflowoven_profile_1", "reflowoven_profile_2":
		return b.getProfile(seriesName)
	default:
		return b.normalBackend.LoadDataWindow(seriesName, start)
	}
}

func (b *backend) CreateSeries(seriesNames []string) error {
	return b.normalBackend.CreateSeries(seriesNames)
}

func (b *backend) InsertValue(seriesName string, timestamp time.Time, value float64) error {
	return b.normalBackend.InsertValue(seriesName, timestamp, value)
}
