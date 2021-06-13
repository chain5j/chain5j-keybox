// description: keybox
//
// @author: xwc1125
// @date: 2020/8/18 0018
package main

import (
	"fmt"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
	"github.com/chain5j/keybox"
	"github.com/chain5j/keybox/bip32"
	"github.com/chain5j/keybox/chain"
	"github.com/chain5j/keybox/chain/btc"
	"github.com/chain5j/keybox/chain/eth"
	"github.com/spf13/cobra"
	"os"
)

var (
	cmd = &cobra.Command{
		Use:   os.Args[0],
		Short: "WalletCtl",
	}
	// 导出主账户
	cmdOprMaster = &cobra.Command{
		Use:   "master",
		Short: "opt the master key",
		Run:   runOprMaster,
	}
	// 生成子账户
	cmdGenChild = &cobra.Command{
		Use:   "geneChild",
		Short: "generate the child account",
		Run:   runGenChild,
	}
	// 导出子账户
	cmdExportChild = &cobra.Command{
		Use:   "exportChild",
		Short: "export the child account",
		Run:   runExportChild,
	}
	// 使用子账户进行签名
	cmdSign = &cobra.Command{
		Use:   "sign",
		Short: "use the childAccount to sign",
		Run:   runSign,
	}
)

var (
	// 主账户部分
	path              string // 钱包路径
	password          string // 密码
	isSaveSubKey      bool   // 是否保存子私钥
	isSaveExtendedKey bool   // 是否保存扩展私钥
	isSaveMnemonic    bool   // 是否保存助记词
	mnemonicType      string // 助记词类型(zh-cn简体中文,zh-tw繁体中文,en English,fr French,it Italian,ja Japanese,ko Korean,es Spanish)
	mnemonic          string // 助记词
	isUsePwdBlur      bool   // 是否使用Password进行混淆
	prvKeyBase58      string // 私钥Base58
	networkType       string // 网络类型
	// 导出主账户信息
	exportMasterMn          bool // 导出主账户助记词
	exportMasterRawKey      bool // 导出主账户基本私钥
	exportMasterExtendedKey bool // 导出主账户扩展私钥
	// 子账户部分
	purposeType  uint32 // 生成类型（44,45）
	org          uint32 // purpose=45时，才使用
	coinType     uint32 // 币种类型
	account      uint32 // 用户空间
	addressIndex uint32 // 子账户的索引
	chainType    string // 链类型（eth,btc）
	childAddress string // 子账户地址
	childKeyPath string // 子账户路径
	// 子账户导出
	exportChildRawKey      bool   // 导出子账户的基础私钥
	exportChildExtendedKey bool   // 导出子账户的扩展私钥
	exportChildKeystore    bool   // 导出子账户的keystore
	childKeystorePwd       string // 子账户导出keystore的加密密码
	// 子账户签名
	signHash string // 交易体Hash
)

func init() {
	// 通用部分
	addFlags := func(cmd *cobra.Command, op string) {
		cmd.Flags().StringVarP(&path, "path", "f", "./wallet.dat", "the wallet file path")
		cmd.Flags().StringVarP(&password, "password", "p", "", "password to encrypt & decrypt wallet")
		cmd.Flags().BoolVar(&isSaveSubKey, "isSaveSubKey", false, "whether the wallet saves the subKey (the default is true) ")
		cmd.Flags().BoolVar(&isSaveExtendedKey, "isSaveExtendedKey", false, "whether the wallet saves the subExtendedKey (the default is false) ")
		cmd.Flags().BoolVar(&isSaveMnemonic, "isSaveMnemonic", true, "whether the wallet saves the mnemonic (the default is true) ")
		cmd.Flags().StringVar(&mnemonicType, "mnemonicType", "en",
			"the mnemonic type, the values is :en[English],zh-cn[简体中文],zh-tw[繁体中文],fr[French],it[Italian],ja[Japanese],ko[Korean],es[Spanish](the default is en) ")
		cmd.Flags().StringVarP(&mnemonic, "mnemonic", "m", "", "if load wallet by mnemonic,please write mnemonic. Words are separated by spaces")
		cmd.Flags().BoolVar(&isUsePwdBlur, "isUsePwdBlur", false, "whether use password to blur the seed.(the default is false)")
		cmd.Flags().StringVarP(&prvKeyBase58, "prvKeyBase58", "k", "", "if load wallet by prvKeyBase58,please write prvKeyBase58")
		cmd.Flags().StringVarP(&networkType, "networkType", "n", "mainnet", "network type,the values is: mainnet,testnet,devnet. (the default is mainnet)")
	}

	// 主账户导出
	{
		cmdOprMaster.Flags().BoolVar(&exportMasterMn, "exportMasterMn", false, "export the master account mnemonic")
		cmdOprMaster.Flags().BoolVar(&exportMasterRawKey, "exportMasterRawKey", false, "export the master account privateKey")
		cmdOprMaster.Flags().BoolVar(&exportMasterExtendedKey, "exportMasterExtendedKey", false, "export the master account base58PrivateKey")
		addFlags(cmdOprMaster, "master")
	}
	// 子账户生成
	{
		cmdGenChild.Flags().StringVarP(&chainType, "chainType", "t", "eth", "choose the chain type.The values is eth,btc(the default is eth)")
		cmdGenChild.Flags().Uint32Var(&purposeType, "purposeType", uint32(44), "choose the purpose type.The values is 44,45(the default is 44)")
		cmdGenChild.Flags().Uint32Var(&org, "org", 0, "if purpose=45,this value is org(the default is 0)")
		cmdGenChild.Flags().Uint32Var(&coinType, "coinType", 0, "coinType(the default is 0)")
		cmdGenChild.Flags().Uint32Var(&account, "account", 0, "the account space(the default is 0)")
		cmdGenChild.Flags().Uint32Var(&addressIndex, "addressIndex", 0, "the address index(the default is 0)")
		addFlags(cmdGenChild, "geneChild")
	}
	// 导出子账户
	{
		cmdExportChild.Flags().StringVarP(&chainType, "chainType", "t", "eth", "choose the chain type.The values is eth,btc(the default is eth)")
		cmdExportChild.Flags().StringVarP(&childAddress, "childAddress", "a", "", "the child address")
		cmdExportChild.Flags().StringVar(&childKeyPath, "childKeyPath", "", "the child account path")
		cmdExportChild.Flags().BoolVar(&exportChildRawKey, "exportChildRawKey", false, "whether export the child rawKey (the default is false) ")
		cmdExportChild.Flags().BoolVar(&exportChildExtendedKey, "exportChildExtendedKey", false, "whether export the child extendedKey (the default is false) ")
		cmdExportChild.Flags().BoolVar(&exportChildKeystore, "exportChildKeystore", false, "whether export the child keystore (the default is false) ")
		cmdExportChild.Flags().StringVar(&childKeystorePwd, "childKeystorePwd", "", "if export the keystore, will use childKeystorePwd to encrypt the privateKey")
		addFlags(cmdExportChild, "exportChild")
	}
	// 签名
	{
		cmdSign.Flags().StringVarP(&chainType, "chainType", "t", "eth", "choose the chain type.The values is eth,btc(the default is eth)")
		cmdSign.Flags().StringVarP(&childAddress, "childAddress", "a", "", "the child address")
		cmdSign.Flags().StringVar(&childKeyPath, "childKeyPath", "", "the child account path")
		cmdSign.Flags().StringVar(&signHash, "signHash", "", "the hash from transaction is need to sign")
		addFlags(cmdSign, "sign")
	}

	cmd.AddCommand(cmdOprMaster, cmdGenChild, cmdExportChild, cmdSign)
}

// 操作主账户
func runOprMaster(cmd *cobra.Command, args []string) {
	wallet, err := loadWallet()
	if err != nil {
		return
	}
	// 导出
	if exportMasterMn {
		masterMnemonic := wallet.ExportMasterMnemonic()
		fmt.Println("masterMnemonic: ", masterMnemonic)
	}
	if exportMasterRawKey {
		masterRawKey := wallet.ExportMasterRawKey()
		fmt.Println("masterRawKey: ", masterRawKey)
	}
	if exportMasterExtendedKey {
		masterExtendedKey := wallet.ExportMasterExtendedKey()
		fmt.Println("masterExtendedKey: ", masterExtendedKey)
	}
}

// 生成子账户
func runGenChild(cmd *cobra.Command, args []string) {
	wallet, err := loadWallet()
	if err != nil {
		return
	}
	chainApi := getChainApi()

	if purposeType != 44 && purposeType != 45 {
		fmt.Println("purpose type is err: ", "purpose type must 44 or 45")
		os.Exit(1)
	}
	subAddr, keyPath, err := wallet.CreateAccount(bip32.ParseHDNum(purposeType), bip32.ParseHDNum(coinType), bip32.ParseHDNum(org), bip32.ParseHDNum(account), uint32(0), addressIndex, chainApi)
	if err != nil {
		fmt.Println("create child account is err: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("subAddress: ", subAddr)
	fmt.Println("childPath: ", keyPath)
}

// 导出子账户内容
func runExportChild(cmd *cobra.Command, args []string) {
	wallet, err := loadWallet()
	if err != nil {
		return
	}
	chainApi := getChainApi()
	// 子私钥导出
	if exportChildRawKey {
		key, err := wallet.ExportRawKey(childAddress, childKeyPath, chainApi)
		if err != nil {
			fmt.Println("export child rawKey is err: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("rawKey: ", key)
	}
	if exportChildExtendedKey {
		key, err := wallet.ExportExtendedKey(childAddress, childKeyPath, chainApi)
		if err != nil {
			fmt.Println("export child extendedKey is err: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("extendedKey: ", key)
	}
	if exportChildKeystore {
		keystore, err := wallet.ExportKeyStore(childAddress, childKeyPath, childKeystorePwd, chainApi)
		if err != nil {
			fmt.Println("export child keystore is err: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("keystore: ", keystore)
	}
}

// 使用子账户进行签名
func runSign(cmd *cobra.Command, args []string) {
	wallet, err := loadWallet()
	if err != nil {
		return
	}
	// 签名
	signHashBytes, err := hexutil.Decode(signHash)
	if err != nil {
		fmt.Println("hex decode signHash is err: ", err.Error())
		os.Exit(1)
	}
	sign, err := wallet.Sign(childAddress, childKeyPath, signHashBytes, getChainApi())
	if err != nil {
		fmt.Println("sign signHash is err: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("signature: ", sign)
}

// 加载Wallet
func loadWallet() (*keybox.Wallet, error) {
	// 设置助记词类型
	if mnemonicType != "" {
		t := keybox.ParseMnemonicType(mnemonicType)
		keybox.SetBip39MnemonicType(t)
	}
	var (
		wallet *keybox.Wallet
		err    error
	)
	if mnemonic != "" {
		wallet, err = keybox.LoadWalletFromMnemonic(path, password, mnemonic, isUsePwdBlur)
	} else if prvKeyBase58 != "" {
		wallet, err = keybox.LoadWalletFromPrvKey(path, password, prvKeyBase58)
	} else {
		wallet, err = keybox.NewWallet(path, password)
	}
	if err != nil {
		fmt.Println("load or new wallet is err: ", err.Error())
		os.Exit(1)
		return nil, err
	}
	// 进行设置
	wallet.SetIsSaveSubKey(isSaveSubKey)
	wallet.SetIsSaveExtendedKey(isSaveExtendedKey)
	if !isSaveMnemonic {
		wallet.DelMnemonic()
	}
	return wallet, err
}

func getChainApi() keybox.ChainAPI {
	// 创建子账户
	var chainApi keybox.ChainAPI
	switch chainType {
	case "ETH":
		chainApi = eth.NewChain(chain.ParseToType(networkType))
	case "BTC":
		chainApi = btc.NewChain(chain.ParseToType(networkType))
	default:
		chainApi = eth.NewChain(chain.ParseToType(networkType))
	}
	return chainApi
}

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
