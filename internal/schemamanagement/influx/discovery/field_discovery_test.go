package discovery

import (
	"fmt"
	"testing"

	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/timescale/outflux/internal/idrf"
	"github.com/timescale/outflux/internal/schemamanagement/influx/influxqueries"
)

func TestDiscoverMeasurementFields(t *testing.T) {
	var mockClient influx.Client
	mockClient = &influxqueries.MockClient{}
	database := "database"
	measure := "measure"

	cases := []testCase{
		{
			expectedError:  true,
			showQueryError: fmt.Errorf("error executing query"),
		}, { // empty result returned, error should be result, must have fields
			expectedError: true,
			showQueryResult: &influxqueries.InfluxShowResult{
				Values: [][]string{},
			},
		}, { // result has more than two columns
			expectedError: true,
			showQueryResult: &influxqueries.InfluxShowResult{
				Values: [][]string{
					{"1", "2", "3"},
				},
			},
		}, {
			expectedError: false,
			showQueryResult: &influxqueries.InfluxShowResult{ // proper result
				Values: [][]string{
					{"1", "boolean"},
					{"2", "float"},
					{"3", "integer"},
					{"4", "string"},
				},
			},
			expectedTags: []*idrf.Column{
				{Name: "1", DataType: idrf.IDRFBoolean},
				{Name: "2", DataType: idrf.IDRFDouble},
				{Name: "3", DataType: idrf.IDRFInteger64},
				{Name: "4", DataType: idrf.IDRFString},
			},
		},
	}

	for _, testCase := range cases {
		fieldExplorer := defaultFieldExplorer{
			queryService: mock(testCase),
		}
		result, err := fieldExplorer.DiscoverMeasurementFields(mockClient, database, measure)
		if err != nil && !testCase.expectedError {
			t.Errorf("unexpected error %v", err)
		} else if err == nil && testCase.expectedError {
			t.Errorf("expected error, none received")
		}

		if testCase.expectedError {
			continue
		}

		expected := testCase.expectedTags
		if len(expected) != len(result) {
			t.Errorf("еxpected result: '%v', got '%v'", expected, result)
		}

		for index, resColumn := range result {
			if resColumn.Name != expected[index].Name || resColumn.DataType != expected[index].DataType {
				t.Errorf("Expected column: %v, got %v", expected[index], resColumn)
			}
		}
	}
}

func TestDiscoverMeasurementFieldsWithRP(t *testing.T) {
	var mockClient influx.Client
	mockClient = &influxqueries.MockClient{}
	database := "database"
	measure := "rp.measure"

	cases := []testCase{
		{
			showQueryResult: &influxqueries.InfluxShowResult{ // proper result
				Values: [][]string{{"1", "boolean"}},
			},
			expectedTags: []*idrf.Column{
				{Name: "1", DataType: idrf.IDRFBoolean},
			},
		},
	}

	for _, testCase := range cases {
		fieldExplorer := defaultFieldExplorer{
			queryService: mock(testCase),
		}
		result, err := fieldExplorer.DiscoverMeasurementFields(mockClient, database, measure)
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}

		expected := testCase.expectedTags
		if len(expected) != len(result) {
			t.Errorf("еxpected result: '%v', got '%v'", expected, result)
		}

		for index, resColumn := range result {
			if resColumn.Name != expected[index].Name || resColumn.DataType != expected[index].DataType {
				t.Errorf("Expected column: %v, got %v", expected[index], resColumn)
			}
		}
	}
}
