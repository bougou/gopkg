package snippet

import "fmt"

// inner client
type Client struct {
	url string
}

func NewClient(url string) *Client {
	return &Client{url}
}

func (c *Client) XXX() {
	fmt.Println("xxx")
}

func (c *Client) YYY() {
	fmt.Println("yyy")
}

// Define SomeClient
type SomeClient struct {
	URL string

	client *Client
}

// NewXXX function
func NewSome(url string) *SomeClient {
	return &SomeClient{
		URL: url,
	}
}

// Use public fields to initialize private fields
func (s *SomeClient) init() error {
	// use the public URL field to init the private client field
	s.client = NewClient(s.URL)

	return nil
}

// exported
func (s *SomeClient) DoSomething() error {
	// !!!!!!!!!!!!!!!!!! NOTE
	// init client if not initialized
	if s.client == nil {
		if err := s.init(); err != nil {
			return nil
		}
	}

	// call other methods
	s.dosomething1()

	s.dosomething2()

	return nil
}

// unexported
func (s *SomeClient) dosomething1() {
	// use client to do something
	s.client.XXX()
}

// unexported
func (s *SomeClient) dosomething2() {
	s.client.YYY()
}
