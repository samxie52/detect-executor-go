package client

import (
	"context"
	"detect-executor-go/pkg/config"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// HttpClient HTTP客户端接口
type HttpClient interface {
	BaseClient
	Get(ctx context.Context, url string, headers map[string]string) (*HttpResponse, error)
	Post(ctx context.Context, url string, body interface{}, headers map[string]string) (*HttpResponse, error)
	Put(ctx context.Context, url string, body interface{}, headers map[string]string) (*HttpResponse, error)
	Delete(ctx context.Context, url string, headers map[string]string) (*HttpResponse, error)
}

// HttpResponse HTTP响应
type HttpResponse struct {
	StatusCode int
	Body       []byte
	Headers    map[string][]string
}

// httpClient HTTP客户端实现
type httpClient struct {
	client *resty.Client
	config config.HTTPClientConfig
}

func NewHTTPClient(config config.HTTPClientConfig) HttpClient {
	//resty 是一个基于 Go 语言的 HTTP 客户端库，提供了丰富的功能和灵活的配置选项
	client := *resty.New()

	//配置超时时间
	client.SetTimeout(config.Timeout)

	//配置重试次数
	client.SetRetryCount(config.RetryCount)

	//配置重试间隔
	client.SetRetryMaxWaitTime(config.RetryInterval)

	return &httpClient{
		client: &client,
		config: config,
	}
}

func (c *httpClient) Get(ctx context.Context, url string, headers map[string]string) (*HttpResponse, error) {
	req := c.client.R().SetContext(ctx)

	if headers != nil {
		req.SetHeaders(headers)
	}

	resp, err := req.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get %s failed: %v", url, err)
	}

	return &HttpResponse{
		StatusCode: resp.StatusCode(),
		Body:       resp.Body(),
		Headers:    c.convertHeaders(resp.Header()),
	}, nil
}

func (c *httpClient) Post(ctx context.Context, url string, body interface{}, headers map[string]string) (*HttpResponse, error) {
	req := c.client.R().SetContext(ctx)

	if headers != nil {
		req.SetHeaders(headers)
	}

	if body != nil {
		req.SetBody(body)
	}

	resp, err := req.Post(url)
	if err != nil {
		return nil, fmt.Errorf("post %s failed: %v", url, err)
	}

	return &HttpResponse{
		StatusCode: resp.StatusCode(),
		Body:       resp.Body(),
		Headers:    c.convertHeaders(resp.Header()),
	}, nil
}

func (c *httpClient) Put(ctx context.Context, url string, body interface{}, headers map[string]string) (*HttpResponse, error) {
	req := c.client.R().SetContext(ctx)

	if headers != nil {
		req.SetHeaders(headers)
	}

	if body != nil {
		req.SetBody(body)
	}

	resp, err := req.Put(url)
	if err != nil {
		return nil, fmt.Errorf("put %s failed: %v", url, err)
	}

	return &HttpResponse{
		StatusCode: resp.StatusCode(),
		Body:       resp.Body(),
		Headers:    c.convertHeaders(resp.Header()),
	}, nil
}

func (c *httpClient) Delete(ctx context.Context, url string, headers map[string]string) (*HttpResponse, error) {
	req := c.client.R().SetContext(ctx)

	if headers != nil {
		req.SetHeaders(headers)
	}

	resp, err := req.Delete(url)
	if err != nil {
		return nil, fmt.Errorf("delete %s failed: %v", url, err)
	}

	return &HttpResponse{
		StatusCode: resp.StatusCode(),
		Body:       resp.Body(),
		Headers:    c.convertHeaders(resp.Header()),
	}, nil
}

func (c *httpClient) Close() error {
	return nil
}

func (c *httpClient) Ping(ctx context.Context) error {
	return nil
}

// convertHeaders 将 http.Header 转换为 map[string][]string
func (c *httpClient) convertHeaders(headers map[string][]string) map[string][]string {
	result := make(map[string][]string)
	for k, v := range headers {
		result[k] = v
	}
	return result
}
