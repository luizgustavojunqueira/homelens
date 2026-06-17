import { useNavigate } from "react-router-dom";

export function Header() {
  const navigate = useNavigate();
  return (
    <header className="px-6 h-12 flex items-center border-b border-(--border)">
      <div className="flex items-center gap-2.5">
        <span className="dot"></span>
        <span
          className="text-lg font-medium text-(--text) tracking-tight hover:cursor-pointer"
          onClick={() => navigate("/")}
        >
          homelens
        </span>
      </div>
    </header>
  );
}
