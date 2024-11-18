# 钱包生成工具

[TOC]

## 项目说明

- 通用参数：

| 参数                  | 说明                                           |
|---------------------|----------------------------------------------|
| -f                  | --path,钱包文件路径                                |
| -p                  | --password,钱包文件加解密密码                         |
| --isSaveSubKey      | 是否保存子账户（默认false）                             |
| --isSaveExtendedKey | 是否保存子账户的扩展私钥（默认false）                        |
| --isSaveMnemonic    | 是否保存主账户的助记词（默认true）                          |
| --mnemonicType      | 助记词类型，类型有en,zh-cn,zh-tw,fr,it,ja,ko,es（默认en） |
| --mnemonic          | 助记词（用于恢复钱包）                                  |
| --isUsePwdBlur      | 是否使用Password进行混淆（默认false）                    |
| --prvKeyBase58      | 扩展私钥（用于恢复钱包）                                 |
| --networkType       | 网络类型，类型有：mainnet,testnet,devnet（默认mainnet）   |

- childPath结构说明

/44/60/0/0/0：其含义为：/BIP44格式/币种类型/组织/地址类型（默认0）/地址索引

### 主账户导出

参数说明：

| 参数                        | 说明                         |
|---------------------------|----------------------------|
| --exportMasterMn          | 是否导出助记词（默认false）           |
| --exportMasterRawKey      | 是否导出16进制私钥（默认false）        |
| --exportMasterExtendedKey | 是否导出扩展私钥，base58进制（默认false） |

- 示例：

```shell script
## 助记词恢复
./walletctl master -f "./wallet1.dat" -p "123456" --mnemonic "security traffic pluck dawn enlist above bunker worth pencil ten garage ribbon"
## 导出扩展私钥
./walletctl master -f "./wallet1.dat" -p "123456" --exportMasterExtendedKey
## 导出助记词
./walletctl master -f "./wallet1.dat" -p "123456" --exportMasterMn
```

### 子账户生成

参数说明：

| 参数             | 说明                                |
|----------------|-----------------------------------|
| -t             | --chainType,链类型，包含有eth、btc（默认eth） |
| --purposeType  | purpose 类型，包含44，45（默认44）          |
| --org          | 当purpose=45时，才被使用（默认0）            |
| --coinType     | 币种类型（默认0）                         |
| --account      | account账户空间（默认0）                  |
| --addressIndex | 地址索引（默认0）                         |

- 示例：

```shell script
## 创建子账户
./walletctl geneChild -f "./wallet1.dat" -p "123456" --chainType "eth" --addressIndex 1

./walletctl geneChild -f "./wallet1.dat" -p "123456" --chainType "eth" --addressIndex 0
```

### 导出子账户

参数说明：

| 参数                       | 说明                                |
|--------------------------|-----------------------------------|
| -t                       | --chainType,链类型，包含有eth、btc（默认eth） |
| --childAddress           | 子账户地址                             |
| --childKeyPath           | 子账户的路径                            |
| --exportChildRawKey      | 是否导出16进制私钥(默认false)               |
| --exportChildExtendedKey | 是否导出扩展私钥(默认false)                 |
| --exportChildKeystore    | 是否导出keystore(默认false)             |
| --childKeystorePwd       | 导出keystore时的加密密码                  |

- 示例：

```shell script
## 导出扩展私钥
./walletctl exportChild -f "./wallet1.dat" -p "123456" --chainType "eth" --childAddress "0xb3d988aFDe88653dc1e2C48f770d7DC5AE93547C" --childKeyPath "/44/0/0/0/0" --exportChildExtendedKey
```

### 签名

参数说明：

| 参数             | 说明                                |
|----------------|-----------------------------------|
| -t             | --chainType,链类型，包含有eth、btc（默认eth） |
| --childAddress | 子账户地址                             |
| --childKeyPath | 子账户的路径                            |
| --signHash     | 交易体Hash                           |

- 示例：

```shell script
## 签名
./walletctl sign -f "./wallet1.dat" -p "123456" --chainType "eth" --childAddress "0xAeff996F0Efb374fCf95Eb6b38fd4aA5E4bbC1b1" --childKeyPath "/44/60/0/0/0" --signHash "0x123456"
```
