import Grid from "../../../components/grid/grid";
import Row from "../../../components/grid/row";
import Cell from "../../../components/grid/cell";
import type { DockerContainer } from "../../../api/models";

export default function DockerContainersTable({ containers, agentIp }: { containers: DockerContainer[], agentIp: string }) {
  const isPrivateIp = (ip: string) => {
    return /^(10\.|172\.(1[6-9]|2[0-9]|3[0-1])\.|192\.168\.|127\.0\.0\.1)/.test(ip);
  };
  const safeIp = isPrivateIp(agentIp) ? agentIp : null;

  if (!containers || containers.length === 0) return null;

  return (
    <div className="flex flex-col gap-4 bg-(--bg-elev) rounded-xl w-full overflow-hidden border border-(--border)">
      <div className="px-3 pt-3">
        <h3 className="text-xl font-medium text-(--text)">
          Docker Containers
        </h3>
      </div>
      <Grid
        columns={["Name", "Image", "State", "Status", "Ports"]}
        widths={["w-[20%]", "w-[30%]", "w-[10%]", "w-[15%]", "w-[25%]"]}
      >
        {containers.map((c) => (
          <Row key={c.name}>
            <Cell>{c.name}</Cell>
            <Cell>{c.image}</Cell>
            <Cell>
              <span
                className={`${c.state === "running" ? "text-green-500" : "text-red-500"}`}
              >
                {c.state}
              </span>
            </Cell>
            <Cell>{c.status}</Cell>

            <Cell>
              <div className="flex flex-wrap gap-2">
                {c.ports
                  ?.filter((p) => p.public_port)
                  .map((p, index) => {
                    if (!safeIp) {
                      return <span key={index} className="px-2 py-1 text-xs font-mono text-(--text) bg-gray-500/10 border border-gray-500/20 rounded">
                        {p.private_port}:{p.public_port}
                      </span>
                    }
                    return (
                    <a
                      key={index}
                      href={`http://${safeIp}:${p.public_port}`}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="px-2 py-1 text-xs font-mono text-(--text) bg-blue-500/10 hover:bg-blue-500/20 border border-blue-500/20 rounded transition-colors"
                      title={`Access port ${p.public_port}`}
                    >
                      {p.private_port}:{p.public_port}
                    </a>
                  )})}

                {(!c.ports ||
                  c.ports.filter((p) => p.public_port).length === 0) && (
                  <span className="text-xs text-gray-500 italic">-</span>
                )}
              </div>
            </Cell>
          </Row>
        ))}
      </Grid>
    </div>
  );
}
