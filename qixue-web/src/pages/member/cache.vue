<template>
	<div class="cache-page">
		<div class="cache-list" v-if="!loading">
			<template v-for="group in groups" :key="group.id">
				<div class="cache-row" @click="toggleGroup(group)">
					<span class="cache-check" :class="{ on: groupChecked(group) }"><van-icon v-if="groupChecked(group)" name="success" size="14" color="#fff" /></span>
					<span class="cache-name">{{ group.name }}</span>
					<span class="cache-size">{{ formatBytes(group.bytes) }}</span>
					<van-icon
						v-if="group.expandable"
						class="cache-arrow"
						:name="group.open ? 'arrow-up' : 'arrow'"
						size="16"
						color="#c8c9cc"
						@click.stop="group.open = !group.open"
					/>
				</div>
				<div v-if="group.expandable && group.open" class="cache-sub-list">
					<div v-if="!group.items.length" class="cache-empty">暂无已缓存文件</div>
					<div class="cache-sub" v-for="item in group.items" :key="item.key" @click="toggleItem(group, item)">
						<span class="cache-check" :class="{ on: item.checked }"><van-icon v-if="item.checked" name="success" size="12" color="#fff" /></span>
						<span class="cache-name">{{ item.title }}</span>
						<span class="cache-size">{{ formatBytes(item.bytes) }}</span>
					</div>
				</div>
			</template>
		</div>
		<div class="cache-loading" v-else>正在计算缓存…</div>
		<div class="cache-foot">
			<button type="button" @click="selectAll">全选</button>
			<button type="button" @click="selectNone">取消选择</button>
			<span class="cache-estimate">预估 {{ formatBytes(selectedBytes) }}</span>
			<button type="button" class="cache-clear" :disabled="clearing || selectedBytes <= 0" @click="clearSelected">清理</button>
		</div>
	</div>
</template>

<script>
import { showModal } from '@/common/dom'
import { clearCacheGroups, formatBytes, listCacheGroups } from '@/common/mediaCache'

export default {
	data() {
		return {
			groups: [],
			loading: true,
			clearing: false
		}
	},
	computed: {
		selectedBytes() {
			return this.groups.reduce((sum, group) => {
				if (group.items && group.items.length) {
					return sum + group.items.filter((item) => item.checked).reduce((inner, item) => inner + item.bytes, 0)
				}
				return sum + (group.checked ? group.bytes : 0)
			}, 0)
		}
	},
	onShow() {
		this.reload()
	},
	methods: {
		formatBytes,
		groupChecked(group) {
			if (group.items && group.items.length)
				return group.items.every((item) => item.checked)
			return !!group.checked
		},
		async reload() {
			this.loading = true
			try {
				const groups = await listCacheGroups()
				groups.forEach((group) => { group.open = false })
				this.groups = groups
			}
			catch (e) {
				this.msg('读取缓存失败')
			}
			finally {
				this.loading = false
			}
		},
		toggleGroup(group) {
			const next = !this.groupChecked(group)
			group.checked = next
			;(group.items || []).forEach((item) => { item.checked = next })
		},
		toggleItem(group, item) {
			item.checked = !item.checked
			group.checked = this.groupChecked(group)
		},
		selectAll() {
			this.groups.forEach((group) => this.setGroup(group, true))
		},
		selectNone() {
			this.groups.forEach((group) => this.setGroup(group, false))
		},
		setGroup(group, on) {
			group.checked = on
			;(group.items || []).forEach((item) => { item.checked = on })
		},
		clearSelected() {
			if (this.selectedBytes <= 0 || this.clearing) return
			showModal({
				title: '清理缓存',
				content: this.groups.some((group) => group.id === 'app' && group.checked)
					? `将清理约 ${formatBytes(this.selectedBytes)}。应用文件会在下次打开时重新下载。`
					: `将清理约 ${formatBytes(this.selectedBytes)}。已缓存的视频和音频需要重新下载。`,
				confirmText: '清理',
				success: (res) => {
					if (res.confirm) this.doClear()
				}
			})
		},
		async doClear() {
			this.clearing = true
			try {
				await clearCacheGroups(this.groups)
				this.msg('清理完成')
				await this.reload()
			}
			catch (e) {
				this.msg('清理失败')
			}
			finally {
				this.clearing = false
			}
		}
	}
}
</script>

<style scoped>
.cache-page {
	min-height: 100vh;
	background: #fff;
	padding-bottom: calc(120 * var(--rpx));
}
.cache-row,
.cache-sub {
	display: flex;
	align-items: center;
	min-height: 56px;
	padding: 0 16px;
	border-bottom: 1px solid #f3f4f6;
}
.cache-sub {
	padding-left: 44px;
	background: #fafafa;
}
.cache-check {
	display: flex;
	align-items: center;
	justify-content: center;
	width: 22px;
	height: 22px;
	margin-right: 12px;
	border: 1px solid #d8d8d8;
	border-radius: 50%;
	background: #fff;
	flex: none;
}
.cache-check.on {
	border-color: var(--van-primary-color);
	background: var(--van-primary-color);
}
.cache-name {
	flex: 1;
	min-width: 0;
	color: #1a1a1a;
	font-size: calc(16px * var(--font-scale, 1));
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}
.cache-size {
	margin-left: 12px;
	color: #8a8f99;
	font-size: calc(15px * var(--font-scale, 1));
}
.cache-arrow {
	margin-left: 8px;
}
.cache-empty,
.cache-loading {
	padding: 18px 16px;
	color: #8a8f99;
	text-align: center;
}
.cache-foot {
	position: fixed;
	left: 0;
	right: 0;
	bottom: 0;
	display: flex;
	align-items: center;
	gap: 8px;
	min-height: 56px;
	padding: 8px 12px calc(8px + env(safe-area-inset-bottom));
	background: #fff;
	border-top: 1px solid #eee;
	box-sizing: border-box;
	z-index: 5;
}
.cache-foot button {
	border: 0;
	background: transparent;
	color: var(--van-primary-color);
	font-size: calc(15px * var(--font-scale, 1));
	padding: 8px 4px;
}
.cache-estimate {
	flex: 1;
	text-align: right;
	color: #8a8f99;
	font-size: calc(14px * var(--font-scale, 1));
}
.cache-clear {
	min-width: 72px;
	height: 36px;
	border-radius: 6px;
	background: var(--van-primary-color) !important;
	color: #fff !important;
}
.cache-clear:disabled {
	opacity: 0.45;
}
</style>
