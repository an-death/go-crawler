package main

import (
	"io"
	"net/http"
	"strings"
)

type responseHandle struct {
	validators []func(*http.Response) error
	bodyHandlers []func(reader io.Reader) error
}

func (rh *responseHandle) HandleResponse(response *http.Response) error {
	if err := rh.runValidators(response); err != nil {
		return err
	}
	defer response.Body.Close()
	return rh.runHandlers(response.Body)
}

func (rh *responseHandle) runValidators(response *http.Response) error {
	for _, validator := range rh.validators{
		if err := validator(response); err != nil {
			return err
		}
	}
	return nil
}

func (rh *responseHandle) runHandlers(reader io.Reader) error{
	for _, handler := range rh.bodyHandlers {
		if err := handler(reader); err != nil {
			return err
		}
	}
	return nil
}


func checkResponseCode(expectedCode int) func(response *http.Response) error {
	return func(response *http.Response) error {
		if response.StatusCode != expectedCode {
			return &DoesNotMatchError{response.StatusCode, expectedCode}
		}
		return nil
	}
}

func checkResponseContentType(expectedContentType string) func(response *http.Response) error {
	return func(response *http.Response) error {
		actualContentType := response.Header.Get("Content-Type")
		if !strings.Contains(strings.ToLower(actualContentType), expectedContentType) {
			return &DoesNotMatchError{actualContentType, expectedContentType}
		}
		return nil
	}
}


