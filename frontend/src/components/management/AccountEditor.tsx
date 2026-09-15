import { useId, useState } from 'react';
import Select from '../Select';
import PasswordInput from '../PasswordInput';
import userService from '../../services/user';
import roleService from '../../services/role';
import permissionService from '../../services/permission';
import { useTranslation } from '../../locales';
import { validatePassword } from '../../utils/password';
import { ErrorMessage, Modal, useList } from './Common';

export type Account = {
  _id: string;
  name: string;
  role?: string;
  permission?: string[];
  created_at?: string;
  updated_at?: string;
};

export default function AccountEditor(
  { kind, value, onClose, onUpdated }: {
    kind: 'user' | 'role';
    value?: Account;
    onClose: () => void;
    onUpdated: () => void;
  },
) {
  const passwordId = useId();
  const { t } = useTranslation();
  const [id, setId] = useState(value?._id || '');
  const [name, setName] = useState(value?.name || value?._id || '');
  const [role, setRole] = useState(value?.role || '');
  const [permissions, setPermissions] = useState<string[]>(value?.permission || []);
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);
  const options = useList<Account | string>(
    kind === 'user' ? roleService.list : permissionService.list,
  );
  const title = `${value ? 'Update' : 'Add'} ${kind}`;
  async function save() {
    const noun = kind === 'user' ? 'User' : 'Role';
    if (id.trim().length < 3) {
      setError(`${noun} ID must be at least 3 characters long`);
      return;
    }
    if (name.trim().length < 3) {
      setError(`${noun} name must be at least 3 characters long`);
      return;
    }
    if (kind === 'user') {
      if (!role) {
        setError('You must select a role');
        return;
      }
      // existing users keep their password when both fields stay blank
      const message = value && !password && !confirmPassword
        ? ''
        : validatePassword(password, confirmPassword);
      if (message) {
        setError(message);
        return;
      }
    }
    setError('');
    setBusy(true);
    try {
      if (kind === 'user') {
        await userService.update({ _id: id, name, role, password, confirmPassword });
      } else await roleService.update({ _id: id, name, permission: permissions });
      onUpdated();
    } catch {
      setError(`Unable to save ${kind}. Please try again.`);
    } finally {
      setBusy(false);
    }
  }
  const roles = options.rows.filter((row): row is Account => typeof row !== 'string');
  const permissionOptions = Array.from(
    new Set([
      ...options.rows.filter((row): row is string => typeof row === 'string'),
      ...permissions,
    ]),
  );
  return (
    <Modal
      title={title}
      onClose={() => {
        if (!busy) onClose();
      }}
    >
      <form
        className='grid gap-4'
        onSubmit={(event) => {
          event.preventDefault();
          void save();
        }}
      >
        <label>
          {t(`${kind}.id`)}
          <input
            className='block w-full'
            value={id}
            disabled={!!value || busy}
            onChange={(e) => setId(e.target.value)}
          />
        </label>
        <label>
          {t(`${kind}.name`)}
          <input
            className='block w-full'
            value={name}
            disabled={busy}
            onChange={(e) => setName(e.target.value)}
          />
        </label>
        <ErrorMessage message={options.error} />
        {options.error && <button type='button' onClick={options.refresh}>Retry options</button>}
        {kind === 'user'
          ? (
            <>
              <label>
                {t('user.role')}
                <Select
                  className='block w-full'
                  value={role}
                  disabled={options.loading || busy}
                  onValueChange={setRole}
                >
                  <option value=''>Select role</option>
                  {roles.filter((r) => r._id !== 'admin').map((r) => (
                    <option key={r._id} value={r._id}>{r.name || r._id}</option>
                  ))}
                  <option value='admin'>admin</option>
                  {role && role !== 'admin' && !roles.some((r) => r._id === role) && (
                    <option value={role}>{role}</option>
                  )}
                </Select>
              </label>
              <div>
                <label htmlFor={`${passwordId}-new`}>{t('user.password')}</label>
                <PasswordInput
                  id={`${passwordId}-new`}
                  className='block w-full'
                  autoComplete='new-password'
                  value={password}
                  disabled={busy}
                  onChange={(e) => setPassword(e.target.value)}
                />
              </div>
              {value && <p>Leave password blank to keep the current password.</p>}
              <div>
                <label htmlFor={`${passwordId}-confirm`}>{t('user.repassword')}</label>
                <PasswordInput
                  id={`${passwordId}-confirm`}
                  className='block w-full'
                  autoComplete='new-password'
                  value={confirmPassword}
                  disabled={busy}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                />
              </div>
            </>
          )
          : (
            <fieldset disabled={busy || options.loading}>
              <legend>{t('role.permission')}</legend>
              {permissionOptions.map((permission) => (
                <label className='flex gap-2 p-2' key={permission}>
                  <input
                    type='checkbox'
                    checked={permissions.includes(permission)}
                    onChange={(e) =>
                      setPermissions((current) =>
                        e.target.checked
                          ? [...current, permission]
                          : current.filter((p) =>
                            p !== permission
                          )
                      )}
                  />
                  {permission}
                </label>
              ))}
            </fieldset>
          )}
        <ErrorMessage message={error} />
        <button
          className='primary'
          type='submit'
          disabled={busy || options.loading || !!options.error}
        >
          Save {kind === 'user' ? 'User' : 'Role'}
        </button>
      </form>
    </Modal>
  );
}
