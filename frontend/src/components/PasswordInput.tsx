import { useId, useState, type ComponentProps } from 'react';
import { useTranslation } from '../locales';
import Icon from './Icon';

export default function PasswordInput({ className = '', id, disabled, ...props }: Omit<ComponentProps<'input'>, 'type'>) {
  const generatedId = useId();
  const [visible, setVisible] = useState(false);
  const { t } = useTranslation();
  const inputId = id || generatedId;
  const toggleLabel = t(visible ? 'passwordInput.hide' : 'passwordInput.show');
  return <span className={`password-input ${className}`}>
    <input {...props} id={inputId} disabled={disabled} type={visible ? 'text' : 'password'} />
    <button type="button" className="password-toggle" aria-label={toggleLabel} title={toggleLabel} aria-controls={inputId} disabled={disabled} onClick={() => setVisible(!visible)}>
      <Icon name={visible ? 'eyeOff' : 'eye'} />
    </button>
  </span>;
}
