import { useId, useState } from 'react';
import auth from '../services/auth';
import { encodeCredentials, getCredentials, setCredentials, type Credentials } from '../auth';
import { useTranslation } from '../locales';
import { validatePassword } from '../utils/password';
import { ErrorMessage, Modal } from './management/Common';
import PasswordInput from './PasswordInput';

export default function Account({ user, onClose }: { user: Credentials; onClose: () => void }) {
  const passwordId = useId();
  const { t } = useTranslation();
  const [name, setName] = useState(user.name);
  const [current, setCurrent] = useState('');
  const [next, setNext] = useState('');
  const [confirm, setConfirm] = useState('');
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);
  const [done, setDone] = useState(false);
  async function save() {
    if (!current) {
      setError('authErrors.currentRequired');
      return;
    }
    const newName = name.trim();
    if (newName.length < 3) {
      setError('account.nameTooShort');
      return;
    }
    const message = next || confirm ? validatePassword(next, confirm) : '';
    if (message) {
      setError(message);
      return;
    }
    if (newName === user.name && !next) {
      setError('account.nothingToChange');
      return;
    }
    setError('');
    setBusy(true);
    try {
      await auth.updateAccount({ _id: user._id, password: current, name: newName, new_password: next || undefined });
      const session = getCredentials();
      if (session === user) {
        setCredentials({ ...session, name: newName, base64encoded: encodeCredentials(user._id, next || current) });
      }
      setCurrent('');
      setNext('');
      setConfirm('');
      setDone(true);
    } catch (err) {
      const status = (err as { response?: { status?: number } }).response?.status;
      setError(status === 401 ? 'account.wrongPassword' : 'account.failed');
    } finally {
      setBusy(false);
    }
  }
  return (
    <Modal title={t('account.title')} onClose={() => { if (!busy) onClose(); }}>
      {done
        ? <p role='status'>{t('account.success')}</p>
        : (
          <form className='grid gap-4' noValidate onSubmit={(event) => { event.preventDefault(); void save(); }}>
            <p>{t('account.userid')}: <strong>{user._id}</strong> · {t('account.role')}: <strong>{user.role}</strong></p>
            <label>
              {t('account.name')}
              <input className='block w-full' autoComplete='name' value={name} disabled={busy} onChange={(e) => setName(e.target.value)} required />
            </label>
            <div>
              <label htmlFor={`${passwordId}-current`}>{t('account.current')}</label>
              <PasswordInput id={`${passwordId}-current`} className='block w-full' autoComplete='current-password' value={current} disabled={busy} onChange={(e) => setCurrent(e.target.value)} required />
            </div>
            <div>
              <label htmlFor={`${passwordId}-new`}>{t('account.new')}</label>
              <PasswordInput id={`${passwordId}-new`} className='block w-full' autoComplete='new-password' value={next} disabled={busy} onChange={(e) => setNext(e.target.value)} />
            </div>
            <div>
              <label htmlFor={`${passwordId}-confirm`}>{t('account.confirm')}</label>
              <PasswordInput id={`${passwordId}-confirm`} className='block w-full' autoComplete='new-password' value={confirm} disabled={busy} onChange={(e) => setConfirm(e.target.value)} />
            </div>
            <p>{t('account.keepPassword')}</p>
            <ErrorMessage message={error ? t(error) : ''} />
            <button className='primary' type='submit' disabled={busy}>{t('account.submit')}</button>
          </form>
        )}
    </Modal>
  );
}
