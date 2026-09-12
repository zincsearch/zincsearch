import { Children, isValidElement, useId, useLayoutEffect, useRef, useState, type ButtonHTMLAttributes, type CSSProperties, type KeyboardEvent, type ReactNode } from 'react';
import { createPortal } from 'react-dom';

type Props = Omit<ButtonHTMLAttributes<HTMLButtonElement>, 'children' | 'onChange' | 'value'> & {
  value: string | number;
  onValueChange: (value: string) => void;
  children: ReactNode;
  displayValue?: ReactNode;
};
type OptionProps = { value?: string | number; children: string | number; disabled?: boolean };

// Single-value controlled select. Option children are declarations, never native DOM controls.
export default function Select({ value, onValueChange, children, displayValue, disabled, className = '', id, ...props }: Props) {
  const generatedId = useId();
  const triggerId = id || generatedId;
  const listId = `${triggerId}-list`;
  const trigger = useRef<HTMLButtonElement>(null);
  const list = useRef<HTMLDivElement>(null);
  const [open, setOpen] = useState(false);
  const [active, setActive] = useState(-1);
  const [style, setStyle] = useState<CSSProperties>({});
  const typed = useRef({ text: '', time: 0 });
  const options = Children.toArray(children).filter(isValidElement<OptionProps>).map(child => ({
    value: String(child.props.value ?? child.props.children),
    label: String(child.props.children),
    disabled: !!child.props.disabled,
  }));
  const selected = options.findIndex(option => option.value === String(value));
  const expanded = open && !disabled;
  const activeOption = options[active];

  useLayoutEffect(() => { if (disabled) setOpen(false); }, [disabled]);

  useLayoutEffect(() => {
    if (!expanded) return;
    const button = trigger.current!;
    const menu = list.current!;
    const position = () => {
      const rect = button.getBoundingClientRect();
      const viewport = window.visualViewport;
      const leftEdge = (viewport?.offsetLeft || 0) + 8;
      const topEdge = (viewport?.offsetTop || 0) + 8;
      const rightEdge = leftEdge + (viewport?.width || window.innerWidth) - 16;
      const bottomEdge = topEdge + (viewport?.height || window.innerHeight) - 16;
      const width = Math.min(Math.max(rect.width, 160), rightEdge - leftEdge);
      const below = Math.max(0, bottomEdge - Math.max(topEdge, rect.bottom + 4));
      const above = Math.max(0, Math.min(bottomEdge, rect.top - 4) - topEdge);
      const desired = Math.min(288, menu.scrollHeight);
      const up = below < desired && above > below;
      const height = Math.min(desired, up ? above : below);
      setStyle({
        position: 'fixed', width,
        left: Math.max(leftEdge, Math.min(rect.left, rightEdge - width)),
        top: Math.max(topEdge, Math.min(up ? rect.top - 4 - height : rect.bottom + 4, bottomEdge - height)),
        maxHeight: Math.max(0, Math.min(288, up ? above : below)),
      });
    };
    const outside = (event: PointerEvent) => {
      if (!button.contains(event.target as Node) && !menu.contains(event.target as Node)) setOpen(false);
    };
    position();
    const observer = new ResizeObserver(position);
    observer.observe(button);
    observer.observe(menu);
    window.addEventListener('resize', position);
    window.addEventListener('scroll', position, true);
    window.visualViewport?.addEventListener('resize', position);
    window.visualViewport?.addEventListener('scroll', position);
    document.addEventListener('pointerdown', outside);
    return () => {
      observer.disconnect();
      window.removeEventListener('resize', position);
      window.removeEventListener('scroll', position, true);
      window.visualViewport?.removeEventListener('resize', position);
      window.visualViewport?.removeEventListener('scroll', position);
      document.removeEventListener('pointerdown', outside);
    };
  }, [expanded, children]);

  useLayoutEffect(() => {
    if (!expanded) return;
    const menu = list.current!;
    const option = menu.children[active] as HTMLElement | undefined;
    if (option) {
      if (option.offsetTop < menu.scrollTop) menu.scrollTop = option.offsetTop;
      else if (option.offsetTop + option.offsetHeight > menu.scrollTop + menu.clientHeight) {
        menu.scrollTop = option.offsetTop + option.offsetHeight - menu.clientHeight;
      }
    }
  }, [active, expanded, style.maxHeight]);

  function show(last = false) {
    typed.current = { text: '', time: 0 };
    const enabled = options.map((option, index) => option.disabled ? -1 : index).filter(index => index >= 0);
    setActive(selected >= 0 && !options[selected].disabled ? selected : (last ? enabled[enabled.length - 1] : enabled[0]) ?? -1);
    setOpen(true);
  }
  function choose(index: number) {
    const option = options[index];
    if (!option || option.disabled || disabled) return;
    setOpen(false);
    trigger.current?.focus();
    if (option.value !== String(value)) onValueChange(option.value);
  }
  function keyDown(event: KeyboardEvent<HTMLButtonElement>) {
    if (disabled) return;
    const key = event.key;
    if (key === 'Escape' && expanded) {
      event.preventDefault();
      event.stopPropagation();
      setOpen(false);
    } else if (key === 'Tab') setOpen(false);
    else if (['ArrowDown', 'ArrowUp', 'Home', 'End'].includes(key)) {
      event.preventDefault();
      if (!expanded && (key === 'ArrowDown' || key === 'ArrowUp')) show(key === 'ArrowUp');
      else {
        const step = key === 'ArrowUp' || key === 'End' ? -1 : 1;
        let next = key === 'Home' ? 0 : key === 'End' ? options.length - 1 : active + step;
        while (next >= 0 && next < options.length && options[next].disabled) next += step;
        if (next >= 0 && next < options.length) setActive(next);
        setOpen(true);
      }
    } else if (key === 'Enter' || (key === ' ' && (!typed.current.text || Date.now() - typed.current.time >= 700))) {
      event.preventDefault();
      if (expanded) choose(active);
      else show();
    } else if (key.length === 1 && !event.ctrlKey && !event.metaKey && !event.altKey) {
      event.preventDefault();
      const now = Date.now();
      const text = (now - typed.current.time < 700 ? typed.current.text : '') + key.toLocaleLowerCase();
      typed.current = { text, time: now };
      const prefix = [...text].every(character => character === text[0]) ? text[0] : text;
      const start = expanded ? active : selected;
      for (let offset = 1; offset <= options.length; offset++) {
        const index = (start + offset + options.length) % options.length;
        if (!options[index].disabled && options[index].label.toLocaleLowerCase().startsWith(prefix)) {
          setActive(index);
          break;
        }
      }
      setOpen(true);
    }
  }
  return <>
    <button {...props} id={triggerId} ref={trigger} type="button" role="combobox" value={value}
      disabled={disabled} className={`ui-select ${className}`} aria-haspopup="listbox" aria-expanded={expanded}
      aria-describedby={[`${triggerId}-value`, props['aria-describedby']].filter(Boolean).join(' ')}
      aria-controls={expanded ? listId : undefined}
      aria-activedescendant={expanded && activeOption && !activeOption.disabled ? `${listId}-${active}` : undefined}
      onClick={() => expanded ? setOpen(false) : show()} onKeyDown={keyDown} onBlur={() => setOpen(false)}>
      <span id={`${triggerId}-value`} className={displayValue === undefined ? 'ui-select-value' : 'sr-only'} aria-hidden="true">{options[selected]?.label || '\u00a0'}</span>
      {displayValue !== undefined && <span className="ui-select-value" aria-hidden="true">{displayValue}</span>}
      <svg className="ui-select-arrow" aria-hidden="true" focusable="false" width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
        <path d="m3 4.5 3 3 3-3" />
      </svg>
    </button>
    {expanded && createPortal(<div id={listId} ref={list} role="listbox" aria-labelledby={triggerId}
      className="ui-select-list" style={style} onMouseDown={event => event.preventDefault()}>
      {options.map((option, index) => <div key={option.value} id={`${listId}-${index}`} role="option"
        aria-selected={index === selected} aria-disabled={option.disabled || undefined}
        data-active={index === active} data-value={option.value} className="ui-select-option"
        onPointerMove={() => { if (!option.disabled) setActive(index); }} onClick={event => { event.stopPropagation(); choose(index); }}>
        <span>{option.label}</span><span aria-hidden="true">{index === selected ? '✓' : ''}</span>
      </div>)}
    </div>, document.body)}
  </>;
}
