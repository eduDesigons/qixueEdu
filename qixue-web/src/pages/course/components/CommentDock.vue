<template>
	<div class="dock">
		<div class="dock__reply" v-if="replyName">
			<span>回复 @{{ replyName }}</span>
			<span class="dock__cancel" @click="$emit('cancel')">取消</span>
		</div>
		<textarea
			ref="box"
			class="dock__input"
			:value="modelValue"
			:placeholder="placeholder"
			maxlength="2000"
			@input="onInput"
			@keydown.enter.exact="onEnter"
		></textarea>
		<div class="dock__send" @click="$emit('send')">发布</div>
	</div>
</template>

<script>
export default {
	name: 'CommentDock',
	props: {
		modelValue: { type: String, default: '' },
		placeholder: { type: String, default: '我来说一说…' },
		replyName: { type: String, default: '' },
	},
	emits: ['update:modelValue', 'send', 'cancel'],
	methods: {
		onInput(e) {
			const detail = e && e.detail
			const value = detail && detail.value != null
				? detail.value
				: (e && e.target ? e.target.value : '')
			this.$emit('update:modelValue', value == null ? '' : value)
		},
		onEnter(e) {
			if (e.isComposing || e.keyCode === 229) return
			e.preventDefault()
			this.$emit('send')
		},
		focus() {
			const el = this.$refs.box
			const node = el && el.focus ? el : (el && el.$el)
			if (node && node.focus) node.focus()
		},
	},
}
</script>

<style scoped>
.dock {
	flex: 1;
	min-width: 0;
	display: flex;
	flex-wrap: wrap;
	align-items: flex-end;
	gap: 8px;
}
.dock__reply {
	flex: 1 0 100%;
	display: flex;
	align-items: center;
	gap: 8px;
	color: #6b7280;
	font-size: 12px;
}
.dock__cancel {
	color: #4286f5;
	cursor: pointer;
}
.dock__input {
	flex: 1;
	min-width: 0;
	height: 40px;
	max-height: 96px;
	padding: 8px 12px;
	border: 0;
	border-radius: 8px;
	outline: none;
	background: #f6f7fb;
	color: #1a1a1a;
	font-family: inherit;
	font-size: 14px;
	line-height: 1.5;
	resize: none;
	box-sizing: border-box;
}
.dock__send {
	flex: none;
	height: 40px;
	padding: 0 16px;
	border-radius: 8px;
	background: #4286f5;
	color: #fff;
	line-height: 40px;
	cursor: pointer;
}
</style>
