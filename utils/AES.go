package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// 通过用户的口令生成 AES 密钥
func DeriveKey(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

// AES 加密函数，使用 AES-GCM 模式加密数据
func Encrypt(plaintext, password string) ([]byte, error) {
	// 生成密钥
	key := DeriveKey(password)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 创建 AES-GCM 模式
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 生成随机 nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// 使用 AES-GCM 模式加密数据
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return ciphertext, nil
}

// AES 解密函数，使用 AES-GCM 模式解密数据
func Decrypt(ciphertext []byte, password string) (string, error) {
	// 生成密钥
	key := DeriveKey(password)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 创建 AES-GCM 模式
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 提取 nonce 和加密数据
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 解密数据
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// 将加密的数据保存到文件
func SaveEncryptedData(filename, data, password string) error {
	// 加密数据
	encryptedData, err := Encrypt(data, password)
	if err != nil {
		return err
	}

	// 设置保存路径
	saveDir := "record"
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return err
	}
	savePath := filepath.Join(saveDir, filename)

	// 将加密数据写入文件
	return os.WriteFile(savePath, encryptedData, 0600)
}

// 从文件加载并解密数据
func LoadEncryptedData(filename, password string) (string, error) {
	// 读取文件内容
	encryptedData, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	// 解密数据
	return Decrypt(encryptedData, password)
}

// 将密文转换为十六进制字符串
func ToHexString(data []byte) string {
	return hex.EncodeToString(data)
}

// 从十六进制字符串解析密文
func FromHexString(hexData string) ([]byte, error) {
	return hex.DecodeString(hexData)
}
