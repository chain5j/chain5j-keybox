// description: keybox 
// 
// @author: xwc1125
// @date: 2020/8/18 0018
package keybox

const (
	MnemonicType_Chinese_Simplified  MnemonicType = "chinese_simplified"
	MnemonicType_Chinese_Traditional MnemonicType = "chinese_traditional"
	MnemonicType_English             MnemonicType = "english"
	MnemonicType_French              MnemonicType = "french"
	MnemonicType_Italian             MnemonicType = "italian"
	MnemonicType_Japanese            MnemonicType = "japanese"
	MnemonicType_Korean              MnemonicType = "korean"
	MnemonicType_Spanish             MnemonicType = "spanish"
)

type MnemonicType string

func ParseMnemonicType(t string) MnemonicType {
	switch t {
	case "zh-cn":
		return MnemonicType_Chinese_Simplified
	case "zh-tw":
		return MnemonicType_Chinese_Traditional
	case "en":
		return MnemonicType_English
	case "fr":
		return MnemonicType_French
	case "it":
		return MnemonicType_Italian
	case "ja":
		return MnemonicType_Japanese
	case "ko":
		return MnemonicType_Korean
	case "es":
		return MnemonicType_Spanish
	default:
		return MnemonicType_English
	}
}
