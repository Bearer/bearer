package cache

import (
	"golang.org/x/exp/slices"

	"github.com/bearer/bearer/internal/scanner/ast/tree"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"
	"github.com/bearer/bearer/internal/util/set"
)

const (
	maxCacheSize = 100
	evictionSize = 20
)

type entry struct {
	Node       *tree.Node
	RuleID     string
	Detections []*detectortypes.Detection
}

type Shared struct {
	ruleIDs set.Set[string]
	data    []entry
}

func NewShared(ruleIDs []string) *Shared {
	idSet := set.New[string]()
	idSet.AddAll(ruleIDs)

	return &Shared{
		ruleIDs: idSet,
	}
}

type Cache struct {
	shared *Shared
	data   []entry
}

func NewCache(sharedCache *Shared) *Cache {
	return &Cache{
		shared: sharedCache,
	}
}

func (cache *Cache) Get(node *tree.Node, ruleID string) ([]*detectortypes.Detection, bool) {
	for _, entry := range *cache.dataFor(ruleID) {
		if entry.Node == node && entry.RuleID == ruleID {
			return entry.Detections, true
		}
	}

	return nil, false
}

func (cache *Cache) Put(node *tree.Node, ruleID string, detections []*detectortypes.Detection) {
	data := cache.dataFor(ruleID)

	if len(*data) > maxCacheSize {
		*data = slices.Delete(*data, 0, evictionSize)
	}

	*data = append(*data, entry{
		Node:       node,
		RuleID:     ruleID,
		Detections: detections,
	})
}

func (cache *Cache) dataFor(ruleID string) *[]entry {
	if cache.shared.ruleIDs.Has(ruleID) {
		return &cache.shared.data
	} else {
		return &cache.data
	}
}
