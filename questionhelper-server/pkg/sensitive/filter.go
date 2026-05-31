package sensitive

import (
	"strings"
	"sync"
)

// Filter 敏感词过滤器
type Filter struct {
	words map[string]bool
	mu    sync.RWMutex
}

// NewFilter 创建过滤器
func NewFilter() *Filter {
	f := &Filter{
		words: make(map[string]bool),
	}
	f.loadDefaultWords()
	return f
}

// loadDefaultWords 加载默认敏感词
func (f *Filter) loadDefaultWords() {
	// TODO: 从文件或数据库加载敏感词库
	defaultWords := []string{
		// 这里只是示例，实际应该从敏感词库文件加载
	}
	f.AddWords(defaultWords)
}

// AddWord 添加敏感词
func (f *Filter) AddWord(word string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.words[word] = true
}

// AddWords 批量添加敏感词
func (f *Filter) AddWords(words []string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, word := range words {
		f.words[word] = true
	}
}

// HasSensitive 检查是否包含敏感词
func (f *Filter) HasSensitive(content string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	for word := range f.words {
		if strings.Contains(content, word) {
			return true
		}
	}
	return false
}

// Replace 替换敏感词
func (f *Filter) Replace(content string, replacement rune) string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	result := content
	for word := range f.words {
		if strings.Contains(result, word) {
			result = strings.ReplaceAll(result, word, strings.Repeat(string(replacement), len(word)))
		}
	}
	return result
}

// FindAll 查找所有敏感词
func (f *Filter) FindAll(content string) []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var found []string
	for word := range f.words {
		if strings.Contains(content, word) {
			found = append(found, word)
		}
	}
	return found
}
