export function Header() {
  return (
    <header className="px-6 h-12 flex items-center border-b border-(--border)">
      <div className="flex items-center gap-2.5">
        <span className="dot"></span>
        <span className="text-lg font-medium text-(--text) tracking-tight">
          homelens
        </span>
      </div>
    </header>
  );
}
