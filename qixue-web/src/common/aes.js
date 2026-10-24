import CryptoJS from 'crypto-js/core'
import 'crypto-js/aes'
import 'crypto-js/mode-ecb'
import 'crypto-js/pad-pkcs7'

export function aesEncrypt(word, keyWord) {
	const key = CryptoJS.enc.Utf8.parse(keyWord)
	const srcs = CryptoJS.enc.Utf8.parse(word)
	const encrypted = CryptoJS.AES.encrypt(srcs, key, {
		mode: CryptoJS.mode.ECB,
		padding: CryptoJS.pad.Pkcs7,
	})
	return encrypted.toString()
}
