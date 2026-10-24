<template>
	<div class="fs-page" :style="pageStyle">
		<div class="fs-preview">
			<div class="pv">
				<div class="pv-brief">
					<div class="pv-brief__left">
						<div class="pv-brief__title">
							<span class="pv-brief__category">预览</span>
							调整文字大小
						</div>
						<div class="pv-brief__subtitle">拖动下方滑块，看看哪种字号更清楚。</div>
					</div>
					<div class="pv-brief__cover">
						<div class="pv-brief__cover-shade">
							<span class="pv-brief__views">字号</span>
						</div>
					</div>
				</div>

				<div class="pv-stats">
					<div class="pv-stats__item">
						<div class="pv-stats__label">6 档</div>
						<div class="pv-stats__tip">从小到超大</div>
					</div>
					<div class="pv-stats__item">
						<div class="pv-stats__label">标准</div>
						<div class="pv-stats__tip">默认大小</div>
					</div>
					<div class="pv-stats__item">
						<div class="pv-stats__label">全站</div>
						<div class="pv-stats__tip">一起变化</div>
					</div>
				</div>

				<div class="pv-cert">选好字号后点「完成」，首页、课程和设置都会按这个大小显示。</div>

				<div class="pv-tabs">
					<div
						class="pv-tabs__item"
						:class="{ 'pv-tabs__item--on': previewTab === 'guide' }"
						@click="previewTab = 'guide'"
					>说明</div>
					<div
						class="pv-tabs__item"
						:class="{ 'pv-tabs__item--on': previewTab === 'effect' }"
						@click="previewTab = 'effect'"
					>效果</div>
				</div>

				<div v-if="previewTab === 'guide'" class="pv-intro">
					<div class="pv-intro__title">怎么调整</div>
					<div class="pv-intro__text">左右拖动下方滑块，可以在小、标准、稍大、大、特大、超大之间切换。上面的标题、正文和列表会马上跟着变，方便对照阅读。</div>
				</div>

				<div v-else class="pv-catalog">
					<div class="pv-catalog__head">调整后的效果</div>
					<div class="pv-chapter">
						<div class="pv-chapter__mark"></div>
						<div class="pv-chapter__info">
							<div class="pv-chapter__title">标题和正文一起变大</div>
							<div class="pv-chapter__meta">列表、提示里的文字也会跟着变化</div>
						</div>
						<div class="pv-chapter__status">即时预览</div>
					</div>
					<div class="pv-chapter pv-chapter--current">
						<div class="pv-chapter__mark pv-chapter__mark--on"></div>
						<div class="pv-chapter__info">
							<div class="pv-chapter__title">点完成后长期有效</div>
							<div class="pv-chapter__meta">设置保存在本机，下次打开仍然是这个字号</div>
						</div>
						<div class="pv-chapter__status pv-chapter__status--on">记住选择</div>
					</div>
				</div>
			</div>
		</div>

		<div class="fs-panel">
			<div class="fs-bar">
				<span class="fs-a fs-a--s">A</span>
				<div class="fs-ctrl">
					<div class="fs-names">
						<div v-for="(item, i) in levels" :key="'name-' + i" class="fs-name">
							<span v-if="i === standardIndex">标准</span>
						</div>
					</div>
					<div
						class="fs-stops"
						@touchstart="onTouch"
						@touchmove.stop.prevent="onTouch"
						@mousedown="onMouseDown"
					>
						<div class="fs-stops__line" :style="lineStyle"></div>
						<div
							v-for="(item, i) in levels"
							:key="'stop-' + i"
							class="fs-stop"
							@click.stop="pick(i)"
						>
							<div class="fs-stop__tick"></div>
							<div v-if="i === levelIndex" class="fs-stop__thumb"></div>
						</div>
					</div>
				</div>
				<span class="fs-a fs-a--l">A</span>
			</div>
			<div class="fs-done" :style="{ background: themeColor }" @click="confirm">完成</div>
		</div>
	</div>
</template>

<script>
import { createSelectorQuery } from '@/common/dom'
import { navigateBack } from '@/common/nav'
import config from '@/config'
import {
	FONT_LEVELS,
	FONT_STANDARD_INDEX,
	readFontScale,
	saveFontScale,
	findLevelIndex
} from '@/common/fontScale'

export default {
	data() {
		return {
			levels: FONT_LEVELS,
			standardIndex: FONT_STANDARD_INDEX,
			levelIndex: findLevelIndex(readFontScale()),
			previewTab: 'guide',
			themeColor: config.primaryColor || '#4286F5'
		}
	},
	computed: {
		previewScale() {
			return this.levels[this.levelIndex].scale
		},
		lineStyle() {
			const inset = 100 / this.levels.length / 2
			return { left: inset + '%', right: inset + '%' }
		},
		pageStyle() {
			return {
				'--font-scale': String(this.previewScale),
				'--fs-theme': this.themeColor
			}
		}
	},
	methods: {
		pick(index) {
			this.levelIndex = index
		},
		goBack() {
			navigateBack(1)
		},
		confirm() {
			saveFontScale(this.previewScale)
			this.goBack()
		},
		pointTo(clientX) {
			createSelectorQuery().in(this).select('.fs-stops').boundingClientRect((rect) => {
				if (!rect || !rect.width) return
				let ratio = (clientX - rect.left) / rect.width
				if (ratio < 0) ratio = 0
				if (ratio > 0.999) ratio = 0.999
				this.levelIndex = Math.floor(ratio * this.levels.length)
			}).exec()
		},
		onTouch(e) {
			const touch = e.touches && e.touches[0]
			if (!touch) return
			this.pointTo(touch.clientX)
		},
		onMouseDown(e) {
			if (!e || typeof window === 'undefined') return
			this.pointTo(e.clientX)
			const move = (ev) => this.pointTo(ev.clientX)
			const up = () => {
				window.removeEventListener('mousemove', move)
				window.removeEventListener('mouseup', up)
			}
			window.addEventListener('mousemove', move)
			window.addEventListener('mouseup', up)
		}
	}
}
</script>

<style lang="scss" scoped>
.fs-page {
	height: calc(100vh - var(--window-top, 0px));
	display: flex;
	flex-direction: column;
	background: #fff;
	box-sizing: border-box;
	overflow: hidden;
}

.fs-panel {
	--font-scale: 1;
}

.fs-preview {
	flex: 1;
	height: 0;
	min-height: 0;
	display: flex;
	flex-direction: column;
	overflow: hidden;
	background: #fff;
}

.pv {
	flex: 1;
	height: 0;
	min-height: 0;
	overflow-x: hidden;
	overflow-y: auto;
	padding-top: 16px;
	padding-bottom: calc(160px + env(safe-area-inset-bottom));
}

.pv-brief {
	display: flex;
	align-items: flex-start;
	padding: 0 20px 16px;
}

.pv-brief__left {
	flex: 1;
	min-width: 0;
	padding-right: 16px;
}

.pv-brief__category {
	display: inline-block;
	margin-right: 6px;
	padding: 0 4px;
	border: 1px solid rgba(250, 137, 25, 0.35);
	border-radius: 4px;
	font-size: calc(13px * var(--font-scale, 1));
	font-weight: 500;
	line-height: 1.4;
	color: #fa8919;
	vertical-align: middle;
}

.pv-brief__title {
	margin-bottom: 8px;
	font-size: calc(18px * var(--font-scale, 1));
	font-weight: 500;
	line-height: 1.4;
	color: #353535;
	word-break: break-word;
}

.pv-brief__subtitle {
	font-size: calc(13px * var(--font-scale, 1));
	line-height: 1.5;
	color: #7f848c;
	word-break: break-word;
}

.pv-brief__cover {
	position: relative;
	flex-shrink: 0;
	width: 90px;
	height: 90px;
	border-radius: 4px;
	background: linear-gradient(160deg, #7eb0ff 0%, var(--fs-theme) 100%);
	overflow: hidden;
}

.pv-brief__cover-shade {
	position: absolute;
	left: 0;
	right: 0;
	bottom: 0;
	min-height: 18px;
	padding: 2px 6px;
	display: flex;
	align-items: center;
	justify-content: flex-end;
	background: linear-gradient(180deg, transparent, rgba(0, 0, 0, 0.45));
	box-sizing: border-box;
}

.pv-brief__views {
	font-size: calc(12px * var(--font-scale, 1));
	line-height: 1.2;
	color: #fff;
}

.pv-stats {
	display: flex;
	align-items: flex-start;
	padding: 12px 10px;
	border-bottom: 8px solid #f6f7fb;
}

.pv-stats__item {
	flex: 1;
	text-align: center;
}

.pv-stats__label {
	font-size: calc(13px * var(--font-scale, 1));
	font-weight: 500;
	line-height: 1.4;
	color: #404040;
	word-break: break-word;
}

.pv-stats__tip {
	margin-top: 2px;
	font-size: calc(12px * var(--font-scale, 1));
	line-height: 1.4;
	color: #7f848c;
	word-break: break-word;
}

.pv-cert {
	margin: 12px 20px 0;
	padding: 10px 12px;
	border-radius: 6px;
	background: #f0f7ff;
	color: #3b6fd9;
	font-size: calc(13px * var(--font-scale, 1));
	line-height: 1.5;
}

.pv-tabs {
	display: flex;
	justify-content: center;
	gap: 48px;
	padding: 16px 20px 8px;
}

.pv-tabs__item {
	position: relative;
	font-size: calc(17px * var(--font-scale, 1));
	line-height: 1.4;
	color: #7f848c;
	cursor: pointer;
}

.pv-tabs__item--on {
	font-weight: 500;
	color: #353535;
}

.pv-tabs__item--on::after {
	content: '';
	position: absolute;
	left: 50%;
	bottom: -6px;
	width: 22px;
	height: 3px;
	margin-left: -11px;
	border-radius: 2px;
	background: var(--fs-theme);
}

.pv-intro {
	padding: 18px 20px 8px;
}

.pv-intro__title {
	position: relative;
	margin-bottom: 10px;
	font-size: calc(17px * var(--font-scale, 1));
	font-weight: 500;
	line-height: 1.4;
	color: #353535;
}

.pv-intro__title::before {
	content: '';
	position: absolute;
	top: -2px;
	left: -8px;
	width: 16px;
	height: 16px;
	border-radius: 50%;
	background: linear-gradient(1turn, #f3f8ff, #b3cefb);
	opacity: 0.5;
}

.pv-intro__text {
	font-size: calc(15px * var(--font-scale, 1));
	line-height: 1.75;
	color: #404040;
	word-break: break-word;
}

.pv-catalog {
	padding: 8px 0 16px;
}

.pv-catalog__head {
	padding: 8px 20px;
	font-size: calc(14px * var(--font-scale, 1));
	font-weight: 500;
	line-height: 1.4;
	color: #404040;
}

.pv-chapter {
	display: flex;
	align-items: flex-start;
	padding: 12px 16px 12px 20px;
	border-top: 1px solid #f6f7fb;
}

.pv-chapter__mark {
	flex-shrink: 0;
	width: 16px;
	height: 16px;
	margin-top: calc(0.15 * 14px * var(--font-scale, 1));
	margin-right: 10px;
	border-radius: 50%;
	border: 2px solid #d5d8de;
	box-sizing: border-box;
}

.pv-chapter__mark--on {
	border-color: var(--fs-theme);
	background: var(--fs-theme);
	box-shadow: inset 0 0 0 3px #fff;
}

.pv-chapter__info {
	flex: 1;
	min-width: 0;
}

.pv-chapter__title {
	font-size: calc(14px * var(--font-scale, 1));
	line-height: 1.45;
	color: #404040;
	word-break: break-word;
}

.pv-chapter--current .pv-chapter__title {
	color: var(--fs-theme);
}

.pv-chapter__meta {
	margin-top: 4px;
	font-size: calc(12px * var(--font-scale, 1));
	line-height: 1.45;
	color: #7f848c;
	word-break: break-word;
}

.pv-chapter__status {
	flex-shrink: 0;
	margin-left: 8px;
	padding: 3px 10px;
	border-radius: 12px;
	background: #f4f5f7;
	font-size: calc(12px * var(--font-scale, 1));
	line-height: 1.3;
	color: #9aa0a6;
	white-space: nowrap;
}

.pv-chapter__status--on {
	background: #f3f8ff;
	color: var(--fs-theme);
}

.fs-panel {
	position: fixed;
	left: 0;
	right: 0;
	bottom: 0;
	z-index: 20;
	background: #fff;
	border-top: 8px solid #f6f7fb;
	padding: 8px 16px calc(12px + env(safe-area-inset-bottom));
	user-select: none;
}

.fs-bar {
	display: flex;
	align-items: flex-end;
	padding: 4px 0 12px;
}

.fs-a {
	flex-shrink: 0;
	min-width: 28px;
	text-align: center;
	color: #353535;
	line-height: 1;
	padding-bottom: 8px;
}

.fs-a--s {
	font-size: calc(16px * var(--font-scale, 1));
}

.fs-a--l {
	font-size: calc(26px * var(--font-scale, 1));
}

.fs-ctrl {
	flex: 1;
	min-width: 0;
	padding: 0 6px;
}

.fs-names,
.fs-stops {
	display: flex;
}

.fs-name {
	flex: 1;
	text-align: center;
	height: 22px;
	font-size: calc(14px * var(--font-scale, 1));
	font-weight: 500;
	color: #353535;
	line-height: 22px;
}

.fs-stops {
	position: relative;
	height: 44px;
	align-items: center;
	touch-action: none;
	cursor: pointer;
}

.fs-stops__line {
	position: absolute;
	top: 50%;
	height: 1px;
	background: #e4e6eb;
}

.fs-stop {
	flex: 1;
	position: relative;
	height: 44px;
	display: flex;
	align-items: center;
	justify-content: center;
}

.fs-stop__tick {
	width: 1px;
	height: 8px;
	background: #d5d8de;
}

.fs-stop__thumb {
	position: absolute;
	left: 50%;
	top: 50%;
	width: 28px;
	height: 28px;
	margin-left: -14px;
	margin-top: -14px;
	border-radius: 50%;
	background: #fff;
	border: 2px solid var(--fs-theme);
	box-sizing: border-box;
	box-shadow: 0 1px 4px rgba(66, 133, 244, 0.25);
	z-index: 1;
}

.fs-done {
	height: 44px;
	border-radius: 8px;
	display: flex;
	align-items: center;
	justify-content: center;
	color: #fff;
	font-size: calc(16px * var(--font-scale, 1));
	font-weight: 500;
	line-height: 22px;
}
</style>
