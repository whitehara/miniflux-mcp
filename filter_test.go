package main

import (
	"testing"
	"time"

	"miniflux.app/v2/client"
)

func TestBuildEntriesFilter(t *testing.T) {
	t.Run("nil_args_returns_nil", func(t *testing.T) {
		f, err := buildEntriesFilter(nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if f != nil {
			t.Errorf("expected nil filter, got %+v", f)
		}
	})

	t.Run("existing_params_still_work", func(t *testing.T) {
		args := map[string]interface{}{
			"status":  "unread",
			"feed_id": float64(42),
			"limit":   float64(10),
			"offset":  float64(20),
		}
		f, err := buildEntriesFilter(args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if f.Status != "unread" {
			t.Errorf("Status: got %q, want %q", f.Status, "unread")
		}
		if f.FeedID != 42 {
			t.Errorf("FeedID: got %d, want 42", f.FeedID)
		}
		if f.Limit != 10 {
			t.Errorf("Limit: got %d, want 10", f.Limit)
		}
		if f.Offset != 20 {
			t.Errorf("Offset: got %d, want 20", f.Offset)
		}
	})

	t.Run("direction_and_order", func(t *testing.T) {
		args := map[string]interface{}{
			"direction": "desc",
			"order":     "published_at",
		}
		f, err := buildEntriesFilter(args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if f.Direction != "desc" {
			t.Errorf("Direction: got %q, want %q", f.Direction, "desc")
		}
		if f.Order != "published_at" {
			t.Errorf("Order: got %q, want %q", f.Order, "published_at")
		}
	})

	t.Run("published_after_valid_RFC3339", func(t *testing.T) {
		ts := "2026-04-22T00:00:00+09:00"
		args := map[string]interface{}{
			"published_after": ts,
		}
		f, err := buildEntriesFilter(args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected, _ := time.Parse(time.RFC3339, ts)
		if f.PublishedAfter != expected.Unix() {
			t.Errorf("PublishedAfter: got %d, want %d", f.PublishedAfter, expected.Unix())
		}
	})

	t.Run("published_before_valid_RFC3339", func(t *testing.T) {
		ts := "2026-04-22T23:59:59Z"
		args := map[string]interface{}{
			"published_before": ts,
		}
		f, err := buildEntriesFilter(args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected, _ := time.Parse(time.RFC3339, ts)
		if f.PublishedBefore != expected.Unix() {
			t.Errorf("PublishedBefore: got %d, want %d", f.PublishedBefore, expected.Unix())
		}
	})

	t.Run("invalid_RFC3339_published_after_returns_error", func(t *testing.T) {
		args := map[string]interface{}{
			"published_after": "not-a-date",
		}
		_, err := buildEntriesFilter(args)
		if err == nil {
			t.Error("expected error for invalid RFC3339, got nil")
		}
	})

	t.Run("invalid_RFC3339_published_before_returns_error", func(t *testing.T) {
		args := map[string]interface{}{
			"published_before": "2026/04/22",
		}
		_, err := buildEntriesFilter(args)
		if err == nil {
			t.Error("expected error for invalid RFC3339, got nil")
		}
	})

	t.Run("category_id", func(t *testing.T) {
		args := map[string]interface{}{
			"category_id": float64(5),
		}
		f, err := buildEntriesFilter(args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if f.CategoryID != 5 {
			t.Errorf("CategoryID: got %d, want 5", f.CategoryID)
		}
	})

	t.Run("search", func(t *testing.T) {
		args := map[string]interface{}{
			"search": "golang",
		}
		f, err := buildEntriesFilter(args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if f.Search != "golang" {
			t.Errorf("Search: got %q, want %q", f.Search, "golang")
		}
	})

	t.Run("starred_true", func(t *testing.T) {
		args := map[string]interface{}{
			"starred": true,
		}
		f, err := buildEntriesFilter(args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if f.Starred != client.FilterOnlyStarred {
			t.Errorf("Starred: got %q, want %q", f.Starred, client.FilterOnlyStarred)
		}
	})

	t.Run("starred_false", func(t *testing.T) {
		args := map[string]interface{}{
			"starred": false,
		}
		f, err := buildEntriesFilter(args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if f.Starred != client.FilterNotStarred {
			t.Errorf("Starred: got %q, want %q", f.Starred, client.FilterNotStarred)
		}
	})

	t.Run("all_params_combined", func(t *testing.T) {
		ts := "2026-04-22T00:00:00Z"
		args := map[string]interface{}{
			"status":          "unread",
			"feed_id":         float64(1),
			"limit":           float64(20),
			"offset":          float64(0),
			"direction":       "desc",
			"order":           "published_at",
			"published_after": ts,
			"category_id":     float64(3),
			"search":          "miniflux",
			"starred":         true,
		}
		f, err := buildEntriesFilter(args)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected, _ := time.Parse(time.RFC3339, ts)
		if f.Status != "unread" {
			t.Errorf("Status: got %q, want unread", f.Status)
		}
		if f.FeedID != 1 {
			t.Errorf("FeedID: got %d, want 1", f.FeedID)
		}
		if f.Limit != 20 {
			t.Errorf("Limit: got %d, want 20", f.Limit)
		}
		if f.Direction != "desc" {
			t.Errorf("Direction: got %q, want desc", f.Direction)
		}
		if f.Order != "published_at" {
			t.Errorf("Order: got %q, want published_at", f.Order)
		}
		if f.PublishedAfter != expected.Unix() {
			t.Errorf("PublishedAfter: got %d, want %d", f.PublishedAfter, expected.Unix())
		}
		if f.CategoryID != 3 {
			t.Errorf("CategoryID: got %d, want 3", f.CategoryID)
		}
		if f.Search != "miniflux" {
			t.Errorf("Search: got %q, want miniflux", f.Search)
		}
		if f.Starred != client.FilterOnlyStarred {
			t.Errorf("Starred: got %q, want %q", f.Starred, client.FilterOnlyStarred)
		}
	})
}
