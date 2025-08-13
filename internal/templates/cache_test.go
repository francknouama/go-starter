package templates

import (
	"testing"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
)

func TestTemplateCache(t *testing.T) {
	cache := NewTemplateCache(10, 5*time.Minute)
	
	t.Run("cache miss and set", func(t *testing.T) {
		key := "test-template"
		content := "Hello {{.Name}}"
		
		// Should be a cache miss
		tmpl, cachedContent, found := cache.Get(key)
		if found {
			t.Error("Expected cache miss, got hit")
		}
		if tmpl != nil {
			t.Error("Expected nil template on cache miss")
		}
		if cachedContent != "" {
			t.Error("Expected empty content on cache miss")
		}
		
		// Parse and cache template
		parsedTmpl := template.New(key).Funcs(sprig.TxtFuncMap())
		parsedTmpl, err := parsedTmpl.Parse(content)
		if err != nil {
			t.Fatalf("Failed to parse template: %v", err)
		}
		
		cache.Set(key, parsedTmpl, content)
		
		// Should be a cache hit now
		cachedTmpl, cachedContent, found := cache.Get(key)
		if !found {
			t.Error("Expected cache hit, got miss")
		}
		if cachedTmpl == nil {
			t.Error("Expected cached template, got nil")
		}
		if cachedContent != content {
			t.Errorf("Expected content %q, got %q", content, cachedContent)
		}
	})
	
	t.Run("cache eviction", func(t *testing.T) {
		smallCache := NewTemplateCache(2, 5*time.Minute)
		
		// Fill cache to capacity
		tmpl1 := template.New("tmpl1")
		tmpl2 := template.New("tmpl2")
		
		smallCache.Set("key1", tmpl1, "content1")
		smallCache.Set("key2", tmpl2, "content2")
		
		// Both should be cached
		_, _, found1 := smallCache.Get("key1")
		_, _, found2 := smallCache.Get("key2")
		
		if !found1 || !found2 {
			t.Error("Both templates should be cached")
		}
		
		// Add third template, should evict least recently used
		tmpl3 := template.New("tmpl3")
		smallCache.Set("key3", tmpl3, "content3")
		
		// key1 should be evicted (oldest)
		_, _, found1 = smallCache.Get("key1")
		_, _, found2 = smallCache.Get("key2")
		_, _, found3 := smallCache.Get("key3")
		
		if found1 {
			t.Error("key1 should have been evicted")
		}
		if !found2 {
			t.Error("key2 should still be cached")
		}
		if !found3 {
			t.Error("key3 should be cached")
		}
	})
	
	t.Run("cache expiration", func(t *testing.T) {
		shortCache := NewTemplateCache(10, 1*time.Millisecond)
		
		tmpl := template.New("test")
		shortCache.Set("expired", tmpl, "content")
		
		// Wait for expiration
		time.Sleep(2 * time.Millisecond)
		
		_, _, found := shortCache.Get("expired")
		if found {
			t.Error("Template should have expired")
		}
	})
	
	t.Run("cache stats", func(t *testing.T) {
		cache := NewTemplateCache(10, 5*time.Minute)
		
		// Initially empty
		size, hits := cache.Stats()
		if size != 0 || hits != 0 {
			t.Errorf("Expected empty cache stats, got size=%d, hits=%d", size, hits)
		}
		
		// Add templates and access them
		tmpl1 := template.New("tmpl1")
		tmpl2 := template.New("tmpl2")
		
		cache.Set("key1", tmpl1, "content1")
		cache.Set("key2", tmpl2, "content2")
		
		// Access templates to increment hit count
		cache.Get("key1")
		cache.Get("key1")
		cache.Get("key2")
		
		size, hits = cache.Stats()
		if size != 2 {
			t.Errorf("Expected cache size 2, got %d", size)
		}
		if hits != 3 {
			t.Errorf("Expected 3 cache hits, got %d", hits)
		}
	})
	
	t.Run("cache clear", func(t *testing.T) {
		cache := NewTemplateCache(10, 5*time.Minute)
		
		tmpl := template.New("test")
		cache.Set("key", tmpl, "content")
		
		// Verify it's cached
		_, _, found := cache.Get("key")
		if !found {
			t.Error("Template should be cached")
		}
		
		// Clear cache
		cache.Clear()
		
		// Should be empty now
		_, _, found = cache.Get("key")
		if found {
			t.Error("Cache should be empty after clear")
		}
		
		size, hits := cache.Stats()
		if size != 0 || hits != 0 {
			t.Errorf("Cache should be empty after clear, got size=%d, hits=%d", size, hits)
		}
	})
}

func TestGetOrParseTemplate(t *testing.T) {
	// Clear global cache before testing
	ClearTemplateCache()
	
	t.Run("parse and cache template", func(t *testing.T) {
		key := "test-template"
		content := "Hello {{.Name}}"
		
		// First call should parse and cache
		tmpl1, err := GetOrParseTemplate(key, content, nil)
		if err != nil {
			t.Fatalf("Failed to parse template: %v", err)
		}
		if tmpl1 == nil {
			t.Error("Expected parsed template, got nil")
		}
		
		// Second call should return cached version
		tmpl2, err := GetOrParseTemplate(key, content, nil)
		if err != nil {
			t.Fatalf("Failed to get cached template: %v", err)
		}
		if tmpl2 == nil {
			t.Error("Expected cached template, got nil")
		}
		
		// Templates should be clones (different instances)
		if tmpl1 == tmpl2 {
			t.Error("Expected different template instances (clones)")
		}
	})
	
	t.Run("content change invalidates cache", func(t *testing.T) {
		key := "changing-template"
		content1 := "Hello {{.Name}}"
		content2 := "Hi {{.Name}}"
		
		// Parse first version
		tmpl1, err := GetOrParseTemplate(key, content1, nil)
		if err != nil {
			t.Fatalf("Failed to parse template: %v", err)
		}
		
		// Parse second version with different content
		tmpl2, err := GetOrParseTemplate(key, content2, nil)
		if err != nil {
			t.Fatalf("Failed to parse template: %v", err)
		}
		
		// Should get fresh template, not cached one
		if tmpl1 == tmpl2 {
			t.Error("Expected different templates for different content")
		}
	})
	
	t.Run("template parsing error", func(t *testing.T) {
		key := "invalid-template"
		content := "Hello {{.Name" // Missing closing braces
		
		_, err := GetOrParseTemplate(key, content, nil)
		if err == nil {
			t.Error("Expected parsing error for invalid template")
		}
	})
	
	t.Run("cache stats integration", func(t *testing.T) {
		ClearTemplateCache()
		
		key1, content1 := "tmpl1", "Content 1"
		key2, content2 := "tmpl2", "Content 2"
		
		// Parse two templates
		GetOrParseTemplate(key1, content1, nil)
		GetOrParseTemplate(key2, content2, nil)
		
		// Access them again (should be cache hits)
		GetOrParseTemplate(key1, content1, nil)
		GetOrParseTemplate(key2, content2, nil)
		
		size, hits := GetTemplateCacheStats()
		if size != 2 {
			t.Errorf("Expected cache size 2, got %d", size)
		}
		if hits != 2 {
			t.Errorf("Expected 2 cache hits, got %d", hits)
		}
	})
}

func BenchmarkTemplateCache(b *testing.B) {
	ClearTemplateCache()
	
	content := `
package {{.PackageName}}

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"{{.ModulePath}}/internal/config"
	"{{.ModulePath}}/internal/logger"
	{{if eq .Framework "gin"}}
	"github.com/gin-gonic/gin"
	{{end}}
)

func main() {
	cfg := config.Load()
	logger := logger.New(cfg.LogLevel)
	
	{{if eq .Framework "gin"}}
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	{{end}}
	
	logger.Info(fmt.Sprintf("Server starting on port %d", cfg.Port))
}
`
	
	b.Run("without cache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tmpl := template.New("test").Funcs(sprig.TxtFuncMap())
			_, err := tmpl.Parse(content)
			if err != nil {
				b.Fatalf("Failed to parse template: %v", err)
			}
		}
	})
	
	b.Run("with cache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := GetOrParseTemplate("benchmark-template", content, nil)
			if err != nil {
				b.Fatalf("Failed to parse template: %v", err)
			}
		}
	})
}