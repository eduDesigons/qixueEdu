<template>
	<div class="feedback-page" :class="{ 'has-footer': composing || ready }">
		<div class="feedback-loading" v-if="!composing && !ready">加载中…</div>
		<div class="feedback-form" v-if="composing">
			<div class="feedback-head">
				<span class="feedback-title">问题或建议</span>
				<span class="feedback-required">必填</span>
				<span class="feedback-cancel" v-if="list.length" @click="composing = false">返回</span>
			</div>
			<div class="feedback-box">
				<textarea
					class="feedback-input"
					v-model="content"
					maxlength="500"
					auto-height
					:show-confirm-bar="false"
					placeholder="请描述遇到的问题，或想提的建议"
					placeholder-class="feedback-placeholder"
				/>
				<span class="feedback-count">{{ content.length }}/500</span>
			</div>
			<div class="feedback-pic-head">
				<span>图片</span>
				<span class="feedback-pic-tip">选填，最多 3 张</span>
			</div>
			<div class="feedback-images">
				<div class="feedback-image" v-for="(url, index) in images" :key="url">
					<img :src="static(url)" @click="preview(index)"  style="object-fit:cover" />
					<div class="feedback-image__del" @click.stop="removeImage(index)">
						<van-icon name="cross" size="10" color="#fff" />
					</div>
				</div>
				<div class="feedback-add" v-if="images.length < 3" @click="chooseImage">
					<van-icon name="plus" size="22" color="#b9bdc2" />
				</div>
			</div>
		</div>

		<div class="feedback-history" v-else-if="ready">
			<div class="feedback-history__bar">
				<span class="feedback-history__title">我的反馈</span>
			</div>
			<div class="feedback-item" v-for="item in list" :key="item.id">
				<div class="feedback-item__head">
					<span class="feedback-item__time">{{ shortTime(item.created_at) }}</span>
					<span class="feedback-item__status" :class="{ done: Number(item.status) === 1 }">
						{{ Number(item.status) === 1 ? '已回复' : '待回复' }}
					</span>
				</div>
				<div class="feedback-item__content">{{ item.content }}</div>
				<div class="feedback-item__pics" v-if="item.images && item.images.length">
					<img
						v-for="(url, index) in item.images"
						:key="url"
						:src="static(url)"
						@click="previewItem(item, index)"
					 style="object-fit:cover" />
				</div>
				<div class="feedback-item__reply" v-if="item.reply">
					<span class="feedback-item__reply-label">回复</span>
					<span class="feedback-item__reply-text">{{ item.reply }}</span>
				</div>
			</div>
		</div>

		<div class="feedback-footer" v-if="composing">
			<div class="feedback-submit" :class="{ disabled: submitting }" @click="submit">
				{{ submitting ? '提交中...' : '提交' }}
			</div>
		</div>
		<div class="feedback-footer" v-else-if="ready">
			<div class="feedback-submit" @click="openForm">提交新的意见建议</div>
		</div>
	</div>
</template>

<script>
import { pickImages, previewImages } from '@/common/media'

export default {
	data() {
		return {
			content: '',
			images: [],
			list: [],
			loading: true,
			submitting: false,
			composing: false,
			ready: false
		}
	},
	onShow() {
		this.loadList()
	},
	methods: {
		loadList() {
			this.loading = true
			this.$http.get('/feedback/list').then(res => {
				this.list = res.data || []
				if (!this.ready) {
					this.composing = this.list.length === 0
					this.ready = true
				}
			}).catch(() => {
				this.msg('加载失败')
				if (!this.ready) {
					this.composing = true
					this.ready = true
				}
			}).finally(() => {
				this.loading = false
			})
		},
		chooseImage() {
			const remain = 3 - this.images.length
			if (remain <= 0) return
			pickImages({ count: remain }).then((list) => {
				list.forEach((item) => {
					this.$http.upload(item.file).then((ret) => {
						const url = ret && ret.data && (ret.data.url || ret.data)
						if (url && this.images.length < 3) this.images.push(url)
						URL.revokeObjectURL(item.url)
					}).catch(() => {
						URL.revokeObjectURL(item.url)
					})
				})
			}).catch(() => {})
		},
		removeImage(index) {
			this.images.splice(index, 1)
		},
		preview(index) {
			const urls = this.images.map((url) => this.static(url))
			previewImages(urls, index)
		},
		previewItem(item, index) {
			const urls = (item.images || []).map((url) => this.static(url))
			previewImages(urls, index)
		},
		openForm() {
			this.content = ''
			this.images = []
			this.composing = true
		},
		shortTime(value) {
			const text = String(value || '')
			if (text.length >= 16) return text.slice(5, 16)
			return text
		},
		submit() {
			if (this.submitting) return
			const content = (this.content || '').trim()
			if (!content) {
				this.msg('请填写内容')
				return
			}
			this.submitting = true
			this.$http.post('/feedback/submit', {
				content,
				images: this.images
			}).then(() => {
				this.msg('提交成功')
				this.content = ''
				this.images = []
				this.composing = false
				this.ready = true
				this.loadList()
			}).catch((err) => {
				this.msg((err && err.message) || '提交失败')
			}).finally(() => {
				this.submitting = false
			})
		}
	}
}
</script>

<style lang="scss" scoped>
.feedback-page {
	min-height: 100vh;
	box-sizing: border-box;
	padding: calc(24 * var(--rpx)) calc(30 * var(--rpx)) calc(calc(32 * var(--rpx)) + env(safe-area-inset-bottom));
	background: #f8f9fa;
}

.feedback-page.has-footer {
	padding-bottom: calc(calc(160 * var(--rpx)) + env(safe-area-inset-bottom));
}

.feedback-loading {
	padding: calc(80 * var(--rpx)) 0;
	text-align: center;
	color: #9ca3af;
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
}

.feedback-form {
	padding: calc(32 * var(--rpx));
	background: #fff;
	border-radius: calc(20 * var(--rpx));
	box-shadow: 0 calc(4 * var(--rpx)) calc(16 * var(--rpx)) rgba(0, 0, 0, 0.04);
}

.feedback-head {
	display: flex;
	align-items: center;
	margin-bottom: calc(24 * var(--rpx));
}

.feedback-cancel {
	margin-left: auto;
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
}

.feedback-title {
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	color: #1a1a1a;
}

.feedback-required {
	margin-left: calc(12 * var(--rpx));
	padding: calc(2 * var(--rpx)) calc(10 * var(--rpx));
	border-radius: calc(6 * var(--rpx));
	background: #fee;
	color: #f56565;
	font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
}

.feedback-box {
	position: relative;
	padding: calc(24 * var(--rpx)) calc(24 * var(--rpx)) calc(56 * var(--rpx));
	background: #f8f9fa;
	border: calc(2 * var(--rpx)) solid #e9ecef;
	border-radius: calc(16 * var(--rpx));
}

.feedback-input {
	width: 100%;
	min-height: calc(200 * var(--rpx));
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	line-height: 1.7;
	color: #1a1a1a;
	background: transparent;
}

.feedback-placeholder {
	color: #adb5bd;
}

.feedback-count {
	position: absolute;
	right: calc(24 * var(--rpx));
	bottom: calc(16 * var(--rpx));
	font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
	color: #adb5bd;
}

.feedback-pic-head {
	display: flex;
	align-items: baseline;
	margin-top: calc(28 * var(--rpx));
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #1a1a1a;
}

.feedback-pic-tip {
	margin-left: calc(12 * var(--rpx));
	font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 400;
	color: #adb5bd;
}

.feedback-images {
	display: flex;
	flex-wrap: wrap;
	gap: calc(16 * var(--rpx));
	margin-top: calc(20 * var(--rpx));
}

.feedback-image,
.feedback-add {
	position: relative;
	width: calc(148 * var(--rpx));
	height: calc(148 * var(--rpx));
	border-radius: calc(16 * var(--rpx));
	overflow: hidden;
}

.feedback-image img {
	width: 100%;
	height: 100%;
}

.feedback-add {
	display: flex;
	align-items: center;
	justify-content: center;
	background: #f8f9fa;
	border: calc(2 * var(--rpx)) dashed #d7dde3;
	box-sizing: border-box;
}

.feedback-image__del {
	position: absolute;
	top: calc(8 * var(--rpx));
	right: calc(8 * var(--rpx));
	width: calc(32 * var(--rpx));
	height: calc(32 * var(--rpx));
	border-radius: 50%;
	background: rgba(0, 0, 0, 0.45);
	display: flex;
	align-items: center;
	justify-content: center;
}

.feedback-history {
	margin-top: calc(28 * var(--rpx));
}

.feedback-history__bar {
	margin: 0 calc(8 * var(--rpx)) calc(16 * var(--rpx));
}

.feedback-history__title {
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	color: #353535;
}

.feedback-empty {
	padding: calc(48 * var(--rpx)) 0;
	text-align: center;
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	color: #adb5bd;
}

.feedback-item {
	margin-bottom: calc(16 * var(--rpx));
	padding: calc(28 * var(--rpx));
	background: #fff;
	border-radius: calc(20 * var(--rpx));
}

.feedback-item__head {
	display: flex;
	align-items: center;
	justify-content: space-between;
}

.feedback-item__time {
	font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
	color: #b9bdc2;
}

.feedback-item__status {
	padding: calc(4 * var(--rpx)) calc(14 * var(--rpx));
	border-radius: calc(20 * var(--rpx));
	background: #fff7e6;
	color: #d48806;
	font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
}

.feedback-item__status.done {
	background: #f6ffed;
	color: #389e0d;
}

.feedback-item__content {
	margin-top: calc(16 * var(--rpx));
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	line-height: 1.7;
	color: #1a1a1a;
	white-space: pre-wrap;
	word-break: break-word;
}

.feedback-item__pics {
	display: flex;
	gap: calc(12 * var(--rpx));
	margin-top: calc(16 * var(--rpx));
}

.feedback-item__pics img {
	width: calc(132 * var(--rpx));
	height: calc(132 * var(--rpx));
	border-radius: calc(12 * var(--rpx));
}

.feedback-item__reply {
	margin-top: calc(20 * var(--rpx));
	padding: calc(16 * var(--rpx)) calc(20 * var(--rpx));
	border-left: calc(6 * var(--rpx)) solid #4286F5;
	border-radius: 0 calc(12 * var(--rpx)) calc(12 * var(--rpx)) 0;
	background: #f4f8ff;
}

.feedback-item__reply-label {
	display: block;
	margin-bottom: calc(6 * var(--rpx));
	font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	color: #4286F5;
}

.feedback-item__reply-text {
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	line-height: 1.7;
	color: #404040;
	white-space: pre-wrap;
	word-break: break-word;
}

.feedback-footer {
	position: fixed;
	left: 0;
	right: 0;
	bottom: 0;
	z-index: 10;
	padding: calc(20 * var(--rpx)) calc(32 * var(--rpx)) calc(calc(20 * var(--rpx)) + env(safe-area-inset-bottom));
	background: #fff;
	border-top: calc(1 * var(--rpx)) solid #e9ecef;
}

.feedback-submit {
	height: calc(88 * var(--rpx));
	border-radius: calc(44 * var(--rpx));
	background: #4286F5;
	color: #fff;
	font-size: calc(calc(30 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	display: flex;
	align-items: center;
	justify-content: center;
}

.feedback-submit.disabled {
	opacity: 0.6;
}
</style>
