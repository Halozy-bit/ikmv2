package cache

import (
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Pagination *pagination

func init() {
	Pagination = new(pagination)
	Pagination.mux = new(sync.RWMutex)
}

var ErrStoreSameID = fmt.Errorf("cache store same id")

// Immutable
type pagination struct {
	page     []primitive.ObjectID
	category map[string][]primitive.ObjectID
	mux      *sync.RWMutex
}

func (p *pagination) Leng() int {
	p.mux.RLock()
	defer p.mux.RUnlock()
	leng := len(p.page)
	return leng
}

func (p *pagination) Page(idx int) primitive.ObjectID {
	p.mux.RLock()
	defer p.mux.RUnlock()
	if idx > len(p.page) {
		return primitive.NilObjectID
	}
	id := p.page[idx-1]
	return id
}

// page starts from page 1
func (p *pagination) CategoryPage(category string, idx int) primitive.ObjectID {
	p.mux.RLock()
	defer p.mux.RUnlock()
	if idx > len(p.category[category]) {
		return primitive.NilObjectID
	}
	id := p.category[category][idx-1]
	return id
}

func (p *pagination) StoreCategory(category string, ids []primitive.ObjectID) error {
	p.mux.Lock()
	defer p.mux.Unlock()
	if len(p.category[category]) > 0 {
		if p.category[category][0] == ids[0] {
			return ErrStoreSameID
		}
	}

	p.category[category] = ids
	return nil
}

func (p *pagination) StorePage(ids []primitive.ObjectID) error {
	p.mux.Lock()
	defer p.mux.Unlock()
	if len(p.page) > 0 {
		if p.page[0] == ids[0] {
			return ErrStoreSameID
		}
	}
	p.page = ids
	return nil
}
