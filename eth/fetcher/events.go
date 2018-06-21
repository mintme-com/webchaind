package fetcher

import "github.com/webchain-network/webchaind/core/types"

type FetcherInsertBlockEvent struct {
	Peer  string
	Block *types.Block
}
