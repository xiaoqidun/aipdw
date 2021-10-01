package aipdw

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

type reqLBS struct {
	SK   string
	Args map[string]string
}

func (r *reqLBS) Encode() string {
	r.Signature()
	req := url.Values{}
	for k, v := range r.Args {
		req.Set(k, v)
	}
	return req.Encode()
}

func (r *reqLBS) Signature() {
	if r.Args == nil {
		r.Args = make(map[string]string)
	}
	var keys []string
	for k := range r.Args {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var keyValue []string
	for i := 0; i < len(keys); i++ {
		keyValue = append(keyValue, fmt.Sprintf("%s=%s", keys[i], r.Args[keys[i]]))
	}
	signStr := strings.Join(keyValue, "&") + r.SK
	signMd5 := md5.Sum([]byte(signStr))
	r.Args["sig"] = hex.EncodeToString(signMd5[:])
}
