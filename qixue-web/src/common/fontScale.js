import { getStorage, setStorage } from '@/common/storage'
export const FONT_SCALE_KEY = 'fontScale'

export const FONT_LEVELS = [
	{ scale: 0.92, label: '小' },
	{ scale: 1, label: '标准' },
	{ scale: 1.12, label: '稍大' },
	{ scale: 1.24, label: '大' },
	{ scale: 1.36, label: '特大' },
	{ scale: 1.48, label: '超大' }
]

export const FONT_STANDARD_INDEX = 1

function clampScale(scale) {
	const n = Number(scale)
	if (!Number.isFinite(n) || n < 0.8 || n > 1.6) return 1
	return FONT_LEVELS[findLevelIndex(n)].scale
}

export function findLevelIndex(scale) {
	const n = Number(scale)
	let best = FONT_STANDARD_INDEX
	let bestDiff = Infinity
	FONT_LEVELS.forEach((item, index) => {
		const diff = Math.abs(item.scale - n)
		if (diff < bestDiff) {
			bestDiff = diff
			best = index
		}
	})
	return best
}

export function readFontScale() {
	try {
		const saved = getStorage(FONT_SCALE_KEY)
		if (saved === '' || saved === null || saved === undefined) return 1
		return clampScale(saved)
	} catch (e) {
		return 1
	}
}

export function fontScaleLabel(scale) {
	return FONT_LEVELS[findLevelIndex(scale == null ? readFontScale() : scale)].label
}

export function applyFontScale(scale) {
	const value = scale == null ? readFontScale() : clampScale(scale)
	const ui = 1 + (value - 1) * 0.65
	if (typeof document !== 'undefined' && document.documentElement) {
		const root = document.documentElement
		root.style.setProperty('--font-scale', String(value))
		root.style.setProperty('--font-scale-ui', String(Math.round(ui * 1000) / 1000))
		root.setAttribute('data-font-custom', Math.abs(value - 1) > 0.001 ? '1' : '0')
	}
	return value
}

export function saveFontScale(scale) {
	const value = clampScale(scale)
	try {
		setStorage(FONT_SCALE_KEY, String(value))
	} catch (e) {}
	return applyFontScale(value)
}
