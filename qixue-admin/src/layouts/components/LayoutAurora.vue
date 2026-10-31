<template>
  <div
    class="layout-shell__aurora"
    aria-hidden="true"
    :style="{ '--aurora-opacity': String(canvasOpacity) }"
  >
    <canvas ref="canvasRef" class="layout-shell__canvas" />
  </div>
</template>

<script setup>
/**
 * DeepSeek 同款 WebGL2 pattern 背景
 * 支持 params / opacity 覆盖；登录页可用 preset="hero" 还原官网动效 + 主题色
 */
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'

const props = defineProps({
  /** app：后台壳层；hero：接近 DeepSeek 首页动效 */
  preset: {
    type: String,
    default: 'app',
  },
  /** 覆盖颜色，如 ['#a8c9cd', '#FFFFFF', '#FFFFFF'] */
  colors: {
    type: Array,
    default: null,
  },
  /** 画布透明度，不传则按 preset */
  opacity: {
    type: Number,
    default: undefined,
  },
  /** 额外覆盖任意 PARAMS 字段 */
  params: {
    type: Object,
    default: null,
  },
})

const canvasRef = ref(null)

const VERT = `#version 300 es
in vec4 a_position;
void main() {
  gl_Position = a_position;
}
`

const FRAG = `#version 300 es
precision mediump float;

uniform float u_time;
uniform float u_pixelRatio;
uniform vec2 u_resolution;

uniform float u_scale;
uniform float u_rotation;
uniform vec4 u_color1;
uniform vec4 u_color2;
uniform vec4 u_color3;
uniform vec4 u_color4;
uniform vec4 u_color5;
uniform vec4 u_color6;
uniform float u_colorCount;
uniform float u_grain;
uniform float u_proportion;
uniform float u_softness;
uniform float u_shape;
uniform float u_shapeScale;
uniform float u_distortion;
uniform float u_swirl;
uniform float u_swirlIterations;
uniform vec2 u_offset;

out vec4 fragColor;

#define TWO_PI 6.28318530718
#define PI 3.14159265358979323846

vec2 rotate(vec2 uv, float th) {
  return mat2(cos(th), sin(th), -sin(th), cos(th)) * uv;
}

float random(vec2 st) {
  return fract(sin(dot(st.xy, vec2(12.9898, 78.233))) * 43758.5453123);
}

float noise(vec2 st) {
  vec2 i = floor(st);
  vec2 f = fract(st);
  float a = random(i);
  float b = random(i + vec2(1.0, 0.0));
  float c = random(i + vec2(0.0, 1.0));
  float d = random(i + vec2(1.0, 1.0));
  vec2 u = f * f * (3.0 - 2.0 * f);
  float x1 = mix(a, b, u.x);
  float x2 = mix(c, d, u.x);
  return mix(x1, x2, u.y);
}

vec3 blend_multi(float mixer, float softness) {
  float edge = 1.0 - softness;
  vec3 col = u_color1.rgb;
  if (u_colorCount > 1.5) {
    float r1 = smoothstep(0.0 + 0.35 * edge, 0.7 - 0.35 * edge, mixer);
    col = mix(col, u_color2.rgb, r1);
  }
  if (u_colorCount > 2.5) {
    float r2 = smoothstep(0.3 + 0.35 * edge, 1.0 - 0.35 * edge, mixer);
    col = mix(col, u_color3.rgb, r2);
  }
  if (u_colorCount > 3.5) {
    col = mix(col, u_color4.rgb, smoothstep(0.4, 0.75, mixer));
  }
  if (u_colorCount > 4.5) {
    col = mix(col, u_color5.rgb, smoothstep(0.55, 0.85, mixer));
  }
  if (u_colorCount > 5.5) {
    col = mix(col, u_color6.rgb, smoothstep(0.7, 0.95, mixer));
  }
  return col;
}

void main() {
  vec2 uv = gl_FragCoord.xy / u_resolution.xy;

  float t = .5 * u_time;
  float noise_scale = .0005 + .006 * u_scale;

  uv -= .5;
  uv *= (noise_scale * u_resolution);
  uv = rotate(uv, u_rotation * .5 * PI);
  uv /= u_pixelRatio;
  uv += .5;
  uv += u_offset;

  float n1 = noise(uv * 1. + t);
  float n2 = noise(uv * 2. - t);
  float angle = n1 * TWO_PI;
  uv.x += 4. * u_distortion * n2 * cos(angle);
  uv.y += 4. * u_distortion * n2 * sin(angle);

  float iterations_number = ceil(clamp(u_swirlIterations, 1., 30.));
  for (float i = 1.; i <= 30.0; i++) {
    if (i > iterations_number) break;
    uv.x += clamp(u_swirl, 0., 2.) / i * cos(t + i * 1.5 * uv.y);
    uv.y += clamp(u_swirl, 0., 2.) / i * cos(t + i * 1. * uv.x);
  }

  float proportion = clamp(u_proportion, 0., 1.);

  float shape = 0.;
  float mixer = 0.;
  if (u_shape < .5) {
    vec2 checks_shape_uv = uv * (.5 + 3.5 * u_shapeScale);
    shape = .5 + .5 * sin(checks_shape_uv.x) * cos(checks_shape_uv.y);
    mixer = shape + .48 * sign(proportion - .5) * pow(abs(proportion - .5), .5);
  } else if (u_shape < 1.5) {
    vec2 stripes_shape_uv = uv * (.25 + 3. * u_shapeScale);
    float f = fract(stripes_shape_uv.y);
    shape = smoothstep(.0, .55, f) * smoothstep(1., .45, f);
    mixer = shape + .48 * sign(proportion - .5) * pow(abs(proportion - .5), .5);
  } else {
    float sh = 1. - uv.y;
    sh -= .5;
    sh /= (noise_scale * u_resolution.y);
    sh += .5;
    float shape_scaling = .2 * (1. - u_shapeScale);
    shape = smoothstep(.45 - shape_scaling, .55 + shape_scaling, sh + .3 * (proportion - .5));
    mixer = shape;
  }

  vec3 col = blend_multi(mixer, clamp(u_softness, 0., 1.));
  fragColor = vec4(col, 1.0);

  if (u_grain > 0.0) {
    float g = random(gl_FragCoord.xy + vec2(u_time * 100.0));
    fragColor.rgb += (g - 0.5) * u_grain;
  }
}
`

const APP_PARAMS = {
  colors: ['#a8c9cd', '#FFFFFF', '#FFFFFF'],
  grain: 0,
  speed: 48,
  shape: 'checks',
  shapeScale: 8,
  rotation: -6,
  proportion: 64,
  softness: 100,
  distortion: 18,
  swirl: 14,
  swirlIterations: 8,
  scale: 0.45,
  offsetX: 0,
  offsetY: 60,
}

const HERO_PARAMS = {
  colors: ['#a8c9cd', '#FFFFFF', '#FFFFFF'],
  grain: 0,
  speed: 36,
  shape: 'checks',
  shapeScale: 10,
  rotation: -5,
  proportion: 50,
  softness: 100,
  distortion: 28,
  swirl: 18,
  swirlIterations: 10,
  scale: 0.5,
  offsetX: 0,
  offsetY: 65,
}

const SHAPE_MAP = { checks: 0, stripes: 1, gradient: 2 }
const FRAME_MS = 1000 / 30

function hexToRgb(hex) {
  const t = String(hex || '').replace('#', '')
  if (t.length < 6)
    return [0, 0, 0]
  return [
    Number.parseInt(t.slice(0, 2), 16) / 255,
    Number.parseInt(t.slice(2, 4), 16) / 255,
    Number.parseInt(t.slice(4, 6), 16) / 255,
  ]
}

const activeParams = computed(() => {
  const base = props.preset === 'hero' ? HERO_PARAMS : APP_PARAMS
  const colors = props.colors?.length
    ? props.colors
    : (props.params?.colors || base.colors)
  return {
    ...base,
    ...(props.params || {}),
    colors,
  }
})

const canvasOpacity = computed(() => {
  if (typeof props.opacity === 'number')
    return props.opacity
  return props.preset === 'hero' ? 1 : 0.76
})

let dispose = null

onMounted(() => {
  const canvas = canvasRef.value
  if (!canvas)
    return

  const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
  const gl = canvas.getContext('webgl2', {
    alpha: true,
    premultipliedAlpha: false,
    powerPreference: 'low-power',
    preserveDrawingBuffer: reduceMotion,
  })
  if (!gl) {
    canvas.style.display = 'none'
    return
  }

  function compile(type, src) {
    const s = gl.createShader(type)
    gl.shaderSource(s, src)
    gl.compileShader(s)
    if (!gl.getShaderParameter(s, gl.COMPILE_STATUS)) {
      console.error('Shader compile error:', gl.getShaderInfoLog(s))
      gl.deleteShader(s)
      return null
    }
    return s
  }

  const vs = compile(gl.VERTEX_SHADER, VERT)
  const fs = compile(gl.FRAGMENT_SHADER, FRAG)
  if (!vs || !fs)
    return

  const prog = gl.createProgram()
  gl.attachShader(prog, vs)
  gl.attachShader(prog, fs)
  gl.linkProgram(prog)
  if (!gl.getProgramParameter(prog, gl.LINK_STATUS)) {
    console.error('Program link error:', gl.getProgramInfoLog(prog))
    return
  }
  gl.useProgram(prog)

  const buf = gl.createBuffer()
  gl.bindBuffer(gl.ARRAY_BUFFER, buf)
  gl.bufferData(gl.ARRAY_BUFFER, new Float32Array([-1, -1, 1, -1, -1, 1, 1, 1]), gl.STATIC_DRAW)
  const aPos = gl.getAttribLocation(prog, 'a_position')
  gl.enableVertexAttribArray(aPos)
  gl.vertexAttribPointer(aPos, 2, gl.FLOAT, false, 0, 0)

  const u = name => gl.getUniformLocation(prog, name)
  const locs = {
    u_time: u('u_time'),
    u_pixelRatio: u('u_pixelRatio'),
    u_resolution: u('u_resolution'),
    u_scale: u('u_scale'),
    u_offset: u('u_offset'),
    u_rotation: u('u_rotation'),
    u_color1: u('u_color1'),
    u_color2: u('u_color2'),
    u_color3: u('u_color3'),
    u_color4: u('u_color4'),
    u_color5: u('u_color5'),
    u_color6: u('u_color6'),
    u_colorCount: u('u_colorCount'),
    u_grain: u('u_grain'),
    u_proportion: u('u_proportion'),
    u_softness: u('u_softness'),
    u_shape: u('u_shape'),
    u_shapeScale: u('u_shapeScale'),
    u_distortion: u('u_distortion'),
    u_swirl: u('u_swirl'),
    u_swirlIterations: u('u_swirlIterations'),
  }

  let dpr = Math.min(window.devicePixelRatio || 1, 1.5)
  const start = performance.now()
  let raf = 0
  let last = 0
  let visible = true
  let disposed = false

  function resize() {
    const rect = canvas.getBoundingClientRect()
    dpr = Math.min(window.devicePixelRatio || 1, 1.5)
    const w = Math.max(1, Math.floor(rect.width * dpr))
    const h = Math.max(1, Math.floor(rect.height * dpr))
    if (canvas.width !== w || canvas.height !== h) {
      canvas.width = w
      canvas.height = h
      gl.viewport(0, 0, w, h)
    }
  }

  function draw() {
    const p = activeParams.value
    const colors = Array.from({ length: 6 }, (_, i) => (i < p.colors.length ? hexToRgb(p.colors[i]) : [0, 0, 0]))
    gl.uniform1f(locs.u_scale, p.scale)
    gl.uniform2f(locs.u_offset, (p.offsetX || 0) / 100, (p.offsetY || 0) / 100)
    gl.uniform1f(locs.u_rotation, p.rotation / 90)
    gl.uniform4f(locs.u_color1, colors[0][0], colors[0][1], colors[0][2], 1)
    gl.uniform4f(locs.u_color2, colors[1][0], colors[1][1], colors[1][2], 1)
    gl.uniform4f(locs.u_color3, colors[2][0], colors[2][1], colors[2][2], 1)
    gl.uniform4f(locs.u_color4, colors[3][0], colors[3][1], colors[3][2], 1)
    gl.uniform4f(locs.u_color5, colors[4][0], colors[4][1], colors[4][2], 1)
    gl.uniform4f(locs.u_color6, colors[5][0], colors[5][1], colors[5][2], 1)
    gl.uniform1f(locs.u_colorCount, p.colors.length)
    gl.uniform1f(locs.u_grain, p.grain)
    gl.uniform1f(locs.u_proportion, p.proportion / 100)
    gl.uniform1f(locs.u_softness, p.softness / 100)
    gl.uniform1f(locs.u_shape, SHAPE_MAP[p.shape] ?? 0)
    gl.uniform1f(locs.u_shapeScale, p.shapeScale / 100)
    gl.uniform1f(locs.u_distortion, p.distortion / 100)
    gl.uniform1f(locs.u_swirl, p.swirl / 50)
    gl.uniform1f(locs.u_swirlIterations, p.swirlIterations)
    const t = (performance.now() - start) * 0.001 * (p.speed / 100)
    gl.uniform1f(locs.u_time, t)
    gl.uniform1f(locs.u_pixelRatio, dpr)
    gl.uniform2f(locs.u_resolution, canvas.width, canvas.height)
    gl.drawArrays(gl.TRIANGLE_STRIP, 0, 4)
  }

  function tick(now) {
    raf = 0
    if (disposed || !visible || document.hidden)
      return
    if (last && now - last < FRAME_MS) {
      schedule()
      return
    }
    last = now
    resize()
    draw()
    if (!reduceMotion)
      schedule()
  }

  function schedule() {
    if (!disposed && !raf && visible && !document.hidden && !reduceMotion)
      raf = requestAnimationFrame(tick)
  }

  resize()
  draw()
  if (!reduceMotion)
    schedule()

  const ro = new ResizeObserver(() => {
    resize()
    draw()
    if (!reduceMotion)
      schedule()
  })
  ro.observe(canvas)

  const io = typeof IntersectionObserver !== 'undefined'
    ? new IntersectionObserver(([entry]) => {
        visible = entry.isIntersecting
        if (visible)
          schedule()
        else if (raf) {
          cancelAnimationFrame(raf)
          raf = 0
        }
      }, { threshold: 0 })
    : null
  io?.observe(canvas)

  const onVis = () => {
    if (document.hidden) {
      if (raf) {
        cancelAnimationFrame(raf)
        raf = 0
      }
    }
    else {
      schedule()
    }
  }
  document.addEventListener('visibilitychange', onVis)

  dispose = () => {
    disposed = true
    if (raf)
      cancelAnimationFrame(raf)
    ro.disconnect()
    io?.disconnect()
    document.removeEventListener('visibilitychange', onVis)
    gl.deleteProgram(prog)
    gl.deleteShader(vs)
    gl.deleteShader(fs)
    gl.deleteBuffer(buf)
  }
})

onBeforeUnmount(() => {
  dispose?.()
  dispose = null
})
</script>

<style scoped>
.layout-shell__aurora {
  position: absolute;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  overflow: hidden;
  background: #fff;
}

.layout-shell__canvas {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  display: block;
  opacity: var(--aurora-opacity, 0.76);
}
</style>
