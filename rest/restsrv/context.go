package restsrv

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/mediacoin-pro/core/chain"
	"github.com/mediacoin-pro/core/chain/assets"
	"github.com/mediacoin-pro/core/chain/txobj"
	"github.com/mediacoin-pro/core/common/bignum"
	"github.com/mediacoin-pro/core/common/bin"
	"github.com/mediacoin-pro/core/common/xlog"
	"github.com/mediacoin-pro/core/crypto"
)

type Context struct {
	*Server
	req      *http.Request
	reqQuery url.Values
	reqBody  *bin.Reader
	rw       http.ResponseWriter
	uriPath  string
	uriParts []string
}

func newContext(
	srv *Server,
	req *http.Request,
	rw http.ResponseWriter,
) *Context {
	path := req.URL.Path
	path = strings.TrimPrefix(path, "/rest")
	if path != "/" {
		path = strings.TrimSuffix(path, "/")
	}
	c := &Context{
		Server:   srv,
		req:      req,
		uriPath:  path,
		reqQuery: req.URL.Query(),
		reqBody:  bin.NewReader(req.Body),
		rw:       rw,
	}
	if (req.Method == "POST" || req.Method == "PUT") && req.ParseForm() == nil {
		c.reqQuery = req.Form
	}
	return c
}

const (
	contentTypeBinary = "binary"
	contentTypeJSON   = "application/json; charset=utf-8"
)

var (
	rePathBlockNum    = regexp.MustCompile(`^/block/(\d+)$`)
	rePathAddressInfo = regexp.MustCompile(`^/address/(@[a-zA-Z0-9\-_]+|MDC[a-zA-Z1-9]+|0x[a-f0-9]+)$`)
	reTxHash          = regexp.MustCompile(`^/tx/([a-f0-9]{64})$`)
	reTxID            = regexp.MustCompile(`^/tx/([a-f0-9]{1,16})$`)

	err404        = errors.New("404 - Not found")
	errUserExists = errors.New("400 - User exists")
)

func (c *Context) Exec() {

	switch {

	case c.uriPath == "/info":
		c.WriteVar(c.bc.Info())

		//	/block/<block-num>
	case c.matchPath(rePathBlockNum):
		num, _ := strconv.ParseUint(c.uriParts[1], 10, 64)
		c.WriteVar(c.bc.GetBlock(num))

		//	/blocks?offset=<block-num>&limit=<count-blocks>
	case c.uriPath == "/blocks":
		offset := c.getUint("offset")
		limit := c.getLimit()
		orderDesc := c.getOrderDesc()
		c.WriteVar(c.bc.GetBlocks(offset, limit, orderDesc))

		//	/tx/<hash:hex>
	case c.matchPath(reTxHash):
		txHash, _ := hex.DecodeString(c.uriParts[1])
		c.WriteVar(c.bc.TransactionByHash(txHash))

		//	/tx/<txID:hex>
	case c.matchPath(reTxID):
		txID, _ := strconv.ParseUint(c.uriParts[1], 16, 64)
		c.WriteVar(c.bc.TransactionByID(txID))

		//	/address/?address=MDC&memo=...
	case c.uriPath == "/address":
		addr, memo := c.getAddress("")
		c.WriteVar(c.bc.AddressInfo(addr, memo, assets.MDC))

		//	/address/MDCxxxxxxxxxxxxx
	case c.matchPath(rePathAddressInfo):
		addr, memo := c.getAddress(c.uriParts[1])
		c.WriteVar(c.bc.AddressInfo(addr, memo, assets.MDC))

	case c.uriPath == "/txs":
		addr, memo := c.getAddress("")
		offset := c.getUint("offset")
		limit := c.getLimit()
		orderDesc := c.getOrderDesc()
		txs, ofst, err := c.bc.TransactionsByAddr(assets.MDC, addr, memo, offset, limit, orderDesc)
		c.WriteVar(NewResponse(txs, ofst, err))

	case c.uriPath == "/put-tx":
		var tx *chain.Transaction
		c.getBinary(&tx)
		err := c.bc.Mempool.Put(tx)
		c.WriteVar(0, err)

	case c.uriPath == "/new-transfer":
		prvKey := c.getPrivateKey()        // private key OR seed
		toAddr, toMemo := c.getAddress("") // address
		amount := c.getAmount("amount")    // amount
		comment := c.getStr("comment", "") // comment
		asset := assets.MDC                //

		tx := txobj.NewSimpleTransfer(c.bc, prvKey, asset, amount, 0, toAddr, toMemo, comment, 0)
		c.assert(tx.Verify(c.bc.Cfg))

		err := c.bc.Mempool.Put(tx)
		c.WriteVar(tx, err)

	case c.uriPath == "/new-user":
		prv := c.getPrivateKey()          // private key OR seed
		nick := c.getStr("login", "")     // user nickname
		referrerID := c.getUint("ref_id") // referral id

		user, err := c.bc.UserByNick(nick)
		c.assert(err)
		if user != nil {
			if user.PublicKey().Equal(prv.PublicKey()) {
				c.WriteVar(user.Tx(), err)
			} else {
				c.WriteError(errUserExists, http.StatusBadRequest)
			}
			return
		}
		tx := txobj.NewUser(c.bc, prv, nick, referrerID)
		c.assert(tx.Verify(c.bc.Cfg))

		err = c.bc.Mempool.Put(tx)
		c.WriteVar(tx, err)

	case c.uriPath == "/new-key":
		prv := c.getPrivateKey() // private key OR seed
		c.WriteVar(struct {
			PrvKey  string `json:"private_key"`
			PubKey  string `json:"public_key"`
			Address string `json:"address"`
			UserID  string `json:"user_id"`
		}{
			prv.String(),
			prv.PublicKey().String(),
			prv.PublicKey().StrAddress(),
			"0x" + prv.PublicKey().HexID(),
		})

	default:
		c.WriteError(err404, http.StatusNotFound)
	}

	return
}

//----------------------- request --------------------------------------
func (c *Context) matchPath(re *regexp.Regexp) bool {
	c.uriParts = re.FindStringSubmatch(c.uriPath)
	return len(c.uriParts) > 0
}

func (c *Context) assert(err error) {
	if err != nil {
		c.WriteError(err, 400)
		panic(err)
	}
}

func (c *Context) exists(name string) bool {
	_, ok := c.reqQuery[name]
	return ok
}

func (c *Context) getLimit() (limit int64) {
	limit = c.getInt("limit")
	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}
	return
}

func (c *Context) getAddress(defaultValue string) (addr []byte, memo uint64) {
	addr, memo, err := c.bc.AddressByStr(c.getStr("address", defaultValue))
	c.assert(err)
	if s := c.getStr("memo", ""); s != "" {
		memo, err = strconv.ParseUint(s, 0, 64)
		c.assert(err)
	}
	return
}

func (c *Context) getPrivateKey() *crypto.PrivateKey {
	if seed := c.getStr("seed", ""); seed != "" {
		return crypto.NewPrivateKeyBySecret(seed)
	}
	if login := c.getStr("login", ""); login != "" {
		return crypto.NewPrivateKeyBySecret(login + "::" + c.getStr("password", ""))
	}
	prv, err := crypto.ParsePrivateKey(c.getStr("private", ""))
	c.assert(err)
	return prv
}

func (c *Context) getOrderDesc() bool {
	return c.getStr("order", "asc") == "desc"
}

func (c *Context) getAmount(name string) (n bignum.Int) {
	v, err := strconv.ParseUint(c.getStr(name, "0"), 10, 64)
	c.assert(err)
	return bignum.NewInt(int64(v))
}

func (c *Context) getStr(name, defaultValue string) string {
	if v := c.reqQuery.Get(name); v != "" {
		return v
	}
	return defaultValue
}

func (c *Context) getInt(name string) int64 {
	n, err := strconv.ParseInt(c.getStr(name, "0"), 0, 64)
	c.assert(err)
	return n
}

func (c *Context) getUint(name string) uint64 {
	n, err := strconv.ParseUint(c.getStr(name, "0"), 0, 64)
	c.assert(err)
	return n
}

func (c *Context) getBinary(v ...interface{}) {
	err := c.reqBody.ReadVar(v...)
	c.assert(err)
}

//----------------------- response -------------------------------------
func (c *Context) WriteError(err error, httpCode int) {
	xlog.Error.Printf("rest> Response-ERROR-%d: %v", httpCode, err)

	var buf io.Reader
	if c.req.Header.Get("Accept") == contentTypeBinary {
		c.rw.Header().Set("Content-Type", contentTypeBinary)
		buf = bytes.NewBufferString(err.Error())
	} else {
		c.rw.Header().Set("Content-Type", contentTypeJSON)
		data, _ := json.Marshal(&Response{Error: err.Error()})
		buf = bytes.NewBuffer(data)
	}
	c.rw.Header().Set("X-Content-Type-Options", "nosniff")
	c.rw.WriteHeader(httpCode)
	io.Copy(c.rw, buf)
}

func (c *Context) WriteVar(v interface{}, ee ...error) {
	if len(ee) > 0 && ee[0] != nil { // error
		c.WriteError(ee[0], 500)
		return
	}
	var buf io.Reader
	if c.req.Header.Get("Accept") == contentTypeBinary {
		// binary-response
		c.rw.Header().Set("Content-Type", contentTypeBinary)
		if r, ok := v.(*Response); ok {
			v = r.Results
			c.rw.Header().Set("X-Next-Offset", r.NextOffset)
		}
		buf = bin.NewBuffer(nil, v)

	} else {
		// json-response
		c.rw.Header().Set("Content-Type", contentTypeJSON)
		var data []byte
		if _, ok := c.reqQuery["pretty"]; ok {
			data, _ = json.MarshalIndent(v, "", "  ")
		} else {
			data, _ = json.Marshal(v)
		}
		buf = bytes.NewBuffer(data)
	}
	_, err := io.Copy(c.rw, buf)
	if err != nil {
		xlog.Error.Printf("rest> http-response-error: %v", err)
	}
}
