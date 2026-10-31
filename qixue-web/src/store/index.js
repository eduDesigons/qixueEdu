import { getStorage, removeStorage, setStorage } from '@/common/storage'
import { reLaunch } from '@/common/nav'
import { createStore } from 'vuex'
import { http } from '@/common/request.js'
import { applySiteBranding } from '@/common/siteTitle.js'

function prefetchAssetOrigin(domain) {
	if (!domain || typeof document === 'undefined' || !document.head) return
	if (document.getElementById('qixue-asset-preconnect')) return
	let href = String(domain).trim()
	if (!href) return
	if (href.indexOf('http') !== 0) href = 'https://' + href.replace(/^\/+/, '')
	href = href.replace(/\/+$/, '')
	const preconnect = document.createElement('link')
	preconnect.id = 'qixue-asset-preconnect'
	preconnect.rel = 'preconnect'
	preconnect.href = href
	preconnect.crossOrigin = 'anonymous'
	document.head.appendChild(preconnect)
	const dns = document.createElement('link')
	dns.rel = 'dns-prefetch'
	dns.href = href
	document.head.appendChild(dns)
}

const store = createStore({
	state: {
		isDark: false,
		token: '',
		nav: {
			top: '0px',
			bottom: '0px',
			height: '20px'
		},
		appInfo: {
			sysConfig: {

			},
			playerConfig: {

			},
			coinConfig: {

			}
		},
		userInfo: {
			id: null,
			name: "未登录",
			avatar: "",
			info_status: 0,
			base_work_city: '',
			mobile: null,
			sex: 1,
			account_type: 1,
		},
		unreadMessageCount: 0,
		badgeFetchedAt: 0
	},
	getters: {
		hasLogin(state){
			return !!state.token;
		}
	},
	mutations: {
		 setDarkTheme(state,  isDark) {
		  state.isDark = isDark

		},
		setStateAttr(state, param){
			if(param instanceof Array){
				for(let item of param){
					state[item.key] = item.val;
				}
			}else{
				state[param.key] = param.val;
			}
		},
		setToken(state, token){
			state.token = token;
			setStorage('token', token)
		},
		setShareInfo(state, shareInfo){
			state.shareInfo = shareInfo
		},
		setUserInfo(state, userInfo){
			state.userInfo = userInfo;
			setStorage('userInfo', userInfo)
		},
		setUnreadMessageCount(state, count) {
			state.unreadMessageCount = Math.max(0, Number(count) || 0)
			state.badgeFetchedAt = Date.now()
		},
		setBadgeFetchedAt(state, at) {
			state.badgeFetchedAt = at || Date.now()
		},
		logout(state){
			state.token = '';
			state.unreadMessageCount = 0;
			state.badgeFetchedAt = 0;
			removeStorage('token');
			setTimeout(()=>{
				state.userInfo = {};
				removeStorage('userInfo')
				reLaunch('/pages/index/login')
			}, 0)
		},
		setSysInfo(state, sysInfo){
			state.appInfo = sysInfo
			setStorage('sysInfo', sysInfo)
			applySiteBranding(sysInfo)
		}
	},
	actions: {
		getSysInfo({commit}){
			http.get('/sys/getSysInfo').then(res=>{
				commit('setSysInfo', res.data)
				prefetchAssetOrigin(res.data && res.data.sysConfig && res.data.sysConfig.assetDomain)
			}).catch(() => {
				try {
					const cached = getStorage('sysInfo')
					if (cached && cached.sysConfig) {
						commit('setSysInfo', cached)
						prefetchAssetOrigin(cached.sysConfig.assetDomain)
					}
				} catch (e) {}
			})
		},
		getUserInfo({state, commit, dispatch}){
			return new Promise((resolve, reject) => {
				http.get(
					'/user/info'
				).then(res =>{
					const userInfo = res.data;
					commit('setUserInfo', userInfo)
					dispatch('refreshMessageBadge')
					return resolve(res)
				}).catch((res)=>{
					 return reject()
				})
			})

		},
		applyMessageBadge({ commit }, count) {
			const n = Math.max(0, Number(count) || 0)
			commit('setUnreadMessageCount', n)
			return n
		},
		refreshMessageBadge({ state, commit, dispatch }) {
			if (!state.token) {
				return dispatch('applyMessageBadge', 0)
			}
			if (state.badgeFetchedAt && Date.now() - state.badgeFetchedAt < 30000) {
				return Promise.resolve(state.unreadMessageCount)
			}
			commit('setBadgeFetchedAt', Date.now())
			return http.get('/user/unreadMessageCount').then(res => {
				const count = (res.data && res.data.count) || 0
				return dispatch('applyMessageBadge', count)
			}).catch(() => dispatch('applyMessageBadge', 0))
		},

	}
})

export default store
