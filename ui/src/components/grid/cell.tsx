import type { ReactNode } from "react";

type Align = "center" | "left" | "right";

interface ICell {
  children: ReactNode;
  align?: Align;
}

const alignClassMapper: Record<Align, string> = {
  center: "text-center",
  left: "text-left",
  right: "text-end",
};

export default function Cell({ children, align }: ICell) {
  return (
    <td
      className={`p-2 px-4 overflow-hidden text-ellipsis ${align ? alignClassMapper[align] : ""}`}
    >
      {children}
    </td>
  );
}
