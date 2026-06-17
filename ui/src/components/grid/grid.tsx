import type { ReactNode } from "react";

interface IGrid {
  widths?: string[];
  columns: string[];
  children: ReactNode;
}

export default function Grid({ columns, children, widths }: IGrid) {
  return (
    <div className="w-full overflow-x-auto">
      <table className="w-full min-w-max table-auto lg:table-fixed">
        <thead className="border-(--border) border-b-3">
          <tr>
            {columns.map((col, index) => (
              <th
                key={index}
                className={`p-2 px-4 text-left ${widths ? widths[index] : ""}`}
              >
                {col}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>{children}</tbody>
      </table>
    </div>
  );
}
