import { useState, useEffect } from "react";
import { DataTable } from "@/components/ui/DataTable";
import { Delinquency, columns } from "@/components/delinquencies/columns";
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetDescription } from "@/components/ui/sheet";
import { DetailsDrawer } from "@/components/DetailsDrawer";

const statusOptions = [
  'Open',
  'Refund',
  'Offset',
  'In Process',
  'Write Off',
  'Referred to Treasury for Collections',
  'Return Credit to Treasury',
  'Waiting on Customer Response',
  'Waiting on GSA Response Pending Payment',
  'Closed - Payment Received',
  'Reverse to Income',
  'Bill as IPAC',
  'Bill as DoD',
  'EIS Issues'
];

const PAGE_SIZE = 500;

export function DelinquenciesPage() {
  const [delinquencies, setDelinquencies] = useState<Delinquency[]>([]);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [isDrawerOpen, setIsDrawerOpen] = useState(false);
  const [selectedDelinquency, setSelectedDelinquency] = useState<Delinquency | null>(null);

  const fetchDelinquencies = async () => {
    try {
      // Assuming a similar API endpoint for delinquencies
      const response = await fetch(`http://10.98.1.142:8080/api/delinquencies?limit=${PAGE_SIZE}&page=${page}`);
      const responseData = await response.json();
      
      if (responseData && Array.isArray(responseData.data)) {
        setDelinquencies(responseData.data);
        setHasMore(responseData.data.length === PAGE_SIZE);
      } else {
        console.error("API response did not contain a 'data' array:", responseData);
        setDelinquencies([]);
        setHasMore(false);
      }

    } catch (error) {
      console.error("Failed to fetch delinquencies:", error);
      setHasMore(false); // Stop trying on error
    }
  };

  useEffect(() => {
    fetchDelinquencies();
  }, [page]);

  const handleRowClick = (delinquency: Delinquency) => {
    setSelectedDelinquency(delinquency);
    setIsDrawerOpen(true);
  };

  const handleSaveDelinquency = async (updatedData: Delinquency) => {
    try {
      if (typeof updatedData.id !== 'number') {
        console.error("Invalid ID for delinquency update:", updatedData.id);
        return;
      }
      const response = await fetch(`http://10.98.1.142:8080/api/delinquencies/${updatedData.id}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(updatedData),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      // Refresh the data after successful update
      fetchDelinquencies();
      setIsDrawerOpen(false);
    } catch (error) {
      console.error("Failed to save delinquency:", error);
    }
  };

  const handleCancelDelinquency = () => {
    setIsDrawerOpen(false);
  };

  return (
    <div className="space-y-4">
      <DataTable 
        columns={columns} 
        data={delinquencies}
        title="Delinquencies"
        description="A list of recent delinquencies from the live API."
        page={page}
        setPage={setPage}
        hasMore={hasMore}
        onRowClick={handleRowClick}
      />

      <Sheet open={isDrawerOpen} onOpenChange={setIsDrawerOpen}>
        <SheetContent>
          <SheetHeader>
            <SheetTitle>Delinquency Details</SheetTitle>
            <SheetDescription>
              View and manage details for this delinquency.
            </SheetDescription>
          </SheetHeader>
          {selectedDelinquency && (
            <DetailsDrawer
              data={selectedDelinquency}
              fields={{
                main: [
                  { key: "business_line", label: "Business Line" },
                  { key: "document_number", label: "Document Number" },
                  { key: "vendor_code", label: "Vendor Code" },
                  { key: "billed_total_amount", label: "Billed Total Amount", type: "currency" },
                  { key: "debit_outstanding_amount", label: "Debit Outstanding Amount", type: "currency" },
                  { key: "credit_outstanding_amount", label: "Credit Outstanding Amount", type: "currency" },
                ],
                status: [
                  { key: "current_status", label: "Current Status", options: statusOptions },
                  { key: "gsa_poc", label: "GSA POC" },
                  { key: "pfs_poc", label: "PFS POC" },
                ],
                comments: [],
              }}
              onSave={handleSaveDelinquency}
              onCancel={handleCancelDelinquency}
            />
          )}
        </SheetContent>
      </Sheet>
    </div>
  );
}