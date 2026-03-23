package main

import (
	"testing"
	"time"

	"github.com/LeeDark/go-web-labs/books/lets-go/snippetbox/internal/assert"
)

func TestHumanDateSingle(t *testing.T) {
	// Init
	tm := time.Date(2026, 3, 23, 13, 28, 0, 0, time.UTC)
	hd := humanDate(tm)

	// Check
	//if hd != "23 Mar 2026 at 13:28" {
	//	t.Errorf("got %q; want %q", hd, "23 Mar 2026 at 13:28")
	//}

	assert.Equal(t, hd, "23 Mar 2026 at 13:28")
}

func TestHumanDateTable(t *testing.T) {
	// Table
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2024 at 10:15",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2024 at 09:15",
		},
	}

	// Loop
	for _, tt := range tests {
		// Sub-test
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			//if hd != tt.want {
			//	t.Errorf("got %q; want %q", hd, tt.want)
			//}

			assert.Equal(t, hd, tt.want)
		})
	}
}
