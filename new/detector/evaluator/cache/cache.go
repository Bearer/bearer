package cache

import (
	"github.com/bearer/bearer/new/detector/detection"
	"github.com/bearer/bearer/new/language/tree"
	"github.com/bearer/bearer/pkg/commands/process/settings"
	"github.com/bearer/bearer/pkg/util/set"
)

type Key struct {
	rootNodeID tree.NodeID
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
		rootNodeID: rootNode.ID(),
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
	cache.dataFor(key)[key] = detections
}

func (cache *Cache) dataFor(key Key) cacheMap {
	if cache.shared.ruleIDs.Has(key.ruleID) {
		return cache.shared.data
	}

	return cache.data
}
