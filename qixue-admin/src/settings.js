/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/13 20:54:36
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

export const defaultPrimaryColor = '#316C72'
/** 与 defaultPrimaryColor 对应的 RGB，供 CSS 变量使用 */
export const defaultPrimaryRgb = '49, 108, 114'

/**
 * js 文件下使用这个做类型提示
 * @type import('naive-ui').GlobalThemeOverrides
 */
export const naiveThemeOverrides = {
  common: {
    primaryColor: '#316C72FF',
    primaryColorHover: '#3A7C83FF',
    primaryColorPressed: '#2B4C59FF',
    primaryColorSuppl: '#3A7C83FF',
    infoColor: '#5B7C86FF',
    successColor: '#2F9B63FF',
    warningColor: '#C9941AFF',
    errorColor: '#D4534AFF',
    textColorBase: '#1C2B2EFF',
    textColor1: '#1C2B2EFF',
    textColor2: '#5A6B6DFF',
    textColor3: '#8A9697FF',
    borderColor: '#316C721A',
    dividerColor: '#316C721A',
    hoverColor: '#316C7214',
    pressedColor: '#316C7226',
    opacityDisabled: '0.42',
    borderRadius: '12px',
    borderRadiusSmall: '10px',
    fontFamily:
      "-apple-system, BlinkMacSystemFont, 'SF Pro Text', 'SF Pro Display', 'PingFang SC', 'Hiragino Sans GB', 'Noto Sans SC', 'Helvetica Neue', 'Microsoft YaHei UI', sans-serif",
    fontFamilyMono:
      "'SF Mono', ui-monospace, 'Cascadia Code', 'JetBrains Mono', Menlo, Monaco, Consolas, monospace",
    fontWeightStrong: '600',
    heightTiny: '26px',
    heightSmall: '30px',
    heightMedium: '36px',
    heightLarge: '44px',
  },
  Button: {
    borderRadiusTiny: '9px',
    borderRadiusSmall: '10px',
    borderRadiusMedium: '12px',
    borderRadiusLarge: '14px',
    heightTiny: '26px',
    heightSmall: '30px',
    heightMedium: '36px',
    heightLarge: '44px',
    fontWeight: '560',
    paddingTiny: '0 10px',
    paddingSmall: '0 12px',
    paddingMedium: '0 16px',
    paddingLarge: '0 22px',
    colorPrimary: '#316C72FF',
    colorHoverPrimary: '#3A7C83FF',
    colorPressedPrimary: '#2B4C59FF',
    colorFocusPrimary: '#3A7C83FF',
    textColorPrimary: '#FFFFFFFF',
    textColorHoverPrimary: '#FFFFFFFF',
    textColorPressedPrimary: '#FFFFFFFF',
    borderPrimary: '1px solid #316C72FF',
    borderHoverPrimary: '1px solid #3A7C83FF',
    borderPressedPrimary: '1px solid #2B4C59FF',
    rippleColorPrimary: '#316C7233',
  },
  Dialog: {
    borderRadius: '20px',
    padding: '22px 24px 18px',
    titleFontWeight: '650',
    titleFontSize: '17px',
    buttonPrimaryColor: '#316C72FF',
    buttonPrimaryColorHover: '#3A7C83FF',
    buttonPrimaryColorPressed: '#2B4C59FF',
  },
  Card: {
    borderRadius: '16px',
    paddingMedium: '18px 20px',
    titleFontWeight: '650',
    titleFontSize: '17px',
  },
  Tabs: {
    tabFontWeightActive: '620',
    tabGapMediumLine: '24px',
    tabPaddingMediumLine: '8px 12px',
    tabBorderRadius: '10px',
    panePaddingMedium: '16px 0 0',
  },
  Input: {
    borderRadius: '12px',
    heightTiny: '26px',
    heightSmall: '30px',
    heightMedium: '36px',
    heightLarge: '42px',
    fontSize: '14px',
  },
  InternalSelection: {
    borderRadius: '12px',
    heightTiny: '26px',
    heightSmall: '30px',
    heightMedium: '36px',
    heightLarge: '42px',
  },
  Tag: {
    borderRadius: '999px',
    heightSmall: '22px',
    heightMedium: '24px',
    fontSizeSmall: '11.5px',
    fontSizeMedium: '12px',
    padding: '0 9px',
  },
  Pagination: {
    itemBorderRadius: '10px',
    itemSizeMedium: '32px',
  },
  DataTable: {
    borderRadius: '14px',
    thFontWeight: '550',
    thPaddingMedium: '10px 12px',
    tdPaddingMedium: '11px 12px',
    fontSizeMedium: '13.5px',
    loadingColor: '#316C72',
    loadingSize: '22px',
    opacityLoading: '0.72',
  },
  Spin: {
    color: '#316C72',
    textColor: '#8A9697',
    fontSize: '13px',
    opacitySpinning: '0.72',
    sizeTiny: '16px',
    sizeSmall: '20px',
    sizeMedium: '24px',
    sizeLarge: '30px',
  },
  Menu: {
    borderRadius: '12px',
    itemBorderRadius: '10px',
    itemHeight: '36px',
  },
  Dropdown: {
    borderRadius: '14px',
    optionHeightMedium: '34px',
    padding: '6px',
  },
  Modal: {
    borderRadius: '20px',
  },
  Drawer: {
    borderRadius: '18px',
  },
  Alert: {
    borderRadius: '14px',
  },
  Message: {
    borderRadius: '16px',
  },
  Notification: {
    borderRadius: '16px',
  },
  Switch: {
    railBorderRadiusMedium: '999px',
    buttonBorderRadiusMedium: '999px',
    railHeightMedium: '28px',
    railWidthMedium: '48px',
    buttonHeightMedium: '24px',
    buttonWidthMedium: '24px',
  },
  Checkbox: {
    borderRadius: '7px',
    sizeMedium: '18px',
  },
  Radio: {
    buttonBorderRadius: '11px',
  },
  Slider: {
    handleSize: '22px',
    railHeight: '6px',
  },
  Progress: {
    railBorderRadiusHorizontal: '999px',
    fillBorderRadiusHorizontal: '999px',
    railHeightHorizontal: '6px',
  },
  Tooltip: {
    borderRadius: '10px',
  },
  Popover: {
    borderRadius: '16px',
    padding: '10px 12px',
  },
  Skeleton: {
    borderRadius: '10px',
  },
  Form: {
    labelFontWeight: '560',
    labelFontSizeTopMedium: '13px',
    feedbackFontSizeMedium: '12px',
    feedbackHeightMedium: '18px',
  },
  Descriptions: {
    titleTextColor: '#1C2B2EFF',
    thTextColor: '#5A6B6DFF',
    tdTextColor: '#1C2B2EFF',
    borderColor: '#316C721A',
  },
  Empty: {
    iconSizeMedium: '56px',
    fontSizeMedium: '13px',
  },
  DatePicker: {
    itemSize: '32px',
    itemCellBorderRadius: '10px',
    panelBorderRadius: '16px',
  },
  Tree: {
    nodeBorderRadius: '10px',
    nodeHeight: '32px',
  },
}

export const basePermissions = [
]
