import { QQ_ICONS } from './qqIcons.js'

export function qqRank(coins, coinConfig) {
	const cfg = coinConfig && typeof coinConfig === 'object' ? coinConfig : {}
	const dayCoin = Math.max(1, Number(cfg.dayCoin) || 5)
	const n = Number(coins)
	if (!isFinite(n) || n < 0) return null
	const days = n / dayCoin
	let level = 0
	if (days >= 5) {
		level = Math.floor(-2 + Math.sqrt(days + 4) + 1e-8)
		if (level < 0) level = 0
		while (level > 0 && days + 1e-8 < level * (level + 4)) level -= 1
		while (days + 1e-8 >= (level + 1) * (level + 5)) level += 1
	}
	const nextDays = (level + 1) * (level + 5)
	const lack = Math.max(0, Math.ceil(nextDays * dayCoin - n - 1e-8))
	const remainDays = Math.max(0, Math.ceil(nextDays - days - 1e-8))
	let rest = level
	const crown = Math.floor(rest / 64)
	rest %= 64
	const sun = Math.floor(rest / 16)
	rest %= 16
	const moon = Math.floor(rest / 4)
	const star = rest % 4
	const marks = []
	const showCrown = Math.min(crown, 6)
	pushMarks(marks, 'crown', showCrown, QQ_ICONS.crown)
	pushMarks(marks, 'sun', sun, QQ_ICONS.sun)
	pushMarks(marks, 'moon', moon, QQ_ICONS.moon)
	pushMarks(marks, 'star', star, QQ_ICONS.star)
	return {
		level: level,
		dayCoin: dayCoin,
		marks: marks,
		crownExtra: Math.max(0, crown - showCrown),
		lack: lack,
		remainDays: remainDays,
		icons: QQ_ICONS,
	}
}

function pushMarks(marks, type, count, src) {
	for (let i = 0; i < count; i++) marks.push({ type: type, icon: src })
}
