import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { toast } from "react-toastify";
import TextInput from "../../components/inputs/TextInput";
import type {
  GetAlertConfigResponse,
  UpdateAlertConfigRequest,
} from "../../api/models";
import { getAlertConfig, saveAlertConfig } from "../../api/alerts";

export default function Settings() {
  const { control, handleSubmit, reset } = useForm<GetAlertConfigResponse>({
    defaultValues: {
      cpu_threshold: 90,
      mem_threshold: 90,
      disk_threshold: 95,
      offline_threshold: 5,
      tolerance_minutes: 5,
      webhook_url: "",
    },
  });

  useEffect(() => {
    getAlertConfig()
      .then((data) => {
        if (data) {
          reset(data);
        }
      })
      .catch(() => {
        toast.error("Failed to load current settings");
      });
  }, [reset]);

  const onSubmit = (data: GetAlertConfigResponse) => {
    const payload: UpdateAlertConfigRequest = {
      ...data,
      cpu_threshold: Number(data.cpu_threshold),
      mem_threshold: Number(data.mem_threshold),
      disk_threshold: Number(data.disk_threshold),
      offline_threshold: Number(data.offline_threshold),
      tolerance_minutes: Number(data.tolerance_minutes),
    };

    saveAlertConfig(payload)
      .then(() => {
        toast.success("Settings saved successfully!");
      })
      .catch(() => {
        toast.error("Failed to save settings");
      });
  };

  return (
    <section className="px-6 py-6 flex-1 overflow-y-auto max-w-screen">
      <div className="mb-6">
        <h2 className="text-2xl font-medium text-(--text)">Alert Settings</h2>
        <p className="text-(--text-dim) mt-1">
          Define the global thresholds for triggering webhook notifications.
        </p>
      </div>

      <div className="border border-(--border) rounded-md bg-(--bg-elev) p-6 max-w-2xl">
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <TextInput
              name="cpu_threshold"
              control={control}
              label="CPU Threshold (%)"
              type="number"
              min="1"
              max="100"
              placeholder="Ex: 90"
            />

            <TextInput
              name="mem_threshold"
              control={control}
              label="Memory Threshold (%)"
              type="number"
              min="1"
              max="100"
              placeholder="Ex: 90"
            />

            <TextInput
              name="disk_threshold"
              control={control}
              label="Disk Threshold (%)"
              type="number"
              min="1"
              max="100"
              placeholder="Ex: 95"
            />

            <TextInput
              name="offline_threshold"
              control={control}
              label="Offline Threshold (minutes)"
              type="number"
              min="1"
              placeholder="Ex: 5"
            />

            <TextInput
              name="tolerance_minutes"
              control={control}
              label="Tolerance Minutes"
              type="number"
              min="1"
              placeholder="Ex: 5"
            />
          </div>

          <div className="pt-2">
            <TextInput
              name="webhook_url"
              control={control}
              label="Webhook URL (Telegram/Discord)"
              type="url"
              placeholder="https://api.telegram.org/bot.../sendMessage"
            />
          </div>

          <div className="pt-4 border-t border-(--border) flex justify-end">
            <button
              type="submit"
              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md font-medium transition-colors cursor-pointer"
            >
              Save Settings
            </button>
          </div>
        </form>
      </div>
    </section>
  );
}
