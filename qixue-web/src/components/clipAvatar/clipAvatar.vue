<template>
	<teleport to="body">
		<div v-if="show" class="avatar-crop">
			<div
				ref="stage"
				class="avatar-crop__stage"
				@pointerdown.prevent="onDown"
				@pointermove.prevent="onMove"
				@pointerup.prevent="onUp"
				@pointercancel.prevent="onUp"
				@wheel.prevent="onWheel"
			>
				<img
					ref="imgEl"
					class="avatar-crop__img"
					:src="src"
					:style="imgStyle"
					draggable="false"
					alt=""
					@load="onLoad"
				/>
				<div class="avatar-crop__veil" :style="veilStyle"></div>
			</div>
			<div class="avatar-crop__bar">
				<button type="button" class="avatar-crop__btn" @click="cancel">取消</button>
				<button type="button" class="avatar-crop__btn avatar-crop__btn--ok" @click="confirm">确定</button>
			</div>
			<canvas ref="canvas" class="avatar-crop__canvas"></canvas>
		</div>
	</teleport>
</template>

<script>
const OUTPUT = 400

export default {
	name: 'clipAvatar',
	emits: ['cropImg'],
	data() {
		return {
			show: false,
			src: '',
			naturalW: 0,
			naturalH: 0,
			scale: 1,
			minScale: 1,
			offsetX: 0,
			offsetY: 0,
			hole: 260,
			dragging: false,
			lastX: 0,
			lastY: 0,
		}
	},
	computed: {
		veilStyle() {
			const s = this.hole
			return {
				width: `${s}px`,
				height: `${s}px`,
				marginLeft: `${-s / 2}px`,
				marginTop: `${-s / 2}px`,
				boxShadow: '0 0 0 9999px rgba(0,0,0,.62)',
			}
		},
		imgStyle() {
			const w = this.naturalW * this.scale
			const h = this.naturalH * this.scale
			return {
				width: `${w}px`,
				height: `${h}px`,
				transform: `translate(${this.offsetX}px, ${this.offsetY}px)`,
			}
		},
	},
	methods: {
		open(src) {
			this.src = src || ''
			this.show = true
			this.naturalW = 0
			this.naturalH = 0
			this.$nextTick(() => this.fitHole())
		},
		fitHole() {
			const stage = this.$refs.stage
			if (!stage) return
			const size = Math.min(stage.clientWidth, stage.clientHeight)
			this.hole = Math.max(180, Math.floor(size * 0.72))
		},
		onLoad(e) {
			const img = e.target
			this.naturalW = img.naturalWidth || 1
			this.naturalH = img.naturalHeight || 1
			this.fitHole()
			const cover = Math.max(this.hole / this.naturalW, this.hole / this.naturalH)
			this.minScale = cover
			this.scale = cover
			const stage = this.$refs.stage
			const sw = stage ? stage.clientWidth : window.innerWidth
			const sh = stage ? stage.clientHeight : window.innerHeight
			this.offsetX = (sw - this.naturalW * this.scale) / 2
			this.offsetY = (sh - this.naturalH * this.scale) / 2
			this.clamp()
		},
		clamp() {
			const stage = this.$refs.stage
			if (!stage || !this.naturalW) return
			const sw = stage.clientWidth
			const sh = stage.clientHeight
			const w = this.naturalW * this.scale
			const h = this.naturalH * this.scale
			const cx = sw / 2
			const cy = sh / 2
			const half = this.hole / 2
			const minX = cx + half - w
			const maxX = cx - half
			const minY = cy + half - h
			const maxY = cy - half
			this.offsetX = Math.min(maxX, Math.max(minX, this.offsetX))
			this.offsetY = Math.min(maxY, Math.max(minY, this.offsetY))
		},
		onDown(e) {
			this.dragging = true
			this.lastX = e.clientX
			this.lastY = e.clientY
			e.currentTarget.setPointerCapture?.(e.pointerId)
		},
		onMove(e) {
			if (!this.dragging) return
			this.offsetX += e.clientX - this.lastX
			this.offsetY += e.clientY - this.lastY
			this.lastX = e.clientX
			this.lastY = e.clientY
			this.clamp()
		},
		onUp() {
			this.dragging = false
		},
		onWheel(e) {
			if (!this.naturalW) return
			const stage = this.$refs.stage
			if (!stage) return
			const rect = stage.getBoundingClientRect()
			const mx = e.clientX - rect.left
			const my = e.clientY - rect.top
			const prev = this.scale
			const next = Math.min(8, Math.max(this.minScale, prev * (e.deltaY > 0 ? 0.92 : 1.08)))
			if (next === prev) return
			this.offsetX = mx - ((mx - this.offsetX) * next) / prev
			this.offsetY = my - ((my - this.offsetY) * next) / prev
			this.scale = next
			this.clamp()
		},
		releaseSrc() {
			if (this.src && this.src.startsWith('blob:')) URL.revokeObjectURL(this.src)
			this.src = ''
			this.show = false
		},
		cancel() {
			this.releaseSrc()
		},
		confirm() {
			const canvas = this.$refs.canvas
			const img = this.$refs.imgEl
			const stage = this.$refs.stage
			if (!canvas || !img || !stage || !this.naturalW) return
			canvas.width = OUTPUT
			canvas.height = OUTPUT
			const ctx = canvas.getContext('2d')
			const sw = stage.clientWidth
			const sh = stage.clientHeight
			const half = this.hole / 2
			const left = sw / 2 - half
			const top = sh / 2 - half
			const sx = (left - this.offsetX) / this.scale
			const sy = (top - this.offsetY) / this.scale
			const sSize = this.hole / this.scale
			ctx.clearRect(0, 0, OUTPUT, OUTPUT)
			ctx.drawImage(img, sx, sy, sSize, sSize, 0, 0, OUTPUT, OUTPUT)
			canvas.toBlob((blob) => {
				if (!blob) return
				const file = new File([blob], 'avatar.jpg', { type: 'image/jpeg' })
				this.$emit('cropImg', file)
				this.releaseSrc()
			}, 'image/jpeg', 0.92)
		},
	},
}
</script>

<style scoped>
.avatar-crop {
	position: fixed;
	inset: 0;
	z-index: 10000;
	display: flex;
	flex-direction: column;
	background: #000;
}
.avatar-crop__stage {
	position: relative;
	flex: 1;
	overflow: hidden;
	touch-action: none;
	cursor: grab;
}
.avatar-crop__stage:active {
	cursor: grabbing;
}
.avatar-crop__img {
	position: absolute;
	left: 0;
	top: 0;
	max-width: none;
	user-select: none;
	-webkit-user-drag: none;
	will-change: transform;
}
.avatar-crop__veil {
	position: absolute;
	left: 50%;
	top: 50%;
	border: 1px solid rgba(255, 255, 255, 0.9);
	border-radius: 50%;
	box-sizing: border-box;
	pointer-events: none;
	z-index: 2;
}
.avatar-crop__bar {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 16px 20px calc(16px + env(safe-area-inset-bottom));
	background: #111;
}
.avatar-crop__btn {
	min-width: 72px;
	height: 36px;
	padding: 0 16px;
	border: 0;
	border-radius: 8px;
	background: transparent;
	color: #fff;
	font-size: 15px;
	cursor: pointer;
}
.avatar-crop__btn--ok {
	background: #4286f5;
}
.avatar-crop__canvas {
	position: fixed;
	left: -9999px;
	top: 0;
	width: 1px;
	height: 1px;
	opacity: 0;
	pointer-events: none;
}
</style>
