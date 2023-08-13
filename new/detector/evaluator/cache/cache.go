package cache

import (
	"github.com/bearer/bearer/new/detector/detection"
	"github.com/bearer/bearer/pkg/ast/tree"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/util/set"
)

const (
	maxCacheSize = 10_000
	evictionSize = 1000
)

// FIXME: get this dynamically
var builtinRuleIDs = []string{
	"datatype",
	"insecure_url",
	"object",
	"string",
	"string_literal",
}

type Key struct {
	rootNode   *tree.Node
	ruleID     string
	scope      settings.RuleReferenceScope
	followFlow bool
}

func NewKey(
	rootNode *tree.Node,
	ruleID string,
	scope settings.RuleReferenceScope,
	followFlow bool,
) Key {
	return Key{
		rootNode:   rootNode,
		ruleID:     ruleID,
		scope:      scope,
		followFlow: followFlow,
	}
}

type cacheMap map[Key][]*detection.Detection

type Shared struct {
	ruleIDs set.Set[string]
	data    cacheMap
}

func NewShared(ruleIDs []string) *Shared {
	idSet := set.New[string]()
	idSet.AddAll(ruleIDs)
	idSet.AddAll(builtinRuleIDs)

	return &Shared{
		ruleIDs: idSet,
		data:    make(cacheMap),
	}
}

type Cache struct {
	shared *Shared
	data   cacheMap
}

func NewCache(sharedCache *Shared) *Cache {
	return &Cache{
		shared: sharedCache,
		data:   make(cacheMap),
	}
}

func (cache *Cache) Get(key Key) ([]*detection.Detection, bool) {
	detections, cached := cache.dataFor(key)[key]
	return detections, cached
}

func (cache *Cache) Put(key Key, detections []*detection.Detection) {
	data := cache.dataFor(key)

	if len(data) > maxCacheSize {
		i := 0
		for evictedKey := range data {
			if i == evictionSize {
				break
			}

			delete(data, evictedKey)
			i++
		}
	}

	data[key] = detections
}

func (cache *Cache) dataFor(key Key) cacheMap {
	if cache.shared.ruleIDs.Has(key.ruleID) {
		return cache.shared.data
	}

	return cache.data
}
