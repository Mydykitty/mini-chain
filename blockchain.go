package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

const blocksBucket = "blocks"

type Blockchain struct {
	Tip []byte
	DB  *bolt.DB
}

// 创建区块链
func CreateBlockchain(address string) *Blockchain {
	db, err := bolt.Open("blockchain.db", 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	var tip []byte

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			genesis := NewGenesisBlock(NewCoinbaseTX(address, "Genesis Block"))
			b, _ = tx.CreateBucket([]byte(blocksBucket))
			b.Put(genesis.Hash, genesis.Serialize())
			b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{tip, db}
}

// 添加新区块
func (bc *Blockchain) AddBlock(transactions []*Transaction) {
	var lastHash []byte

	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash)

	err = bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		b.Put(newBlock.Hash, newBlock.Serialize())
		b.Put([]byte("l"), newBlock.Hash)
		bc.Tip = newBlock.Hash
		return nil
	})
}

// 区块链迭代器
type BlockchainIterator struct {
	CurrentHash []byte
	DB          *bolt.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.Tip, bc.DB}
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block
	err := i.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.CurrentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	i.CurrentHash = block.PrevHash
	return block
}

// 查找UTXO
func (bc *Blockchain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	accumulated := 0
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
		Outputs:
			for outIdx, out := range tx.Vout {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
					accumulated += out.Value
					unspentOuts[txID] = append(unspentOuts[txID], outIdx)
					if accumulated >= amount {
						break
					}
				}
			}

			if !tx.IsCoinbase() {
				for _, in := range tx.Vin {
					if bytes.Equal(HashPubKey(in.PubKey), pubKeyHash) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.OutIndex)
					}
				}
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}

	return accumulated, unspentOuts
}

// 根据交易ID查找交易
func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Equal(tx.ID, ID) {
				return *tx, nil
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return Transaction{}, fmt.Errorf("交易 %x 未找到", ID)
}

func (bc *Blockchain) PrintBlockchain() {
	bci := bc.Iterator()
	fmt.Println("=== 开始遍历区块链 ===")

	for {
		block := bci.Next()
		fmt.Printf("\n--- 区块 ---\n")
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("PrevHash: %x\n", block.PrevHash)
		fmt.Printf("时间戳: %d\n", block.Timestamp)

		for i, tx := range block.Transactions {
			fmt.Printf("  交易 %d:\n", i)
			fmt.Printf("    ID: %x\n", tx.ID)
			if tx.IsCoinbase() {
				fmt.Println("    Coinbase交易（挖矿奖励）")
			}

			for j, in := range tx.Vin {
				fmt.Printf("    输入 %d:\n", j)
				fmt.Printf("      TxID: %x\n", in.Txid)
				fmt.Printf("      OutIndex: %d\n", in.OutIndex)
				fmt.Printf("      PubKey: %x\n", in.PubKey)
				fmt.Printf("      Signature: %x\n", in.Signature)
			}

			for j, out := range tx.Vout {
				fmt.Printf("    输出 %d:\n", j)
				fmt.Printf("      金额: %d\n", out.Value)
				fmt.Printf("      公钥哈希: %x\n", out.PubKeyHash)
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	fmt.Println("=== 遍历结束 ===")
}
