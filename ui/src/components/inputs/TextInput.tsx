import { useRef } from "react";
import {
  Controller,
  type Control,
  type FieldValues,
  type Path,
} from "react-hook-form";

interface InputProps<
  T extends FieldValues,
> extends React.InputHTMLAttributes<HTMLInputElement> {
  name: Path<T>;
  control: Control<T>;
  label?: string;
  debounceTime?: number;
  onDebounce?: (value: string) => void;
}

export default function TextInput<T extends FieldValues>({
  name,
  control,
  label,
  debounceTime = 500,
  onDebounce,
  className = "",
  ...rest
}: InputProps<T>) {
  const timerRef = useRef<number | null>(null);

  return (
    <Controller
      name={name}
      control={control}
      render={({ field, fieldState }) => {
        const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
          field.onChange(e);

          if (onDebounce) {
            if (timerRef.current) clearTimeout(timerRef.current);
            timerRef.current = setTimeout(() => {
              onDebounce(e.target.value);
            }, debounceTime);
          }
        };

        return (
          <div className="flex flex-col w-full">
            {label && (
              <label className="text-sm text-(--text) mb-1 font-medium">
                {label}
              </label>
            )}
            <input
              {...field}
              {...rest}
              onChange={handleChange}
              className={`px-3 py-1.5 bg-(--bg-elev) border border-(--border) rounded-md text-(--text) focus:outline-none focus:ring-2 focus:ring-blue-500 transition-shadow ${
                fieldState.error ? "border-red-500" : ""
              } ${className}`}
            />
            {fieldState.error && (
              <span className="text-xs text-red-500 mt-1">
                {fieldState.error.message}
              </span>
            )}
          </div>
        );
      }}
    />
  );
}
