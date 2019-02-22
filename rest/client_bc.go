package rest

import (
	"fmt"
	"net/url"

	"github.com/mediacoin-pro/core/chain"
	"github.com/mediacoin-pro/core/crypto"
)

func (c *Client) GetBlock(num uint64) (block *chain.Block, err error) {
	err = c.httpGet(fmt.Sprintf("/block/%d", num), &block)
	return
}

func (c *Client) GetBlocks(offset uint64, limit int) (blocks []*chain.Block, err error) {
	err = c.httpGet("/blocks?"+url.Values{
		"offset": {fmt.Sprint(offset)},
		"limit":  {fmt.Sprint(limit)},
	}.Encode(), &blocks)
	return
}

func (c *Client) PutTx(tx *chain.Transaction) (err error) {
	return c.httpPut("/put-tx", tx)
}

func (c *Client) AddressInfo(addr []byte, memo uint64) (info *chain.AddressInfo, err error) {
	err = c.httpGet("/address/"+crypto.EncodeAddress(addr, memo), &info)
	return
}
