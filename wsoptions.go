package dgo2poc

type WSOpts struct {
	id wsIdentify
}

// Options for WSClient.
type WSOpt func(opts *WSOpts)

// Connect as a shard. By running multiple shards, you can split your bot across multiple processes.
func WithShards(num, of int) WSOpt {
	return WSOpt(func(opts *WSOpts) {
		opts.id.Shard = [2]int{num, of}
	})
}

// Set the threshold before a guild is considered "large" and will not have offline members
// returned. Must be in the range of 50-250.
func WithLargeThreshold(num int) WSOpt {
	return WSOpt(func(opts *WSOpts) {
		opts.id.LargeThreshold = num
	})
}

// Set the initial status. By default, it will be "online" with no game playing.
func WithStatus(s WSStatus) WSOpt {
	return WSOpt(func(opts *WSOpts) {
		opts.id.Presence = s
	})
}
