export const courseTypes = [
  { label: '视频课', value: 1 },
  { label: '音频课', value: 2 },
  { label: '图文课', value: 3 },
]
export function getCourseTypeName(type) {
  return courseTypes.find(item => item.value === type)?.label
}
