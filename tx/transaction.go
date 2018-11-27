package tx

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
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
	ScriptSig string
}

type TXOutput struct {
	Value        int
	ScriptPubkey string
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

func (txo *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return (txo.ScriptPubkey == unlockingData)
}

func (txi *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return (txi.ScriptSig == unlockingData)
}

func NewConinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{subsidy, to}
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
			inputs = append(inputs, TXInput{txID, out, from})
		}
	}

	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from})
	}

	tx := Transaction{[]byte{}, inputs, outputs}
	tx.setID()

	return &tx
}
