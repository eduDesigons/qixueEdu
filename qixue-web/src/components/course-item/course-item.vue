<template>
	<div class="LessonItem_lesson-wrap_3WSpt LessonList_lesson-item_2rWPz">
		<img class="LessonItem_img-wrap_3IUmH"
			src="/static/img/video-img-1.png"
			alt="MySQL实战45讲"  style="object-fit:cover" />
		<div class="LessonItem_content_16t2y">
			<div class="">
				<p class="LessonItem_title_2XyiJ" item-title=""><span class="LessonItem_icon-sort_2PwIC"
						data-icon="1"></span><span>MySQL实战45讲</span></p>
				<p class="LessonItem_desc_3U0cu" item-desc=""><span class="LessonItem_name_3SiLy">林晓斌</span><span
						class=""></span></p>
			</div>
			<p class="LessonItem_count_22tHn">23.4w 人已学习</p>
		</div>
	</div>

</template>

<script>
	export default {
		name: "user-item",
		data() {
			return {

			};
		},
		props: {
			info: {
				type: Object,
				default: () => ({})
			}
		}
	}
</script>

<style lang="scss">
	.LessonItem_lesson-wrap_3WSpt{display:flex;width:313px;height:90px;}
	.LessonItem_img-wrap_3IUmH{display:block;width:120px;height:90px;border-radius:4px;}
	.LessonItem_content_16t2y{display:flex;flex-direction:column;justify-content:space-between;margin-left:13px;}
	.LessonItem_title_2XyiJ{position:relative;top:-3px;width:180px;line-height:21px;overflow:hidden;display:-webkit-box;-webkit-line-clamp:2;-webkit-box-orient:vertical;font-size: calc(15px * var(--font-scale, 1));font-weight:500;color:#353535;}
	.LessonItem_title_2XyiJ .LessonItem_icon-sort_2PwIC{position:relative;top:4px;display:inline-block;width:16px;height:19px;margin-right:6px;background:url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAmCAYAAAClI5npAAAAAXNSR0IArs4c6QAAAERlWElmTU0AKgAAAAgAAYdpAAQAAAABAAAAGgAAAAAAA6ABAAMAAAABAAEAAKACAAQAAAABAAAAIKADAAQAAAABAAAAJgAAAAAjxpdTAAAFKklEQVRYCbVYPY8cRRDtnp3Z2727tQMkHB3/gAyJBAlfQGASREBAREJCABmBEYKTQAhiSAgIiJBwgOzElhxwK12CREZCQgIRREgkILwzvFfVb6Zm9s53i+SyPFVTXf3qs3vHzgl09+zX5/Ls3zvVrDqqqjTrOmpTyjklyuRG1FMu672sdbcKTxkXA4B1Xd6Afus2zWuvvPDMj5nOq/rRD3BSTQEn2wOwiwpOAW8ZTAHCBoht+6h5vu7y33dSmlUGoswKR7zAzHi6Iuw3XxlRt22HCpUMpxFM8GzZHJlU0XcNB0cdQIwmEXcwJja5DFx2h13XQo0AjReTEZsAloTUQxT9CAF0s00JYGpO710LJyVDOmK2cqjgyrIF4+AeIDb73lGASobm3axuN54lA2ccBFNBYFDwfBNhW1aFxiDFRk7KsGexyEmSVUAP2H1ovWZJZSAluZHjeBBQTF7dREquF1mcBlH24DE3RZmRSt0qfINjbkRRjpLlZby+nZEPZJyZKZ4jFTzA1TaAfRPlSJxRBXnSUx9BtIJmIGsf+tfDmXbYz8A8aNfxfKECeCk9M/v4oJ2SpgiZAQ8OCCEjzoRXwDCJM5khmnbS2TorwH4Q4yIKa+qn+LBxMBrKD0CpxekjyMQpQ+haP2KehdnCQjq+SyY/b92Uo4fK5/Zefq+k9nsLSkr9jRfClM42GB40jgcVZItdaY0dDun6+rRFtMY9wMuF8G4uCL7zmDAAcpJdRNTpYiHuqKdQSGcbuM6NfAFB9qH3VypqO5vFwPZibbD3/BmGyDV6C8ZUyUw86qJc1ollx7C3Z4R88YQh4wXl0WBdPgNKtwBYRNIxAsm+zsr7KeAayI8YHPJogswhLqp+6Mrda78PbgDbct/bDqUiTmWUxzcHvg2mNyErEDZAtJIHlZfIvEFEcPgznongc5ywxzKZGTsFXhAHjc8W7rlGfh5pbbQ/mkoWJ0iQmavPQEEg47oANQ+xKMQQCUtc+p6fB6jLgEaQ/SKKPYc39Zw20Tnbo0H0/T6g0Z76nhRZzyFYUEUB2X4Ndc75y0iw/hdykgHt2tJ3OpFMTrJ5sLZxI4M/J+Aw1PgA1jFUiBqoAABYDSak0iLZs2Xb8kgXSqiPGfK6nqXlcp8VGAB8Aix2y8DhVQa+uezcWmgtYlttF6BiiwvIiM2aOi0We6kBJ43ugZFl/3JxgFt3u7XPy95vL0INh8vlonesdZsBnwyq6GzIWEYD15qnzKfKSpvpTBCrntfpYH+55ViY5YuovLJf5iNmLdPYEF8fh0M7Tonfds28SfsHizQHfxzZjxFCNxvrH+QLewqPNlP0DJJs3DWp2WvS4eH+pY7NHA+vQHnTtU5OYiDSmUKFEaeyyMz0cHWQ5ghgF7LvgT5l4gVwyeLuTYV35/NFk1bXDtPe3nwXv71txSHiOedfXkDilKd/ed+YDSvUpjV+V45v3HjqGMdq3SPuKJTvgZI2kht9sVAdEqYM1XqW08m7b710Gnzd/Ob+TzcxPidI4MWgv1TMn3zxoHhn+f0Mk5N079tdn/M6V9XJ7bHjLQe7BpI//vzB70j7aUdSuuSePNi6Qsa33751St1V6UqB5PRHjSN4D/m+SWD7wMDo2wcGMq5z+vC9d15eX9VptHv91rOneH9sa1DZe/mzr85W//z159eo+qvIucW/bL9PVfXRB//TcQwiyqzIpuvexzf9MVLF/8bk71bXmjd6Gwby6ZcPr/eKJyR8+/CX63fPfl49IfjdYf8DP/c1u0oFYa0AAAAASUVORK5CYII=) no-repeat;background-size:contain;}
	.LessonItem_title_2XyiJ .LessonItem_icon-sort_2PwIC:before{display:block;content:attr(data-icon);position:absolute;top:-2px;left:5px;font-size: calc(12px * var(--font-scale, 1));font-weight:500;color:#fff;}
	.LessonItem_title_2XyiJ .LessonItem_icon-sort_2PwIC[data-icon="1"]{background:url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAmCAMAAACS/WnbAAAAeFBMVEUAAAD3cV77SUf3cmD8S0n2dmT/eGj/Tkn0dGP0cWD7SEb2cV73cV/7R0T3cV76Rkb2cF72XFPycWD3V0/4VE7zalzxdGL1Y1f2X1T3WlHyb1/4UUzzbV30ZVnzaVv0Z1r5Tkr5TEj1YVb0YFX0ZFj6SEX6Skf3aFmjN2ZxAAAAEHRSTlMAwMA/PzgQEOHh7Ou0tKCg+1VhygAAAPVJREFUOMuFiotygjAURLdWQfs0BIXwqgmP/v8fdidxGgNCziyTvewBDsfrKscD9+smB5yyTU7IIqBaklUZcz/Ar6kahsV3/w9NBFwizIXu0jEPHV0E/ERAGwGDpx1axr7E9QFFBJiiMIWxFK4zfjcwfW96w/y/AQZlBJRClMLfPCwsrkNEiAsqAqRSUkmLsj18BeQm5Rm5Q+b+Lw+Gr/g9U1hFcaZQu6vOa2Y2O+Epdl4XpJ0duC2wsxe0X/RNa51zDoWA3f715bkw6WmadnuQUAH/j9PI3Oe5grfR4eaF8o7vYF4oX0g//bxUPhKWNMUKSYI4f9U/Y8HhyWGrAAAAAElFTkSuQmCC) no-repeat;background-size:contain;}
	.LessonItem_desc_3U0cu{width:180px;margin-top:4px;line-height:17px;font-size: calc(12px * var(--font-scale, 1));font-weight:400;color:#7f848c;overflow:hidden;display:-webkit-box;-webkit-line-clamp:1;-webkit-box-orient:vertical;}
	.LessonItem_desc_3U0cu .LessonItem_name_3SiLy{margin-right:2px;}
	.LessonItem_count_22tHn{line-height:17px;font-size: calc(12px * var(--font-scale, 1));font-weight:400;color:#7f848c;}
	/*! CSS Used from: Embedded */
	.LessonList_list_TK6EU .LessonList_lesson-item_2rWPz{margin:0 13px 16px 0;}
</style>