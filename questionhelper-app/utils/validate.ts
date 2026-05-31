// 手机号验证
export const isValidPhone = (phone: string): boolean => {
  return /^1[3-9]\d{9}$/.test(phone)
}

// 邮箱验证
export const isValidEmail = (email: string): boolean => {
  return /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/.test(email)
}

// 密码强度验证（至少6位，包含字母和数字）
export const isValidPassword = (password: string): boolean => {
  return /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d@$!%*#?&]{6,}$/.test(password)
}

// 身份证号验证
export const isValidIdCard = (idCard: string): boolean => {
  return /^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$/.test(idCard)
}

// 用户名验证（2-20位，字母数字中文）
export const isValidUsername = (username: string): boolean => {
  return /^[一-龥a-zA-Z0-9]{2,20}$/.test(username)
}

// 非空验证
export const isNotEmpty = (value: string): boolean => {
  return value !== null && value !== undefined && value.trim() !== ''
}

// 长度验证
export const isLengthInRange = (value: string, min: number, max: number): boolean => {
  return value.length >= min && value.length <= max
}
