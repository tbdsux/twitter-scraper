package twitterscraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"
)

const bearerToken string = "AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA"

// RequestAPI get JSON from frontend API and decodes it
func (s *Scraper) RequestAPI(req *http.Request, target interface{}) error {
	s.wg.Wait()
	if s.delay > 0 {
		defer s.delayRequest()
	}

	if err := s.prepareRequest(req); err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return s.handleResponse(resp, target)
}

func (s *Scraper) delayRequest() {
	s.wg.Add(1)
	go func() {
		time.Sleep(time.Second * time.Duration(s.delay))
		s.wg.Done()
	}()
}

func (s *Scraper) prepareRequest(req *http.Request) error {
	req.Header.Set("User-Agent", s.userAgent)

	if !s.isLogged {
		if err := s.setGuestToken(req); err != nil {
			return err
		}
	}

	s.setAuthorizationHeader(req)
	s.setCSRFToken(req)

	return nil
}

func (s *Scraper) setGuestToken(req *http.Request) error {
	if !s.IsGuestToken() || s.guestCreatedAt.Before(time.Now().Add(-time.Hour*3)) {
		if err := s.GetGuestToken(); err != nil {
			return err
		}
	}
	req.Header.Set("X-Guest-Token", s.guestToken)
	return nil
}

func (s *Scraper) setAuthorizationHeader(req *http.Request) {
	if s.oAuthToken != "" && s.oAuthSecret != "" {
		req.Header.Set("Authorization", s.sign(req.Method, req.URL))
	} else {
		req.Header.Set("Authorization", "Bearer "+s.bearerToken)
	}
}

func (s *Scraper) setCSRFToken(req *http.Request) {
	for _, cookie := range s.client.Jar.Cookies(req.URL) {
		if cookie.Name == "ct0" {
			req.Header.Set("X-CSRF-Token", cookie.Value)
			break
		}
	}
}

func (s *Scraper) handleResponse(resp *http.Response, target interface{}) error {
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	okResponses := []int{http.StatusOK, http.StatusAccepted, http.StatusNoContent, http.StatusCreated}

	if !slices.Contains(okResponses, resp.StatusCode) {
		return fmt.Errorf("response status %s: %s", resp.Status, content)
	}

	if resp.Header.Get("X-Rate-Limit-Remaining") == "0" {
		s.guestToken = ""
	}

	if target == nil {
		return nil
	}

	return json.Unmarshal(content, target)
}

// GetGuestToken from Twitter API
func (s *Scraper) GetGuestToken() error {
	req, err := http.NewRequest("POST", "https://api.x.com/1.1/guest/activate.json", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.bearerToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response status %s: %s", resp.Status, body)
	}

	var jsn map[string]interface{}
	if err := json.Unmarshal(body, &jsn); err != nil {
		return err
	}
	var ok bool
	if s.guestToken, ok = jsn["guest_token"].(string); !ok {
		return fmt.Errorf("guest_token not found")
	}
	s.guestCreatedAt = time.Now()

	return nil
}

func (s *Scraper) ClearGuestToken() error {
	s.guestToken = ""
	s.guestCreatedAt = time.Time{}

	return nil
}
