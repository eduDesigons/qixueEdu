import { showToast } from '@/common/toast'
import { openLocation, previewImage, setClipboard, timeFrom } from '@/common/dom'
import { navigateBack, navigateTo, redirectTo, switchTab } from '@/common/nav'
import { rewriteHtml } from '@/common/html'
import {mapState} from 'vuex'
import { applySiteBranding } from '@/common/siteTitle.js'

export default{
	computed: {
		...mapState({
		    isDark: 'isDark',
			userInfo: 'userInfo',
			appInfo: 'appInfo',
			nav: 'nav',
			unreadMessageCount: 'unreadMessageCount'
		}),
		assetDomain() {
			return (this.appInfo?.sysConfig?.assetDomain || '')
		},
		siteName() {
			return String(this.appInfo?.sysConfig?.name || '').trim()
		},
	},
	onShow() {
		this.$nextTick(() => {
			applySiteBranding(this.appInfo)
		})
	},
	methods: {
		str2arr(str){
			return str ? str.split(',') : []
		},
		copyIt(str){
			setClipboard({ data: str,
				success: () => {
					this.toast('已经复制到剪切板')
				} })
		},
		static(src){
			if (!src) return null
			if (typeof src !== 'string') return src
			if (src.indexOf('http') === 0 || src.indexOf('data:') === 0 || src.indexOf('blob:') === 0) return src
			if (src.indexOf('/static/') === 0 || src.indexOf('/src/') === 0 || src.indexOf('/assets/') === 0 || src.indexOf('/@') === 0) return src
			if (!this.assetDomain) return src
			return this.assetDomain + (src.startsWith('/') ? src : '/' + src)
		},
		richHtml(html) {
			return rewriteHtml(html, this.assetDomain)
		},
		get_unix_time(dateStr) {
			if(!dateStr) return '';
		    var newstr = dateStr.replace(/-/g, '/');
		    var date = new Date(newstr);
		    return date.getTime();
		},
		date(dateTime, formatStr = 'yyyy-mm-dd') {
			return this.formatDate(dateTime, formatStr)
		},
		formatDate(dateTime, formatStr = 'yyyy-mm-dd') {
			if (dateTime === undefined || dateTime === null || dateTime === '') return ''
			const raw = String(dateTime).trim()
			let date
			if (/^\d{10}$/.test(raw)) date = new Date(Number(raw) * 1000)
			else if (/^\d+$/.test(raw)) date = new Date(Number(raw))
			else if (raw.indexOf('-') !== -1 && raw.indexOf('T') === -1) date = new Date(raw.replace(/-/g, '/'))
			else date = new Date(dateTime)
			if (Number.isNaN(date.getTime())) return ''
			const pad = (n) => String(n).padStart(2, '0')
			const year = String(date.getFullYear())
			const month = pad(date.getMonth() + 1)
			const day = pad(date.getDate())
			const hour = pad(date.getHours())
			const minute = pad(date.getMinutes())
			const second = pad(date.getSeconds())
			return String(formatStr)
				.replace(/YYYY|yyyy/g, year)
				.replace(/YY|yy/g, year.slice(2))
				.replace(/MM/g, month)
				.replace(/DD|dd/g, day)
				.replace(/HH|hh/g, hour)
				.replace(/ss/g, second)
				.replace(/mm/g, (token, offset, source) => (source.charAt(offset - 1) === ':' ? minute : month))
		},
		timeFrom(dateStr, format = false) {
			if(!dateStr) return '';
			return timeFrom(this.get_unix_time(dateStr), format);
		},
		showDuration(duration) {
			if(!duration) return '';
			const days = Math.floor(duration / 86400);
			const hours = Math.floor(duration / 3600);
			if(days > 0){
				return `${hours}小时`;
			}			
			const minutes = Math.floor((duration % 3600) / 60);
			const seconds = duration % 60;
			if(hours > 0){
				return `${hours}时${minutes}分`;
			}else{
				return `${minutes}分${seconds}秒`;
			}
		},
		showSize(size) {
			if(!size) return '';
			if(size < 1024){
				return `${size}B`;
			}else if(size < 1024 * 1024){
				return `${(size / 1024).toFixed(2)}KB`;
			}else{
				return `${(size / 1024 / 1024).toFixed(2)}MB`;
			}
		},
		previewPics(urls =[], i = 0) {
			previewImage({
				showmenu: false,
				current: Array.isArray(urls) ? this.static(urls[i]) : this.static(urls),
				urls: Array.isArray(urls) ? urls.map(e=>this.static(e)) : [this.static(urls)]
			})
		},
		msg(msg){
			showToast(msg)
		},
		navTo(url, close=false){
			if(url.startsWith('/pages/index/home') || url.startsWith('/pages/index/my') || url.startsWith('/pages/index/course') || url.startsWith('/pages/index/learn')){
				switchTab(url)
			}else{
				if(close){
					redirectTo(url)
				}else{
					navigateTo(url)
				}
			}
		},
		back(dely=0){
			console.log('back')
			if(dely){
				setTimeout(()=>{
					navigateBack(1)
				}, dely)
			}else{
				navigateBack(1)
			}
			
		},
		toMap(item){
			openLocation({
				latitude: Number(item.lat),
				longitude: Number(item.lng),
				name: item.name,
				address: item.address
			})
		}
	},
}
