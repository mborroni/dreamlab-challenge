package conversion

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConversion_IPv4ToDecimal(t *testing.T) {
	type fields struct {
		ipv4 string
	}

	type want struct {
		decimal int64
		err     error
	}

	tests := []struct {
		name         string
		fields       fields
		expectations func(fields fields)
		want         want
	}{
		{name: "ok a",
			fields: fields{
				ipv4: "181.44.9.182",
			},
			want: want{
				decimal: 3039562166,
				err:     nil,
			},
		},
		{name: "ok b",
			fields: fields{
				ipv4: "181.167.120.9",
			},
			want: want{
				decimal: 3047651337,
				err:     nil,
			},
		},
		{name: "error missing dots",
			fields: fields{
				ipv4: "181449182",
			},
			want: want{
				decimal: 0,
				err:     NotIPv4{},
			},
		},
		{name: "error chars in string",
			fields: fields{
				ipv4: "181.abc.120.9",
			},
			want: want{
				decimal: 0,
				err:     NotIPv4{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IPv4ToDecimal(tt.fields.ipv4)
			assert.EqualValues(t, tt.want.decimal, got)
			assert.IsType(t, tt.want.err, err)
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestConversion_DecimalToIPv4(t *testing.T) {
	type fields struct {
		decimal int64
	}

	type want struct {
		ipv4 string
	}

	tests := []struct {
		name         string
		fields       fields
		expectations func(fields fields)
		want         want
	}{
		{name: "ok a",
			fields: fields{
				decimal: 3039562166,
			},
			want: want{
				ipv4: "181.44.9.182",
			},
		},
		{name: "ok b",
			fields: fields{
				decimal: 3047651337,
			},
			want: want{
				ipv4: "181.167.120.9",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DecimalToIPv4(tt.fields.decimal)
			assert.EqualValues(t, tt.want.ipv4, got)
		})
	}
}
