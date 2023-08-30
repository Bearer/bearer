package cache

import (
	"github.com/rs/zerolog/log"

	"github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/detectorset"
	"github.com/bearer/bearer/internal/util/set"
)

const (
	maxCacheSize = 1000
	evictionSize = 100
)

type entry struct {
	RuleID string
	Result *detectorset.Result
}

type Shared struct {
	ruleIDs set.Set[string]
	data    map[*tree.Node][]entry
}

func NewShared(ruleIDs []string) *Shared {
	idSet := set.New[string]()
	idSet.AddAll(ruleIDs)

	return &Shared{
		ruleIDs: idSet,
		data:    make(map[*tree.Node][]entry),
	}
}

type Cache struct {
	shared *Shared
	data   map[*tree.Node][]entry
}

func NewCache(sharedCache *Shared) *Cache {
	return &Cache{
		shared: sharedCache,
		data:   make(map[*tree.Node][]entry),
	}
}

func (cache *Cache) Get(node *tree.Node, ruleID string) (*detectorset.Result, bool) {
	if cache == nil {
		return nil, false
	}

	for _, entry := range cache.dataFor(ruleID)[node] {
		if entry.RuleID == ruleID {
			return entry.Result, true
		}
	}

	return nil, false
}

func (cache *Cache) Put(node *tree.Node, ruleID string, result *detectorset.Result) {
	if cache == nil {
		return
	}

	data := cache.dataFor(ruleID)

	if len(data) > maxCacheSize {
		log.Trace().Msg("detection cache full, evicting entries")

		i := 0
		for evictedNode := range data {
			if i == evictionSize {
				break
			}

			data[evictedNode] = nil
			delete(data, evictedNode)

			i++
		}
	}

	data[node] = append(data[node], entry{
		RuleID: ruleID,
		Result: result,
	})
}

func (cache *Cache) dataFor(ruleID string) map[*tree.Node][]entry {
	if cache.shared.ruleIDs.Has(ruleID) {
		return cache.shared.data
	} else {
		return cache.data
	}
}
