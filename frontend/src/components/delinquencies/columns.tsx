import { ColumnDef } from "@tanstack/react-table";
import { DataTableColumnHeader } from "@/components/ui/DataTableColumnHeader";

export type Delinquency = {
  id: number;
  business_line: string;
  document_number: string; // Document Number
  vendor_code: string; // Vendor Code
  status: string; // Current Status as Status
  billed_total_amount: number;
  debit_outstanding_amount: number; // Debit Outstanding Amount
  credit_outstanding_amount: number; // Credit Outstanding Amount
};

export const columns: ColumnDef<Delinquency>[] = [
  {
    accessorKey: "business_line",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Business Line" />
    ),
  },
  {
    accessorKey: "document_number",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Document Number" />
    ),
  },
  {
    accessorKey: "vendor_code",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Vendor Code" />
    ),
  },
  {
    accessorKey: "status",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Status" />
    ),
  },
  {
    accessorKey: "billed_total_amount",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Billed Total Amount" />
    ),
  },
  {
    accessorKey: "debit_outstanding_amount",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Debit Outstanding Amount" />
    ),
  },
  {
    accessorKey: "credit_outstanding_amount",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title="Credit Outstanding Amount" />
    ),
  },
];
