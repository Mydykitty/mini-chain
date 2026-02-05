# mini-chain

一个基于 Go 实现的简易区块链节点，支持 **工作量证明 (PoW)**、**UTXO交易系统**、**P2P节点通信** 和 **数字签名验证**。  
项目旨在学习和演示区块链核心机制，包含从交易创建到区块挖矿、节点同步的完整流程。

---

## 📦 功能特性

### 核心功能
- **区块链存储**
    - 链上每个区块包含前置哈希、交易列表、时间戳、Nonce。
    - 支持 PoW 挖矿和区块验证。
- **交易系统**
    - 支持 **Coinbase** 奖励交易。
    - 支持 **UTXO 转账**交易，自动找零。
    - 使用 **ECDSA** 数字签名保证交易安全。
- **网络通信**
    - TCP P2P 网络节点，可广播交易和区块。
    - 简单命令协议：`version`, `inv`, `getdata`, `block`, `tx`。
- **钱包管理**
    - 创建和管理多个钱包。
    - 公私钥生成和地址管理。
    - 钱包数据保存至 `wallets.json`。

### 可选功能/优化方向
- 动态难度调整，模拟真实区块链挖矿。
- Merkle 树交易哈希，提升交易完整性校验。
- 节点发现与连接池，提高 P2P 网络健壮性。
- 钱包文件加密，保护私钥安全。
- CLI 操作：查询余额、创建交易、挖矿、同步区块链。

---

## ⚙️ 技术细节

- **语言**：Go
- **哈希算法**：SHA-256 + RIPEMD160
- **签名算法**：ECDSA (P256)
- **数据序列化**：Gob
- **网络协议**：TCP

### 数据结构
- `Block`：包含 `PrevHash`, `Transactions`, `Timestamp`, `Nonce`, `Hash`
- `Transaction`：包含 `ID`, `Vin`, `Vout`，支持 Coinbase 与 UTXO 转账
- `TXInput/TXOutput`：输入输出结构，支持锁定/解锁逻辑
- `Wallet`：包含私钥、公共地址
- `Blockchain`：区块链操作集合，包括查找 UTXO、添加区块

---

## 🚀 快速开始

### 1. 克隆项目
```bash
$ git clone https://github.com/你的用户名/mini-chain.git
$ cd mini-chain
```

### 2. 编译和运行
```bash
$ go run . 3000

$ go run . 3001

$ curl -X POST http://localhost:43000/tx \
-H "Content-Type: application/json" \
-d '{
"from": "miner",
"to": "alice",
"amount": 10
}'

$ go run . 3000 getchain

$ go run . 3001 getchain
```