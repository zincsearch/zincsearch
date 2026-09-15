import { useEffect, type RefObject } from 'react';

// Closes a popover when pointing outside `ref`, including only its own portaled Select menus as inside.
export default function useClickOutside(ref: RefObject<HTMLElement | null>, active: boolean, onOutside: () => void) {
  useEffect(() => {
    if (!active) return;
    const outside = (event: PointerEvent) => {
      const target = event.target as Element;
      if (ref.current?.contains(target)) return;
      const listbox = target.closest('[role="listbox"]');
      if (listbox?.id && Array.from(ref.current?.querySelectorAll('[role="combobox"][aria-controls]') ?? [])
        .some(trigger => trigger.getAttribute('aria-controls') === listbox.id)) return;
      onOutside();
    };
    document.addEventListener('pointerdown', outside);
    return () => document.removeEventListener('pointerdown', outside);
  }, [ref, active, onOutside]);
}
