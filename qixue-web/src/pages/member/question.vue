<template>
	<div class="question-page">

		<div class="question-content">
			<van-collapse :accordion="true" :border="false" v-model="activeNames">
				<van-collapse-item
					v-for="(item, index) in questionList"
					:key="item.id || index"
					:name="item.id || index"
					:title="item.title"
				>
					<div class="question-answer">
						<div class="qx-parse" style="user-select:text" v-html="richHtml(item.answer)"></div>
					</div>
				</van-collapse-item>
			</van-collapse>

			<div class="question-empty" v-if="!questionLoading && questionList.length === 0">
				<van-empty description="暂无常见问题" />
			</div>

			<div class="question-loading" v-if="questionLoading">
				<van-loading type="circular" />
				<span>加载中...</span>
			</div>
		</div>
		<div class="question-entry" @click="navTo('/pages/member/feedback')">
			<div class="question-entry__main">
				<span class="question-entry__title">没有找到答案？</span>
				<span class="question-entry__desc">把问题或建议发给我们</span>
			</div>
			<div class="question-entry__go">
				<span>问题/建议</span>
				<van-icon name="arrow" size="14" color="#4286F5" />
			</div>
		</div>
	</div>
</template>

<script>
export default {
	data() {
		return {
			questionList: [],
			questionLoading: false,
			activeNames: null
		}
	},
	onLoad() {
		this.getQuestionList()
	},
	methods: {
		getQuestionList() {
			this.questionLoading = true
			this.$http.get('/sys/getQuestion').then(res => {
				this.questionLoading = false
				this.questionList = res.data || []
			}).catch(err => {
				this.questionLoading = false
				console.log('获取问题列表失败:', err)
				this.msg('获取问题列表失败')
			})
		}
	}
}
</script>

<style lang="scss" scoped>
.question-page {
	min-height: 100vh;
	padding-bottom: calc(calc(32 * var(--rpx)) + env(safe-area-inset-bottom));
	box-sizing: border-box;
	background: #f8f9fa;
}

.question-content {
	padding: calc(20 * var(--rpx)) 0 0;
}

.question-entry {
	display: flex;
	align-items: center;
	margin: calc(24 * var(--rpx)) calc(30 * var(--rpx)) 0;
	padding: calc(28 * var(--rpx)) calc(28 * var(--rpx));
	background: #fff;
	border-radius: calc(20 * var(--rpx));
	box-shadow: 0 calc(4 * var(--rpx)) calc(16 * var(--rpx)) rgba(0, 0, 0, 0.04);
}

.question-entry__main {
	flex: 1;
	min-width: 0;
}

.question-entry__title {
	display: block;
	font-size: calc(calc(30 * var(--rpx)) * var(--font-scale, 1));
	font-weight: 600;
	color: #1a1a1a;
	line-height: calc(42 * var(--rpx));
}

.question-entry__desc {
	display: block;
	margin-top: calc(6 * var(--rpx));
	font-size: calc(calc(24 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
	line-height: calc(34 * var(--rpx));
}

.question-entry__go {
	display: flex;
	align-items: center;
	flex-shrink: 0;
	margin-left: calc(16 * var(--rpx));
	font-size: calc(calc(26 * var(--rpx)) * var(--font-scale, 1));
	color: #4286F5;
	font-weight: 500;
}

.question-entry__go span {
	margin-right: calc(4 * var(--rpx));
}

.question-answer {
	padding: calc(24 * var(--rpx)) calc(32 * var(--rpx));
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #404040;
	line-height: 1.8;
}

.question-empty {
	padding: calc(200 * var(--rpx)) 0;
}

.question-loading {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	gap: calc(24 * var(--rpx));
	padding: calc(200 * var(--rpx)) 0;
	font-size: calc(calc(28 * var(--rpx)) * var(--font-scale, 1));
	color: #7f848c;
}
</style>