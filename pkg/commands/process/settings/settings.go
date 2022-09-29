package settings

var DefaultSettings = TypeSettings{
	MaximumMemoryMb:      uint64(2046),
	MemoryCheckEachFiles: 100,
}

type TypeSettings struct {
	MaximumMemoryMb      uint64
	MemoryCheckEachFiles int
}

func Default() TypeSettings {
	return DefaultSettings
}
