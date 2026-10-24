import { request } from '@/utils'

export default {
  // 轮播图
  getBannerList: () => request.get('/admin/decoration/banner/getBannerList'),
  saveBannerInfo: data => request.post('/admin/decoration/banner/saveBannerInfo', data),
  delBanner: id => request.get(`/admin/decoration/banner/delBanner?id=${id}`),
  toggleStatus: id => request.get(`/admin/decoration/banner/toggleStatus?id=${id}`),
  updateBannerSort: ids => request.post('/admin/decoration/banner/updateBannerSort', { ids }),

  // 通知
  // 获取通知列表
  getNoticeList: data => request.post('/admin/decoration/notice/getNoticeList', data),
  // 获取通知详情
  getNoticeDetail: id => request.get(`/admin/decoration/notice/getNoticeDetail?id=${id}`),
  // 保存/创建通知
  saveNotice: data => request.post('/admin/decoration/notice/saveNotice', data),
  // 更新通知发布状态
  updateNoticeStatus: data => request.post('/admin/decoration/notice/updateNoticeStatus', data),
  // 删除通知
  deleteNotice: id => request.get(`/admin/decoration/notice/deleteNotice?id=${id}`),
}
