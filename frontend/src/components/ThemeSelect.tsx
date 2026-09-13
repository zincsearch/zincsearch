import Select from './Select';
import { useTranslation } from '../locales';
import { setTheme, useTheme, type Theme } from '../theme';

export default function ThemeSelect() {
  const theme = useTheme();
  const { t } = useTranslation();
  const icon = <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.7" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
    {theme === 'system' ? <><rect x="3" y="4" width="18" height="13" rx="2" /><path d="M8 21h8m-4-4v4" /></> : theme === 'light' ? <><circle cx="12" cy="12" r="4" /><path d="M12 2v2m0 16v2M2 12h2m16 0h2M5 5l1.5 1.5m11 11L19 19M5 19l1.5-1.5m11-11L19 5" /></> : <path d="M20.5 13A8.5 8.5 0 0 1 11 3.5 8.5 8.5 0 1 0 20.5 13Z" />}
  </svg>;
  return <Select className="theme-select" aria-label={t('theme.label')} title={`${t('theme.label')}: ${t(`theme.${theme}`)}`} displayValue={icon} value={theme} onValueChange={value => setTheme(value as Theme)}>
    {(['system', 'light', 'dark'] as const).map(value => <option key={value} value={value}>{t(`theme.${value}`)}</option>)}
  </Select>;
}
