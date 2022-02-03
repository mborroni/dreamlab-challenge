package ips

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/mborroni/dreamlab-challenge/internal/conversion"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newMockAddressesService(ctrl *gomock.Controller) *AddressesService {
	return NewAddressesService(NewMockrepository(ctrl))
}

func TestAddressesService_Get(t *testing.T) {
	ass := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := newMockAddressesService(ctrl)

	type fields struct {
		ip string
	}

	type want struct {
		ip  *IP
		err error
	}

	tests := []struct {
		name         string
		fields       fields
		expectations func(fields fields)
		want         want
	}{
		{name: "ok",
			fields: fields{
				ip: "8.243.138.218",
			},
			expectations: func(fields fields) {
				service.repository.(*Mockrepository).
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(dummyIP1, nil)
			},
			want: want{
				ip:  dummyIP1,
				err: nil,
			},
		},
		{name: "invalid IP",
			fields: fields{
				ip: "181.ab.9.182",
			},
			expectations: func(fields fields) {
				service.repository.(*Mockrepository).
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(nil, conversion.NotIPv4{})
			},
			want: want{
				ip:  nil,
				err: conversion.NotIPv4{},
			},
		},
		{name: "no content",
			fields: fields{
				ip: "181.44.9.182",
			},
			expectations: func(fields fields) {
				service.repository.(*Mockrepository).
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrNoRows)
			},
			want: want{
				ip:  nil,
				err: nil,
			},
		},
		{name: "error",
			fields: fields{
				ip: "8.243.138.218",
			},
			expectations: func(fields fields) {
				service.repository.(*Mockrepository).
					EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone)
			},
			want: want{
				ip:  nil,
				err: sql.ErrConnDone,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expectations(tt.fields)
			got, err := service.Get(context.Background(), tt.fields.ip)
			ass.EqualValues(tt.want.ip, got)
			ass.IsType(tt.want.err, err)
			ass.Equal(tt.want.err, err)
		})
	}
}

func TestAddressesService_List(t *testing.T) {
	ass := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := newMockAddressesService(ctrl)

	type fields struct {
		limit   int
		filters map[string]interface{}
	}

	type want struct {
		ips []*IP
		err error
	}

	tests := []struct {
		name         string
		fields       fields
		expectations func(fields fields)
		want         want
	}{
		{name: "ok",
			fields: fields{
				limit:   2,
				filters: map[string]interface{}{"country": "Argentina"},
			},
			expectations: func(fields fields) {
				service.repository.(*Mockrepository).
					EXPECT().
					List(gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*IP{dummyIP1, dummyIP2}, nil)
			},
			want: want{
				ips: []*IP{dummyIP1, dummyIP2},
				err: nil,
			},
		},
		{name: "ok",
			fields: fields{
				limit:   2,
				filters: map[string]interface{}{"country": "Argentina"},
			},
			expectations: func(fields fields) {
				service.repository.(*Mockrepository).
					EXPECT().
					List(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("error"))
			},
			want: want{
				ips: nil,
				err: errors.New("error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expectations(tt.fields)
			got, err := service.List(context.Background(), tt.fields.limit, tt.fields.filters)
			ass.EqualValues(tt.want.ips, got)
			ass.IsType(tt.want.err, err)
			ass.Equal(tt.want.err, err)
		})
	}
}

func TestAddressesService_GetIPQuantityByCountry(t *testing.T) {
	ass := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := newMockAddressesService(ctrl)

	type fields struct {
		country string
	}

	type want struct {
		quantity int
		err      error
	}

	tests := []struct {
		name         string
		fields       fields
		expectations func(fields fields)
		want         want
	}{
		{name: "ok",
			fields: fields{
				country: "Switzerland",
			},
			expectations: func(fields fields) {
				service.repository.(*Mockrepository).
					EXPECT().
					GetIPQuantityByCountry(gomock.Any(), gomock.Any()).
					Return(10, nil)
			},
			want: want{
				quantity: 10,
				err:      nil,
			},
		},
		{name: "no content",
			fields: fields{
				country: "Switzerland",
			},
			expectations: func(fields fields) {
				service.repository.(*Mockrepository).
					EXPECT().
					GetIPQuantityByCountry(gomock.Any(), gomock.Any()).
					Return(0, sql.ErrNoRows)
			},
			want: want{
				quantity: 0,
				err:      nil,
			},
		},
		{name: "error",
			fields: fields{
				country: "Switzerland",
			},
			expectations: func(fields fields) {
				service.repository.(*Mockrepository).
					EXPECT().
					GetIPQuantityByCountry(gomock.Any(), gomock.Any()).
					Return(0, sql.ErrConnDone)
			},
			want: want{
				quantity: 0,
				err:      sql.ErrConnDone,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expectations(tt.fields)
			got, err := service.GetIPQuantityByCountry(context.Background(), tt.fields.country)
			ass.EqualValues(tt.want.quantity, got)
			ass.IsType(tt.want.err, err)
			ass.Equal(tt.want.err, err)
		})
	}
}

func TestAddressesService_GetTopNISPByCountry(t *testing.T) {
	ass := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := newMockAddressesService(ctrl)

	type fields struct {
		top     int
		country string
	}

	type want struct {
		quantity []*ISPCount
		err      error
	}

	tests := []struct {
		name         string
		fields       fields
		expectations func(fields fields)
		want         want
	}{
		{name: "ok",
			fields: fields{
				top:     2,
				country: "Switzerland",
			},
			expectations: func(fields fields) {
				service.repository.(*Mockrepository).
					EXPECT().
					GetTopNISPByCountry(gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]*ISPCount{
						{
							ISP:   "Rook Media GmbH",
							Total: 180,
						},
						{
							ISP:   "RapidSeedbox Ltd",
							Total: 139,
						},
					}, nil)
			},
			want: want{
				quantity: []*ISPCount{
					{
						ISP:   "Rook Media GmbH",
						Total: 180,
					},
					{
						ISP:   "RapidSeedbox Ltd",
						Total: 139,
					},
				},
				err: nil,
			},
		},
		{name: "error",
			fields: fields{
				top:     2,
				country: "Switzerland",
			},
			expectations: func(fields fields) {
				service.repository.(*Mockrepository).
					EXPECT().
					GetTopNISPByCountry(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, sql.ErrConnDone)
			},
			want: want{
				quantity: nil,
				err:      sql.ErrConnDone,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expectations(tt.fields)
			got, err := service.GetTopNISPByCountry(context.Background(), tt.fields.top, tt.fields.country)
			ass.EqualValues(tt.want.quantity, got)
			ass.IsType(tt.want.err, err)
			ass.Equal(tt.want.err, err)
		})
	}
}

var (
	dummyIP1 = &IP{
		From:      150178522,
		To:        150178522,
		ProxyType: "PUB",
		Country: Country{
			Code:   "AR",
			Name:   "Argentina",
			Region: "Ciudad Autonoma de Buenos Aires",
			City:   "Buenos Aires",
		},
		ISP:    "CTL LATAM",
		Domain: "centurylink.com",
		Usage:  "ISP",
		ASN:    3356,
		AS:     "Level 3 Parent LLC",
	}

	dummyIP2 = &IP{
		From:      417862038,
		To:        417862038,
		ProxyType: "PUB",
		Country: Country{
			Code:   "AR",
			Name:   "Argentina",
			Region: "Cordoba",
			City:   "Cordoba",
		},
		ISP:    "Telecom Argentina S.A.",
		Domain: "telecom.com.ar",
		Usage:  "ISP/MOB",
		ASN:    7303,
		AS:     "Latin American and Caribbean IP address Regional Registry",
	}
)