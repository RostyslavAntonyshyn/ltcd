// Copyright (c) 2013 Conformal Systems LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package btcws_test

import (
	"github.com/conformal/btcjson"
	"github.com/conformal/btcws"
	"github.com/davecgh/go-spew/spew"
	"reflect"
	"testing"
)

var ntfntests = []struct {
	name   string
	f      func() btcjson.Cmd
	result btcjson.Cmd // after marshal and unmarshal
}{
	{
		name: "accountbalance",
		f: func() btcjson.Cmd {
			return btcws.NewAccountBalanceNtfn("abcde", 1.2345, true)
		},
		result: &btcws.AccountBalanceNtfn{
			Account:   "abcde",
			Balance:   1.2345,
			Confirmed: true,
		},
	},
	{
		name: "blockconnected",
		f: func() btcjson.Cmd {
			return btcws.NewBlockConnectedNtfn(
				"000000004811dda1c320ad5d0ea184a20a53acd92292c5f1cb926c3ee82abf70",
				153469)
		},
		result: &btcws.BlockConnectedNtfn{
			Hash:   "000000004811dda1c320ad5d0ea184a20a53acd92292c5f1cb926c3ee82abf70",
			Height: 153469,
		},
	},
	{
		name: "blockdisconnected",
		f: func() btcjson.Cmd {
			return btcws.NewBlockDisconnectedNtfn(
				"000000004811dda1c320ad5d0ea184a20a53acd92292c5f1cb926c3ee82abf70",
				153469)
		},
		result: &btcws.BlockDisconnectedNtfn{
			Hash:   "000000004811dda1c320ad5d0ea184a20a53acd92292c5f1cb926c3ee82abf70",
			Height: 153469,
		},
	},
	{
		name: "btcdconnected",
		f: func() btcjson.Cmd {
			return btcws.NewBtcdConnectedNtfn(true)
		},
		result: &btcws.BtcdConnectedNtfn{
			Connected: true,
		},
	},
	{
		name: "recvtx no block",
		f: func() btcjson.Cmd {
			return btcws.NewRecvTxNtfn("lalala the hex tx", nil)
		},
		result: &btcws.RecvTxNtfn{
			HexTx: "lalala the hex tx",
			Block: nil,
		},
	},
	{
		name: "recvtx with block",
		f: func() btcjson.Cmd {
			block := &btcws.BlockDetails{
				Height: 153469,
				Hash:   "000000004811dda1c320ad5d0ea184a20a53acd92292c5f1cb926c3ee82abf70",
				Index:  1,
				Time:   1386944019,
			}
			return btcws.NewRecvTxNtfn("lalala the hex tx", block)
		},
		result: &btcws.RecvTxNtfn{
			HexTx: "lalala the hex tx",
			Block: &btcws.BlockDetails{
				Height: 153469,
				Hash:   "000000004811dda1c320ad5d0ea184a20a53acd92292c5f1cb926c3ee82abf70",
				Index:  1,
				Time:   1386944019,
			},
		},
	},
	{
		name: "redeemingtx",
		f: func() btcjson.Cmd {
			return btcws.NewRedeemingTxNtfn("lalala the hex tx", nil)
		},
		result: &btcws.RedeemingTxNtfn{
			HexTx: "lalala the hex tx",
			Block: nil,
		},
	},
	{
		name: "redeemingtx with block",
		f: func() btcjson.Cmd {
			block := &btcws.BlockDetails{
				Height: 153469,
				Hash:   "000000004811dda1c320ad5d0ea184a20a53acd92292c5f1cb926c3ee82abf70",
				Index:  1,
				Time:   1386944019,
			}
			return btcws.NewRedeemingTxNtfn("lalala the hex tx", block)
		},
		result: &btcws.RedeemingTxNtfn{
			HexTx: "lalala the hex tx",
			Block: &btcws.BlockDetails{
				Height: 153469,
				Hash:   "000000004811dda1c320ad5d0ea184a20a53acd92292c5f1cb926c3ee82abf70",
				Index:  1,
				Time:   1386944019,
			},
		},
	},
	{
		name: "newtx",
		f: func() btcjson.Cmd {
			details := map[string]interface{}{
				"key1": float64(12345),
				"key2": true,
				"key3": "lalala",
				"key4": []interface{}{"abcde", float64(12345)},
			}
			return btcws.NewTxNtfn("abcde", details)
		},
		result: &btcws.TxNtfn{
			Account: "abcde",
			Details: map[string]interface{}{
				"key1": float64(12345),
				"key2": true,
				"key3": "lalala",
				"key4": []interface{}{"abcde", float64(12345)},
			},
		},
	},
	{
		name: "walletlockstate",
		f: func() btcjson.Cmd {
			return btcws.NewWalletLockStateNtfn("abcde", true)
		},
		result: &btcws.WalletLockStateNtfn{
			Account: "abcde",
			Locked:  true,
		},
	},
	{
		name: "alltx",
		f: func() btcjson.Cmd {
			return btcws.NewAllTxNtfn(
				"062f2b5f7d28c787e0f3aee382132241cd590efb7b83bd2c7f506de5aa4ef275",
				34567765)
		},
		result: &btcws.AllTxNtfn{
			TxID:   "062f2b5f7d28c787e0f3aee382132241cd590efb7b83bd2c7f506de5aa4ef275",
			Amount: 34567765,
		},
	},
	{
		name: "allverbosetx",
		f: func() btcjson.Cmd {
			return btcws.NewAllVerboseTxNtfn(&btcjson.TxRawResult{
				Hex:      "01000000010cdf900074a3622499a2f28f44a94476f27a8900a2bdd60e042754b6cab09741000000008a473044022012e11012fad1eb21ba1c82deb8da98778b08e714b72f281293064528343fae0502204294d7520f469f9673087a55395de0ce0e9074dce236db9fe7f30013b5fd00b90141047b6ff7832b4a763666e5481a0bd9eedb656d9f882d215c16fe9563d7b191cd67b2a41601a853a9f9d92773ae6f912ef451a089148e510623759cf55c408efdefffffffff02f4063f00000000001976a914b269e0ceec5d5b5e192cf580ae42341e0f79b0b588aca8c84b02000000001976a91439233c0d43a1411e547c60bad8985bae3530b6af88ac00000000",
				Txid:     "0cfeb968fb5d0f6b9a2a1de37c0607a1964dd3e335f203377cec90e03b20869e",
				Version:  0x1,
				LockTime: 0x0,
			})
		},
		result: &btcws.AllVerboseTxNtfn{
			RawTx: &btcjson.TxRawResult{
				Hex:      "01000000010cdf900074a3622499a2f28f44a94476f27a8900a2bdd60e042754b6cab09741000000008a473044022012e11012fad1eb21ba1c82deb8da98778b08e714b72f281293064528343fae0502204294d7520f469f9673087a55395de0ce0e9074dce236db9fe7f30013b5fd00b90141047b6ff7832b4a763666e5481a0bd9eedb656d9f882d215c16fe9563d7b191cd67b2a41601a853a9f9d92773ae6f912ef451a089148e510623759cf55c408efdefffffffff02f4063f00000000001976a914b269e0ceec5d5b5e192cf580ae42341e0f79b0b588aca8c84b02000000001976a91439233c0d43a1411e547c60bad8985bae3530b6af88ac00000000",
				Txid:     "0cfeb968fb5d0f6b9a2a1de37c0607a1964dd3e335f203377cec90e03b20869e",
				Version:  0x1,
				LockTime: 0x0,
			},
		},
	},
}

func TestNtfns(t *testing.T) {
	for _, test := range ntfntests {
		// create notification.
		n := test.f()

		// verify that id is nil.
		if n.Id() != nil {
			t.Errorf("%s: notification should not have non-nil id %v",
				test.name, n.Id())
			continue
		}

		mn, err := n.MarshalJSON()
		if err != nil {
			t.Errorf("%s: failed to marshal notification: %v",
				test.name, err)
			continue
		}

		n2, err := btcjson.ParseMarshaledCmd(mn)
		if err != nil {
			t.Errorf("%s: failed to ummarshal cmd: %v",
				test.name, err)
			continue
		}

		if !reflect.DeepEqual(test.result, n2) {
			t.Errorf("%s: unmarshal not as expected. "+
				"got %v wanted %v", test.name, spew.Sdump(n2),
				spew.Sdump(test.result))
		}
		if !reflect.DeepEqual(n, n2) {
			t.Errorf("%s: unmarshal not as we started with. "+
				"got %v wanted %v", test.name, spew.Sdump(n2),
				spew.Sdump(n))
		}

		// Read marshaled notification back into n.  Should still
		// match result.
		n.UnmarshalJSON(mn)
		if !reflect.DeepEqual(test.result, n) {
			t.Errorf("%s: unmarshal not as expected. "+
				"got %v wanted %v", test.name, spew.Sdump(n),
				spew.Sdump(test.result))
		}
	}
}
