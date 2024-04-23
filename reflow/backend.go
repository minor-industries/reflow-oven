package reflow

import (
	"github.com/minor-industries/rtgraph/schema"
	"github.com/minor-industries/rtgraph/storage"
	"time"
)

type Backend struct {
	T0            time.Time
	NormalBackend storage.StorageBackend
	Profile       Schedule
}

func (b *Backend) getProfile(profileName string) (schema.Series, error) {
	profile := map[string]Schedule{
		"reflowoven_profile":   b.Profile,
		"reflowoven_profile_1": Profile1,
		"reflowoven_profile_2": Profile2,
	}[profileName]

	var values []schema.Value

	for _, point := range profile {
		t := point.T()
		ts := b.T0.Add(t)
		values = append(values, schema.Value{
			Timestamp: ts,
			Value:     point.Val,
		})
	}

	return schema.Series{
		SeriesName: profileName,
		Values:     values,
	}, nil
}

func (b *Backend) LoadDataWindow(seriesName string, start time.Time) (schema.Series, error) {
	switch seriesName {
	case "reflowoven_profile_1", "reflowoven_profile_2":
		return b.getProfile(seriesName)
	default:
		return b.NormalBackend.LoadDataWindow(seriesName, start)
	}
}

func (b *Backend) CreateSeries(seriesNames []string) error {
	return b.NormalBackend.CreateSeries(seriesNames)
}

func (b *Backend) InsertValue(seriesName string, timestamp time.Time, value float64) error {
	return b.NormalBackend.InsertValue(seriesName, timestamp, value)
}
