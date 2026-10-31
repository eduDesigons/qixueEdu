<template>
	<div class="course-video-player course-video-content">
		<div ref="host" class="course-video-host"></div>
		<div v-if="showResumeBar" class="video-resume">上次看到 {{ resumeLabel }}</div>
	</div>
</template>

<script>
import { getStorage, setStorage } from '@/common/storage'
	import { cacheFullVideo, getVideoCacheStatus, offlinePlaybackSrc } from '@/common/mediaCache'
	import Player from 'xgplayer'
	import 'xgplayer/dist/index.min.css'

	const SPEED_LIST = [
		{ rate: 2, text: '2.0x', iconText: '2.0x' },
		{ rate: 1.5, text: '1.5x', iconText: '1.5x' },
		{ rate: 1.25, text: '1.25x', iconText: '1.25x' },
		{ rate: 1, text: '1.0x', iconText: '倍速' },
		{ rate: 0.75, text: '0.75x', iconText: '0.75x' }
	]

	function finiteSeconds(value) {
		const n = Number(value)
		return Number.isFinite(n) && n > 0 ? n : 0
	}

	export default {
		name: 'CourseVideoPlayer',
		emits: ['play', 'pause', 'timeupdate', 'ended', 'canplay'],
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
				cacheBusy: false,
				cacheState: 'idle',
				cacheAbort: null,
				blobUrl: '',
				playbackUrl: '',
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
				this.unbindMediaWatch()
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
				if (!src) return src
				try {
					return (await offlinePlaybackSrc(src)) || src
				} catch (e) {
					return src
				}
			},
			mountPlayer() {
				if (this.cacheAbort) this.cacheAbort.abort()
				const gen = ++this.mountGen
				this._playbackRetried = false
				this.playbackUrl = ''
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
				this.playbackUrl = url
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
						'webkit-playsinline': true,
						'x5-playsinline': true,
						'x5-video-player-type': 'h5'
					}
				})
				this.player = player
				this.bindMediaWatch(player)
				this.setupResumeBar(start)
				player.on('play', () => {
					this.showResumeBar = false
					this.maybeAutoCache()
					this.startProgressWatch()
					this.$emit('play')
				})
				player.on('pause', () => {
					this.stopProgressWatch()
					this.emitProgress()
					this.$emit('pause')
				})
				player.on('error', () => {
					if (this._playbackRetried || !this.playSrc) return
					this._playbackRetried = true
					const usingLocal = this.playbackUrl && this.playbackUrl !== this.playSrc
					const next = usingLocal
						? Promise.resolve(this.playSrc)
						: offlinePlaybackSrc(this.playSrc).then((offline) => (offline && offline !== this.playbackUrl) ? offline : '')
					next.then((url) => {
						if (!url || url === this.playbackUrl) return
						const time = player.currentTime || 0
						this.destroyPlayer()
						if (String(url).startsWith('blob:')) this.blobUrl = url
						this.createPlayer(url)
						this.$nextTick(() => {
							if (this.player && time > 1) this.player.currentTime = time
						})
					}).catch(() => {})
				})
				player.on('ended', () => {
					this.stopProgressWatch()
					this.emitProgress()
					this.$emit('ended')
				})
				player.on('timeupdate', () => this.emitProgress())
				player.on('durationchange', () => this.emitProgress())
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
				player.on('ready', () => {
					bindRate()
					if (!this._mediaVideo) this.bindMediaWatch(player)
				})
				bindRate()
			},
			mediaEl() {
				const player = this.player
				if (!player) return null
				return player.media || player.video || null
			},
			playbackPoint() {
				const video = this.mediaEl()
				const player = this.player
				return {
					currentTime: finiteSeconds(video && video.currentTime) || finiteSeconds(player && player.currentTime),
					duration: finiteSeconds(video && video.duration) || finiteSeconds(player && player.duration)
				}
			},
			emitProgress() {
				if (!this.player) return
				this.$emit('timeupdate', this.playbackPoint())
			},
			bindMediaWatch(player) {
				this.unbindMediaWatch()
				const video = player.media || player.video
				if (!video) return
				video.setAttribute('playsinline', 'true')
				video.setAttribute('webkit-playsinline', 'true')
				video.setAttribute('x5-playsinline', 'true')
				video.setAttribute('x5-video-player-type', 'h5')
				this._mediaVideo = video
				this._onMediaProgress = () => this.emitProgress()
				video.addEventListener('loadedmetadata', this._onMediaProgress)
				video.addEventListener('durationchange', this._onMediaProgress)
				video.addEventListener('webkitendfullscreen', this._onMediaProgress)
				video.addEventListener('x5videoexitfullscreen', this._onMediaProgress)
				this._onVisible = () => this.emitProgress()
				document.addEventListener('visibilitychange', this._onVisible)
			},
			unbindMediaWatch() {
				this.stopProgressWatch()
				const video = this._mediaVideo
				if (video && this._onMediaProgress) {
					video.removeEventListener('loadedmetadata', this._onMediaProgress)
					video.removeEventListener('durationchange', this._onMediaProgress)
					video.removeEventListener('webkitendfullscreen', this._onMediaProgress)
					video.removeEventListener('x5videoexitfullscreen', this._onMediaProgress)
				}
				if (this._onVisible) {
					document.removeEventListener('visibilitychange', this._onVisible)
				}
				this._mediaVideo = null
				this._onMediaProgress = null
				this._onVisible = null
			},
			startProgressWatch() {
				if (this._watchTimer) return
				this._watchTimer = setInterval(() => this.emitProgress(), 1000)
			},
			stopProgressWatch() {
				if (!this._watchTimer) return
				clearInterval(this._watchTimer)
				this._watchTimer = null
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
				if (!this.playSrc) {
					this.cacheState = 'idle'
					return
				}
				try {
					const status = await getVideoCacheStatus(this.playSrc)
					this.cacheState = status.cached ? 'cached' : 'idle'
				} catch (e) {
					this.cacheState = 'idle'
				}
			},
			async maybeAutoCache() {
				if (this.cacheBusy || this.cacheState === 'cached' || !this.playSrc) return
				await this.runCache()
			},
			async runCache() {
				this.cacheBusy = true
				this.cacheState = 'downloading'
				const abort = new AbortController()
				this.cacheAbort = abort
				try {
					await cacheFullVideo(this.playSrc, {
						title: this.title,
						signal: abort.signal,
					})
					this.cacheState = 'cached'
				} catch (e) {
					if (this.cacheState !== 'cached') this.cacheState = 'idle'
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
