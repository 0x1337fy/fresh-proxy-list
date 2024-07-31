package entity

import "encoding/xml"

type ProxyXMLClassicView struct {
	XMLName xml.Name `xml:"proxies"`
	Proxies []string `xml:"proxy"`
}

type ProxyXMLAdvancedView struct {
	XMLName xml.Name `xml:"proxies"`
	Proxies []Proxy  `xml:"proxy"`
}
