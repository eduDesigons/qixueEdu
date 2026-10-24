import assert from 'node:assert/strict'
import {
  OFFLINE_VIDEO_PATH,
  assertDownloadComplete,
  isMp4Url,
  offlinePlaybackUrl,
  parseContentRangeTotal,
  stableMediaKey,
} from '../src/common/mediaCacheCore.js'

const signed = 'https://cdn.example.com/course/a.mp4?token=abc&expires=170000&x-amz-signature=zzz&v=2'
const stable = stableMediaKey(signed)
assert.equal(stable, 'https://cdn.example.com/course/a.mp4?v=2')
assert.equal(stableMediaKey('https://cdn.example.com/course/a.mp4?token=other&v=2'), stable)
assert.equal(isMp4Url(signed), true)
assert.equal(isMp4Url('https://cdn.example.com/a.jpg?token=1'), false)
assert.equal(parseContentRangeTotal('bytes 0-1/98765'), 98765)
assert.equal(parseContentRangeTotal('bytes */100'), 100)
assert.equal(parseContentRangeTotal('bytes 0-1/*'), null)

const partial = assertDownloadComplete({ status: 206, received: 2, contentLength: 2, contentRangeTotal: 1000 })
assert.equal(partial.ok, false)

const mismatch = assertDownloadComplete({ status: 200, received: 10, contentLength: 100, contentRangeTotal: null })
assert.equal(mismatch.ok, false)

const rangeMismatch = assertDownloadComplete({ status: 200, received: 50, contentLength: 50, contentRangeTotal: 80 })
assert.equal(rangeMismatch.ok, false)

const full = assertDownloadComplete({ status: 200, received: 80, contentLength: 80, contentRangeTotal: 80 })
assert.equal(full.ok, true)

const chunked = assertDownloadComplete({ status: 200, received: 80, contentLength: null, contentRangeTotal: null })
assert.equal(chunked.ok, true)

const offline = offlinePlaybackUrl(stable)
assert.ok(offline.startsWith(OFFLINE_VIDEO_PATH))
assert.equal(new URL(offline, 'https://app.local').searchParams.get('key'), stable)
assert.equal(isMp4Url(offline), false)

console.log('media-cache core ok')
