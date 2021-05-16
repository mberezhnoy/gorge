package sepa

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mattn/go-nulltype"

	"github.com/stretchr/testify/assert"
	"github.com/whitewater-guide/gorge/core"
	"github.com/whitewater-guide/gorge/testutils"
)

func setupTestServer() *httptest.Server {
	// return testutils.SetupFileServer(nil, nil)
	return testutils.SetupFileServer(map[string]string{
		"/10048":                     "10048.csv",
		"/SEPA_River_Levels_Web.csv": "SEPA_River_Levels_Web.csv",
	}, nil)
}

func TestSepa_ListGauges(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()
	s := scriptSepa{
		name:         "sepa",
		listURL:      ts.URL,
		gaugeURLBase: ts.URL,
	}
	actual, err := s.ListGauges()
	expected := core.Gauges{
		core.Gauge{
			GaugeID: core.GaugeID{
				Code:   "10048",
				Script: "sepa",
			},
			Name:      "Tay - Perth",
			LevelUnit: "m",
			FlowUnit:  "",
			Location: &core.Location{
				Latitude:  56.41191,
				Longitude: -3.4342,
				Altitude:  2,
			},
			URL: "http://apps.sepa.org.uk/waterlevels/default.aspx?sd=t&lc=10048",
		},
		core.Gauge{
			GaugeID: core.GaugeID{
				Script: "sepa",
				Code:   "115371",
			},
			Name:      "Lochindorb - Lochindorb Level",
			URL:       "http://apps.sepa.org.uk/waterlevels/default.aspx?sd=t&lc=115371",
			LevelUnit: "m",
			Location: &core.Location{
				Latitude:  57.39798,
				Longitude: -3.71467,
				Altitude:  294,
			},
		},
	}
	if assert.NoError(t, err) {
		assert.Equal(t, expected, actual)
	}
}

func TestSepa_Harvest(t *testing.T) {
	ts := setupTestServer()
	defer ts.Close()
	s := scriptSepa{
		name:         "sepa",
		listURL:      ts.URL,
		gaugeURLBase: ts.URL,
	}
	expected := core.Measurements{
		&core.Measurement{
			GaugeID: core.GaugeID{
				Script: "sepa",
				Code:   "10048",
			},
			Level:     nulltype.NullFloat64Of(2.042),
			Timestamp: core.HTime{Time: time.Date(2020, time.January, 14, 19, 30, 0, 0, time.UTC)},
		},
		&core.Measurement{
			GaugeID: core.GaugeID{
				Script: "sepa",
				Code:   "10048",
			},
			Level:     nulltype.NullFloat64Of(1.984),
			Timestamp: core.HTime{Time: time.Date(2020, time.January, 14, 19, 45, 0, 0, time.UTC)},
		},
	}
	actual, err := core.HarvestSlice(&s, core.StringSet{"10048": {}}, 0)
	if assert.NoError(t, err) {
		assert.Equal(t, expected, actual)
	}
}
