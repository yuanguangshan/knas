package ssh

import (
	"testing"
)

func TestExtractContentPrefix(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		n        int
		expected string
	}{
		{
			name:     "simple text",
			content:  "Hello World",
			n:        10,
			expected: "Hello_World",
		},
		{
			name:     "chinese characters",
			content:  "你好世界测试",
			n:        10,
			expected: "你好世界测试",
		},
		{
			name:     "mixed content",
			content:  "Test测试123",
			n:        10,
			expected: "Test测试123",
		},
		{
			name:     "with special characters",
			content:  "Hello/World\\Test",
			n:        20,
			expected: "Hello_World_Test",
		},
		{
			name:     "with newlines",
			content:  "Line1\nLine2\tLine3",
			n:        20,
			expected: "Line1_Line2_Line3",
		},
		{
			name:     "only special characters",
			content:  "!@#$%^&*()",
			n:        10,
			expected: "untitled",
		},
		{
			name:     "truncate long content",
			content:  "This is a very long content that should be truncated",
			n:        10,
			expected: "This_is_a_",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractContentPrefix(tt.content, tt.n)
			if result != tt.expected {
				t.Errorf("extractContentPrefix() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExpandPath(t *testing.T) {
	client := &Client{
		config: &Config{
			User: "testuser",
		},
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "path with tilde",
			input:    "~/test/path",
			expected: "/home/testuser/test/path",
		},
		{
			name:     "absolute path",
			input:    "/absolute/path",
			expected: "/absolute/path",
		},
		{
			name:     "relative path",
			input:    "relative/path",
			expected: "relative/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.expandPath(tt.input)
			if result != tt.expected {
				t.Errorf("expandPath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected *Config
	}{
		{
			name: "default port",
			config: &Config{
				Host:    "localhost",
				User:    "root",
				KeyPath: "~/.ssh/id_rsa",
			},
			expected: &Config{
				Host:    "localhost",
				Port:    "22",
				User:    "root",
				KeyPath: "~/.ssh/id_rsa",
			},
		},
		{
			name: "custom port",
			config: &Config{
				Host:    "localhost",
				Port:    "2222",
				User:    "root",
				KeyPath: "~/.ssh/id_rsa",
			},
			expected: &Config{
				Host:    "localhost",
				Port:    "2222",
				User:    "root",
				KeyPath: "~/.ssh/id_rsa",
			},
		},
		{
			name: "default base path",
			config: &Config{
				Host:    "localhost",
				User:    "root",
				KeyPath: "~/.ssh/id_rsa",
			},
			expected: &Config{
				Host:     "localhost",
				Port:     "22",
				User:     "root",
				KeyPath:  "~/.ssh/id_rsa",
				BasePath: "~/knas_archive",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.config)
			if client.config.Host != tt.expected.Host {
				t.Errorf("Host = %v, want %v", client.config.Host, tt.expected.Host)
			}
			if client.config.Port != tt.expected.Port {
				t.Errorf("Port = %v, want %v", client.config.Port, tt.expected.Port)
			}
			if client.config.User != tt.expected.User {
				t.Errorf("User = %v, want %v", client.config.User, tt.expected.User)
			}
			if tt.expected.BasePath != "" && client.config.BasePath != tt.expected.BasePath {
				t.Errorf("BasePath = %v, want %v", client.config.BasePath, tt.expected.BasePath)
			}
		})
	}
}
