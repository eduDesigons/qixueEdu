import { NAvatar, NTag } from 'naive-ui'
import { h } from 'vue'
import { useRouter } from 'vue-router'
import { getAssetUrl } from '@/utils'

export function useOpsUserNav() {
  const router = useRouter()
  function openUser(userId) {
    if (userId)
      router.push(`/member/user/${userId}`)
  }
  function userCell(row) {
    const user = row.user
    return h(
      'div',
      {
        class: 'flex items-center gap-8 cursor-pointer',
        onClick: () => openUser(row.user_id),
      },
      [
        h(NAvatar, { size: 36, src: getAssetUrl(user?.avatar) || undefined, round: true }),
        h('div', { class: 'min-w-0' }, [
          h('div', { class: 'font-500 truncate' }, user?.nick_name || '-'),
          h('div', { class: 'text-12 text-secondary' }, user?.mobile || ''),
        ]),
      ],
    )
  }
  return { openUser, userCell }
}

export function endTag(isEnd, yes = '已结业', no = '学习中') {
  return h(
    NTag,
    { size: 'small', round: true, bordered: false, type: Number(isEnd) === 1 ? 'success' : 'default' },
    { default: () => (Number(isEnd) === 1 ? yes : no) },
  )
}

export function sumCell(sum) {
  const n = Number(sum) || 0
  return h('span', { class: n > 0 ? 'text-success' : (n < 0 ? 'text-error' : '') }, n > 0 ? `+${n}` : String(n))
}

export const endOptions = [
  { label: '全部', value: null },
  { label: '已结业', value: 1 },
  { label: '未结业', value: 0 },
]
