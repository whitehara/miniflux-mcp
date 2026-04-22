package main

import (
	"fmt"
	"time"

	"miniflux.app/v2/client"
)

// buildEntriesFilter constructs a client.Filter from MCP tool argument map.
// Returns (nil, nil) when argsMap is nil.
// Returns an error if published_after or published_before contains an invalid RFC3339 string.
func buildEntriesFilter(argsMap map[string]interface{}) (*client.Filter, error) {
	if argsMap == nil {
		return nil, nil
	}

	filter := &client.Filter{}

	if v, ok := argsMap["status"].(string); ok {
		filter.Status = v
	}
	if v, ok := argsMap["feed_id"].(float64); ok {
		filter.FeedID = int64(v)
	}
	if v, ok := argsMap["limit"].(float64); ok {
		filter.Limit = int(v)
	}
	if v, ok := argsMap["offset"].(float64); ok {
		filter.Offset = int(v)
	}
	if v, ok := argsMap["direction"].(string); ok {
		filter.Direction = v
	}
	if v, ok := argsMap["order"].(string); ok {
		filter.Order = v
	}
	if v, ok := argsMap["published_after"].(string); ok {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return nil, fmt.Errorf("published_after: invalid RFC3339 value %q: %w", v, err)
		}
		filter.PublishedAfter = t.Unix()
	}
	if v, ok := argsMap["published_before"].(string); ok {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return nil, fmt.Errorf("published_before: invalid RFC3339 value %q: %w", v, err)
		}
		filter.PublishedBefore = t.Unix()
	}
	if v, ok := argsMap["category_id"].(float64); ok {
		filter.CategoryID = int64(v)
	}
	if v, ok := argsMap["search"].(string); ok {
		filter.Search = v
	}
	if v, ok := argsMap["starred"].(bool); ok {
		if v {
			filter.Starred = client.FilterOnlyStarred
		} else {
			filter.Starred = client.FilterNotStarred
		}
	}

	return filter, nil
}
