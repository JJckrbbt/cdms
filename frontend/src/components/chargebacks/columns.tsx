"use client"

import { ColumnDef } from "@tanstack/react-table"
import { DataTableColumnHeader } from "@/components/ui/DataTableColumnHeader"

export type Chargeback = {
  id: number;
  bd_doc_num: string;
  customer_name: string;
  current_status: string;
};

export const columns: ColumnDef<Chargeback>[] = [
  {
    accessorKey: "current_status",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Status" />
    ),
  },
  {
    accessorKey: "bd_doc_num",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Document Number" />
    ),
  },
  {
    accessorKey: "customer_name",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Customer Name" />
    ),
  },
]