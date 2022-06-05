package partner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/andrysds/colly-example/product"
	"github.com/stretchr/testify/mock"
)

func TestNewPartner(t *testing.T) {
	mockUsername := "sample username"
	mockPassword := "sample password"
	mockLoginUrl := "https://example.com/login"
	mockGetProductBaseUrl := "https://example.com/product/"

	os.Setenv(usernameEnvKey, mockUsername)
	os.Setenv(passwordEnvKey, mockPassword)
	os.Setenv(loginUrlEnvKey, mockLoginUrl)
	os.Setenv(getProductBaseUrlEnvKey, mockGetProductBaseUrl)
	defer os.Unsetenv(usernameEnvKey)
	defer os.Unsetenv(passwordEnvKey)
	defer os.Unsetenv(loginUrlEnvKey)
	defer os.Unsetenv(getProductBaseUrlEnvKey)

	want := &Partner{
		httpClient:        &http.Client{},
		authToken:         "",
		username:          mockUsername,
		password:          mockPassword,
		loginUrl:          mockLoginUrl,
		getProductBaseUrl: mockGetProductBaseUrl,
	}
	if got := NewPartner(); !reflect.DeepEqual(got, want) {
		t.Errorf("NewPartner() = %v, want %v", got, want)
	}
}

func TestPartner_Login(t *testing.T) {
	mockUsername := "sample username"
	mockPassword := "sample password"
	mockAuthToken := "sample auth token"
	mockLoginUrl := "https://example.com/login"

	mockReqMatcher := func(req *http.Request) bool {
		mockReqJsonBody, _ := json.Marshal(
			map[string]string{
				"username": mockUsername,
				"password": mockPassword,
			},
		)
		mockReqBody := bytes.NewBuffer(mockReqJsonBody)
		mockReq, _ := http.NewRequest(http.MethodPost, mockLoginUrl, mockReqBody)

		return req.Method == mockReq.Method &&
			reflect.DeepEqual(req.URL, mockReq.URL) &&
			reflect.DeepEqual(req.Body, mockReq.Body)
	}

	tests := []struct {
		name          string
		httpClient    func() httpClient
		wantErr       bool
		wantAuthToken string
	}{
		{
			name: "got error from httpClient",
			httpClient: func() httpClient {
				c := &mockHttpClient{}
				c.On("Do", mock.MatchedBy(mockReqMatcher)).Return(nil, fmt.Errorf("sample error"))
				return c
			},
			wantErr:       true,
			wantAuthToken: "",
		},
		{
			name: "got non 2xx response from httpClient",
			httpClient: func() httpClient {
				c := &mockHttpClient{}
				mockRes := &http.Response{StatusCode: http.StatusBadRequest}
				c.On("Do", mock.MatchedBy(mockReqMatcher)).Return(mockRes, nil)
				return c
			},
			wantErr:       true,
			wantAuthToken: "",
		},
		{
			name: "happy path",
			httpClient: func() httpClient {
				c := &mockHttpClient{}

				mockResJsonBody, _ := json.Marshal(
					map[string]map[string]string{
						"data": {
							"token": mockAuthToken,
						},
					},
				)
				mockResBody := bytes.NewBuffer(mockResJsonBody)
				mockRes := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(mockResBody),
				}

				c.On("Do", mock.MatchedBy(mockReqMatcher)).Return(mockRes, nil)
				return c
			},
			wantErr:       false,
			wantAuthToken: mockAuthToken,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Partner{
				httpClient: tt.httpClient(),
				username:   mockUsername,
				password:   mockPassword,
				loginUrl:   mockLoginUrl,
			}
			err := p.Login()
			if (err != nil) != tt.wantErr {
				t.Errorf("Partner.Login() error = %v, wantErr %v", err, tt.wantErr)
			}
			if p.authToken != tt.wantAuthToken {
				t.Errorf("Partner.Login() Partner.auth = %v, wantErr %v", p.authToken, tt.wantAuthToken)
			}
		})
	}
}

func TestPartner_GetProduct(t *testing.T) {
	mockAuthToken := "sample auth token"
	mockGetProductBaseUrl := "https://example.com/product/"
	mockSlug := "sample-slug"

	mockReqMatcher := func(req *http.Request) bool {
		mockReq, _ := http.NewRequest(http.MethodGet, mockGetProductBaseUrl+mockSlug, nil)
		mockReq.Header.Set("Authorization", mockAuthToken)
		return req.Method == mockReq.Method &&
			reflect.DeepEqual(req.URL, mockReq.URL) &&
			reflect.DeepEqual(req.Header, mockReq.Header)
	}

	mockProduct := &product.Product{
		Name:        "sample name",
		Description: "sample description",
		Variants: []product.Variant{
			{
				Name:  "sample name",
				Price: 1000,
				Stock: 0,
			},
		},
	}

	tests := []struct {
		name       string
		httpClient func() httpClient
		want       *product.Product
		wantErr    bool
	}{
		{
			name: "got error from httpClient",
			httpClient: func() httpClient {
				c := &mockHttpClient{}
				c.On("Do", mock.MatchedBy(mockReqMatcher)).Return(nil, fmt.Errorf("sample error"))
				return c
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "got non 2xx response from httpClient",
			httpClient: func() httpClient {
				c := &mockHttpClient{}
				mockRes := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBuffer([]byte(""))),
				}
				c.On("Do", mock.MatchedBy(mockReqMatcher)).Return(mockRes, nil)
				return c
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy path",
			httpClient: func() httpClient {
				c := &mockHttpClient{}

				mockResJsonBody, _ := json.Marshal(
					map[string]interface{}{
						"name":        mockProduct.Name,
						"description": mockProduct.Description,
						"variants": []map[string]interface{}{
							{
								"variants_name": mockProduct.Variants[0].Name,
								"price":         mockProduct.Variants[0].Price,
								"stock":         mockProduct.Variants[0].Stock,
							},
						},
					},
				)
				mockResBody := bytes.NewBuffer(mockResJsonBody)
				mockRes := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(mockResBody),
				}

				c.On("Do", mock.MatchedBy(mockReqMatcher)).Return(mockRes, nil)
				return c
			},
			want:    mockProduct,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Partner{
				httpClient:        tt.httpClient(),
				authToken:         mockAuthToken,
				getProductBaseUrl: mockGetProductBaseUrl,
			}
			got, err := p.GetProduct(mockSlug)
			if (err != nil) != tt.wantErr {
				t.Errorf("Partner.GetProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Partner.GetProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}
