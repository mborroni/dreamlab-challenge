package models

import (
	"github.com/mborroni/dreamlab-challenge/internal/conversion"
	ips "github.com/mborroni/dreamlab-challenge/internal/ipAddresses"
)

type IP struct {
	IP        string  `json:"ip,omitempty"`
	From      int64   `json:"from,omitempty"`
	To        int64   `json:"to,omitempty"`
	ProxyType string  `json:"proxy_type,omitempty"`
	Country   Country `json:"country,omitempty"`
	ISP       string  `json:"isp,omitempty"`
	Domain    string  `json:"domain,omitempty"`
	Usage     string  `json:"usage,omitempty"`
	ASN       int     `json:"asn,omitempty"`
	AS        string  `json:"as,omitempty"`
}

type Country struct {
	Code   string `json:"code,omitempty"`
	Name   string `json:"name,omitempty"`
	Region string `json:"region,omitempty"`
	City   string `json:"city,omitempty"`
}

func ToIPModel(ip string, entity *ips.IP) *IP {
	return &IP{
		IP:        ip,
		ProxyType: entity.ProxyType,
		Country: Country{
			Code:   entity.Country.Code,
			Name:   entity.Country.Name,
			Region: entity.Country.Region,
			City:   entity.Country.City,
		},
		ISP:    entity.ISP,
		Domain: entity.Domain,
		Usage:  entity.Usage,
		ASN:    entity.ASN,
		AS:     entity.AS,
	}
}

func ToIPsModel(entities []*ips.IP) []*IP {
	output := make([]*IP, 0)
	for _, ip := range entities {
		tmp := &IP{
			IP: conversion.DecimalToIPv4(ip.From),
			Country: Country{
				Name: ip.Country.Name,
				City: ip.Country.City,
			},
		}
		output = append(output, tmp)
	}
	return output
}
