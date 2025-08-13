package templates

import (
	"sync"
	"text/template"
	"time"
	
	"github.com/Masterminds/sprig/v3"
)

// TemplateCache provides thread-safe caching for parsed templates
type TemplateCache struct {
	cache    map[string]*cachedTemplate
	mu       sync.RWMutex
	maxSize  int
	maxAge   time.Duration
}

// cachedTemplate stores a parsed template with metadata
type cachedTemplate struct {
	template  *template.Template
	content   string
	timestamp time.Time
	hits      int64
}

// NewTemplateCache creates a new template cache with specified limits
func NewTemplateCache(maxSize int, maxAge time.Duration) *TemplateCache {
	return &TemplateCache{
		cache:   make(map[string]*cachedTemplate),
		maxSize: maxSize,
		maxAge:  maxAge,
	}
}

// Get retrieves a cached template or returns nil if not found/expired
func (tc *TemplateCache) Get(key string) (*template.Template, string, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	
	cached, exists := tc.cache[key]
	if !exists {
		return nil, "", false
	}
	
	// Check if template is expired
	if time.Since(cached.timestamp) > tc.maxAge {
		return nil, "", false
	}
	
	cached.hits++
	return cached.template, cached.content, true
}

// Set stores a parsed template in the cache
func (tc *TemplateCache) Set(key string, tmpl *template.Template, content string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	
	// Evict least recently used if at capacity
	if len(tc.cache) >= tc.maxSize {
		tc.evictLRU()
	}
	
	tc.cache[key] = &cachedTemplate{
		template:  tmpl,
		content:   content,
		timestamp: time.Now(),
		hits:      0,
	}
}

// evictLRU removes the least recently used template
func (tc *TemplateCache) evictLRU() {
	var lruKey string
	var lruTime time.Time
	
	for key, cached := range tc.cache {
		if lruKey == "" || cached.timestamp.Before(lruTime) {
			lruKey = key
			lruTime = cached.timestamp
		}
	}
	
	if lruKey != "" {
		delete(tc.cache, lruKey)
	}
}

// Clear removes all cached templates
func (tc *TemplateCache) Clear() {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	
	tc.cache = make(map[string]*cachedTemplate)
}

// Stats returns cache statistics
func (tc *TemplateCache) Stats() (size int, totalHits int64) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	
	size = len(tc.cache)
	for _, cached := range tc.cache {
		totalHits += cached.hits
	}
	return
}

// ParsedTemplateCache is a global cache for parsed templates
var globalTemplateCache = NewTemplateCache(100, 30*time.Minute)

// GetOrParseTemplate retrieves a cached template or parses it if not cached
func GetOrParseTemplate(key string, content string, funcMap template.FuncMap) (*template.Template, error) {
	// Check cache first
	if cachedTmpl, cachedContent, found := globalTemplateCache.Get(key); found && cachedContent == content {
		// Clone the template to avoid concurrent modification
		cloned, err := cachedTmpl.Clone()
		if err == nil {
			return cloned, nil
		}
		// If clone fails, continue to parse fresh
	}
	
	// Parse template
	tmpl := template.New(key).Funcs(sprig.TxtFuncMap())
	if funcMap != nil {
		tmpl = tmpl.Funcs(funcMap)
	}
	
	parsed, err := tmpl.Parse(content)
	if err != nil {
		return nil, err
	}
	
	// Cache the parsed template
	globalTemplateCache.Set(key, parsed, content)
	
	// Return a clone to avoid concurrent modification
	return parsed.Clone()
}

// ClearTemplateCache clears the global template cache
func ClearTemplateCache() {
	globalTemplateCache.Clear()
}

// GetTemplateCacheStats returns global cache statistics
func GetTemplateCacheStats() (size int, totalHits int64) {
	return globalTemplateCache.Stats()
}