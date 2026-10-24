<template>
	<div class="course-video-player course-video-content">
		<div ref="host" class="course-video-host"></div>
		<div v-if="showResumeBar" class="video-resume">上次看到 {{ resumeLabel }}</div>
		<button
			v-if="playSrc && (!playing || cacheBusy)"
			type="button"
			class="video-offline"
			:disabled="cacheBusy"
			:title="cacheHint || cacheLabel"
			@click.stop="onCacheClick"
		>{{ cacheLabel }}</button>
	</div>
</template>

<script>
import { getStorage, setStorage } from '@/common/storage'
	import { cacheFullVideo, formatBytes, getVideoCacheStatus, offlinePlaybackSrc, removeCachedVideo } from '@/common/mediaCache'
	import Player from 'xgplayer'
	import 'xgplayer/dist/index.min.css'

	const SPEED_LIST = [
		{ rate: 2, text: '2.0x', iconText: '2.0x' },
		{ rate: 1.5, text: '1.5x', iconText: '1.5x' },
		{ rate: 1.25, text: '1.25x', iconText: '1.25x' },
		{ rate: 1, text: '1.0x', iconText: '倍速' },
		{ rate: 0.75, text: '0.75x', iconText: '0.75x' }
	]

	export default {
		name: 'CourseVideoPlayer',
		props: {
			src: { type: String, default: '' },
			poster: { type: String, default: '' },
			title: { type: String, default: '' },
			initialTime: { type: Number, default: 0 },
			videoId: { type: String, default: 'courseVideoPlayer' }
		},
		data() {
			return {
				player: null,
				showResumeBar: false,
				resumeTimer: null,
				playing: false,
				cacheBusy: false,
				cacheState: 'idle',
				cacheHint: '',
				cacheAbort: null,
				blobUrl: '',
				mountGen: 0
			}
		},
		computed: {
			playSrc() {
				return this.static(this.src)
			},
			posterSrc() {
				return this.poster ? this.static(this.poster) : ''
			},
			playerConfig() {
				return (this.appInfo && this.appInfo.playerConfig) || {}
			},
			allowSpeed() {
				return this.playerConfig.allow_speed !== false
			},
			resumeLabel() {
				return this.formatClock(this.initialTime)
			},
			cacheLabel() {
				if (this.cacheBusy) return '缓存中'
				if (this.cacheState === 'cached') return '已缓存'
				return '离线'
			}
		},
		watch: {
			playSrc() {
				this.mountPlayer()
				this.refreshCacheStatus()
			},
			posterSrc() {
				if (this.player && this.posterSrc) this.player.poster = this.posterSrc
			}
		},
		mounted() {
			this.mountPlayer()
			this.refreshCacheStatus()
		},
		beforeUnmount() {
			if (this.cacheAbort) this.cacheAbort.abort()
			this.destroyPlayer()
		},
		methods: {
			savedRate() {
				const saved = Number(getStorage('video_playback_rate'))
				return SPEED_LIST.some((item) => item.rate === saved) ? saved : 1
			},
			formatClock(seconds) {
				const s = Math.max(0, Math.floor(Number(seconds) || 0))
				const m = Math.floor(s / 60)
				const r = s % 60
				return (m < 10 ? '0' : '') + m + ':' + (r < 10 ? '0' : '') + r
			},
			setupResumeBar(seconds) {
				if (this.resumeTimer) {
					clearTimeout(this.resumeTimer)
					this.resumeTimer = null
				}
				if (Number(seconds) > 2) {
					this.showResumeBar = true
					this.resumeTimer = setTimeout(() => {
						this.showResumeBar = false
						this.resumeTimer = null
					}, 3000)
				} else {
					this.showResumeBar = false
				}
			},
			hostEl() {
				const el = this.$refs.host
				return el && el.$el ? el.$el : el
			},
			destroyPlayer() {
				if (this.resumeTimer) {
					clearTimeout(this.resumeTimer)
					this.resumeTimer = null
				}
				if (this.player) {
					try { this.player.destroy() } catch (e) {}
					this.player = null
				}
				if (this.blobUrl) {
					URL.revokeObjectURL(this.blobUrl)
					this.blobUrl = ''
				}
				const host = this.hostEl()
				if (host) host.innerHTML = ''
			},
			async resolvePlaybackUrl(src) {
				if (!src || navigator.onLine !== false) return src
				try {
					return (await offlinePlaybackSrc(src)) || src
				} catch (e) {
					return src
				}
			},
			mountPlayer() {
				const gen = ++this.mountGen
				this.destroyPlayer()
				const host = this.hostEl()
				const src = this.playSrc
				if (!host || !src) return
				this.resolvePlaybackUrl(src).then((url) => {
					if (gen !== this.mountGen || !url) {
						if (url && String(url).startsWith('blob:')) URL.revokeObjectURL(url)
						return
					}
					if (String(url).startsWith('blob:')) this.blobUrl = url
					this.createPlayer(url)
				})
			},
			createPlayer(url) {
				const host = this.hostEl()
				if (!host || !url) return
				const start = Number(this.initialTime) || 0
				const rate = this.allowSpeed ? this.savedRate() : 1
				const player = new Player({
					el: host,
					url,
					poster: this.posterSrc,
					width: '100%',
					height: '100%',
					fluid: false,
					fitVideoSize: 'fixed',
					videoFillMode: 'contain',
					autoplay: false,
					videoInit: true,
					playsinline: true,
					lang: 'zh',
					startTime: start > 1 ? start : 0,
					defaultPlaybackRate: rate,
					playbackRate: this.allowSpeed ? SPEED_LIST : false,
					volume: 0.8,
					keyShortcut: true,
					cssFullscreen: false,
					minWaitDelay: 0,
					commonStyle: {
						progressColor: '#4285f4',
						playedColor: '#4285f4',
						cachedColor: 'rgba(255,255,255,.28)',
						volumeColor: '#4285f4'
					},
					ignores: ['download', 'screenshot', 'pip', 'miniscreen', 'definition', 'cssfullscreen', 'dynamicbg'].concat(this.allowSpeed ? [] : ['playbackrate']),
					videoAttributes: {
						preload: 'auto',
						playsinline: true,
						'webkit-playsinline': true
					}
				})
				this.player = player
				this.setupResumeBar(start)
				player.on('play', () => {
					this.playing = true
					this.showResumeBar = false
					if (this.cacheAbort) this.cacheAbort.abort()
					this.$emit('play')
				})
				player.on('pause', () => {
					this.playing = false
					this.$emit('pause')
				})
				player.on('error', () => {
					if (navigator.onLine !== false || this._offlineRetried || !this.playSrc) return
					this._offlineRetried = true
					offlinePlaybackSrc(this.playSrc).then((offline) => {
						if (!offline) return
						const time = player.currentTime || 0
						this.destroyPlayer()
						if (offline.startsWith('blob:')) this.blobUrl = offline
						this.createPlayer(offline)
						this.$nextTick(() => {
							if (this.player && time > 1) this.player.currentTime = time
						})
					}).catch(() => {})
				})
				player.on('ended', () => this.$emit('ended'))
				player.on('timeupdate', () => {
					this.$emit('timeupdate', {
						currentTime: player.currentTime || 0,
						duration: player.duration || 0
					})
				})
				player.on('canplay', () => this.$emit('canplay'))
				const bindRate = () => {
					const video = player.video || player.media
					if (!video || video.__qxRateBound) return
					video.__qxRateBound = true
					video.addEventListener('ratechange', () => {
						const next = Number(video.playbackRate)
						if (next > 0) setStorage('video_playback_rate', next)
					})
				}
				player.on('ready', bindRate)
				bindRate()
			},
			pause() {
				try {
					if (this.player) this.player.pause()
				} catch (e) {}
			},
			seek(seconds) {
				try {
					if (this.player) this.player.currentTime = Number(seconds) || 0
				} catch (e) {}
			},
			async refreshCacheStatus() {
				this.cacheHint = ''
				if (!this.playSrc) {
					this.cacheState = 'idle'
					return
				}
				try {
					const status = await getVideoCacheStatus(this.playSrc)
					this.cacheState = status.cached ? 'cached' : 'idle'
					if (status.cached) this.cacheHint = formatBytes(status.bytes)
				} catch (e) {
					this.cacheState = 'idle'
				}
			},
			async onCacheClick() {
				if (this.cacheBusy || !this.playSrc) return
				if (this.cacheState === 'cached') {
					this.cacheBusy = true
					try {
						await removeCachedVideo(this.playSrc)
						this.cacheState = 'idle'
						this.cacheHint = '已清除本地视频'
					} catch (e) {
						this.cacheHint = e.message || '清除失败'
					} finally {
						this.cacheBusy = false
					}
					return
				}
				this.cacheBusy = true
				this.cacheState = 'downloading'
				this.cacheHint = '正在下载完整文件…'
				const abort = new AbortController()
				this.cacheAbort = abort
				try {
					const result = await cacheFullVideo(this.playSrc, {
						title: this.title,
						signal: abort.signal,
						onProgress: ({ received, total }) => {
							if (total > 0) this.cacheHint = `正在缓存 ${Math.min(100, Math.floor(received / total * 100))}%`
							else this.cacheHint = `已下载 ${formatBytes(received)}`
						}
					})
					this.cacheState = 'cached'
					this.cacheHint = result.warning || `已缓存 ${formatBytes(result.bytes)}，可离线播放`
				} catch (e) {
					if (e && e.name === 'AbortError') {
						this.cacheHint = '已取消缓存'
					} else {
						this.cacheHint = (e && e.message) || '缓存失败'
					}
					this.cacheState = 'idle'
				} finally {
					this.cacheBusy = false
					this.cacheAbort = null
				}
			}
		}
	}
</script>

<style lang="scss" scoped>
	.course-video-player {
		position: relative;
		width: 100%;
		aspect-ratio: 16 / 9;
		background: #000;
		overflow: hidden;
	}

	.course-video-host {
		width: 100%;
		height: 100%;
	}

	.video-offline {
		position: absolute;
		top: 8px;
		right: 8px;
		z-index: 20;
		border-radius: 999px;
		padding: 3px 8px;
		background: rgba(17, 24, 39, 0.45);
		color: rgba(255, 255, 255, 0.88);
		font-size: 11px;
		line-height: 1.4;
	}

	.video-offline:disabled {
		opacity: 0.7;
	}

	.video-resume {
		position: absolute;
		left: 16px;
		top: 16px;
		z-index: 30;
		padding: 6px 12px;
		border-radius: 16px;
		background: rgba(17, 24, 39, 0.72);
		color: #fff;
		font-size: 12px;
		pointer-events: none;
	}

	.course-video-player :deep(.xgplayer-start) {
		width: 58px;
		height: 58px;
	}

	.course-video-player :deep(xg-start-inner) {
		background: rgba(66, 133, 244, 0.92);
		box-shadow: 0 8px 24px rgba(0, 0, 0, 0.28);
	}

	.course-video-player :deep(.xgplayer-controls) {
		height: 48px;
		background-image: linear-gradient(180deg, transparent, rgba(0, 0, 0, 0.72));
	}

	.course-video-player :deep(.xgplayer-progress-outer) {
		height: 3px;
		background: transparent;
	}

	.course-video-player :deep(.xgplayer-progress-inner) {
		background: rgba(255, 255, 255, 0.28) !important;
	}

	.course-video-player :deep(.xgplayer-progress-played),
	.course-video-player :deep(xg-mini-progress-played),
	.course-video-player :deep(.xgplayer-drag),
	.course-video-player :deep(.xg-curbar) {
		background: #4285f4 !important;
		background-image: none !important;
	}

	.course-video-player :deep(.xgplayer-progress-cache),
	.course-video-player :deep(xg-mini-progress-cache) {
		background: rgba(255, 255, 255, 0.38) !important;
	}

	.course-video-player :deep(.xgplayer-progress-btn) {
		width: 12px;
		height: 12px;
		background: #fff !important;
		border: 0 !important;
		box-shadow: 0 0 0 4px rgba(66, 133, 244, 0.35) !important;
	}

	.course-video-player :deep(.xgplayer-progress-btn:before) {
		display: none;
	}

	.course-video-player :deep(.xgplayer-time) {
		font-size: 12px;
		font-variant-numeric: tabular-nums;
	}

	.course-video-player :deep(.xgplayer-playbackrate .icon-text) {
		min-width: 40px;
		height: 22px;
		padding: 0 8px;
		border-radius: 11px;
		background: rgba(255, 255, 255, 0.16);
		line-height: 22px;
		font-size: 12px;
		font-weight: 600;
	}

	.course-video-player :deep(.xg-options-list) {
		width: 92px;
		padding: 6px 0;
		border-radius: 8px;
		background: rgba(17, 24, 39, 0.92);
	}

	.course-video-player :deep(.xg-options-list li) {
		height: 32px;
		margin: 0;
		padding: 0;
		line-height: 32px;
		font-size: 13px;
	}

	.course-video-player :deep(.xg-options-list li:hover),
	.course-video-player :deep(.xg-options-list li.selected) {
		color: #8eb7ff !important;
		background: rgba(66, 133, 244, 0.16);
	}

	.course-video-player :deep(.xgplayer-slider) {
		border-radius: 8px;
		background: rgba(17, 24, 39, 0.92);
	}

	.course-video-player :deep(xg-trigger .time-preview .xg-cur),
	.course-video-player :deep(xg-trigger .xg-top-note i) {
		color: #8eb7ff !important;
	}
</style>
