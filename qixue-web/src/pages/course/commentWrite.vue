<template>
	<div>
		<div g-input="" class="Input_g-input_1Z4Gs" modelmodifiers="[object Object]" autofocus=""><!----><textarea
				type="textarea" modelmodifiers="[object Object]" autofocus="" class="Input_g-input-inner_2YRiO"
				g-input-inner="" placeholder="我来说一说..."></textarea><!----><!----><span
				class="Input_g-input-word-limit_gCB4A Input_g-input-word-limit-position_2hb9q"
				g-input-word-limit="">0/2000 字符</span></div>
	</div>
</template>

<script>
	export default {
		data() {
			return {

			}
		},
		methods: {

		}
	}
</script>

<style>
/*! CSS Used from: Embedded */
.Input_g-input_1Z4Gs{box-sizing:border-box;overflow:hidden;position:relative;display:flex;justify-content:center;align-items:center;background-color:#fff;border:1px solid #e9eaed;}
.Input_g-input_1Z4Gs textarea{height:100%;}
.Input_g-input-inner_2YRiO{box-sizing:border-box;padding:6px 10px;flex:1;background-color:#fff;border:none;font-size: calc(14px * var(--font-scale, 1));font-weight:400;}
.Input_g-input-word-limit_gCB4A{color:#b9bdc2;font-size: calc(14px * var(--font-scale, 1));font-weight:400;}
.Input_g-input-inner_2YRiO::placeholder{line-height:24px;font-size: calc(15px * var(--font-scale, 1));color:#b9bdc2;}
.Input_g-input-inner_2YRiO:focus{outline:none;}
.Input_g-input-word-limit-position_2hb9q{position:absolute;bottom:5px;right:10px;}
/*! CSS Used from: Embedded */
.WriteComment_leave-word-release_3RooL [g-input]{padding:2px 10px 32px 10px;border:none;}
.WriteComment_leave-word-release_3RooL [g-input] [g-input-inner]{min-height:246px;padding-bottom:8px;line-height:24px;color:#404040;font-size: calc(15px * var(--font-scale, 1));font-weight:400;}
.WriteComment_leave-word-release_3RooL [g-input-word-limit]{bottom:12px;right:20px;}
</style>