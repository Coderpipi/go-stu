package distributed_cache_stu

type (
	PeerPicker interface {
		PickPeer(key string) (PeerGetter, bool)
	}

	PeerGetter interface {
		Get(group, key string) ([]byte, error)
	}
)
