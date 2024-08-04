package repository

import (
	"fresh-proxy-list/internal/entity"
	"reflect"
	"testing"
)

func TestNewProxyRepository(t *testing.T) {
	tests := []struct {
		name string
		want ProxyRepositoryInterface
	}{
		{
			name: "New Proxy Repository",
			want: &ProxyRepository{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewProxyRepository()
			if got == nil {
				t.Errorf("NewProxyRepository() = nil, want non-nil")
			}
			if _, ok := got.(*ProxyRepository); !ok {
				t.Errorf("NewProxyRepository() = %T, want *ProxyRepository", got)
			}
		})
	}
}

func TestProxyRepository(t *testing.T) {
	type fields struct {
		allClassicView     []string
		httpClassicView    []string
		httpsClassicView   []string
		socks4ClassicView  []string
		socks5ClassicView  []string
		allAdvancedView    []entity.AdvancedProxy
		httpAdvancedView   []entity.Proxy
		httpsAdvancedView  []entity.Proxy
		socks4AdvancedView []entity.Proxy
		socks5AdvancedView []entity.Proxy
	}
	type args struct {
		proxy entity.Proxy
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    fields
		wantErr error
	}{
		{
			name: "Store HTTP proxy",
			fields: fields{
				allClassicView:     []string{},
				httpClassicView:    []string{},
				httpsClassicView:   []string{},
				socks4ClassicView:  []string{},
				socks5ClassicView:  []string{},
				allAdvancedView:    []entity.AdvancedProxy{},
				httpAdvancedView:   []entity.Proxy{},
				httpsAdvancedView:  []entity.Proxy{},
				socks4AdvancedView: []entity.Proxy{},
				socks5AdvancedView: []entity.Proxy{},
			},
			args: args{
				proxy: entity.Proxy{
					Category:  testHTTPCategory,
					IP:        testIP1,
					Port:      testPort,
					Proxy:     testProxy1,
					TimeTaken: testTimeTaken,
					CheckedAt: testCheckedAt,
				},
			},
			want: fields{
				allClassicView:    []string{testProxy1},
				httpClassicView:   []string{testProxy1},
				httpsClassicView:  []string{},
				socks4ClassicView: []string{},
				socks5ClassicView: []string{},
				allAdvancedView: []entity.AdvancedProxy{
					{
						Proxy:      testProxy1,
						IP:         testIP1,
						Port:       testPort,
						TimeTaken:  testTimeTaken,
						CheckedAt:  testCheckedAt,
						Categories: []string{testHTTPCategory},
					},
				},
				httpAdvancedView: []entity.Proxy{
					{
						Category:  testHTTPCategory,
						IP:        testIP1,
						Port:      testPort,
						Proxy:     testProxy1,
						TimeTaken: testTimeTaken,
						CheckedAt: testCheckedAt,
					},
				},
				httpsAdvancedView:  []entity.Proxy{},
				socks4AdvancedView: []entity.Proxy{},
				socks5AdvancedView: []entity.Proxy{},
			},
			wantErr: nil,
		},
		{
			name: "Store HTTPS proxy",
			fields: fields{
				allClassicView:     []string{},
				httpClassicView:    []string{},
				httpsClassicView:   []string{},
				socks4ClassicView:  []string{},
				socks5ClassicView:  []string{},
				allAdvancedView:    []entity.AdvancedProxy{},
				httpAdvancedView:   []entity.Proxy{},
				httpsAdvancedView:  []entity.Proxy{},
				socks4AdvancedView: []entity.Proxy{},
				socks5AdvancedView: []entity.Proxy{},
			},
			args: args{
				proxy: entity.Proxy{
					Category:  testHTTPSCategory,
					IP:        testIP1,
					Port:      testPort,
					Proxy:     testProxy1,
					TimeTaken: testTimeTaken,
					CheckedAt: testCheckedAt,
				},
			},
			want: fields{
				allClassicView:    []string{testProxy1},
				httpClassicView:   []string{},
				httpsClassicView:  []string{testProxy1},
				socks4ClassicView: []string{},
				socks5ClassicView: []string{},
				allAdvancedView: []entity.AdvancedProxy{
					{
						Proxy:      testProxy1,
						IP:         testIP1,
						Port:       testPort,
						TimeTaken:  testTimeTaken,
						CheckedAt:  testCheckedAt,
						Categories: []string{testHTTPSCategory},
					},
				},
				httpAdvancedView: []entity.Proxy{},
				httpsAdvancedView: []entity.Proxy{
					{
						Category:  testHTTPSCategory,
						IP:        testIP1,
						Port:      testPort,
						Proxy:     testProxy1,
						TimeTaken: testTimeTaken,
						CheckedAt: testCheckedAt,
					},
				},
				socks4AdvancedView: []entity.Proxy{},
				socks5AdvancedView: []entity.Proxy{},
			},
			wantErr: nil,
		},
		{
			name: "Store SOCKS4 proxy",
			fields: fields{
				allClassicView:     []string{},
				httpClassicView:    []string{},
				httpsClassicView:   []string{},
				socks4ClassicView:  []string{},
				socks5ClassicView:  []string{},
				allAdvancedView:    []entity.AdvancedProxy{},
				httpAdvancedView:   []entity.Proxy{},
				httpsAdvancedView:  []entity.Proxy{},
				socks4AdvancedView: []entity.Proxy{},
				socks5AdvancedView: []entity.Proxy{},
			},
			args: args{
				proxy: entity.Proxy{
					Category:  testSOCKS4Category,
					IP:        testIP1,
					Port:      testPort,
					Proxy:     testProxy1,
					TimeTaken: testTimeTaken,
					CheckedAt: testCheckedAt,
				},
			},
			want: fields{
				allClassicView:    []string{testProxy1},
				httpClassicView:   []string{},
				httpsClassicView:  []string{},
				socks4ClassicView: []string{testProxy1},
				socks5ClassicView: []string{},
				allAdvancedView: []entity.AdvancedProxy{
					{
						Proxy:      testProxy1,
						IP:         testIP1,
						Port:       testPort,
						TimeTaken:  testTimeTaken,
						CheckedAt:  testCheckedAt,
						Categories: []string{testSOCKS4Category},
					},
				},
				httpAdvancedView:  []entity.Proxy{},
				httpsAdvancedView: []entity.Proxy{},
				socks4AdvancedView: []entity.Proxy{
					{
						Category:  testSOCKS4Category,
						IP:        testIP1,
						Port:      testPort,
						Proxy:     testProxy1,
						TimeTaken: testTimeTaken,
						CheckedAt: testCheckedAt,
					},
				},
				socks5AdvancedView: []entity.Proxy{},
			},
			wantErr: nil,
		},
		{
			name: "Store SOCKS5 proxy",
			fields: fields{
				allClassicView:     []string{},
				httpClassicView:    []string{},
				httpsClassicView:   []string{},
				socks4ClassicView:  []string{},
				socks5ClassicView:  []string{},
				allAdvancedView:    []entity.AdvancedProxy{},
				httpAdvancedView:   []entity.Proxy{},
				httpsAdvancedView:  []entity.Proxy{},
				socks4AdvancedView: []entity.Proxy{},
				socks5AdvancedView: []entity.Proxy{},
			},
			args: args{
				proxy: entity.Proxy{
					Category:  testSOCKS5Category,
					IP:        testIP1,
					Port:      testPort,
					Proxy:     testProxy1,
					TimeTaken: testTimeTaken,
					CheckedAt: testCheckedAt,
				},
			},
			want: fields{
				allClassicView:    []string{testProxy1},
				httpClassicView:   []string{},
				httpsClassicView:  []string{},
				socks4ClassicView: []string{},
				socks5ClassicView: []string{testProxy1},
				allAdvancedView: []entity.AdvancedProxy{
					{
						Proxy:      testProxy1,
						IP:         testIP1,
						Port:       testPort,
						TimeTaken:  testTimeTaken,
						CheckedAt:  testCheckedAt,
						Categories: []string{testSOCKS5Category},
					},
				},
				httpAdvancedView:   []entity.Proxy{},
				httpsAdvancedView:  []entity.Proxy{},
				socks4AdvancedView: []entity.Proxy{},
				socks5AdvancedView: []entity.Proxy{
					{
						Category:  testSOCKS5Category,
						IP:        testIP1,
						Port:      testPort,
						Proxy:     testProxy1,
						TimeTaken: testTimeTaken,
						CheckedAt: testCheckedAt,
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "Duplicate proxy",
			fields: fields{
				allClassicView:    []string{testProxy1, testProxy2},
				httpClassicView:   []string{testProxy1, testProxy2},
				httpsClassicView:  []string{},
				socks4ClassicView: []string{},
				socks5ClassicView: []string{},
				allAdvancedView: []entity.AdvancedProxy{
					{
						Proxy:      testProxy1,
						IP:         testIP1,
						Port:       testPort,
						TimeTaken:  testTimeTaken,
						CheckedAt:  testCheckedAt,
						Categories: []string{testHTTPCategory},
					},
					{
						Proxy:      testProxy2,
						IP:         testIP2,
						Port:       testPort,
						TimeTaken:  testTimeTaken,
						CheckedAt:  testCheckedAt,
						Categories: []string{testHTTPCategory},
					},
				},
				httpAdvancedView: []entity.Proxy{
					{
						Category:  testHTTPCategory,
						IP:        testIP1,
						Port:      testPort,
						Proxy:     testProxy1,
						TimeTaken: testTimeTaken,
						CheckedAt: testCheckedAt,
					},
					{
						Category:  testHTTPCategory,
						IP:        testIP2,
						Port:      testPort,
						Proxy:     testProxy2,
						TimeTaken: testTimeTaken,
						CheckedAt: testCheckedAt,
					},
				},
				httpsAdvancedView:  []entity.Proxy{},
				socks4AdvancedView: []entity.Proxy{},
				socks5AdvancedView: []entity.Proxy{},
			},
			args: args{
				proxy: entity.Proxy{
					Category:  testHTTPCategory,
					IP:        testIP1,
					Port:      testPort,
					Proxy:     testProxy1,
					TimeTaken: testTimeTaken,
					CheckedAt: testCheckedAt,
				},
			},
			want: fields{
				allClassicView:    []string{testProxy1, testProxy2},
				httpClassicView:   []string{testProxy1, testProxy2, testProxy1},
				httpsClassicView:  []string{},
				socks4ClassicView: []string{},
				socks5ClassicView: []string{},
				allAdvancedView: []entity.AdvancedProxy{
					{
						Proxy:      testProxy1,
						IP:         testIP1,
						Port:       testPort,
						TimeTaken:  testTimeTaken,
						CheckedAt:  testCheckedAt,
						Categories: []string{testHTTPCategory},
					},
					{
						Proxy:      testProxy2,
						IP:         testIP2,
						Port:       testPort,
						TimeTaken:  testTimeTaken,
						CheckedAt:  testCheckedAt,
						Categories: []string{testHTTPCategory},
					},
				},
				httpAdvancedView: []entity.Proxy{
					{
						Category:  testHTTPCategory,
						IP:        testIP1,
						Port:      testPort,
						Proxy:     testProxy1,
						TimeTaken: testTimeTaken,
						CheckedAt: testCheckedAt,
					},
					{
						Category:  testHTTPCategory,
						IP:        testIP2,
						Port:      testPort,
						Proxy:     testProxy2,
						TimeTaken: testTimeTaken,
						CheckedAt: testCheckedAt,
					},
					{
						Category:  testHTTPCategory,
						IP:        testIP1,
						Port:      testPort,
						Proxy:     testProxy1,
						TimeTaken: testTimeTaken,
						CheckedAt: testCheckedAt,
					},
				},
				httpsAdvancedView:  []entity.Proxy{},
				socks4AdvancedView: []entity.Proxy{},
				socks5AdvancedView: []entity.Proxy{},
			},
			wantErr: nil,
		},
		{
			name: "Duplicate proxy with different proxy category",
			fields: fields{
				allClassicView:    []string{testProxy1},
				httpClassicView:   []string{testProxy1},
				httpsClassicView:  []string{},
				socks4ClassicView: []string{},
				socks5ClassicView: []string{},
				allAdvancedView: []entity.AdvancedProxy{
					{
						Proxy:      testProxy1,
						IP:         testIP1,
						Port:       testPort,
						TimeTaken:  testTimeTaken,
						CheckedAt:  testCheckedAt,
						Categories: []string{testHTTPCategory},
					},
				},
				httpAdvancedView: []entity.Proxy{
					{
						Category:  testHTTPCategory,
						IP:        testIP1,
						Port:      testPort,
						Proxy:     testProxy1,
						TimeTaken: testTimeTaken,
						CheckedAt: testCheckedAt,
					},
				},
				httpsAdvancedView:  []entity.Proxy{},
				socks4AdvancedView: []entity.Proxy{},
				socks5AdvancedView: []entity.Proxy{},
			},
			args: args{
				proxy: entity.Proxy{
					Category:  testSOCKS5Category,
					IP:        testIP1,
					Port:      testPort,
					Proxy:     testProxy1,
					TimeTaken: testTimeTaken,
					CheckedAt: testCheckedAt,
				},
			},
			want: fields{
				allClassicView:    []string{testProxy1},
				httpClassicView:   []string{testProxy1},
				httpsClassicView:  []string{},
				socks4ClassicView: []string{},
				socks5ClassicView: []string{testProxy1},
				allAdvancedView: []entity.AdvancedProxy{
					{
						Proxy:      testProxy1,
						IP:         testIP1,
						Port:       testPort,
						TimeTaken:  testTimeTaken,
						CheckedAt:  testCheckedAt,
						Categories: []string{testHTTPCategory, testSOCKS5Category},
					},
				},
				httpAdvancedView: []entity.Proxy{
					{
						Category:  testHTTPCategory,
						IP:        testIP1,
						Port:      testPort,
						Proxy:     testProxy1,
						TimeTaken: testTimeTaken,
						CheckedAt: testCheckedAt,
					},
				},
				httpsAdvancedView:  []entity.Proxy{},
				socks4AdvancedView: []entity.Proxy{},
				socks5AdvancedView: []entity.Proxy{
					{
						Category:  testSOCKS5Category,
						IP:        testIP1,
						Port:      testPort,
						Proxy:     testProxy1,
						TimeTaken: testTimeTaken,
						CheckedAt: testCheckedAt,
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ProxyRepository{
				allClassicView:     tt.fields.allClassicView,
				httpClassicView:    tt.fields.httpClassicView,
				httpsClassicView:   tt.fields.httpsClassicView,
				socks4ClassicView:  tt.fields.socks4ClassicView,
				socks5ClassicView:  tt.fields.socks5ClassicView,
				allAdvancedView:    tt.fields.allAdvancedView,
				httpAdvancedView:   tt.fields.httpAdvancedView,
				httpsAdvancedView:  tt.fields.httpsAdvancedView,
				socks4AdvancedView: tt.fields.socks4AdvancedView,
				socks5AdvancedView: tt.fields.socks5AdvancedView,
			}
			r.Store(&tt.args.proxy)

			if !reflect.DeepEqual(r.allClassicView, tt.want.allClassicView) {
				t.Errorf("GetAllClassicView() = %v, want %v", r.allClassicView, tt.want.allClassicView)
			}
			if !reflect.DeepEqual(r.httpClassicView, tt.want.httpClassicView) {
				t.Errorf("GetHTTPClassicView() = %v, want %v", r.httpClassicView, tt.want.httpClassicView)
			}
			if !reflect.DeepEqual(r.httpsClassicView, tt.want.httpsClassicView) {
				t.Errorf("GetHTTPSClassicView() = %v, want %v", r.httpsClassicView, tt.want.httpsClassicView)
			}
			if !reflect.DeepEqual(r.socks4ClassicView, tt.want.socks4ClassicView) {
				t.Errorf("GetSOCKS4ClassicView() = %v, want %v", r.socks4ClassicView, tt.want.socks4ClassicView)
			}
			if !reflect.DeepEqual(r.socks5ClassicView, tt.want.socks5ClassicView) {
				t.Errorf("GetSOCKS5ClassicView() = %v, want %v", r.socks5ClassicView, tt.want.socks5ClassicView)
			}
			if !reflect.DeepEqual(r.allAdvancedView, tt.want.allAdvancedView) {
				t.Errorf("GetAllAdvancedView() = %v, want %v", r.allAdvancedView, tt.want.allAdvancedView)
			}
			if !reflect.DeepEqual(r.httpAdvancedView, tt.want.httpAdvancedView) {
				t.Errorf("GetHTTPAdvancedView() = %v, want %v", r.httpAdvancedView, tt.want.httpAdvancedView)
			}
			if !reflect.DeepEqual(r.httpsAdvancedView, tt.want.httpsAdvancedView) {
				t.Errorf("GetHTTPSAdvancedView() = %v, want %v", r.httpsAdvancedView, tt.want.httpsAdvancedView)
			}
			if !reflect.DeepEqual(r.socks4AdvancedView, tt.want.socks4AdvancedView) {
				t.Errorf("GetSOCKS4AdvancedView() = %v, want %v", r.socks4AdvancedView, tt.want.socks4AdvancedView)
			}
			if !reflect.DeepEqual(r.socks5AdvancedView, tt.want.socks5AdvancedView) {
				t.Errorf("GetSOCKS5AdvancedView() = %v, want %v", r.socks5AdvancedView, tt.want.socks5AdvancedView)
			}
		})
	}
}

func TestGetAllClassicView(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *ProxyRepository
		want  []string
	}{
		{
			name: "Empty All proxies",
			setup: func() *ProxyRepository {
				return &ProxyRepository{
					allClassicView: []string{},
				}
			},
			want: []string{},
		},
		{
			name: "With All proxies",
			setup: func() *ProxyRepository {
				r := &ProxyRepository{}
				r.allClassicView = []string{testProxy1, testProxy2}
				return r
			},
			want: []string{testProxy1, testProxy2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setup()
			got := r.GetAllClassicView()
			if !reflect.DeepEqual(got, tt.want) {
				t.Logf("got: %v, want: %v", got, tt.want)
				t.Errorf("GetAllClassicView() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHTTPClassicView(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *ProxyRepository
		want  []string
	}{
		{
			name: "Empty HTTP proxies",
			setup: func() *ProxyRepository {
				return &ProxyRepository{
					httpClassicView: []string{},
				}
			},
			want: []string{},
		},
		{
			name: "With HTTP proxies",
			setup: func() *ProxyRepository {
				r := &ProxyRepository{}
				r.httpClassicView = []string{testProxy1, testProxy2}
				return r
			},
			want: []string{testProxy1, testProxy2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setup()
			got := r.GetHTTPClassicView()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHTTPClassicView() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHTTPSClassicView(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *ProxyRepository
		want  []string
	}{
		{
			name: "Empty HTTPS proxies",
			setup: func() *ProxyRepository {
				return &ProxyRepository{
					httpsClassicView: []string{},
				}
			},
			want: []string{},
		},
		{
			name: "With HTTPS proxies",
			setup: func() *ProxyRepository {
				r := &ProxyRepository{}
				r.httpsClassicView = []string{testProxy1, testProxy2}
				return r
			},
			want: []string{testProxy1, testProxy2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setup()
			got := r.GetHTTPSClassicView()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHTTPSClassicView() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSOCKS4ClassicView(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *ProxyRepository
		want  []string
	}{
		{
			name: "Empty SOCKS4 proxies",
			setup: func() *ProxyRepository {
				r := &ProxyRepository{}
				r.socks4ClassicView = []string{}
				return r
			},
			want: []string{},
		},
		{
			name: "With SOCKS4 proxies",
			setup: func() *ProxyRepository {
				return &ProxyRepository{
					socks4ClassicView: []string{testProxy1, testProxy2},
				}
			},
			want: []string{testProxy1, testProxy2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setup()
			got := r.GetSOCKS4ClassicView()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSOCKS4ClassicView() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSOCKS5ClassicView(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *ProxyRepository
		want  []string
	}{
		{
			name: "Empty SOCKS5 proxies",
			setup: func() *ProxyRepository {
				return &ProxyRepository{
					socks5ClassicView: []string{},
				}
			},
			want: []string{},
		},
		{
			name: "With SOCKS5 proxies",
			setup: func() *ProxyRepository {
				r := &ProxyRepository{}
				r.socks5ClassicView = []string{testProxy1, testProxy2}
				return r
			},
			want: []string{testProxy1, testProxy2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setup()
			got := r.GetSOCKS5ClassicView()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSOCKS5ClassicView() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllAdvancedView(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *ProxyRepository
		want  []entity.AdvancedProxy
	}{
		{
			name: "Empty All proxies",
			setup: func() *ProxyRepository {
				return &ProxyRepository{
					allAdvancedView: []entity.AdvancedProxy{},
				}
			},
			want: []entity.AdvancedProxy{},
		},
		{
			name: "With All proxies",
			setup: func() *ProxyRepository {
				r := &ProxyRepository{}
				r.allAdvancedView = []entity.AdvancedProxy{
					{
						Proxy: testProxy1,
						IP:    testIP1,
						Port:  testPort,
					},
					{
						Proxy: testProxy2,
						IP:    testIP2,
						Port:  testPort,
					},
				}
				return r
			},
			want: []entity.AdvancedProxy{
				{
					Proxy: testProxy1,
					IP:    testIP1,
					Port:  testPort,
				},
				{
					Proxy: testProxy2,
					IP:    testIP2,
					Port:  testPort,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setup()
			got := r.GetAllAdvancedView()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllAdvancedView() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHTTPAdvancedView(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *ProxyRepository
		want  []entity.Proxy
	}{
		{
			name: "Empty HTTP proxies",
			setup: func() *ProxyRepository {
				return &ProxyRepository{
					httpAdvancedView: []entity.Proxy{},
				}
			},
			want: []entity.Proxy{},
		},
		{
			name: "With HTTP proxies",
			setup: func() *ProxyRepository {
				r := &ProxyRepository{}
				r.httpAdvancedView = []entity.Proxy{
					{
						Proxy: testProxy1,
						IP:    testIP1,
						Port:  testPort,
					},
					{
						Proxy: testProxy2,
						IP:    testIP2,
						Port:  testPort,
					},
				}
				return r
			},
			want: []entity.Proxy{
				{
					Proxy: testProxy1,
					IP:    testIP1,
					Port:  testPort,
				},
				{
					Proxy: testProxy2,
					IP:    testIP2,
					Port:  testPort,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setup()
			got := r.GetHTTPAdvancedView()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHTTPAdvancedView() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHTTPSAdvancedView(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *ProxyRepository
		want  []entity.Proxy
	}{
		{
			name: "Empty HTTPS proxies",
			setup: func() *ProxyRepository {
				return &ProxyRepository{
					httpsAdvancedView: []entity.Proxy{},
				}
			},
			want: []entity.Proxy{},
		},
		{
			name: "With HTTPS proxies",
			setup: func() *ProxyRepository {
				r := &ProxyRepository{}
				r.httpsAdvancedView = []entity.Proxy{
					{
						Proxy: "192.168.1.3:8080",
						IP:    "192.168.1.3",
						Port:  testPort,
					},
					{
						Proxy: "192.168.1.4:8080",
						IP:    "192.168.1.4",
						Port:  testPort,
					},
				}
				return r
			},
			want: []entity.Proxy{
				{
					Proxy: "192.168.1.3:8080",
					IP:    "192.168.1.3",
					Port:  testPort,
				},
				{
					Proxy: "192.168.1.4:8080",
					IP:    "192.168.1.4",
					Port:  testPort,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setup()
			got := r.GetHTTPSAdvancedView()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHTTPSAdvancedView() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSOCKS4AdvancedView(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *ProxyRepository
		want  []entity.Proxy
	}{
		{
			name: "Empty SOCKS4 proxies",
			setup: func() *ProxyRepository {
				return &ProxyRepository{
					socks4AdvancedView: []entity.Proxy{},
				}
			},
			want: []entity.Proxy{},
		},
		{
			name: "With SOCKS proxies",
			setup: func() *ProxyRepository {
				r := &ProxyRepository{}
				r.socks4AdvancedView = []entity.Proxy{
					{
						Proxy: "192.168.1.5:8080",
						IP:    "192.168.1.5",
						Port:  testPort,
					},
					{
						Proxy: "192.168.1.6:8080",
						IP:    "192.168.1.6",
						Port:  testPort,
					},
				}
				return r
			},
			want: []entity.Proxy{
				{
					Proxy: "192.168.1.5:8080",
					IP:    "192.168.1.5",
					Port:  testPort,
				},
				{
					Proxy: "192.168.1.6:8080",
					IP:    "192.168.1.6",
					Port:  testPort,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setup()
			got := r.GetSOCKS4AdvancedView()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSOCKS4AdvancedView() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSOCKS5AdvancedView(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *ProxyRepository
		want  []entity.Proxy
	}{
		{
			name: "Empty SOCKS5 proxies",
			setup: func() *ProxyRepository {
				return &ProxyRepository{
					socks5AdvancedView: []entity.Proxy{},
				}
			},
			want: []entity.Proxy{},
		},
		{
			name: "With SOCKS5 proxies",
			setup: func() *ProxyRepository {
				r := &ProxyRepository{}
				r.socks5AdvancedView = []entity.Proxy{
					{
						Proxy: "192.168.1.7:8080",
						IP:    "192.168.1.7",
						Port:  testPort,
					},
					{
						Proxy: "192.168.1.8:8080",
						IP:    "192.168.1.8",
						Port:  testPort,
					},
				}
				return r
			},
			want: []entity.Proxy{
				{
					Proxy: "192.168.1.7:8080",
					IP:    "192.168.1.7",
					Port:  testPort,
				},
				{
					Proxy: "192.168.1.8:8080",
					IP:    "192.168.1.8",
					Port:  testPort,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.setup()
			got := r.GetSOCKS5AdvancedView()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSOCKS5AdvancedView() = %v, want %v", got, tt.want)
			}
		})
	}
}
