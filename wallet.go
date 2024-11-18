package keybox

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chain5j/keybox/bip32"
	"github.com/chain5j/keybox/bip39"
	"github.com/chain5j/keybox/bip39/wordlists"
	"github.com/chain5j/keybox/bip44"
	"github.com/chain5j/keybox/crypto/scrypt"
	"github.com/chain5j/keybox/util/dateutil"
	log "github.com/chain5j/log15"
	"github.com/pborman/uuid"
)

type ChildKeyPropertyInfo struct {
	Purpose       uint32 `json:"purpose"`
	ChainType     uint32 `json:"chainType"`
	AlgorithmType uint32 `json:"algorithmType"`
	Org           uint32 `json:"org"`
	CoinType      uint32 `json:"coinType"`
	Time          uint32 `json:"time"`
	Key           []byte `json:"key"`
}

type ExtendedKey struct {
	Key       []byte `json:"key"`        // 33 bytes
	ChainCode []byte `json:"chain_code"` // 32 bytes
}

var isLog bool // 是否打印日志

// Wallet 管理钱包文件
type Wallet struct {
	mu             sync.RWMutex
	Path           string                           `json:"path"`
	Mnemonic       string                           `json:"mnemonic"`
	Password       string                           `json:"password"`
	Key            *bip32.Key                       `json:"key"`
	Time           uint32                           `json:"time"`
	AddrLinkPubkey map[string]string                `json:"addrLinkPubkey"` // 地址和公钥的配置
	ChildKeyInfo   map[string]*ChildKeyPropertyInfo `json:"childKeyInfo"`   // 公钥对应的子私钥内容

	IsSaveSubKey      bool `json:"isSaveSubKey"`      // 是否保存子私钥
	IsSaveExtendedKey bool `json:"isSaveExtendedKey"` // 是否保存扩展私钥
}

// 将wallet进行scrypt加密，并写入文件中
func writeContentToWalletFile(wallet *Wallet, path string, password string) error {
	if nil == wallet || len(path) == 0 {
		return fmt.Errorf("writeContentToWalletFile parameter error")
	}
	data, err := json.Marshal(wallet)
	if err != nil {
		return fmt.Errorf("NewWallet json.Marshal err:%v", err.Error())
	}
	if len(password) == 0 {
		err = ioutil.WriteFile(path, data, 0644)
		if err != nil {
			return fmt.Errorf("NewWallet ioutil.WriteFile err:%v", err.Error())
		}
		return nil
	}

	startTime := getLogCurrentTime()
	// 使用scrypt加密
	dataEnc, err := scrypt.EncryptKey(&scrypt.Key{Id: uuid.NewRandom(), PrivateKey: data}, password, scrypt.StandardScryptN, scrypt.StandardScryptP)
	printMsg("scrypt.EncryptKey", startTime)
	if err != nil {
		return fmt.Errorf("NewWallet aes.AesEncrypt err:%v", err.Error())
	}
	// 转换为base64
	dataBase64 := base64.StdEncoding.EncodeToString(dataEnc)
	wallet.mu.Lock()
	defer wallet.mu.Unlock()
	err = ioutil.WriteFile(path, []byte(dataBase64), 0644)
	if err != nil {
		return fmt.Errorf("NewWallet ioutil.WriteFile err:%v", err.Error())
	}
	return nil
}

// 读取文件内容，并使用scrypt解密
func readContentToWalletFile(path string, password string) (wallet *Wallet, err error) {
	if len(path) == 0 || len(password) == 0 {
		return nil, fmt.Errorf("readContentToWalletFile parameter error")
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("readContentToWalletFile ioutil.ReadFile err:%v", err.Error())
	}
	walletData, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, fmt.Errorf("readContentToWalletFile base64.StdEncoding.DecodeString err:%v", err.Error())
	}
	w := new(Wallet)
	if len(password) == 0 {
		err = json.Unmarshal(walletData, w)
		if err != nil {
			return nil, fmt.Errorf("readContentToWalletFile json.Unmarshal err:%v", err.Error())
		}
		return w, nil
	}
	// 使用scrypt解密
	startTime := getLogCurrentTime()
	key, err := scrypt.DecryptKey(walletData, password)
	printMsg("scrypt.DecryptKey", startTime)
	if err != nil {
		return nil, fmt.Errorf("readContentToWalletFile aes.AesDecrypt err:%v", err.Error())
	}
	dataDec := key.PrivateKey
	err = json.Unmarshal(dataDec, w)
	if err != nil {
		return nil, fmt.Errorf("readContentToWalletFile json.Unmarshal err:%v", err.Error())
	}
	return w, nil
}

// ==========================主账户============================
// 设置助记词的类型【默认是使用English助记词，如果更换，需要在最前面初始化】
func SetBip39MnemonicType(mnemonicType MnemonicType) {
	var wordList []string
	switch mnemonicType {
	case MnemonicType_Chinese_Simplified:
		wordList = wordlists.ChineseSimplified
	case MnemonicType_Chinese_Traditional:
		wordList = wordlists.ChineseTraditional
	case MnemonicType_English:
		wordList = wordlists.English
	case MnemonicType_French:
		wordList = wordlists.French
	case MnemonicType_Italian:
		wordList = wordlists.Italian
	case MnemonicType_Japanese:
		wordList = wordlists.Japanese
	case MnemonicType_Korean:
		wordList = wordlists.Korean
	case MnemonicType_Spanish:
		wordList = wordlists.Spanish
	default:
		wordList = wordlists.English
	}
	bip39.SetWordList(wordList)
}

func newWallet() *Wallet {
	return &Wallet{
		IsSaveSubKey:      false,
		IsSaveExtendedKey: false,
	}
}

// NewWallet 创建钱包文件实例
func NewWallet(path string, password string) (*Wallet, error) {
	startTime := getLogCurrentTime()
	// 参数检查
	if len(path) == 0 {
		return nil, fmt.Errorf("NewWallet path parameter error")
	}
	wallet := newWallet()
	// 判断钱包文件是否存在
	_, err := ioutil.ReadFile(path)
	printMsg("ioutil.ReadFile", startTime)
	if err == nil {
		return readWalletFromFile(path, password, wallet, err)
	}
	// 不存在的话则创建一个钱包文件
	// 使用bip39处理，生成助记词
	startTime = getLogCurrentTime()
	entropy, err := bip39.NewEntropy(128)
	printMsg("bip39.NewEntropy", startTime)
	if err != nil {
		return nil, fmt.Errorf("NewWallet bip39.NewEntropy err:%v", err.Error())
	}

	startTime = getLogCurrentTime()
	mnemonic, err := bip39.NewMnemonic(entropy)
	printMsg("bip39.NewMnemonic", startTime)
	if err != nil {
		return nil, fmt.Errorf("NewWallet bip39.NewMnemonic err:%v", err.Error())
	}
	wallet.Mnemonic = mnemonic
	// 创建主私钥
	startTime = getLogCurrentTime()
	seed := bip39.NewSeed(mnemonic, password)
	printMsg("bip39.NewSeed", startTime)
	// 创建主私钥
	startTime = getLogCurrentTime()
	mKey, err := bip32.NewMasterKey(seed)
	printMsg("bip32.NewMasterKey", startTime)
	if err != nil {
		return nil, fmt.Errorf("NewWallet bip32.NewMasterKey err:%v", err.Error())
	}
	wallet.Key = mKey
	wallet.Time = uint32(time.Now().Unix())
	wallet.Path = path
	wallet.Password = password
	wallet.ChildKeyInfo = make(map[string]*ChildKeyPropertyInfo, 0)
	wallet.AddrLinkPubkey = make(map[string]string, 0)
	writeContentToWalletFile(wallet, path, password)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

// 从助记词中恢复主钱包
// isUsePwdBlur 是否使用Password进行混淆
func LoadWalletFromMnemonic(path string, password string, mnemonic string, isUsePwdBlur bool) (*Wallet, error) {
	// 参数检查
	if len(path) == 0 {
		return nil, fmt.Errorf("LoadWalletFromMnemonic path parameter error")
	}
	// 判断钱包文件是否存在
	_, err := ioutil.ReadFile(path)

	wallet := newWallet()
	if err == nil {
		return readWalletFromFile(path, password, wallet, err)
	}

	// 助记词判断
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("LoadWalletFromMnemonic mnemonic is err")
	}

	startTime := getLogCurrentTime()
	// 不存在的话则创建一个钱包文件
	wallet.Mnemonic = mnemonic
	// 创建主私钥
	var seed []byte
	if isUsePwdBlur {
		// 使用password作为混淆因子
		seed, err = bip39.NewSeedWithErrorChecking(mnemonic, password)
	} else {
		// ETH，BTC都没有添加混淆因子
		seed, err = bip39.NewSeedWithErrorChecking(mnemonic, "")
	}
	printMsg("bip39.NewSeedWithErrorChecking", startTime)

	// 创建主私钥
	startTime = getLogCurrentTime()
	mKey, err := bip32.NewMasterKey(seed)
	printMsg("bip32.NewMasterKey", startTime)
	if err != nil {
		return nil, fmt.Errorf("LoadWalletFromMnemonic bip32.NewMasterKey err:%v", err.Error())
	}
	wallet.Key = mKey
	wallet.Time = uint32(time.Now().Unix())
	wallet.Path = path
	wallet.Password = password
	wallet.ChildKeyInfo = make(map[string]*ChildKeyPropertyInfo, 0)
	wallet.AddrLinkPubkey = make(map[string]string, 0)
	writeContentToWalletFile(wallet, path, password)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

// 从私钥中恢复钱包
func LoadWalletFromPrvKey(path string, password string, prvKeyBase58 string) (*Wallet, error) {
	// 参数检查
	if len(path) == 0 {
		return nil, fmt.Errorf("LoadWalletFromPrvKey path parameter error")
	}
	wallet := newWallet()
	// 判断钱包文件是否存在
	_, err := ioutil.ReadFile(path)
	if err == nil {
		return readWalletFromFile(path, password, wallet, err)
	}
	if prvKeyBase58 == "" {
		return nil, errors.New("LoadWalletFromPrvKey prvKeyBase58 is empty")
	}
	startTime := getLogCurrentTime()
	// 文件不存在，解析prvKeyBase58
	mKey, err := bip32.B58Deserialize(prvKeyBase58)
	printMsg("bip32.B58Deserialize", startTime)
	if err != nil {
		return nil, fmt.Errorf("LoadWalletFromPrvKey bip32.Deserialize err:%v", err.Error())
	}
	wallet.Key = mKey
	wallet.Time = uint32(time.Now().Unix())
	wallet.Path = path
	wallet.Password = password
	wallet.ChildKeyInfo = make(map[string]*ChildKeyPropertyInfo, 0)
	wallet.AddrLinkPubkey = make(map[string]string, 0)
	writeContentToWalletFile(wallet, path, password)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

// 从文件中读取wallet
func readWalletFromFile(path string, password string, wallet *Wallet, err error) (*Wallet, error) {
	wallet, err = readContentToWalletFile(path, password)
	if err != nil {
		return nil, fmt.Errorf("LoadWalletFromMnemonic json.Unmarshal err:%v", err.Error())
	}
	wallet.Path = path
	wallet.Password = password
	if nil == wallet.ChildKeyInfo {
		wallet.ChildKeyInfo = make(map[string]*ChildKeyPropertyInfo, 0)
		wallet.AddrLinkPubkey = make(map[string]string, 0)
	}
	return wallet, nil
}

// ==========================主账户============================
// 导出主账户的助记词
func (w *Wallet) ExportMasterMnemonic() string {
	return w.Mnemonic
}

// 导出主账户的扩展私钥
func (w *Wallet) ExportMasterExtendedKey() string {
	return w.Key.String()
}

// 导出主账户的扩展私钥
func (w *Wallet) ExportMasterRawKey() string {
	return hex.EncodeToString(w.Key.Key)
}

// 删除助记词
func (w *Wallet) DelMnemonic() error {
	if w.Mnemonic == "" {
		return nil
	}
	w.Mnemonic = ""
	err := writeContentToWalletFile(w, w.Path, w.Password)
	return err
}

// ==========================设置============================
// 设置是否打印日志
func SetIsLog(b bool) {
	isLog = b
}

// 设置是否保存子私钥
func (w *Wallet) SetIsSaveSubKey(b bool) {
	w.IsSaveSubKey = b
}

// 是否保存ExtendedKey
func (w *Wallet) SetIsSaveExtendedKey(b bool) {
	w.IsSaveExtendedKey = b
}

func getLogCurrentTime() int64 {
	if isLog {
		return dateutil.CurrentTime()
	}
	return 0
}

func printMsg(msg string, startTime int64) {
	if isLog {
		log.Info(msg, "耗时", dateutil.GetDistanceTimeToCurrent(startTime))
	}
}

// ==========================子账户============================

// build child key path
// purpose 默认44，默认bip44
// algorithmType 算法类型（s256，p256，gm2）
// orgOrCoinType purpose=44,代表coinType,否则代表组织
// _account 账户空间
// change 除了btc，其他的都为0
// addressIndex 地址索引
func buildChildKeyPath(purpose, coinType, org, _account, change, addressIndex uint32) (string, error) {
	if coinType < bip32.FirstHardenedChild || _account < bip32.FirstHardenedChild {
		return "", fmt.Errorf("wallet buildChildKeyPath parameter error")
	}
	path := "/" + strconv.FormatUint(uint64(purpose)-uint64(bip32.FirstHardenedChild), 10)
	path += "/" + strconv.FormatUint(uint64(coinType)-uint64(bip32.FirstHardenedChild), 10) // 币种
	if purpose != bip44.Purpose {
		if org < bip32.FirstHardenedChild {
			return "", fmt.Errorf("wallet buildChildKeyPath parameter error")
		}
		path += "/" + strconv.FormatUint(uint64(org)-uint64(bip32.FirstHardenedChild), 10) // 组织
	}
	path += "/" + strconv.FormatUint(uint64(_account)-uint64(bip32.FirstHardenedChild), 10) // 账户空间
	path += "/" + strconv.FormatUint(uint64(change), 10)
	path += "/" + strconv.FormatUint(uint64(addressIndex), 10)
	return path, nil
}

// 通过路径获取相应的参数
func parseChildKeyPath(path string) ([]uint32, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("wallet getAccountTypeFromKeyPath parameter error")
	}
	arrPath := strings.Split(path, "/")
	arrPath = arrPath[1:]
	hardenedIndexStart := len(arrPath) - 2 // 最后两位保持int数据
	params := make([]uint32, 0)
	for i, p := range arrPath {
		typeTemp, _ := strconv.ParseUint(p, 10, 32)
		if i < hardenedIndexStart {
			typeTemp = typeTemp + uint64(bip32.FirstHardenedChild)
			if typeTemp > (uint64(1)<<uint(32) - 1) {
				return nil, fmt.Errorf("type beyond the limit err")
			}
		}
		params = append(params, uint32(typeTemp))
	}
	return params, nil
}

// 创建账户[同一机构下，同一中签名算法，的同一用户只会保留一个私钥]
// purpose purpose=44时，不使用org
// org：组织
// coinType：币种
// _account：将密钥空间划分为独立的用户身份[每一个用户对应一个地址空间]
// change：0用于外部接收地址 1用于找零地址
// addressIndex：地址索引[官方推荐不超过20]
func (w *Wallet) CreateAccount(purpose, coinType, org, _account, change, addressIndex uint32, api ChainAPI) (addr string, keyPath string, err error) {
	// 参数校验
	if nil == w.Key {
		return "", "", fmt.Errorf("wallet CreateAccount should create wallet first")
	}
	if api == nil {
		return "", "", fmt.Errorf("wallet CreateAccount chainApi is nil")
	}

	if coinType < bip32.FirstHardenedChild || _account < bip32.FirstHardenedChild {
		return "", "", fmt.Errorf("wallet CreateAccount _account should more than the %x", bip32.FirstHardenedChild)
	}
	if purpose != bip44.Purpose {
		if api.ChainInfo().Algorithm < bip32.FirstHardenedChild {
			return "", "", fmt.Errorf("wallet CreateAccount algorithmType should more than the %x", bip32.FirstHardenedChild)
		}
		if org < bip32.FirstHardenedChild {
			return "", "", fmt.Errorf("wallet CreateAccount org should more than the %x", bip32.FirstHardenedChild)
		}
	}

	// Generate sub-private key from master private key
	startTime := getLogCurrentTime()
	key, err := bip44.NewKeyFromMasterKeyWithOrg(w.Key, purpose, coinType, org, _account, change, addressIndex)
	printMsg("bip44.NewKeyFromMasterKey", startTime)
	if err != nil {
		return "", "", fmt.Errorf("wallet CreateAccount bip44.NewKeyFromMasterKey err:%v", err.Error())
	}
	// Get the public key from the private key
	startTime = getLogCurrentTime()
	pubKey, err := api.GetPubKeyFromPriKey(key.Key)
	printMsg("api.GetPubKeyFromPriKey", startTime)
	if err != nil {
		return "", "", fmt.Errorf("wallet CreateAccount getPubKeyFromPriKey err:%v", err.Error())
	}

	startTime = getLogCurrentTime()
	addr, err = api.GetAddressFromPubKey(pubKey)
	printMsg("api.GetAddressFromPubKey", startTime)
	if err != nil {
		return "", "", err
	}
	if len(addr) == 0 {
		return "", "", fmt.Errorf("wallet CreateAccount accountObj.GetAddressFromPubKey return value err")
	}

	// 生成路径
	startTime = getLogCurrentTime()
	keyPath, err = buildChildKeyPath(purpose, coinType, org, _account, change, addressIndex)
	printMsg("buildChildKeyPath", startTime)
	if err != nil {
		return "", "", err
	}
	if w.IsSaveSubKey {
		// Generate sub-private key propertyInfo
		childKeyPropertyInfo := new(ChildKeyPropertyInfo)
		childKeyPropertyInfo.Purpose = purpose
		childKeyPropertyInfo.ChainType = api.ChainInfo().ChainType
		childKeyPropertyInfo.AlgorithmType = api.ChainInfo().Algorithm
		childKeyPropertyInfo.Org = org
		childKeyPropertyInfo.CoinType = coinType
		childKeyPropertyInfo.Time = uint32(time.Now().Unix())
		// 将subKey进行加密
		subPwd := w.getSubPwd(keyPath)

		subPivKey := key.Key
		if w.IsSaveExtendedKey {
			subPivKey, err = key.Serialize()
			if err != nil {
				return "", "", err
			}
		}

		startTime = getLogCurrentTime()
		encryptKey, err := scrypt.EncryptKey(&scrypt.Key{Id: uuid.NewRandom(), Path: keyPath, PrivateKey: subPivKey}, subPwd, scrypt.StandardScryptN, scrypt.StandardScryptP)
		printMsg("scrypt.EncryptKey sub", startTime)
		if err != nil {
			return "", "", err
		}

		childKeyPropertyInfo.Key = encryptKey

		startTime = getLogCurrentTime()
		pubKeyStr := w.getPubKey(pubKey)
		w.mu.Lock()
		w.AddrLinkPubkey[addr] = pubKeyStr
		w.ChildKeyInfo[pubKeyStr] = childKeyPropertyInfo
		w.mu.Unlock()
		printMsg("AddrLinkPubkey", startTime)
		err = writeContentToWalletFile(w, w.Path, w.Password)
		if err != nil {
			return "", "", fmt.Errorf("wallet CreateAccount ioutil.WriteFile err:%v", err.Error())
		}
	}

	if nil == w.ChildKeyInfo {
		w.ChildKeyInfo = make(map[string]*ChildKeyPropertyInfo, 0)
		w.AddrLinkPubkey = make(map[string]string, 0)
	}
	return addr, keyPath, nil
}

// 通过路径恢复私钥
func (w *Wallet) createAccountByPath(path string, api ChainAPI) (addr string, key *bip32.Key, err error) {
	if api == nil {
		return "", nil, fmt.Errorf("wallet CreateAccount chainApi is nil")
	}
	if path == "" {
		return "", nil, errors.New("path is empty")
	}
	childKeyPath, err := parseChildKeyPath(path)
	if err != nil {
		return "", nil, err
	}

	org := uint32(0)
	offset := 0
	if len(childKeyPath) != 5 {
		org = childKeyPath[2]
		offset = 1
	}

	startTime := getLogCurrentTime()
	// Generate sub-private key from master private key
	key, err = bip44.NewKeyFromMasterKeyWithOrg(w.Key, childKeyPath[0], childKeyPath[1], org, childKeyPath[offset+2], childKeyPath[offset+3], childKeyPath[offset+4])
	printMsg("bip44.NewKeyFromMasterKey", startTime)
	if err != nil {
		return "", nil, fmt.Errorf("wallet CreateAccount bip44.NewKeyFromMasterKey err:%v", err.Error())
	}
	// Get the public key from the private key
	startTime = getLogCurrentTime()
	pubKey, err := api.GetPubKeyFromPriKey(key.Key)
	printMsg("api.GetPubKeyFromPriKey", startTime)
	if err != nil {
		return "", nil, fmt.Errorf("wallet CreateAccount getPubKeyFromPriKey err:%v", err.Error())
	}

	startTime = getLogCurrentTime()
	addr, err = api.GetAddressFromPubKey(pubKey)
	printMsg("api.GetAddressFromPubKey", startTime)
	if err != nil {
		return "", nil, err
	}
	if len(addr) == 0 {
		return "", nil, fmt.Errorf("addr is empty")
	}

	return addr, key, nil
}

// ListAccount list all account
func (w *Wallet) ListAccount() ([]string, error) {
	if nil == w.AddrLinkPubkey {
		return nil, nil
	}
	accounts := make([]string, 0)
	for key := range w.AddrLinkPubkey {
		accounts = append(accounts, key)
	}
	return accounts, nil
}

// 将地址对应的私钥转换成keystore
func (w *Wallet) ExportKeyStore(address, keyPath string, keystorePwd string, api ChainAPI) (keyStore string, err error) {
	// 参数校验
	if len(address) == 0 || len(keystorePwd) == 0 {
		return "", fmt.Errorf("wallet ExportKeyStore parameter error")
	}

	bip32Key, err := w.getRawPrivateKey(address, keyPath, api)
	if err != nil {
		return "", err
	}

	keyStoreBytes, err := scrypt.EncryptKey(
		&scrypt.Key{
			Id:         uuid.NewRandom(),
			Path:       keyPath,
			Addr:       address,
			PrivateKey: bip32Key.Key,
		},
		keystorePwd,
		scrypt.StandardScryptN,
		scrypt.StandardScryptP)
	return string(keyStoreBytes), err
}

func (w *Wallet) getRawPrivateKey(address string, keyPath string, api ChainAPI) (bip32Key *bip32.Key, err error) {
	pubKeyStr := w.AddrLinkPubkey[address]
	startTime := getLogCurrentTime()
	if pubKeyStr != "" {
		childKeyInfo := w.ChildKeyInfo[pubKeyStr]
		if nil == childKeyInfo {
			return nil, fmt.Errorf("wallet ExportKeyStore key store not exist")
		}
		priKeyBytes1 := childKeyInfo.Key

		// 将subKeyEnc进行解密
		// 外层的密码+addr作为
		subPwd := w.getSubPwd(keyPath)
		priKey, err := scrypt.DecryptKey(priKeyBytes1, subPwd)
		printMsg("scrypt.DecryptKey sub", startTime)
		if err != nil {
			return nil, err
		}
		priKeyBytes := priKey.PrivateKey
		if len(priKeyBytes) != 32 {
			bip32Key, err = bip32.Deserialize(priKeyBytes)
			if err != nil {
				return nil, err
			}
		}
	}
	startTime = getLogCurrentTime()
	addr, prvKey, err := w.createAccountByPath(keyPath, api)
	printMsg("createAccountByPath", startTime)
	if err != nil {
		return nil, err
	}
	if addr != address {
		return nil, errors.New("address is diff")
	}
	bip32Key = prvKey
	return
}

// 将地址对应的私钥直接输出
func (w *Wallet) ExportRawKey(address, keyPath string, api ChainAPI) (key string, err error) {
	// 参数校验
	if len(address) == 0 {
		return "", fmt.Errorf("wallet ExportKeyStore parameter error")
	}
	if nil == w.ChildKeyInfo {
		return "", fmt.Errorf("wallet ExportKeyStore key store not exist")
	}

	bip32Key, err := w.getRawPrivateKey(address, keyPath, api)
	if err != nil {
		return "", err
	}
	return api.ExportPrivateKey(bip32Key.Key, false)
}

// 导出扩展私钥
func (w *Wallet) ExportExtendedKey(address, keyPath string, api ChainAPI) (extendedKey string, err error) {
	// 参数校验
	if len(address) == 0 {
		return "", fmt.Errorf("wallet ExportKeyStore parameter error")
	}
	bip32Key, err := w.getRawPrivateKey(address, keyPath, api)
	if err != nil {
		return "", err
	}
	return bip32Key.String(), nil
}

// 导入keystore
func (w *Wallet) ImportKeyStore(key []byte, password string) (address string, err error) {
	return "", fmt.Errorf("not support")
}

// 导入私钥
func (w *Wallet) ImportRawKey(key string, password string) (address string, err error) {
	return "", fmt.Errorf("not support")
}

func (w *Wallet) getSubPwd(keyPath string) string {
	// 外层的密码+addr作为
	keccak256Hash := scrypt.Keccak256([]byte(w.Password + keyPath))
	subPwd := hex.EncodeToString(keccak256Hash)
	return subPwd
}

// 获取w中存储的key值
func (w *Wallet) getPubKey(pubKey []byte) string {
	// 将地址进行hash
	pubKeyStr := hex.EncodeToString(pubKey)
	return pubKeyStr
}

// GetPriKeyFromAddress 获取某个地址的私钥
func (w *Wallet) GetPriKeyFromAddress(address, keyPath string, api ChainAPI) ([]byte, error) {
	strKey, err := w.ExportRawKey(address, keyPath, api)
	if err != nil {
		return nil, fmt.Errorf("wallet GetPriKeyFromAddress err:%v", err.Error())
	}
	return hex.DecodeString(strKey)
}

// Sign
func (w *Wallet) Sign(address, keyPath string, hash []byte, api ChainAPI) (string, error) {
	// 参数校验
	if len(address) == 0 {
		return "", fmt.Errorf("wallet ExportKeyStore parameter error")
	}
	if nil == w.ChildKeyInfo {
		return "", fmt.Errorf("wallet ExportKeyStore key store not exist")
	}

	startTime := getLogCurrentTime()
	bip32Key, err := w.getRawPrivateKey(address, keyPath, api)
	printMsg("w.getRawPrivateKey(address, keyPath, api)", startTime)
	startTime = getLogCurrentTime()
	sign, err := api.SignToStr(bip32Key.Key, hash)
	printMsg("api.SignToStr", startTime)
	if err != nil {
		return "", err
	}
	return sign, nil
}
