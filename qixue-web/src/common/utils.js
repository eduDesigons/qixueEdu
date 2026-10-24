import { showToast } from '@/common/toast'
let _debounceTimeout = null,
	_throttleRunning = false

/**
 * 防抖
 * @param {Function} 执行函数
 * @param {Number} delay 延时ms   
 */
export const debounce = (fn, delay = 500) => {
	clearTimeout(_debounceTimeout);
	_debounceTimeout = setTimeout(() => {
		fn();
	}, delay);
}
/**
 * 节流
 * @param {Function} 执行函数
 * @param {Number} delay 延时ms  
 */
export const throttle = (fn, delay = 500) => {
	if (_throttleRunning) {
		return;
	}
	_throttleRunning = true;
	fn();
	setTimeout(() => {
		_throttleRunning = false;
	}, delay);
}
/**
 * toast
 */
export const msg = (title = '', param={}) => {
	if(!title) return;
	showToast('');
}

/**
 * 获取页面栈
 * @param {Number} preIndex为1时获取上一页
 * @return {Object} 
 */
export const prePage = (preIndex = 1) => {
	const pages = getCurrentPages();
	const prePage = pages[pages.length - (preIndex + 1)];

	return prePage.$vm;
}