<template>
	<van-popup
		:show="show"
	 position="bottom"
	 :close-on-click-overlay="true"
		round="calc(40 * var(--rpx))"
		@update:show="(v)=>{if(!v)handleClose()}"
	 :safe-area-inset-bottom="true"
		:overlay="true"
	>
		<div class="comment-input-popup" :class="{ 'is-expanded': isExpand }" :style="popupStyle">
			<div class="comment-input-container">
				<div class="comment-input-header">
					<div class="comment-input-expand" @click.stop="handleExpand">
						<img src="/static/img/closepand.png" width="18" height="18" alt="" style="width:18px;height:18px;" v-if="isExpand" />
						<img src="/static/img/expand.png" width="18" height="18" alt="" style="width:18px;height:18px;" v-else />
					</div>
				</div>
				<div class="comment-input-box">
					<textarea
						ref="textarea"
						v-model="commentText"
						class="comment-input-textarea"
						:placeholder="placeholder"
						:maxlength="2000"
						:auto-height="true"
						:show-confirm-bar="false"
						:adjust-position="true"
						:hold-keyboard="true"
						:focus="isFocused"
						@input="handleInput"
						@focus="handleFocus"
						@blur="handleBlur"
						@keyboardheightchange="handleKeyboardHeightChange"
					></textarea>
					<div class="comment-input-footer">
						<span class="comment-input-count">{{ commentText.length }}/2000字符</span>
						<div class="comment-input-publish" @click="handlePublish">
							<span>发布</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	</van-popup>
</template>

<script>
import { getSystemInfo } from '@/common/dom'
export default {
	name: 'CommentInputPopup',
	props: {
		show: {
			type: Boolean,
			default: false
		},
		placeholder: {
			type: String,
			default: '我来说一说...'
		}
	},
	data() {
		return {
			commentText: '',
			isExpand: false,
			keyboardHeight: 0,
			windowHeight: 0,
			isFocused: false
		}
	},
	mounted() {
		getSystemInfo({
			success: (res) => {
				this.windowHeight = res.windowHeight;
			}
		});
	},
	computed: {
		popupStyle() {
			let height = 'auto';
			if (this.isExpand) {
				if (this.keyboardHeight > 0 && this.windowHeight > 0) {
					const heightPx = this.windowHeight - this.keyboardHeight;
					height = `${heightPx}px`;
				} else {
					height = '100vh';
				}
			}
			return {
				height: height
			};
		}
	},
	watch: {
		show(newVal) {
			if (newVal) {
				this.$nextTick(() => {
					setTimeout(() => {
						this.isFocused = true;
					}, 300);
				});
			} else {
				this.commentText = '';
				this.isExpand = false;
				this.keyboardHeight = 0;
				this.isFocused = false;
			}
		},
		isExpand(newVal) {
			if (newVal && this.show) {
				this.$nextTick(() => {
					this.isFocused = true;
				});
			}
		}
	},
	methods: {
		handleClose() {
			this.$emit('close');
		},
		handleInput(e) {
			const detail = e && e.detail
			const value = detail && detail.value != null
				? detail.value
				: (e && e.target ? e.target.value : '')
			this.commentText = value == null ? '' : value
		},
		handlePublish() {
			if (!this.commentText.trim()) {
				this.msg('请输入评论内容');
				return;
			}
			if (this.commentText.length > 2000) {
				this.msg('评论内容不能超过2000字符');
				return;
			}
			this.$emit('publish', this.commentText);
		},
		handleExpand() {
			const wasFocused = this.isFocused;
			this.isExpand = !this.isExpand;
			this.$nextTick(() => {
				setTimeout(() => {
					this.isFocused = true;
				}, 50);
			});
			this.$emit('expand', this.isExpand);
		},
		handleFocus(e) {
			this.isFocused = true;
			if (e.detail && e.detail.height) {
				this.keyboardHeight = e.detail.height;
			}
		},
		handleBlur() {
			// 延迟处理，避免与展开操作冲突
			setTimeout(() => {
				if (!this.show) {
					this.keyboardHeight = 0;
					this.isFocused = false;
				} else {
					this.$nextTick(() => {
						if (this.show) {
							this.isFocused = true;
						}
					});
				}
			}, 150);
		},
		handleKeyboardHeightChange(e) {
			if (e.detail && e.detail.height !== undefined) {
				this.keyboardHeight = e.detail.height;
			}
		}
	}
}
</script>

<style lang="scss" scoped>
.comment-input-popup {
	background: #fff;
	border-radius: calc(40 * var(--rpx)) calc(40 * var(--rpx)) 0 0;
	transition: height 0.3s ease;
}

.comment-input-popup.is-expanded {

	overflow: hidden;
}

.comment-input-container {
	position: relative;
	padding: calc(20 * var(--rpx));
	height: 100%;
	display: flex;
	flex-direction: column;
	overflow: hidden;
}

.comment-input-header {
	position: absolute;
	top: calc(20 * var(--rpx));
	right: calc(20 * var(--rpx));
	z-index: 10;
}

.comment-input-expand {
	width: calc(60 * var(--rpx));
	height: calc(60 * var(--rpx));
	display: flex;
	align-items: center;
	justify-content: center;
}

.expand-icon {
	width: calc(40 * var(--rpx));
	height: calc(40 * var(--rpx));
	position: relative;
	display: flex;
	align-items: center;
	justify-content: center;
}

.expand-arrow {
	position: absolute;
	width: 0;
	height: 0;
	border: calc(6 * var(--rpx)) solid transparent;
}

.expand-arrow:nth-child(1) {
	top: 0;
	left: 50%;
	transform: translateX(-50%);
	border-bottom-color: #7f848c;
	border-top: none;
}

.expand-arrow:nth-child(2) {
	right: 0;
	top: 50%;
	transform: translateY(-50%);
	border-left-color: #7f848c;
	border-right: none;
}

.expand-arrow:nth-child(3) {
	bottom: 0;
	left: 50%;
	transform: translateX(-50%);
	border-top-color: #7f848c;
	border-bottom: none;
}

.expand-arrow:nth-child(4) {
	left: 0;
	top: 50%;
	transform: translateY(-50%);
	border-right-color: #7f848c;
	border-left: none;
}

.comment-input-box {
	position: relative;
	background: #fff;
	border-radius: calc(8 * var(--rpx));
	padding: calc(20 * var(--rpx));
	padding-bottom: calc(80 * var(--rpx));
	min-height: calc(300 * var(--rpx));
	height: 100%;
	display: flex;
	flex-direction: column;
	overflow: hidden;
}

.comment-input-textarea {
	width: 100%;
	min-height: calc(260 * var(--rpx));
	font-size: calc(calc(30 * var(--rpx)) * var(--font-scale, 1));
	line-height: 1.6;
	color: #404040;
	background: transparent;
	border: none;
	outline: none;
	box-sizing: border-box;
	flex: 1;
	overflow-y: auto;

	scrollbar-width: none;
	-ms-overflow-style: none;
}

.comment-input-textarea::-webkit-scrollbar {
	display: none;
}

.comment-input-textarea::placeholder {
	color: #b9bdc2;
}

.comment-input-footer {
	position: absolute;
	bottom: calc(20 * var(--rpx));
	right: calc(20 * var(--rpx));
	left: calc(20 * var(--rpx));
	display: flex;
	align-items: center;
	justify-content: flex-end;
}

.comment-input-count {
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #b9bdc2;
}

.comment-input-publish {
	padding: calc(10 * var(--rpx)) calc(24 * var(--rpx));
	background: #4285f4;
    border-radius: calc(30 * var(--rpx));
    margin-left: calc(20 * var(--rpx));
}

.comment-input-publish span {
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #fff;
	font-weight: 500;
}
</style>

