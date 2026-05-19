import { useRef, useState } from "react";
import { createPortal } from "react-dom";

interface ITooltip {
  content: React.ReactNode;
  children?: React.ReactNode;
}

export default function Tooltip({ content, children }: ITooltip) {
  const [show, setShow] = useState(false);
  const [pos, setPos] = useState({ left: 0, top: 0 });
  const triggerRef = useRef<HTMLDivElement>(null);

  const handleEnter = () => {
    const rect = triggerRef.current?.getBoundingClientRect();
    if (rect) {
      setPos({ left: rect.left + rect.width / 2, top: rect.bottom });
    }
    setShow(true);
  };

  return (
    <>
      <div
        ref={triggerRef}
        className="relative inline-block h-full"
        onMouseEnter={handleEnter}
        onMouseLeave={() => setShow(false)}
      >
        {children}
      </div>
      {createPortal(
        <div
          style={{ left: pos.left, top: pos.top }}
          className={[
            'fixed -translate-x-1/2 mt-2 w-max z-50',
            'bg-black/75 text-(--text) text-sm rounded px-3 py-2',
            'origin-top transition-[scale,opacity] duration-150 ease-out',
            show ? 'scale-y-100 opacity-100' : 'scale-y-0 opacity-0 pointer-events-none',
          ].join(' ')}
        >
          {content}
        </div>,
        document.body
      )}
    </>
  );
}
