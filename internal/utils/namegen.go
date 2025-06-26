package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// Adjectives for project name generation (GitHub-style)
var adjectives = []string{
	"amazing", "awesome", "brilliant", "clever", "creative", "elegant", "fantastic",
	"friendly", "gentle", "glorious", "graceful", "incredible", "inspiring", "magical",
	"marvelous", "mighty", "perfect", "powerful", "radiant", "shining", "sparkling",
	"stellar", "stunning", "superb", "supreme", "vibrant", "wonderful", "zesty",
	"bold", "brave", "bright", "calm", "charming", "cheerful", "clean", "cool",
	"cozy", "crisp", "dazzling", "delightful", "divine", "dynamic", "energetic",
	"epic", "exquisite", "fabulous", "fresh", "glowing", "golden", "grand",
	"happy", "harmonious", "honest", "luminous", "majestic", "noble", "optimal",
	"peaceful", "pleasant", "polished", "precious", "pristine", "pure", "quick",
	"quiet", "rapid", "refined", "robust", "serene", "sharp", "sleek", "smart",
	"smooth", "solid", "stable", "strong", "swift", "tranquil", "trusty", "warm",
}

// Nouns for project name generation (tech/general themed)
var nouns = []string{
	"app", "api", "bot", "cli", "code", "core", "dash", "edge", "flow", "forge",
	"gate", "grid", "hub", "kit", "lab", "link", "mesh", "node", "path", "pilot",
	"port", "pulse", "quest", "rail", "sage", "sync", "tool", "wave", "zone",
	"atlas", "beacon", "bridge", "castle", "cloud", "craft", "crystal", "engine",
	"falcon", "fiber", "flame", "galaxy", "garden", "harbor", "helmet", "island",
	"journey", "kernel", "ladder", "matrix", "mirror", "nexus", "ocean", "palace",
	"planet", "portal", "prism", "rocket", "shadow", "sphere", "summit", "temple",
	"tower", "vault", "vessel", "village", "vision", "widget", "wizard", "wonder",
	"anchor", "arrow", "badge", "barrel", "basket", "beacon", "bicycle", "binder",
	"bottle", "branch", "bucket", "button", "cabinet", "camera", "canvas", "carpet",
	"circle", "compass", "corner", "cutter", "diamond", "drawer", "driver", "folder",
	"hammer", "handle", "helper", "hockey", "holder", "hunter", "jacket", "ladder",
	"laser", "lever", "marble", "marker", "master", "meter", "motor", "needle",
	"panel", "paper", "pencil", "picker", "player", "pocket", "printer", "puzzle",
	"reader", "robot", "ruler", "scanner", "server", "slider", "socket", "spider",
	"spring", "switch", "tablet", "ticker", "timer", "tracker", "tunnel", "turtle",
	"vector", "viewer", "wallet", "window", "worker", "writer",
}

var rng *rand.Rand

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// GenerateRandomProjectName generates a random project name in the format "adjective-noun"
// Similar to GitHub's repository name suggestions
func GenerateRandomProjectName() string {
	adjective := adjectives[rng.Intn(len(adjectives))]
	noun := nouns[rng.Intn(len(nouns))]
	return fmt.Sprintf("%s-%s", adjective, noun)
}

// GenerateMultipleNames generates multiple random project name suggestions
func GenerateMultipleNames(count int) []string {
	if count <= 0 {
		count = 3
	}

	names := make([]string, count)
	used := make(map[string]bool)

	for i := 0; i < count; i++ {
		// Ensure uniqueness within the generated set
		var name string
		attempts := 0
		for {
			name = GenerateRandomProjectName()
			if !used[name] || attempts > 50 { // Prevent infinite loop
				break
			}
			attempts++
		}
		names[i] = name
		used[name] = true
	}

	return names
}

// GenerateWithPrefix generates a random name with a specific prefix
func GenerateWithPrefix(prefix string) string {
	if prefix == "" {
		return GenerateRandomProjectName()
	}

	noun := nouns[rng.Intn(len(nouns))]
	return fmt.Sprintf("%s-%s", prefix, noun)
}

// GenerateWithSuffix generates a random name with a specific suffix
func GenerateWithSuffix(suffix string) string {
	if suffix == "" {
		return GenerateRandomProjectName()
	}

	adjective := adjectives[rng.Intn(len(adjectives))]
	return fmt.Sprintf("%s-%s", adjective, suffix)
}

// IsValidProjectNameChar checks if a character is valid for project names
func IsValidProjectNameChar(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') ||
		ch == '-' || ch == '_'
}

// SanitizeProjectName removes invalid characters from a project name
func SanitizeProjectName(name string) string {
	result := make([]rune, 0, len(name))

	for _, ch := range name {
		if IsValidProjectNameChar(ch) {
			result = append(result, ch)
		}
	}

	return string(result)
}
