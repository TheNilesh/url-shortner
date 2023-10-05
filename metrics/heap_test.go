package metrics

import (
	"testing"
)

func TestHeap(t *testing.T) {
	h := NewHeap()

	h.IncOrPush("google.com")
	h.IncOrPush("github.com")
	h.IncOrPush("google.com")
	h.IncOrPush("facebook.com")
	h.IncOrPush("github.com")
	h.IncOrPush("github.com")

	top3 := h.GetMaxValuePairs(3)
	if len(top3) != 3 {
		t.Errorf("Expected 3 top value pairs, but got %d", len(top3))
	}
	if top3[0].Key != "github.com" || top3[0].Value != 3 {
		t.Errorf("Expected top value pair to be {Key: 'github.com', Value: 3}, but got %v", top3[0])
	}
	if top3[1].Key != "google.com" || top3[1].Value != 2 {
		t.Errorf("Expected second top value pair to be {Key: 'google.com', Value: 2}, but got %v", top3[1])
	}
	if top3[2].Key != "facebook.com" || top3[2].Value != 1 {
		t.Errorf("Expected third top value pair to be {Key: 'facebook.com', Value: 1}, but got %v", top3[2])
	}
}
