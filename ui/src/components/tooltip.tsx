import { useState } from "react";

interface ITooltip {
  content: React.ReactNode;
  children?: React.ReactNode;
}

export default function Tooltip({ content, children }: ITooltip) {
  const [show, setShow] = useState(false);

  return (
    <div
      className="relative inline-block h-full"
      onMouseEnter={() => setShow(true)}
      onMouseLeave={() => setShow(false)}
    >
      <div
        className={[
          'absolute top-full left-1/2 -translate-x-1/2 mt-2 w-max z-10',
          'bg-black/75 text-(--text) text-sm rounded px-3 py-2',
          'origin-top transition-[scale,opacity] duration-150 ease-out',
          show ? 'scale-y-100 opacity-100' : 'scale-y-0 opacity-0 pointer-events-none',
        ].join(' ')}
      >
        {content}
      </div>
      {children}
    </div>
  )
}
