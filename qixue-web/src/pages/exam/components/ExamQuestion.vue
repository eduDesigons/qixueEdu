<template>
	<div>
		<div v-if="question.type === 'completion'" class="completion-question">
			<div class="question-title-wrapper">
				<template v-for="(part, partIndex) in completionParts" :key="partIndex">
					<span v-if="part.type === 'text'" class="title-text">{{ part.text }}</span>
					<div v-else-if="part.type === 'blank'" class="blank-wrapper">
						<div
							class="blank-input"
							v-for="(char, charIndex) in part.chars"
							:key="charIndex"
							:class="{ filled: char, active: isActiveBlank(part.blankIndex, charIndex) }"
							@click="$emit('blank', question, part.blankIndex, charIndex)"
						>
							{{ char || ' ' }}
						</div>
					</div>
				</template>
			</div>
		</div>

		<div v-else-if="question.type === 'single'" class="single-question">
			<div class="question-title-wrapper">
				<span class="question-type-tag single-tag">单选</span>
				<span class="question-title-text">{{ question.title }}</span>
			</div>
			<div class="single-options">
				<div
					class="single-option-item"
					v-for="(option, optIndex) in question.options"
					:key="optIndex"
					:class="{ selected: singleAnswer === optIndex }"
					@click="$emit('select-single', question, optIndex)"
				>
					<div class="option-index">{{ String.fromCharCode(65 + optIndex) }}</div>
					<span class="option-text">{{ option }}</span>
					<div class="option-indicator" v-if="singleAnswer === optIndex"></div>
				</div>
			</div>
		</div>

		<div v-else-if="question.type === 'multiple'" class="multiple-question">
			<div class="question-title-wrapper">
				<div class="question-type-tag multiple-tag">多选</div>
				<span class="question-title-text">{{ question.title }}</span>
			</div>
			<div class="multiple-options">
				<div
					class="multiple-option-item"
					v-for="(option, optIndex) in question.options"
					:key="optIndex"
					:class="{ selected: isMultipleSelected(optIndex) }"
					@click="$emit('toggle-multiple', question, optIndex)"
				>
					<div class="option-index">{{ String.fromCharCode(65 + optIndex) }}</div>
					<span class="option-text">{{ option }}</span>
					<div class="option-checkmark" v-if="isMultipleSelected(optIndex)">✓</div>
				</div>
			</div>
		</div>

		<div v-else-if="question.type === 'text'" class="text-question">
			<div class="question-title-text">{{ question.title }}</div>
			<div class="text-input-wrapper">
				<textarea
					:value="textAnswer"
					class="text-input"
					placeholder="请输入您的答案..."
					:maxlength="500"
					@input="$emit('update-text', question, $event.target.value)"
					@blur="$emit('blur-text')"
				></textarea>
				<div class="text-count">{{ textAnswer.length }}/500</div>
			</div>
		</div>

		<div v-else-if="question.type === 'uploadImage'" class="upload-question">
			<div class="question-title-wrapper">
				<div class="question-type-tag upload-tag">图片</div>
				<span class="question-title-text">{{ question.title }}</span>
			</div>
			<div class="upload-tip">
				<span>最多上传 {{ question.options || 1 }} 张图片</span>
			</div>
			<div class="upload-list">
				<div class="upload-item" v-for="(item, idx) in uploads" :key="idx">
					<img :src="static(item)" class="upload-image" @click="$emit('preview-image', question.id, idx)" style="object-fit:cover" />
					<div class="upload-delete" @click="$emit('remove-image', question.id, idx)">
						<van-icon name="cross" size="20" color="#fff" />
					</div>
				</div>
				<div
					class="upload-add-btn"
					v-if="uploads.length < (question.options || 1)"
					@click="$emit('choose-image', question)"
				>
					<van-icon name="plus" size="32" color="#999" />
					<span class="upload-add-text">上传图片</span>
				</div>
			</div>
		</div>

		<div v-else-if="question.type === 'uploadDocument'" class="upload-question">
			<div class="question-title-wrapper">
				<div class="question-type-tag upload-tag">文档</div>
				<span class="question-title-text">{{ question.title }}</span>
			</div>
			<div class="upload-tip">
				<span>最多上传 {{ question.options || 1 }} 个文档</span>
			</div>
			<div class="upload-list">
				<div class="upload-file-item" v-for="(item, idx) in uploads" :key="idx">
					<div class="upload-file-info">
						<van-icon name="description" size="32" color="#4285f4" />
						<div class="upload-file-detail">
							<span class="upload-file-name">{{ fileName(item) }}</span>
						</div>
					</div>
					<div class="upload-file-delete" @click="$emit('remove-file', question.id, idx)">
						<van-icon name="cross" size="20" color="#999" />
					</div>
				</div>
				<div
					class="upload-add-btn upload-file-btn"
					v-if="uploads.length < (question.options || 1)"
					@click="$emit('choose-file', question)"
				>
					<van-icon name="plus" size="32" color="#999" />
					<span class="upload-add-text">上传文档</span>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
import { getCompletionParts, getFileName, getUploadAnswers, getUserAnswers } from '../helpers'

export default {
	name: 'ExamQuestion',
	props: {
		question: { type: Object, required: true },
		userAnswers: { type: Object, default: () => ({}) },
		optionPopup: { type: Object, default: () => ({}) },
	},
	emits: [
		'blank', 'select-single', 'toggle-multiple', 'update-text', 'blur-text',
		'preview-image', 'remove-image', 'choose-image', 'remove-file', 'choose-file',
	],
	computed: {
		completionParts() {
			return getCompletionParts(this.question, this.userAnswers)
		},
		singleAnswer() {
			return getUserAnswers(this.userAnswers[this.question.id])
		},
		textAnswer() {
			return this.userAnswers[this.question.id] || ''
		},
		uploads() {
			return getUploadAnswers(this.userAnswers[this.question.id])
		},
	},
	methods: {
		isActiveBlank(blankIndex, charIndex) {
			return !!(
				this.optionPopup.show &&
				String(this.optionPopup.questionId) === String(this.question.id) &&
				this.optionPopup.blankIndex === blankIndex &&
				this.optionPopup.charIndex === charIndex
			)
		},
		isMultipleSelected(optionIndex) {
			const answers = this.userAnswers[this.question.id] || []
			return Array.isArray(answers) && answers.includes(optionIndex)
		},
		fileName(url) {
			return getFileName(url)
		},
	},
}
</script>

<style lang="scss" scoped>
$primary-color: #4286F5;
$primary-light: #e8f2ff;
$primary-dark: #357ae8;
$warning-color: #ff9800;
$text-primary: #1a1a1a;
$text-secondary: #666;
$text-tertiary: #999;
$bg-primary: #ffffff;
$bg-secondary: #f5f7fa;
$border-color: #e5e7eb;
$shadow-sm: 0 calc(2 * var(--rpx)) calc(8 * var(--rpx)) rgba(0, 0, 0, 0.04);
$shadow-md: 0 calc(4 * var(--rpx)) calc(16 * var(--rpx)) rgba(0, 0, 0, 0.08);

.question-title-text {
	font-size: calc(calc(34 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	color: $text-primary;
	line-height: 1.8;
	margin-bottom: calc(32 * var(--rpx));
	white-space: pre-wrap;
	letter-spacing: calc(0.5 * var(--rpx));
}

.completion-question .question-title-wrapper {
	font-size: calc(calc(32 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	color: #1a1a1a;
	line-height: 2;
	white-space: pre-wrap;
	display: block;
	margin-bottom: 0;
}

.title-text {
	display: inline;
}

.blank-wrapper {
	display: inline-flex;
	flex-wrap: wrap;
	align-items: center;
	gap: calc(8 * var(--rpx));
	margin: 0 calc(4 * var(--rpx));
	vertical-align: middle;
}

.blank-input {
	min-width: calc(64 * var(--rpx));
	height: calc(68 * var(--rpx));
	line-height: calc(68 * var(--rpx));
	text-align: center;
	background: $bg-primary;
	border: calc(2 * var(--rpx)) solid $primary-color;
	border-radius: calc(12 * var(--rpx));
	font-size: calc(calc(30 * var(--rpx)) * var(--font-scale, 1));
	color: $primary-color;
	font-weight: 600;
	box-shadow: $shadow-sm;
	transition: all 0.3s;

	&.filled {
		background: linear-gradient(135deg, $primary-light 0%, rgba(66, 134, 245, 0.1) 100%);
		color: $text-primary;
		border-color: $primary-dark;
	}

	&.active {
		border-color: $primary-dark;
		box-shadow: 0 0 0 calc(4 * var(--rpx)) rgba(66, 134, 245, 0.25);
		background: $primary-light;
	}

	&:active {
		background: $primary-light;
		transform: scale(0.95);
	}
}

.single-question .question-title-wrapper,
.multiple-question .question-title-wrapper {
	display: inline-block;
	align-items: flex-start;
	margin-bottom: calc(32 * var(--rpx));
}

.single-question .question-title-text,
.multiple-question .question-title-text {
	margin-bottom: 0;
	flex: 1;
}

.question-type-tag {
	padding: calc(6 * var(--rpx)) calc(16 * var(--rpx));
	border-radius: calc(6 * var(--rpx));
	font-size: calc(calc(22 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 500;
	margin-right: calc(16 * var(--rpx));
	flex-shrink: 0;
	margin-top: calc(4 * var(--rpx));

	&.single-tag,
	&.upload-tag {
		background: #e8f4ff;
		color: #4285f4;
	}

	&.multiple-tag {
		background: #fff3e0;
		color: #ff9800;
	}
}

.single-options,
.multiple-options {
	display: flex;
	flex-direction: column;
	gap: calc(20 * var(--rpx));
}

.single-option-item,
.multiple-option-item {
	display: flex;
	align-items: center;
	padding: calc(32 * var(--rpx)) calc(28 * var(--rpx));
	background: $bg-primary;
	border: calc(2 * var(--rpx)) solid $border-color;
	border-radius: calc(20 * var(--rpx));
	transition: all 0.3s;
	position: relative;
	box-shadow: $shadow-sm;

	&:active {
		transform: scale(0.98);
	}
}

.single-option-item.selected {
	background: linear-gradient(135deg, $primary-light 0%, rgba(66, 134, 245, 0.05) 100%);
	border-color: $primary-color;
	box-shadow: 0 calc(6 * var(--rpx)) calc(20 * var(--rpx)) rgba(66, 134, 245, 0.2);
	transform: translateY(calc(-2 * var(--rpx)));
}

.multiple-option-item.selected {
	background: linear-gradient(135deg, #fff3e0 0%, rgba(255, 152, 0, 0.05) 100%);
	border-color: $warning-color;
	box-shadow: 0 calc(6 * var(--rpx)) calc(20 * var(--rpx)) rgba(255, 152, 0, 0.2);
	transform: translateY(calc(-2 * var(--rpx)));
}

.option-index {
	width: calc(52 * var(--rpx));
	height: calc(52 * var(--rpx));
	display: flex;
	align-items: center;
	justify-content: center;
	background: $bg-secondary;
	border-radius: 50%;
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 700;
	color: $text-secondary;
	margin-right: calc(24 * var(--rpx));
	flex-shrink: 0;
	transition: all 0.3s;
	box-shadow: $shadow-sm;
}

.single-option-item.selected .option-index {
	background: linear-gradient(135deg, $primary-color 0%, $primary-dark 100%);
	color: #fff;
	box-shadow: 0 calc(4 * var(--rpx)) calc(12 * var(--rpx)) rgba(66, 134, 245, 0.3);
	transform: scale(1.1);
}

.multiple-option-item.selected .option-index {
	background: linear-gradient(135deg, $warning-color 0%, darken($warning-color, 10%) 100%);
	color: #fff;
	box-shadow: 0 calc(4 * var(--rpx)) calc(12 * var(--rpx)) rgba(255, 152, 0, 0.3);
	transform: scale(1.1);
}

.option-text {
	flex: 1;
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #404040;
	line-height: 1.8;
}

.option-indicator {
	width: calc(28 * var(--rpx));
	height: calc(28 * var(--rpx));
	border-radius: 50%;
	background: linear-gradient(135deg, $primary-color 0%, $primary-dark 100%);
	border: calc(4 * var(--rpx)) solid #fff;
	box-shadow: 0 calc(4 * var(--rpx)) calc(12 * var(--rpx)) rgba(66, 134, 245, 0.4);
	flex-shrink: 0;
	margin-left: calc(16 * var(--rpx));
	animation: pulse 2s infinite;
}

@keyframes pulse {
	0%, 100% { transform: scale(1); opacity: 1; }
	50% { transform: scale(1.1); opacity: 0.8; }
}

.option-checkmark {
	width: calc(32 * var(--rpx));
	height: calc(32 * var(--rpx));
	display: flex;
	align-items: center;
	justify-content: center;
	background: #ff9800;
	border-radius: calc(6 * var(--rpx));
	color: #fff;
	font-size: calc(calc(20 * var(--rpx)) * var(--font-scale, 1));
	font-weight: bold;
	flex-shrink: 0;
	margin-left: calc(16 * var(--rpx));
}

.text-input-wrapper {
	position: relative;
	background: $bg-secondary;
	border-radius: calc(20 * var(--rpx));
	padding: calc(32 * var(--rpx));
	min-height: calc(300 * var(--rpx));
	border: calc(2 * var(--rpx)) solid $border-color;
	box-shadow: $shadow-sm;
}

.text-input {
	width: 100%;
	min-height: calc(300 * var(--rpx));
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #404040;
	line-height: 1.8;
}

.text-count {
	position: absolute;
	bottom: calc(16 * var(--rpx));
	right: calc(24 * var(--rpx));
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #999;
}

.upload-question {
	margin-bottom: calc(32 * var(--rpx));
}

.upload-tip {
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	color: $text-tertiary;
	margin-bottom: calc(24 * var(--rpx));
	font-weight: 500;
	padding: calc(16 * var(--rpx)) calc(24 * var(--rpx));
	background: linear-gradient(135deg, $primary-light 0%, rgba(66, 134, 245, 0.03) 100%);
	border-radius: calc(12 * var(--rpx));
	border-left: calc(4 * var(--rpx)) solid $primary-color;
}

.upload-list {
	display: flex;
	flex-wrap: wrap;
	gap: calc(20 * var(--rpx));
}

.upload-item {
	position: relative;
	width: calc(200 * var(--rpx));
	height: calc(200 * var(--rpx));
	border-radius: calc(16 * var(--rpx));
	overflow: hidden;
	box-shadow: $shadow-md;
	border: calc(2 * var(--rpx)) solid $border-color;
}

.upload-image {
	width: 100%;
	height: 100%;
	border-radius: calc(16 * var(--rpx));
}

.upload-delete {
	position: absolute;
	top: calc(-10 * var(--rpx));
	right: calc(-10 * var(--rpx));
	width: calc(40 * var(--rpx));
	height: calc(40 * var(--rpx));
	background: #ff4757;
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
}

.upload-add-btn {
	width: calc(200 * var(--rpx));
	height: calc(200 * var(--rpx));
	border: calc(2 * var(--rpx)) dashed $primary-color;
	border-radius: calc(16 * var(--rpx));
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	background: linear-gradient(135deg, $primary-light 0%, rgba(66, 134, 245, 0.03) 100%);
}

.upload-file-btn {
	height: calc(120 * var(--rpx));
}

.upload-add-text {
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #999;
	margin-top: calc(12 * var(--rpx));
}

.upload-file-item {
	width: 100%;
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: calc(28 * var(--rpx));
	background: $bg-secondary;
	border-radius: calc(16 * var(--rpx));
	margin-bottom: calc(16 * var(--rpx));
	border: calc(2 * var(--rpx)) solid $border-color;
}

.upload-file-info {
	display: flex;
	align-items: center;
	flex: 1;
	margin-right: calc(16 * var(--rpx));
}

.upload-file-detail {
	flex: 1;
	margin-left: calc(20 * var(--rpx));
}

.upload-file-name {
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #404040;
	font-weight: 500;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}

.upload-file-delete {
	padding: calc(8 * var(--rpx));
}
</style>
