<template>
	<div v-if="visible" class="cap-mask" @click="onMask">
		<div class="cap-box" :style="themeStyle" @click.stop>
			<div class="cap-hd">
				<span>请完成安全验证</span>
				<span class="cap-close" @click.stop="dismiss">×</span>
			</div>
			<div class="cap-img" :style="{ width: imgW + 'px', height: imgH + 'px' }">
				<img v-if="bgSrc" class="cap-bg" :src="bgSrc" alt="" draggable="false" />
				<img
					v-if="pieceSrc"
					class="cap-piece"
					:src="pieceSrc"
					alt=""
					draggable="false"
					:style="{
						width: pieceW + 'px',
						height: pieceH + 'px',
						left: left + 'px',
						top: pieceTop + 'px',
					}"
				/>
				<div class="cap-refresh" @click="load">↻</div>
			</div>
			<div
				class="cap-bar"
				:style="{ width: imgW + 'px' }"
				@touchstart="onStart"
				@touchmove.stop.prevent="onMove"
				@touchend="onEnd"
				@mousedown="onMouseDown"
			>
				<div class="cap-fill" :class="state" :style="{ width: (left > 0 ? left : 0) + 'px' }" />
				<span v-if="left === startLeft && !state" class="cap-tip">向右滑动完成验证</span>
				<div class="cap-knob" :class="state" :style="{ left: left + 'px' }">
					<span>{{ state === 'ok' ? '✓' : state === 'fail' ? '×' : '›' }}</span>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
import { showToast } from '@/common/toast'
	import { http } from '@/common/request.js'
	import config from '@/config.js'

	const IMG_W = 300
	const IMG_H = 220
	const KNOB = 40

	function asDataUrl(raw) {
		if (!raw) return ''
		if (String(raw).startsWith('data:')) return raw
		return 'data:image/png;base64,' + raw
	}

	export default {
		computed: {
			themeStyle() {
				const c = config.primaryColor || '#4286F5'
				const rgb = config.primaryRgb || '66,134,245'
				return {
					'--cap-primary': c,
					'--cap-primary-rgb': rgb,
				}
			},
		},
		data() {
			return {
				visible: false,
				imgW: IMG_W,
				imgH: IMG_H,
				knob: KNOB,
				pieceW: 47,
				pieceH: 47,
				pieceTop: 0,
				startLeft: 0,
				bgSrc: '',
				pieceSrc: '',
				token: '',
				left: 0,
				state: '',
				dragging: false,
				startX: 0,
				dragOrigin: 0,
				checking: false,
				openedAt: 0,
			}
		},
		beforeUnmount() {
			this.unbindMouse()
		},
		methods: {
			open() {
				this.openedAt = Date.now()
				this.visible = true
				this.checking = false
				this.state = ''
				this.left = 0
				this.load()
			},
			onMask() {
				if (Date.now() - this.openedAt < 500) return
				this.dismiss()
			},
			dismiss() {
				this.visible = false
				this.dragging = false
				this.unbindMouse()
			},
			load() {
				this.left = 0
				this.state = ''
				this.dragging = false
				this.checking = false
				http.get('/auth/captcha').then((res) => {
					const data = res.data || {}
					if (!data.image || !data.thumb || !data.token) {
						showToast('验证码加载失败')
						return
					}
					this.bgSrc = asDataUrl(data.image)
					this.pieceSrc = asDataUrl(data.thumb)
					this.token = data.token
					this.pieceW = Number(data.thumbWidth) || 47
					this.pieceH = Number(data.thumbHeight) || 47
					this.pieceTop = Number(data.thumbY) || 0
					this.startLeft = Number(data.thumbX) || 0
					this.left = this.startLeft
				}).catch(() => {
					showToast('验证码加载失败')
				})
			},
			maxLeft() {
				return Math.max(0, this.imgW - this.pieceW)
			},
			eventX(e) {
				if (!e) return null
				const touch = (e.touches && e.touches[0]) || (e.changedTouches && e.changedTouches[0])
				if (touch && typeof touch.clientX === 'number') return touch.clientX
				if (typeof e.clientX === 'number') return e.clientX
				if (e.detail && typeof e.detail.x === 'number') return e.detail.x
				return null
			},
			onStart(e) {
				if (this.state === 'ok') return
				const x = this.eventX(e)
				if (x == null) return
				this.dragging = true
				this.startX = x
				this.dragOrigin = this.left
				this.state = 'drag'
			},
			onMove(e) {
				if (!this.dragging) return
				const x = this.eventX(e)
				if (x == null) return
				this.applyMove(x)
			},
			onEnd() {
				if (!this.dragging) return
				this.dragging = false
				this.unbindMouse()
				if (Math.abs(this.left - this.startLeft) <= 2) {
					this.left = this.startLeft
					this.state = ''
					return
				}
				this.verify()
			},
			onMouseDown(e) {
				this.onStart(e)
				if (typeof window === 'undefined') return
				window.addEventListener('mousemove', this.onMouseMove)
				window.addEventListener('mouseup', this.onMouseUp)
			},
			onMouseMove(e) {
				this.onMove(e)
			},
			onMouseUp() {
				this.unbindMouse()
				this.onEnd()
			},
			unbindMouse() {
				if (typeof window === 'undefined') return
				window.removeEventListener('mousemove', this.onMouseMove)
				window.removeEventListener('mouseup', this.onMouseUp)
			},
			applyMove(x) {
				let next = this.dragOrigin + (x - this.startX)
				next = Math.max(0, Math.min(this.maxLeft(), next))
				this.left = next
			},
			verify() {
				if (this.checking || !this.token) return
				this.checking = true
				http.post('/auth/captchaCheck', {
					token: this.token,
					x: Math.round(this.left),
					y: Math.round(this.pieceTop),
				}).then((res) => {
					const code = res.data && res.data.captchaVerification
					if (!code) {
						this.state = 'fail'
						this.checking = false
						setTimeout(() => this.load(), 700)
						return
					}
					this.state = 'ok'
					setTimeout(() => {
						this.dismiss()
						this.$emit('success', code)
					}, 400)
				}).catch(() => {
					this.state = 'fail'
					this.checking = false
					setTimeout(() => this.load(), 700)
				})
			},
		},
	}
</script>

<style scoped>
	.cap-mask {
		position: fixed;
		left: 0;
		top: 0;
		right: 0;
		bottom: 0;
		z-index: 1000;
		background: rgba(18, 28, 30, 0.55);
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.cap-box {
		background: #fff;
		border-radius: 16px;
		padding: 0 15px 16px;
		box-sizing: border-box;
	}

	.cap-hd {
		height: 48px;
		display: flex;
		align-items: center;
		justify-content: space-between;
		font-size: 16px;
		font-weight: 600;
		color: #1c2b2e;
	}

	.cap-close {
		width: 28px;
		text-align: center;
		font-size: 22px;
		color: #5a6b6d;
	}

	.cap-img {
		position: relative;
		border-radius: 12px;
		overflow: hidden;
		background: #e7eeef;
	}

	.cap-bg,
	.cap-piece {
		display: block;
		pointer-events: none;
		-webkit-user-drag: none;
	}

	.cap-bg {
		width: 100%;
		height: 100%;
	}

	.cap-piece {
		position: absolute;
	}

	.cap-refresh {
		position: absolute;
		top: 8px;
		right: 8px;
		width: 28px;
		height: 28px;
		border-radius: 8px;
		background: rgba(18, 28, 30, 0.45);
		color: #fff;
		text-align: center;
		line-height: 28px;
		font-size: 16px;
	}

	.cap-bar {
		position: relative;
		height: 40px;
		margin-top: 12px;
		background: #e7eeef;
		border-radius: 12px;
		overflow: hidden;
	}

	.cap-fill {
		position: absolute;
		left: 0;
		top: 0;
		bottom: 0;
		background: rgba(var(--cap-primary-rgb, 66, 134, 245), 0.55);
	}

	.cap-fill.ok {
		background: linear-gradient(90deg, #e5f6ee 0%, #9dcfb8 58%, #3d8b74 100%);
	}

	.cap-fill.fail {
		background: linear-gradient(90deg, #fbecea 0%, #efc0b4 58%, #c97b63 100%);
	}

	.cap-tip {
		position: absolute;
		left: 0;
		right: 0;
		top: 0;
		height: 40px;
		line-height: 40px;
		text-align: center;
		font-size: 14px;
		color: #666;
	}

	.cap-knob {
		position: absolute;
		top: 0;
		width: 40px;
		height: 40px;
		border-radius: 10px;
		background: #fff;
		border: 1px solid rgba(var(--cap-primary-rgb, 66, 134, 245), 0.2);
		box-sizing: border-box;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 20px;
		color: var(--cap-primary, #4286f5);
		transition: background 0.3s ease, border-radius 0.3s ease, color 0.3s ease;
	}

	.cap-knob.drag {
		background: var(--cap-primary, #4286f5);
		color: #fff;
		border-color: var(--cap-primary, #4286f5);
		border-radius: 0 10px 10px 0;
	}

	.cap-knob.ok {
		background: #3d8b74;
		border-color: #3d8b74;
		color: #fff;
	}

	.cap-knob.fail {
		background: #c97b63;
		border-color: #c97b63;
		color: #fff;
	}
</style>
