import { useState, useMemo } from "react";
import Grid from "../../../components/grid/grid";
import Row from "../../../components/grid/row";
import Cell from "../../../components/grid/cell";
import type { Process } from "../../../api/models";
import { Search, ArrowUpDown } from "lucide-react";

type SortField = "pid" | "cpu" | "memory" | "name";
type SortOrder = "asc" | "desc";

export default function ProcessesTable({ processes }: { processes: Process[] }) {
  const [search, setSearch] = useState("");
  const [sortField, setSortField] = useState<SortField>("cpu");
  const [sortOrder, setSortOrder] = useState<SortOrder>("desc");

  const filteredAndSortedProcesses = useMemo(() => {
    if (!processes) return [];
    
    let result = processes.filter((p) => {
      const term = search.toLowerCase();
      return (
        p.name.toLowerCase().includes(term) ||
        p.command.toLowerCase().includes(term) ||
        p.user.toLowerCase().includes(term) ||
        String(p.pid).includes(term)
      );
    });

    result.sort((a, b) => {
      let valA: number | string = a[sortField];
      let valB: number | string = b[sortField];
      if (typeof valA === "string") valA = valA.toLowerCase();
      if (typeof valB === "string") valB = valB.toLowerCase();

      if (valA < valB) return sortOrder === "asc" ? -1 : 1;
      if (valA > valB) return sortOrder === "asc" ? 1 : -1;
      return 0;
    });

    return result;
  }, [processes, search, sortField, sortOrder]);

  const toggleSort = (field: SortField) => {
    if (sortField === field) {
      setSortOrder(sortOrder === "asc" ? "desc" : "asc");
    } else {
      setSortField(field);
      setSortOrder("desc");
    }
  };

  if (!processes || processes.length === 0) return null;

  return (
    <div className="flex flex-col gap-4 bg-(--bg-elev) rounded-xl w-full overflow-hidden border border-(--border)">
      <div className="px-4 pt-4 pb-2 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
        <h3 className="text-xl font-medium text-(--text)">Top Processes</h3>

        <div className="relative w-full sm:w-64">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-(--text-dim)" />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search PID, name, command..."
            className="w-full pl-9 pr-3 py-1.5 text-sm bg-(--bg) border border-(--border) rounded-md text-(--text) placeholder:text-(--text-dim) focus:outline-none focus:border-sky-500"
          />
        </div>
      </div>

      <Grid
        columns={[
          <button key="pid" onClick={() => toggleSort("pid")} className="flex items-center gap-1 hover:text-sky-400">
            PID <ArrowUpDown className="w-3 h-3" />
          </button>,
          <button key="cpu" onClick={() => toggleSort("cpu")} className="flex items-center gap-1 hover:text-sky-400">
            CPU <ArrowUpDown className="w-3 h-3" />
          </button>,
          <button key="mem" onClick={() => toggleSort("memory")} className="flex items-center gap-1 hover:text-sky-400">
            Mem <ArrowUpDown className="w-3 h-3" />
          </button>,
          "User",
          <button key="name" onClick={() => toggleSort("name")} className="flex items-center gap-1 hover:text-sky-400">
            Name <ArrowUpDown className="w-3 h-3" />
          </button>,
          "Command",
        ]}
        widths={[
          "w-[10%]",
          "w-[10%]",
          "w-[10%]",
          "w-[10%]",
          "w-[15%]",
          "w-[45%]",
        ]}
      >
        {filteredAndSortedProcesses.map((p) => (
          <Row key={p.pid}>
            <Cell>
              <span className="font-mono text-gray-400">{p.pid}</span>
            </Cell>
            <Cell>
              <span className={p.cpu > 50 ? "text-red-500 font-medium" : ""}>
                {p.cpu.toFixed(1)} %
              </span>
            </Cell>
            <Cell>{p.memory.toFixed(1)} %</Cell>
            <Cell>{p.user}</Cell>
            <Cell>
              <span className="block truncate">{p.name}</span>
            </Cell>
            <Cell>
              <span
                title={p.command}
                className="block truncate text-gray-400 font-mono text-sm"
              >
                {p.command}
              </span>
            </Cell>
          </Row>
        ))}
      </Grid>
    </div>
  );
}
