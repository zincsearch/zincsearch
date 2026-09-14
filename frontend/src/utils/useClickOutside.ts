import { useEffect, type RefObject } from 'react';

// Closes a popover when pointing outside `ref`. Portaled Select menus are treated as inside.
export default function useClickOutside(ref: RefObject<HTMLElement | null>, active: boolean, onOutside: () => void) {
  useEffect(() => {
    if (!active) return;
    const outside = (event: PointerEvent) => {
      const target = event.target as Element;
      if (!ref.current?.contains(target) && !target.closest('[role="listbox"]')) onOutside();
    };
    document.addEventListener('pointerdown', outside);
    return () => document.removeEventListener('pointerdown', outside);
  }, [ref, active, onOutside]);
}
