"use client"

import { ColumnDef } from "@tanstack/react-table"
import { DataTableColumnHeader } from "@/components/ui/DataTableColumnHeader"

export type Chargeback = {
  id: number;
  bd_doc_num: string;
  customer_name: string;
  current_status: string;
  region: string;
  vendor: string;
  alc: string;
  customer_tas: string;
  org_code: string;
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
  {
    accessorKey: "region",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Region" />
    ),
  },
  {
    accessorKey: "vendor",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Vendor" />
    ),
  },
  {
    accessorKey: "alc",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="ALC" />
    ),
  },
  {
    accessorKey: "customer_tas",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Customer TAS" />
    ),
  },
  {
    accessorKey: "org_code",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Org Code" />
    ),
  },
]