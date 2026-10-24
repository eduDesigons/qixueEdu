import store from '@/store'
import config from '@/config'
import md5 from 'js-md5'
import { showToast, showLoading, hideLoading } from '@/common/toast'
import { toUploadBlob } from '@/common/media'

const REQUEST_TIMEOUT = 120000

const createHeader = () => {
	const chars = 'ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678'
	let random = ''
	for (let i = 0; i < 32; i++) random += chars.charAt(Math.floor(Math.random() * chars.length))
	const header = {
		Authorization: 'Bearer ' + store.state.token,
		'version-xcx': config.version,
		'timestmp-xcx': Date.now(),
		'random-xcx': random,
	}
	header['sign-xcx'] = md5('xq' + header['version-xcx'] + header['timestmp-xcx'] + header['random-xcx'])
	return header
}

function isLoginPage() {
	try {
		return location.pathname === '/pages/index/login'
	} catch (e) {
		return false
	}
}

function forceLogin(header, body) {
	const sent = String((header && header.Authorization) || '')
	const currentToken = store.state.token || ''
	if (currentToken && sent && sent !== 'Bearer ' + currentToken) return false
	if (isLoginPage()) return false
	showToast((body && (body.message || body.msg)) || '登录已过期，请重新登录')
	store.commit('logout')
	return true
}

function requestViaFetch(url, data, method, header) {
	const m = (method || 'GET').toUpperCase()
	const opts = {
		method: m,
		headers: Object.assign({ 'Content-Type': 'application/json' }, header || {}),
	}
	if (typeof AbortSignal !== 'undefined' && AbortSignal.timeout) {
		opts.signal = AbortSignal.timeout(REQUEST_TIMEOUT)
	}
	let finalUrl = url
	if (m === 'GET' || m === 'HEAD') {
		const qs = Object.keys(data || {})
			.filter((k) => data[k] !== undefined && data[k] !== null)
			.map((k) => `${encodeURIComponent(k)}=${encodeURIComponent(data[k])}`)
			.join('&')
		if (qs) finalUrl += (finalUrl.includes('?') ? '&' : '?') + qs
	} else {
		opts.body = JSON.stringify(data || {})
	}
	return fetch(finalUrl, opts).then(async (res) => {
		const ct = res.headers.get('content-type') || ''
		const body = ct.includes('application/json') ? await res.json() : await res.text()
		return { statusCode: res.status, data: body }
	})
}

const inflight = new Map()

function inflightKey(method, url, data) {
	const m = (method || 'GET').toUpperCase()
	if (m !== 'GET' && m !== 'HEAD') return ''
	try {
		return m + ' ' + url + '?' + JSON.stringify(data || {})
	} catch (e) {
		return m + ' ' + url
	}
}

export const request = ({ url = '', data = {}, method = 'get', showLoading: loading = false } = {}) => {
	if (loading) showLoading(typeof loading === 'string' ? loading : '加载中')
	const header = createHeader()
	const finalUrl = url.includes('://') ? url : config.baseUrl + url
	const key = inflightKey(method, finalUrl, data)
	if (key && inflight.has(key)) {
		return inflight.get(key).finally(() => {
			if (loading) hideLoading()
		})
	}
	const pending = requestViaFetch(finalUrl, data, method || 'POST', header)
		.then((res) => {
			hideLoading()
			const body = res.data && typeof res.data === 'object' && !(res.data instanceof ArrayBuffer) ? res.data : {}
			if (res.statusCode === 401 || body.code == 401) {
				if (!forceLogin(header, body)) {
					return Promise.reject(body.code != null ? body : { code: 401, message: '登录已过期，请重新登录' })
				}
				return Promise.reject(body)
			}
			if (res.statusCode !== 200) {
				showToast('网络错误')
				return Promise.reject('网络错误')
			}
			if (body.code == 200) return res.data
			if (body.code == 400) return Promise.reject(res.data)
			showToast(body.message || body.msg || '请求失败')
			return Promise.reject(res.data)
		})
		.catch((err) => {
			hideLoading()
			return Promise.reject(err)
		})
		.finally(() => {
			if (key) inflight.delete(key)
		})
	if (key) inflight.set(key, pending)
	return pending
}

export const upload = async (file) => {
	showLoading('上传中')
	try {
		const blob = await toUploadBlob(file)
		const form = new FormData()
		const filename = (file && file.name) || (blob.type && blob.type.includes('png') ? 'upload.png' : 'upload.jpg')
		form.append('file', blob, filename)
		const header = createHeader()
		// multipart 由浏览器带 boundary
		delete header['Content-Type']
		const res = await fetch(config.baseUrl + '/common/upload', {
			method: 'POST',
			headers: header,
			body: form,
		})
		const ct = res.headers.get('content-type') || ''
		let body = ct.includes('application/json') ? await res.json() : {}
		if (typeof body === 'string') {
			try { body = JSON.parse(body) } catch (e) { body = {} }
		}
		if (res.status === 401 || (body && body.code == 401)) {
			if (!forceLogin(header, body || {})) {
				throw body && body.code != null ? body : { code: 401, message: '登录已过期，请重新登录' }
			}
			throw body
		}
		if (res.status !== 200 || !body || body.code != 200) {
			showToast((body && (body.message || body.msg)) || '上传失败')
			throw body || { message: '上传失败' }
		}
		return body
	} finally {
		hideLoading()
	}
}

export const http = {
	get: (url, data = {}, showLoadingFlag = false) => request({ url, data, method: 'get', showLoading: showLoadingFlag }),
	post: (url, data = {}, showLoadingFlag = false) => request({ url, data, method: 'post', showLoading: showLoadingFlag }),
	upload: (file) => upload(file),
}
