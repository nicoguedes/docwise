package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port         string
	DatabaseURL  string
	OllamaURL   string
	EmbedModel   string
	ChatModel    string
	ChunkSize    int
	ChunkOverlap int
}

func Load() *Config {
	return &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://docwise:docwise@localhost:5432/docwise?sslmode=disable"),
		OllamaURL:   getEnv("OLLAMA_URL", "http://localhost:11434"),
		EmbedModel:   getEnv("EMBED_MODEL", "nomic-embed-text"),
		ChatModel:    getEnv("CHAT_MODEL", "llama3.1:8b"),
		ChunkSize:    getEnvInt("CHUNK_SIZE", 1000),
		ChunkOverlap: getEnvInt("CHUNK_OVERLAP", 200),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
