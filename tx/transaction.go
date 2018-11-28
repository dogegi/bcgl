package tx

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"

	"../utils"
	"../wallet"
)

const subsidy = 10

type IBlockChain interface {
	FindSpendableOutputs(address string, amount int) (int, map[string][]int)
}

type Transaction struct {
	ID    []byte
	Vins  []TXInput
	Vouts []TXOutput
}

type TXInput struct {
	Txid      []byte
	VoutIdx   int
	Signature []byte
	PubKey    []byte
}

type TXOutput struct {
	Value      int
	PubKeyHash []byte
}

func (tx *Transaction) setID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vins) == 1 && len(tx.Vins[0].Txid) == 0 && tx.Vins[0].VoutIdx == -1
}

func (txi *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	//return (txi.ScriptSig == unlockingData)
	return true
}

func (txi *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.HashPubKey(txi.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

func (txo *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	//return (txo.ScriptPubkey == unlockingData)
	return true
}

func (txo *TXOutput) Lock(address []byte) {
	pubKeyHash := utils.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	txo.PubKeyHash = pubKeyHash
}

func (txo *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(txo.PubKeyHash, pubKeyHash) == 0
}

func NewConinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TXInput{[]byte{}, -1, []byte{}, []byte{}}
	txout := TXOutput{subsidy, []byte{}}
	tx := Transaction{[]byte{}, []TXInput{txin}, []TXOutput{txout}}
	tx.setID()

	return &tx
}

func NewUTXOTransaction(from, to string, amount int, bc IBlockChain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	acc, vaildOutputs := bc.FindSpendableOutputs(from, amount)
	if acc < amount {
		log.Panic("ERRPR: NOt enough funds")
	}

	for txid, outs := range vaildOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			inputs = append(inputs, TXInput{txID, out, []byte{}, []byte{}})
		}
	}

	outputs = append(outputs, TXOutput{amount, []byte{}})
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, []byte{}})
	}

	tx := Transaction{[]byte{}, inputs, outputs}
	tx.setID()

	return &tx
}
