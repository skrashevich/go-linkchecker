package main

import (
	"reflect"
	"sync"
	"testing"
)

func TestCheckLink(t *testing.T) {
	// Test case 1: Testing a valid link that returns a 200 status code
	ch1 := make(chan checkResult)
	wg1 := sync.WaitGroup{}
	wg1.Add(1)
	go checkLink("https://example.com", ch1, &wg1)
	result1 := <-ch1
	wg1.Wait()
	if result1.OK != true {
		t.Errorf("Expected OK to be true, but got false")
	}

	// Test case 2: Testing a valid link that returns a non-200 status code
	ch2 := make(chan checkResult)
	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	go checkLink("https://example.com/404", ch2, &wg2)
	result2 := <-ch2
	wg2.Wait()
	if result2.OK != false {
		t.Errorf("Expected OK to be false, but got true")
	}

	// Test case 3: Testing an invalid link
	ch3 := make(chan checkResult)
	wg3 := sync.WaitGroup{}
	wg3.Add(1)
	go checkLink("https://invalid-link.com", ch3, &wg3)
	result3 := <-ch3
	wg3.Wait()
	if result3.OK != false {
		t.Errorf("Expected OK to be false, but got true")
	}
}

func TestCreateWhiteListSet(t *testing.T) {
	tests := []struct {
		whiteList string
		expected  map[string]bool
	}{
		{
			whiteList: "example.com, google.com",
			expected: map[string]bool{
				"example.com": true,
				"google.com":  true,
			},
		},
		{
			whiteList: "example.com,google.com,go.dev",
			expected: map[string]bool{
				"example.com": true,
				"google.com":  true,
				"go.dev":      true,
			},
		},
		{
			whiteList: "github.com",
			expected: map[string]bool{
				"github.com": true,
			},
		},
		{
			whiteList: "",
			expected:  map[string]bool{},
		},
	}

	for _, test := range tests {
		result := createWhiteListSet(test.whiteList)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Expected %v, but got %v", test.expected, result)
		}
	}
}
