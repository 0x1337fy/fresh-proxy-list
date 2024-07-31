package repository

import (
	"encoding/json"
	"fresh-proxy-list/internal/entity"
)

const (
	testHTTPCategory       = "HTTP"
	testHTTPSCategory      = "HTTPS"
	testSOCKS4Category     = "SOCKS4"
	testSOCKS5Category     = "SOCKS5"
	testIP1                = "192.168.0.1"
	testIP2                = "192.168.0.2"
	testIP3                = "192.168.0.2"
	testPort               = "8080"
	testProxy1             = testIP1 + ":" + testPort
	testProxy2             = testIP2 + ":" + testPort
	testProxy3             = testIP3 + ":" + testPort
	testTimeTaken          = 10.4243
	testCheckedAt          = "2024-07-27T00:00:00Z"
	testDirectory          = "/tmp"
	testFilePath           = "/tmp/test_file"
	testErrCreateFile      = "error creating file %s: %s"
	testErrCreateDirectory = "error creating directory %s: %s"
	testErrWriting         = "error writing"
	testErrWritingTXT      = "error writing TXT: %v"
	testErrEncode          = "error encoding %s: %s"
	errUnexpected          = "unexpected error: %v"
	errExpected            = "expected %q, got %q"
	errExpectedContain     = "expected error containing %q, got %v"
	testUnsupportedFormat  = "unsupported format: %s"
)

var (
	testIPs = []string{
		testIP1,
		testIP2,
		testIP3,
	}
	testProxies = []entity.Proxy{
		{
			Proxy:     testProxy1,
			IP:        testIP1,
			Port:      testPort,
			TimeTaken: testTimeTaken,
			CheckedAt: testCheckedAt,
		},
	}
	testAdvancedProxies = []entity.AdvancedProxy{
		{
			Proxy:      testProxy2,
			IP:         testIP2,
			Port:       testPort,
			Categories: []string{testHTTPCategory, testSOCKS5Category},
			TimeTaken:  testTimeTaken,
			CheckedAt:  testCheckedAt,
		},
	}
	testIPsToText, _             = json.Marshal(testIPs)
	testProxiesToString, _       = json.Marshal(testProxies)
	testAdvancedProxiesToText, _ = json.Marshal(testAdvancedProxies)
)
