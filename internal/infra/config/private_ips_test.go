package config

import (
	"net"
	"testing"
)

func TestParseCIDR(t *testing.T) {
	tests := []struct {
		args struct {
			cidr string
		}
		want *net.IPNet
	}{
		{
			args: struct {
				cidr string
			}{
				cidr: "10.0.0.0/8",
			},
			want: &net.IPNet{
				IP:   net.IP{10, 0, 0, 0},
				Mask: net.CIDRMask(8, 32),
			},
		},
		{
			args: struct {
				cidr string
			}{
				cidr: "172.16.0.0/12",
			},
			want: &net.IPNet{
				IP:   net.IP{172, 16, 0, 0},
				Mask: net.CIDRMask(12, 32),
			},
		},
		{
			args: struct {
				cidr string
			}{
				cidr: "192.168.0.0/16",
			},
			want: &net.IPNet{
				IP:   net.IP{192, 168, 0, 0},
				Mask: net.CIDRMask(16, 32),
			},
		},
		{
			args: struct {
				cidr string
			}{
				cidr: "invalid-cidr",
			},
			want: nil, // We expect this case to panic
		},
	}

	for _, test := range tests {
		if test.want != nil {
			result := ParseCIDR(test.args.cidr)
			if result.String() != test.want.String() {
				t.Errorf("ParseCIDR(%q) = %v; want %v", test.args.cidr, result, test.want)
			}
		} else {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("ParseCIDR(%q) did not panic", test.args.cidr)
				}
			}()
			ParseCIDR(test.args.cidr)
		}
	}
}
