import { useEffect } from "react";
import { useParams } from "react-router-dom";
import { getSnapshots } from "../../api/snapshots";
import Gauge from "../../components/charts/gauge";
import { convertByteToMetric, formatByteStr } from "../../utils";
import Line from "../../components/charts/line";
import { useAgents } from "../../store/agentsStore";
import { getMultiSeries, getSeries } from "./agentDetailsUtils";
import Grid from "../../components/grid/grid";
import Row from "../../components/grid/row";
import Cell from "../../components/grid/cell";
import { useForm } from "react-hook-form";
import TextInput from "../../components/inputs/TextInput";
import { updateName } from "../../api/agents";
import { toast } from "react-toastify";

interface AgentForm {
  agentName: string;
}

export default function AgentDetails() {
  const params = useParams();
  const agentGuid = params.guid;

  const agent = useAgents((state) => state.getAgentState(agentGuid!));

  const { control, reset } = useForm<AgentForm>({
    defaultValues: { agentName: agent?.name || "" },
  });

  useEffect(() => {
    if (agent?.name) {
      reset({ agentName: agent.name });
    }
  }, [agent?.name, reset]);

  useEffect(() => {
    if (!agentGuid) return;
    getSnapshots(agentGuid).then((res) => {
      useAgents.getState().insertHistory(agentGuid, res.snapshots);
    });
  }, [agentGuid]);

  if (agent === undefined) return;

  const handleNameUpdate = (newName: string) => {
    if (newName.trim() === "" || newName === agent.name) return;

    updateName({ name: newName, guid: agent.guid }).then((res) => {
      if (res) {
        toast.success("Agent name changed");
      } else {
        toast.error("Error changing agent name");
      }
    });
  };

  const latestData = agent.latest_snapshot.data;
  const agentIp = latestData.agent_ip;

  const timestamps = agent.history.map((snap) => snap.timestamp) ?? [];

  const currentDiskUsage = latestData.disk.disk_space.usage_percent;
  const currentDiskUsed = latestData.disk.disk_space.used;
  const currentDiskTotal = latestData.disk.disk_space.total;

  const diskTotalIoHistory = getMultiSeries(agent.history, (snap) =>
    snap.data.disk.disk_io_usage.map((io) => io.read_mbps + io.write_mbps),
  );

  const diskNames = latestData.disk.disk_io_usage.map((disk) => disk.name);

  const currentCpuAvgUsage =
    latestData.cpu.reduce((cum, curr) => cum + curr.usage_percent, 0) /
    latestData.cpu.length;
  const cpusHistory = getMultiSeries(agent.history, (snap) =>
    snap.data.cpu.map((cpu) => cpu.usage_percent),
  );
  const cpuAvgHistory = getSeries(
    agent.history,
    (snap) =>
      snap.data.cpu.reduce((cum, curr) => cum + curr.usage_percent, 0) /
      snap.data.cpu.length,
  );

  const currentMemUsed = latestData.memory.used ?? 0;
  const currentMemTotal = latestData.memory.total ?? 0;
  const currentMemUsage = (currentMemUsed / currentMemTotal) * 100;
  const currentTemp = latestData.temperature
    ? latestData.temperature.reduce((cum, curr) => cum + curr.temp_c, 0) /
      latestData.temperature.length
    : undefined;
  const memUsedHistory = agent.history.map((snap) =>
    convertByteToMetric(snap.data.memory.used, "GB", "KB"),
  );

  const netRxHistory = getSeries(agent.history, (snap) =>
    convertByteToMetric(
      snap.data.network.reduce((sum, net) => sum + net.rx_bps, 0),
      "MB",
    ),
  );

  const netTxHistory = getSeries(agent.history, (snap) =>
    convertByteToMetric(
      snap.data.network.reduce((sum, net) => sum + net.tx_bps, 0),
      "MB",
    ),
  );

  return (
    <section className="px-6 py-6 flex-1 overflow-y-auto max-w-screen">
      <div className="flex items-baseline gap-3 mb-6 w-full max-w-sm">
        <TextInput
          name="agentName"
          control={control}
          onDebounce={handleNameUpdate}
          debounceTime={800}
          placeholder="Nome do Agente"
          className="text-xl font-medium w-full"
        />
      </div>

      <div className="space-y-6">
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-4">
          <div className="border border-(--border) rounded-md bg-(--bg-elev) p-4 h-64">
            <Gauge
              value={currentDiskUsage}
              label="Disk Usage"
              total={formatByteStr(currentDiskTotal)}
              used={formatByteStr(currentDiskUsed)}
            />
          </div>

          <div className="border border-(--border) rounded-md bg-(--bg-elev) p-4 h-64">
            <Gauge
              value={currentCpuAvgUsage}
              label="CPU Usage"
              total="100 %"
              used={`${currentCpuAvgUsage.toFixed(2)} %`}
            />
          </div>

          <div className="border border-(--border) rounded-md bg-(--bg-elev) p-4 h-64">
            <Gauge
              value={currentMemUsage}
              label="RAM Usage"
              total={formatByteStr(currentMemTotal, "KB")}
              used={formatByteStr(currentMemUsed, "KB")}
            />
          </div>

          {currentTemp && (
            <div className="border border-(--border) rounded-md bg-(--bg-elev) p-4 h-64">
              <Gauge
                value={currentTemp}
                label="Temperature"
                symbol="C°"
                total={"100 C°"}
                used={`${currentTemp.toFixed(2)} C°`}
              />
            </div>
          )}
        </div>
        <div className="border border-(--border) rounded-md bg-(--bg-elev)"></div>
        <div className="grid grid-cols-1 xl:grid-cols-2 grid-rows-2 gap-4">
          <div className="border border-(--border) rounded-md bg-(--bg-elev)">
            <div className="px-4 py-3 border-b border-(--border)">
              <h3 className="text-sm font-medium text-(--text)">
                CPU Usage History
              </h3>
            </div>

            <div className="h-96 p-2">
              <Line
                isTotalAverage={true}
                timestamps={timestamps}
                label="CPU Usage (%)"
                valueFormatter={(v) => `${v.toFixed(2)} %`}
                series={[
                  {
                    name: "Average",
                    values: cpuAvgHistory,
                  },
                  ...cpusHistory.map((cpu, i) => ({
                    name: `CPU ${i}`,
                    values: cpu,
                    visible: false,
                  })),
                ]}
              />
            </div>
          </div>

          <div className="border border-(--border) rounded-md bg-(--bg-elev)">
            <div className="px-4 py-3 border-b border-(--border)">
              <h3 className="text-sm font-medium text-(--text)">
                Memory Usage History
              </h3>
            </div>

            <div className="h-80 p-2">
              <Line
                valueFormatter={(value) => formatByteStr(value, "GB")}
                timestamps={timestamps}
                label="RAM Used (GB)"
                series={[
                  {
                    values: memUsedHistory,
                    name: "RAM",
                  },
                ]}
              />
            </div>
          </div>

          <div className="border border-(--border) rounded-md bg-(--bg-elev)">
            <div className="px-4 py-3 border-b border-(--border)">
              <h3 className="text-sm font-medium text-(--text)">
                Disk Usage History
              </h3>
            </div>

            <div className="h-80 p-2">
              <Line
                timestamps={timestamps}
                label="Disk IO (MB/s)"
                valueFormatter={(v) => `${v.toFixed(2)} MB/s`}
                series={diskTotalIoHistory.map((disk, i) => ({
                  name: diskNames[i],
                  values: disk,
                }))}
              />
            </div>
          </div>

          <div className="border border-(--border) rounded-md bg-(--bg-elev)">
            <div className="px-4 py-3 border-b border-(--border)">
              <h3 className="text-sm font-medium text-(--text)">
                Network Throughput
              </h3>
            </div>

            <div className="h-80 p-2">
              <Line
                timestamps={timestamps}
                label="Network Throughput"
                valueFormatter={(v) => `${v.toFixed(2)} MB/s`}
                series={[
                  {
                    name: "RX",
                    values: netRxHistory,
                  },
                  {
                    name: "TX",
                    values: netTxHistory,
                  },
                ]}
              />
            </div>
          </div>
        </div>
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
            {latestData.containers?.map((c) => (
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
                      .map((p, index) => (
                        <a
                          key={index}
                          href={`http://${agentIp}:${p.public_port}`}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="px-2 py-1 text-xs font-mono text-(--text) bg-blue-500/10 hover:bg-blue-500/20 border border-blue-500/20 rounded transition-colors"
                          title={`Access port ${p.public_port}`}
                        >
                          {p.private_port}:{p.public_port}
                        </a>
                      ))}

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
            {latestData.processes?.map((p) => (
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
        </div>{" "}
      </div>
    </section>
  );
}
