package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	bitcoinWallet "github.com/Amirilidan78/bitcoin-wallet"
	"github.com/Amirilidan78/bitcoin-wallet/enums"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"log"
)

var privateKeyHex = "88414dbb373a211bc157265a267f3de6a4cec210f3a5da12e89630f2c447ad27"
var toAddressHex = "tb1q9dkhf8vxlvujxjmnslxsv97nseg9pjmxqsku3v"
var chain = &chaincfg.TestNet3Params

func createWallet() *bitcoinWallet.BitcoinWallet {
	w, _ := bitcoinWallet.CreateBitcoinWallet(enums.TEST_NODE, privateKeyHex)
	return w
}

func main() {

	tx()

}

func tx() {

	amount := int64(5000)
	fee := int64(1000)

	wallet := createWallet()

	priv, _ := wallet.PrivateKeyBTCE()

	tx, err := createTransaction(chain, priv, wallet.Address, toAddressHex, amount, fee)

	fmt.Println(tx, err)

}

func test() {

	totalAmount := int64(1200000)
	amount := int64(100000)
	fee := int64(10000)

	privKey := "cS9Zef6XdN3jHTFJFSsyJAtmDgCCdnygyVUJsLoyB8neuwhidUNJ"
	spendAddrStr := "tb1qppv790u4dz48ctnk3p7ss7fmspckagp3wrfyp0"
	destAddrStr := "n3xmiD8kMU42seVZEc5icPKxTrqATm5t6x"
	chain := &chaincfg.TestNet3Params
	txHash := "4de74b4af672742f331eb3e2712ab54eca9a49d7e34d132b5dd39b98186057d7"
	position := 0

	spendAddr, err := btcutil.DecodeAddress(spendAddrStr, chain)
	if err != nil {
		log.Println("DecodeAddress spendAddr err", err)
		return
	}

	destAddr, err := btcutil.DecodeAddress(destAddrStr, chain)
	if err != nil {
		log.Println("DecodeAddress destAddrStr err", err)
		return
	}

	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		log.Println("wif err", err)
		return
	}

	spenderAddrByte, err := txscript.PayToAddrScript(spendAddr)
	if err != nil {
		log.Println("spendAddr PayToAddrScript err", err)
		return
	}

	destAddrByte, err := txscript.PayToAddrScript(destAddr)
	if err != nil {
		log.Println("destAddr PayToAddrScript err", err)
		return
	}

	// == //

	utxoHash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		log.Println("NewHashFromStr err", err)
		return
	}

	outPoint := wire.NewOutPoint(utxoHash, uint32(position))
	redeemTx := wire.NewMsgTx(2)

	txIn := wire.NewTxIn(outPoint, nil, [][]byte{})
	txIn.Sequence = txIn.Sequence - 2
	redeemTx.AddTxIn(txIn)

	redeemTxOut0 := wire.NewTxOut(amount, destAddrByte)
	redeemTxOut1 := wire.NewTxOut(totalAmount-amount-fee, spenderAddrByte)

	redeemTx.AddTxOut(redeemTxOut0)
	redeemTx.AddTxOut(redeemTxOut1)
	redeemTx.LockTime = 2407372

	if err != nil {
		log.Println("DecodeString pkScript err", err)
		return
	}

	sigHashes := txscript.NewTxSigHashes(redeemTx, txscript.NewMultiPrevOutFetcher(map[wire.OutPoint]*wire.TxOut{
		*outPoint: {},
	}))

	signature, err := txscript.WitnessSignature(redeemTx, sigHashes, 0, totalAmount, spenderAddrByte, txscript.SigHashAll, wif.PrivKey, true)
	if err != nil {
		log.Println("WitnessSignature err", err)
		return
	}
	redeemTx.TxIn[0].Witness = signature

	var signedTx bytes.Buffer
	redeemTx.Serialize(&signedTx)

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	log.Println(hexSignedTx)

}
