package rest

import (
	"testing"

	"github.com/mediacoin-pro/core/common/enc"
	"github.com/mediacoin-pro/core/crypto"
	"github.com/stretchr/testify/assert"
)

var testClient = NewClient("")

func TestClient_GetBlock(t *testing.T) {

	block, err := testClient.GetBlock(123)

	println(enc.IndentJSON(block))

	assert.NoError(t, err)
	assert.EqualValues(t, 123, block.Num)
}

func TestClient_GetBlocks(t *testing.T) {

	blocks, err := testClient.GetBlocks(122, 3)

	println(enc.IndentJSON(blocks))

	assert.NoError(t, err)
	assert.Equal(t, 3, len(blocks))
	assert.EqualValues(t, 123, blocks[0].Num)
	assert.EqualValues(t, 124, blocks[1].Num)
	assert.EqualValues(t, 125, blocks[2].Num)
}

func TestClient_AddressInfo(t *testing.T) {

	addr, memo, _ := crypto.DecodeAddress("MDC6ZKGnnz4g2y8eRoZKhaPjPbjsUGCUrUC") // alice address
	info, err := testClient.AddressInfo(addr, memo)

	println(enc.IndentJSON(info))

	assert.NoError(t, err)
	assert.Equal(t, memo, info.Memo)
	assert.Equal(t, addr, info.Address)
	assert.True(t, !info.Balance.IsZero())
}

//func TestClient_PutTx(t *testing.T) {
//
//	prvA := crypto.NewPrivateKeyBySecret("alice::*********")                // alice private key by password
//	addrB := crypto.MustParseAddress("MDCAr1nxktEAA8tun8KmF8KqMTvoXWQCvct") // bob address
//	amount := bignum.NewInt(22 * assets.MilliCoin)                          //
//	tx := txobj.NewSimpleTransfer(nil, prvA, assets.MDC, amount, 0, addrB, 0, "test-transfer")
//
//	err := testClient.PutTx(tx)
//
//	assert.NoError(t, err)
//}
