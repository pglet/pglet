package page

import (
	"fmt"
	"strings"
	"sync"
)

type pages struct {
	sync.RWMutex
	items map[string]*Page
}

var (
	pr   *pages
	once sync.Once
)

// Pages holds the index of all pages.
func Pages() *pages {
	once.Do(func() {
		pr = &pages{
			items: make(map[string]*Page),
		}
	})

	return pr
}

func (pr *pages) AddPage(p *Page) error {
	pr.Lock()
	defer pr.Unlock()

	if _, exists := pr.items[p.Name]; exists {
		return fmt.Errorf("page with '%s' name already exists", p.Name)
	}

	pr.items[p.Name] = p
	return nil
}

func (pr *pages) GetPage(name string) *Page {
	pr.RLock()
	defer pr.RUnlock()
	return pr.items[name]
}

func (pr *pages) RemovePage(name string) {
	pr.Lock()
	defer pr.Lock()
	delete(pr.items, name)
}

func (pr *pages) String() string {
	keys := make([]string, 0, len(pr.items))
	for key := range pr.items {
		keys = append(keys, key)
	}

	return fmt.Sprintf("[%s]", strings.Join(keys, ", "))
}
