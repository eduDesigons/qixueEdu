const primaryColor = import.meta.env.VITE_PRIMARY_COLOR || '#4286F5'
const primaryRgb = import.meta.env.VITE_PRIMARY_RGB || '66,134,245'

export default {
  baseUrl: import.meta.env.VITE_BASE_URL || '/api/webapp',
  version: import.meta.env.VITE_VERSION || 'v2.5.0',
  title: import.meta.env.VITE_TITLE || '',
  envName: import.meta.env.VITE_ENV_NAME || 'production',
  primaryColor,
  primaryRgb,
}
