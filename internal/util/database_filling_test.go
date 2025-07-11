package util

import (
	"SearchService/internal/model"
	"reflect"
	"testing"
)

func TestParseCSVRecord(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    model.Advertisement
		wantErr bool
	}{
		{
			name:  "valid record",
			input: []string{"1", "Name", "Description", "Brand", "Category", "123.45", "USD", "10", "Red", "L", "In Stock"},
			want: model.Advertisement{
				Index:        1,
				Name:         "Name",
				Description:  "Description",
				Brand:        "Brand",
				Category:     "Category",
				Price:        123.45,
				Currency:     "USD",
				Stock:        10,
				Color:        "Red",
				Size:         "L",
				Availability: "In Stock",
			},
			wantErr: false,
		},
		{
			name:    "invalid index",
			input:   []string{"notanumber", "Name", "Description"},
			wantErr: true,
		},
		{
			name:    "not enough fields",
			input:   []string{"1", "Name"},
			wantErr: true,
		},
		{
			name:    "invalid price",
			input:   []string{"1", "Name", "Description", "Brand", "Category", "notafloat", "USD", "10", "Red", "L", "In Stock"},
			wantErr: true,
		},
		{
			name:    "invalid stock",
			input:   []string{"1", "Name", "Description", "Brand", "Category", "123.45", "USD", "notanint", "Red", "L", "In Stock"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCSVRecord(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseCSVRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCSVRecord() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
