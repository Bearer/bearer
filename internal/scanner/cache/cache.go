package cache

import (
	"github.com/rs/zerolog/log"

	treepkg "github.com/bearer/bearer/internal/scanner/ast/tree"
	"github.com/bearer/bearer/internal/scanner/detectorset"
	"github.com/bearer/bearer/internal/util/set"
)

const (
	minNodeCount = 20_000 // minimum number of AST nodes for cache to be enabled
	maxCacheSize = 1000   // maximum number of AST nodes per cache/shared-cache
	evictionSize = 100    // number of AST nodes to evict when max is reached
)

type entry struct {
	DetectorID int
	Result     *detectorset.Result
}

type Shared struct {
	detectorIDs set.Set[int]
	data        map[*treepkg.Node][]entry
}

func NewShared(detectorIDs []int) *Shared {
	idSet := set.New[int]()
	idSet.AddAll(detectorIDs)

	return &Shared{
		detectorIDs: idSet,
		data:        make(map[*treepkg.Node][]entry),
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

func (cache *Cache) Get(node *treepkg.Node, detectorID int) (*detectorset.Result, bool) {
	if cache == nil || !cache.enabled {
		return nil, false
	}

	for _, entry := range cache.dataFor(detectorID)[node] {
		if entry.DetectorID == detectorID {
			return entry.Result, true
		}
	}

	return nil, false
}

func (cache *Cache) Put(node *treepkg.Node, detectorID int, result *detectorset.Result) {
	if cache == nil || !cache.enabled {
		return
	}

	data := cache.dataFor(detectorID)

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
		DetectorID: detectorID,
		Result:     result,
	})
}

func (cache *Cache) dataFor(detectorID int) map[*treepkg.Node][]entry {
	if cache.shared.detectorIDs.Has(detectorID) {
		return cache.shared.data
	} else {
		return cache.data
	}
}
