// Pico HTTP Client
// this is a simple http client that retries on error
// Thanks to:
// + https://medium.com/@nitishkr88/http-retries-in-go-e622e51d249f
// + https://github.com/hashicorp/go-retryablehttp
package httpdo

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var (
	defaultRetryMaxAttempts = 3
	defaultRetryWaitMin     = 1 * time.Second
	defaultRetryWaitMax     = 20 * time.Second
)

type Client struct {
	HTTPClient *http.Client

	RetryWaitMin     time.Duration
	RetryWaitMax     time.Duration
	RetryMaxAttempts int

	CheckForRetry CheckForRetry
	Backoff       Backoff
}

func NewClient() *Client {
	return &Client{
		HTTPClient:       &http.Client{},
		RetryWaitMin:     defaultRetryWaitMin,
		RetryWaitMax:     defaultRetryWaitMax,
		RetryMaxAttempts: defaultRetryMaxAttempts,
		CheckForRetry:    DefaultRetryPolicy,
		Backoff:          DefaultBackoff,
	}
}

type Request struct {
	*http.Request
	body io.ReadSeeker
}

func NewRequest(method, url string, body io.ReadSeeker) (*Request, error) {
	var rcBody io.ReadCloser
	if body != nil {
		rcBody = ioutil.NopCloser(body)
	}

	req, err := http.NewRequest(method, url, rcBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	return &Request{
		Request: req,
		body:    body,
	}, nil
}

func (c *Client) Do(req *Request) (*http.Response, error) {
	for i := 0; i < c.RetryMaxAttempts; i++ {

		if req.Body != nil {
			if _, err := req.body.Seek(0, 0); err != nil {
				return nil, fmt.Errorf("error seeking body: %v", err)
			}
		}

		var code int
		resp, err := c.HTTPClient.Do(req.Request)
		if err != nil {
			code = resp.StatusCode
		}

		retryOK, retryErr := c.CheckForRetry(resp, err)

		if !retryOK {
			if retryErr != nil {
				err = retryErr
			}
			return resp, err
		}

		if err == nil {
			c.drainBody(resp.Body)
		}

		remain := c.RetryMaxAttempts - i
		if remain <= 0 {
			break
		}

		wait := c.Backoff(c.RetryWaitMin, c.RetryWaitMax, i, resp)
		desc := fmt.Sprintf("%s %s", req.Method, req.URL)

		if code > 0 {
			desc = fmt.Sprintf("%s (status: %d)", desc, code)
		}

		fmt.Printf("retrying %s in %s (%d left)\n", desc, wait, remain)

		time.Sleep(wait)
	}

	return nil, fmt.Errorf("%s %s given up after %d attempts", req.Method, req.URL, c.RetryMaxAttempts+1)
}

type CheckForRetry func(resp *http.Response, err error) (bool, error)

func DefaultRetryPolicy(resp *http.Response, err error) (bool, error) {
	if err != nil {
		return true, err
	}

	if resp.StatusCode == 0 && resp.StatusCode >= 500 {
		return true, nil
	}

	return false, nil
}

type Backoff func(min, max time.Duration, attemps int, resp *http.Response) time.Duration

func DefaultBackoff(min, max time.Duration, attemps int, resp *http.Response) time.Duration {
	if resp != nil {
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
			if s, ok := resp.Header["Retry-After"]; ok {
				if sleep, err := strconv.Atoi(s[0]); err == nil {
					return time.Duration(sleep) * time.Second
				}
			}
		}
	}

	t := math.Pow(2, float64(attemps)) * float64(min)
	sleep := time.Duration(t)

	if float64(sleep) != t || sleep > max {
		sleep = max
	}

	return sleep
}

// ExponentialJitterBackoff returns a Backoff function which uses an exponential jitter algorithm
// References:
// + https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
// + https://www.oreilly.com/library/view/cloud-native-go/9781492076322/
func ExponentialJitterBackoff(min, max time.Duration, attemps int, resp *http.Response) time.Duration {
	if resp != nil {
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
			if s, ok := resp.Header["Retry-After"]; ok {
				if sleep, err := strconv.Atoi(s[0]); err == nil {
					return time.Duration(sleep) * time.Second
				}
			}
		}
	}

	backoff := math.Pow(2, float64(attemps)) * float64(min)
	jitter := rand.Int63n(int64(backoff * 3))

	sleep := min + time.Duration(jitter)

	if sleep > max {
		sleep = max
	}

	return sleep
}

func (c *Client) drainBody(body io.ReadCloser) {
	defer body.Close()
	_, err := io.Copy(ioutil.Discard, io.LimitReader(body, 4096))
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func (c *Client) Get(url string) (*http.Response, error) {
	req, err := NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *Client) Post(url string, body io.ReadSeeker) (*http.Response, error) {
	req, err := NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.Do(req)
}
