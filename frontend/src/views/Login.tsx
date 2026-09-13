import { useState, type FormEvent } from 'react';
import { Navigate, useLocation, useNavigate } from 'react-router';
import auth from '../services/auth';
import Select from '../components/Select';
import ThemeSelect from '../components/ThemeSelect';
import { encodeCredentials, setCredentials, useAuth } from '../auth';
import { languages, useTranslation, type Locale } from '../locales';

export default function Login() {
  const user = useAuth();
  const { t, locale, changeLanguage } = useTranslation();
  const navigate = useNavigate();
  const location = useLocation();
  const [id, setId] = useState('');
  const [password, setPassword] = useState('');
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  const destination = location.state?.from || '/search';
  if (user) return <Navigate to={destination} replace />;
  async function submit(event: FormEvent) {
    event.preventDefault();
    setBusy(true);
    setError('');
    try {
      const base64encoded = encodeCredentials(id, password);
      const { data } = await auth.login({ _id: id, password, base64encoded });
      if (!data.validated) throw new Error('Invalid credentials');
      setCredentials({ _id: id, name: data.user.name || id, role: data.user.role, base64encoded });
      setPassword('');
      navigate(destination, { replace: true });
    } catch (error) {
      setError(error instanceof Error ? error.message : 'Unable to sign in');
    } finally {
      setBusy(false);
    }
  }
  return <main className="login-page">
    <form className="card login-card w-full max-w-sm space-y-5" onSubmit={submit}>
      <div className="workspace-avatar" aria-hidden="true">ZS</div>
      <h1>Zinc Search</h1>
      <ThemeSelect />
      <Select className="language-select" aria-label={t('menu.language')} value={locale} onValueChange={value => changeLanguage(value as Locale)}>{Object.entries(languages).map(([value, label]) => <option key={value} value={value}>{label}</option>)}</Select>
      <label className="block">{t('login.userid')}<input className="mt-1 w-full" data-cy="login-user-id" autoComplete="username" value={id} onChange={e => setId(e.target.value)} required /></label>
      <label className="block">{t('login.password')}<input className="mt-1 w-full" data-cy="login-password" type="password" autoComplete="current-password" value={password} onChange={e => setPassword(e.target.value)} required /></label>
      {error && <p className="error" role="alert">{error}</p>}
      <button className="primary w-full" data-cy="login-sign-in" disabled={busy}>{busy ? 'Signing in…' : t('login.signIn')}</button>
    </form>
  </main>;
}
