import Grid from "../../../components/grid/grid";
import Row from "../../../components/grid/row";
import Cell from "../../../components/grid/cell";
import type { Process } from "../../../api/models";

export default function ProcessesTable({ processes }: { processes: Process[] }) {
  if (!processes || processes.length === 0) return null;

  return (
    <div className="flex flex-col gap-4 bg-(--bg-elev) rounded-xl w-full overflow-hidden border border-(--border)">
      <div className="px-4 pt-4 pb-2">
        <h3 className="text-xl font-medium text-(--text)">Top Processes</h3>
      </div>

      <Grid
        columns={["PID", "CPU", "Mem", "User", "Name", "Command"]}
        widths={[
          "w-[8%]",
          "w-[10%]",
          "w-[10%]",
          "w-[10%]",
          "w-[12%]",
          "w-[50%]",
        ]}
      >
        {processes.map((p) => (
          <Row key={p.pid}>
            <Cell>
              <span className="font-mono text-gray-400">{p.pid}</span>
            </Cell>
            <Cell>
              <span
                className={p.cpu > 50 ? "text-red-500 font-medium" : ""}
              >
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
