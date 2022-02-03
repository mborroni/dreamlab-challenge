package ips

type IP struct {
	From      int64
	To        int64
	ProxyType string
	Country   Country
	ISP       string
	Domain    string
	Usage     string
	ASN       int
	AS        string
}

type Country struct {
	Code   string
	Name   string
	Region string
	City   string
}
