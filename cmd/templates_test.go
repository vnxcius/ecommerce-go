package main

import (
	"testing"
	"time"

	"github.com/vnxcius/ecommerce-go/internal/assert"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{
			name: "sum positive numbers",
			a:    1,
			b:    2,
			want: 3,
		},
		{
			name: "sum 0",
			a:    0,
			b:    0,
			want: 0,
		},
		{
			name: "sum negative numbers",
			a:    -1,
			b:    2,
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, add(tt.a, tt.b), tt.want)
		})
	}
}

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 6, 24, 10, 15, 0, 0, time.UTC),
			want: "24 Jun 2024 às 10:15",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "GMT-3",
			tm:   time.Date(2024, 6, 24, 10, 15, 0, 0, time.FixedZone("GMT-03", -3*60*60)),
			want: "24 Jun 2024 às 10:15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, humanDate(tt.tm), tt.want)
		})
	}
}

func TestHumanDateShort(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 6, 24, 10, 15, 0, 0, time.UTC),
			want: "24/06/2024",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "GMT-3",
			tm:   time.Date(2024, 6, 24, 10, 15, 0, 0, time.FixedZone("GMT-03", -3*60*60)),
			want: "24/06/2024",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, humanDateShort(tt.tm), tt.want)
		})
	}
}

func TestDateFormat(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 6, 24, 10, 15, 0, 0, time.UTC),
			want: "2024-06-24",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "GMT-3",
			tm:   time.Date(2024, 6, 24, 10, 15, 0, 0, time.FixedZone("GMT-03", -3*60*60)),
			want: "2024-06-24",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, dateFormat(tt.tm), tt.want)
		})
	}
}

func TestConvertPrice(t *testing.T) {
	tests := []struct {
		name string
		p    float64
		want float64
	}{
		{
			name: "price 5.99",
			p:    5.99,
			want: 5.99,
		},
		{
			name: "price 0",
			p:    0,
			want: 0,
		},
		{
			name: "price 0.00",
			p:    0.00,
			want: 0,
		},
		{
			name: "price 0.01",
			p:    0.01,
			want: 0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, convertPrice(tt.p), tt.want)
		})
	}
}

func TestCalculateDiscount(t *testing.T) {
	tests := []struct {
		name string
		p    float64
		d    float64
		want float64
	}{
		{
			name: "price 10, discount 10",
			p:    10,
			d:    10,
			want: 9,
		},
		{
			name: "price 10, discount 0",
			p:    10,
			d:    0,
			want: 10,
		},
		{
			name: "price 10, discount 100",
			p:    5.99,
			d:    100,
			want: 0,
		},
		{
			name: "price 10, discount -10",
			p:    10,
			d:    -10,
			want: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, calculateDiscount(tt.p, tt.d), tt.want)
		})
	}
}
