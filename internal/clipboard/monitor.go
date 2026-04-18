package clipboard

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/atotto/clipboard"
)

type ClipboardItem struct {
	Content   string
	Timestamp time.Time
	Hash      string
}

type Monitor struct {
	lastHash      string
	lastContent   string
	minLength     int
	stopChan      chan struct{}
	wg            sync.WaitGroup
	itemChan      chan ClipboardItem
	pollInterval  time.Duration
	excludeWords  []string
}

type MonitorConfig struct {
	MinLength     int
	PollInterval  time.Duration
	ExcludeWords  []string
	BufferSize    int
}

func NewMonitor(config MonitorConfig) *Monitor {
	if config.MinLength == 0 {
		config.MinLength = 100
	}
	if config.PollInterval == 0 {
		config.PollInterval = 500 * time.Millisecond
	}
	if config.BufferSize == 0 {
		config.BufferSize = 10
	}

	return &Monitor{
		minLength:    config.MinLength,
		pollInterval: config.PollInterval,
		excludeWords: config.ExcludeWords,
		stopChan:     make(chan struct{}),
		itemChan:     make(chan ClipboardItem, config.BufferSize),
	}
}

func (m *Monitor) hashContent(content string) string {
	hash := md5.Sum([]byte(content))
	return hex.EncodeToString(hash[:])
}

func (m *Monitor) shouldSync(content string) bool {
	// 检查长度
	if len(content) < m.minLength {
		return false
	}

	// 检查是否包含排除词
	for _, word := range m.excludeWords {
		if strings.Contains(content, word) {
			log.Printf("[DEBUG] Content excluded due to keyword: %s", word)
			return false
		}
	}

	// 检查是否是重复内容
	currentHash := m.hashContent(content)
	if currentHash == m.lastHash {
		return false
	}

	return true
}

func (m *Monitor) readClipboard() (string, error) {
	content, err := clipboard.ReadAll()
	if err != nil {
		return "", fmt.Errorf("failed to read clipboard: %w", err)
	}
	return content, nil
}

func (m *Monitor) Start() {
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		log.Println("[INFO] Clipboard monitor started")

		ticker := time.NewTicker(m.pollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-m.stopChan:
				log.Println("[INFO] Clipboard monitor stopped")
				return
			case <-ticker.C:
				content, err := m.readClipboard()
				if err != nil {
					log.Printf("[ERROR] Failed to read clipboard: %v", err)
					continue
				}

				if m.shouldSync(content) {
					m.lastHash = m.hashContent(content)
					m.lastContent = content

					item := ClipboardItem{
						Content:   content,
						Timestamp: time.Now(),
						Hash:      m.lastHash,
					}

					select {
					case m.itemChan <- item:
						log.Printf("[INFO] New clipboard item detected: %s (length: %d)", m.lastHash[:8], len(content))
					default:
						log.Printf("[WARN] Item channel full, dropping item")
					}
				}
			}
		}
	}()
}

func (m *Monitor) Stop() {
	close(m.stopChan)
	m.wg.Wait()
	close(m.itemChan)
}

func (m *Monitor) Items() <-chan ClipboardItem {
	return m.itemChan
}

func (m *Monitor) GetLastContent() string {
	return m.lastContent
}
