package cache

import (
	"github.com/rs/zerolog/log"

	treepkg "github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/detectorset"
	"github.com/bearer/bearer/internal/scanner/ruleset"
	"github.com/bearer/bearer/internal/util/set"
)

const (
	minNodeCount = 20_000 // minimum number of AST nodes for cache to be enabled
	maxCacheSize = 1000   // maximum number of AST nodes per cache/shared-cache
	evictionSize = 100    // number of AST nodes to evict when max is reached
)

type entry struct {
	rule   *ruleset.Rule
	result *detectorset.Result
}

type Shared struct {
	ruleIndexSet set.Set[int]
	data         map[*treepkg.Node][]entry
}

func NewShared(rules []*ruleset.Rule) *Shared {
	ruleIndexSet := set.New[int]()
	for _, rule := range rules {
		if rule.Type() == ruleset.RuleTypeBuiltin || rule.Type() == ruleset.RuleTypeShared {
			ruleIndexSet.Add(rule.Index())
		}
	}

	return &Shared{
		ruleIndexSet: ruleIndexSet,
		data:         make(map[*treepkg.Node][]entry),
	}
}

type Cache struct {
	enabled bool
	shared  *Shared
	data    map[*treepkg.Node][]entry
}

func NewCache(tree *treepkg.Tree, sharedCache *Shared) *Cache {
	enabled := tree.NodeCount() > minNodeCount
	if enabled {
		log.Trace().Msg("cache enabled")
	}

	return &Cache{
		enabled: enabled,
		shared:  sharedCache,
		data:    make(map[*treepkg.Node][]entry),
	}
}

func (cache *Cache) Clear() {
	if cache == nil || !cache.enabled {
		return
	}

	clear(cache.data)
}

func (cache *Cache) Get(node *treepkg.Node, rule *ruleset.Rule) (*detectorset.Result, bool) {
	if cache == nil || !cache.enabled {
		return nil, false
	}

	for _, entry := range cache.dataFor(rule)[node] {
		if entry.rule == rule {
			return entry.result, true
		}
	}

	return nil, false
}

func (cache *Cache) Put(node *treepkg.Node, rule *ruleset.Rule, result *detectorset.Result) {
	if cache == nil || !cache.enabled {
		return
	}

	data := cache.dataFor(rule)

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
		rule:   rule,
		result: result,
	})
}

func (cache *Cache) dataFor(rule *ruleset.Rule) map[*treepkg.Node][]entry {
	if cache.shared.ruleIndexSet.Has(rule.Index()) {
		return cache.shared.data
	} else {
		return cache.data
	}
}
