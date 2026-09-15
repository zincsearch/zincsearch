import { translate } from '../locales';

export function validatePassword(password: string, confirmPassword: string): string {
  if (password.length < 8) return translate('passwordValidation.tooShort');
  if (!/[a-z]/i.test(password)) return translate('passwordValidation.missingLetter');
  if (!/[0-9]/.test(password)) return translate('passwordValidation.missingDigit');
  if (password !== confirmPassword) return translate('passwordValidation.mismatch');
  return '';
}
