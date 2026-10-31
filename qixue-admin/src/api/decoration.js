import { request } from '@/utils'

export default {
  getBannerList: () => request.get('/admin/decoration/banner/getBannerList'),
  saveBannerInfo: data => request.post('/admin/decoration/banner/saveBannerInfo', data),
  delBanner: id => request.get(`/admin/decoration/banner/delBanner?id=${id}`),
  toggleStatus: id => request.get(`/admin/decoration/banner/toggleStatus?id=${id}`),
  updateBannerSort: ids => request.post('/admin/decoration/banner/updateBannerSort', { ids }),

  getNoticeList: data => request.post('/admin/decoration/notice/getNoticeList', data),
  getNoticeDetail: id => request.get(`/admin/decoration/notice/getNoticeDetail?id=${id}`),
  saveNotice: data => request.post('/admin/decoration/notice/saveNotice', data),
  updateNoticeStatus: data => request.post('/admin/decoration/notice/updateNoticeStatus', data),
  deleteNotice: id => request.get(`/admin/decoration/notice/deleteNotice?id=${id}`),
}
