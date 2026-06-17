import type { ReactNode } from "react";

interface IRow {
  children: ReactNode;
}

export default function Row({ children }: IRow) {
  return <tr className="border-b-1 border-(--border) ">{children}</tr>;
}
