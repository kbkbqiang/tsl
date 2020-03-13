/*
 @Time : 2020/3/13 9:43 AM
 @Author : chenye
 @File : eas_go
 @Software : GoLand
 @Remark : 
*/

package EAS

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

func GoEncrypt(encodeByte []byte) (cipherText string, err error){

	ckey, err := aes.NewCipher(key)
	if nil != err {
		fmt.Println("钥匙创建错误:", err)
		return "", err
	}

	blockSize := ckey.BlockSize()

	fmt.Println("加密的字符串", string(encodeByte), "\n加密钥匙", key, "\n向量IV", string(iv))

	fmt.Println("加密前的字节：", encodeByte, "\n")

	encrypter := cipher.NewCBCEncrypter(ckey, iv)

	// PKCS7补码
	encodeByte = PKCS7Padding(encodeByte, blockSize)
	out := make([]byte, len(encodeByte))

	encrypter.CryptBlocks(out, encodeByte)
	fmt.Println("加密后字节：", out)


	// hex 兼容nodejs cropty-js包
	//cipherText = hex.EncodeToString(out)
	cipherText = base64.URLEncoding.EncodeToString(out)
	return cipherText, nil
}

func GoDecrypt(encodeStr string) (origByte []byte, err error) {
	ckey, err := aes.NewCipher(key)
	if nil != err {
		fmt.Println("钥匙创建错误:", err)
		return origByte, err
	}

	//base64Str,err := hex.DecodeString(encodeStr)
	if err != nil {
		return origByte, err
	}
	//base64Out := base64.URLEncoding.EncodeToString(base64Str)


	//fmt.Println("\n开始解码")
	decrypter := cipher.NewCBCDecrypter(ckey, iv)

	base64In, err := base64.URLEncoding.DecodeString(encodeStr)

	if err != nil {
		return origByte, err
	}

	in := make([]byte, len(base64In))

	decrypter.CryptBlocks(in, base64In)

	//fmt.Println("解密后的字节：", in)

	// 去除PKCS7补码
	in = UnPKCS7Padding(in)

	//fmt.Println("去PKCS7补码：", in)
	//fmt.Println("解密：", string(in))
	return in,nil
}


