import { useNavigate } from "react-router-dom";
import { Settings } from "lucide-react";

export function Header() {
  const navigate = useNavigate();
  return (
    <header className="px-6 h-12 flex items-center justify-between border-b border-(--border)">
      <div className="flex items-center gap-2.5">
        <span className="dot"></span>
        <span
          className="text-lg font-medium text-(--text) tracking-tight hover:cursor-pointer"
          onClick={() => navigate("/")}
        >
          homelens
        </span>
      </div>

      <button
        onClick={() => navigate("/settings")}
        className="text-(--text-dim) hover:text-(--text) transition-colors cursor-pointer"
        title="Configurações"
      >
        <Settings size={20} />
      </button>
    </header>
  );
}
