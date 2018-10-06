package api

import (
	"context"
	"time"

	"errors"

	"vest.com/vest/go/internal/cloud/dataset"
	"vest.com/vest/go/internal/errutil"
)

type DatasetCode string

const (
	Gold DatasetCode = "WGC"
)

var (
	ErrNotFound = errors.New("not found")
)

func (c DatasetCode) toDatasetCode() (dataset.Code, error) {
	switch c {
	case Gold:
		return dataset.Gold, nil
	default:
		return "", errutil.New("bad dataset code")
	}
}

type APIer interface {
	Chart(_ context.Context, _ DatasetCode, MinTime, MaxTime time.Time) (Chart, error)
}

type API struct {
	DatasetService dataset.QuandlService
}

type Chart struct {
	Dataset
	Metrics
}

type Dataset struct {
	Labels []string
	Values []float64
}

type Metrics struct {
	Start, End time.Time
}

func (a *API) Chart(ctx context.Context, code DatasetCode, MinTime, MaxTime time.Time) (Chart, error) {
	dc, err := code.toDatasetCode()
	if err != nil {
		return Chart{}, errutil.Wrap(err)
	}

	cloudChart, err := a.DatasetService.Get(ctx, dc, MinTime, MaxTime)
	switch err {
	case nil:
	case dataset.ErrNotFound:
		return Chart{}, ErrNotFound
	default:
		return Chart{}, errutil.Wrap(err)
	}

	c := Chart{}
	for k, v := range cloudChart.Data {
		c.Labels = append(c.Labels, k)
		c.Values = append(c.Values, v)
	}
	c.Start = cloudChart.Start
	c.End = cloudChart.End

	return c, nil
}