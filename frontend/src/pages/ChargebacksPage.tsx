import { useState, useEffect } from "react";
import { DataTable } from "@/components/ui/DataTable";
import { Chargeback, columns } from "@/components/chargebacks/columns";
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetDescription } from "@/components/ui/sheet";
import { DetailsDrawer } from "@/components/DetailsDrawer";

const PAGE_SIZE = 500;

const chargebackStatusOptions = [
  'Open',
  'Hold Pending External Action',
  'Hold Pending Internal Action',
  'In Research',
  'Passed to PFS',
  'Completed by PFS',
  'PFS Return to GSA',
  'New'
];

interface ChargebacksPageProps {
  onUploadSuccess?: () => void; // Optional, as it's not directly used here for refresh
}

export function ChargebacksPage({ onUploadSuccess }: ChargebacksPageProps) {
  const [chargebacks, setChargebacks] = useState<Chargeback[]>([]);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [isDrawerOpen, setIsDrawerOpen] = useState(false);
  const [selectedChargeback, setSelectedChargeback] = useState<Chargeback | null>(null);

  const fetchChargebacks = async () => {
    try {
      const response = await fetch(`http://10.98.1.142:8080/api/chargebacks?limit=${PAGE_SIZE}&page=${page}`);
      const responseData = await response.json();
      
      if (responseData && Array.isArray(responseData.data)) {
        setChargebacks(responseData.data);
        setHasMore(responseData.data.length === PAGE_SIZE);
      } else {
        console.error("API response did not contain a 'data' array:", responseData);
        setChargebacks([]);
        setHasMore(false);
      }

    } catch (error) {
      console.error("Failed to fetch chargebacks:", error);
      setHasMore(false); // Stop trying on error
    }
  };

  useEffect(() => {
    fetchChargebacks();
  }, [page]);

  const handleRowClick = (chargeback: Chargeback) => {
    setSelectedChargeback(chargeback);
    setIsDrawerOpen(true);
  };

  // This effect will run when onUploadSuccess is called, triggering a refresh
  useEffect(() => {
    if (onUploadSuccess) {
      // This is a placeholder. In a real app, you might want a more explicit way to trigger refresh
      // e.g., a refresh button or a state variable that changes.
      // For now, we'll just re-fetch if the prop changes (which it won't, but demonstrates the intent)
      // A better approach would be to pass a 'refreshTrigger' state from App.tsx
      // and have this useEffect depend on it.
    }
  }, [onUploadSuccess]);

  const handleSaveChargeback = async (updatedData: Chargeback) => {
    try {
      if (typeof updatedData.id !== 'number') {
        console.error("Invalid ID for chargeback update:", updatedData.id);
        return;
      }
      const response = await fetch(`http://10.98.1.142:8080/api/chargebacks/${updatedData.id}`, {
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
      fetchChargebacks();
      setIsDrawerOpen(false);
    } catch (error) {
      console.error("Failed to save chargeback:", error);
    }
  };

  const handleCancelChargeback = () => {
    setIsDrawerOpen(false);
  };

  return (
    <div className="space-y-4">
      <DataTable 
        columns={columns} 
        data={chargebacks}
        title="Chargebacks"
        description="A list of recent chargebacks from the live API."
        page={page}
        setPage={setPage}
        hasMore={hasMore}
        onRowClick={handleRowClick}
      />

      <Sheet open={isDrawerOpen} onOpenChange={setIsDrawerOpen}>
        <SheetContent>
          <SheetHeader>
            <SheetTitle>Chargeback Details</SheetTitle>
            <SheetDescription>
              View and manage details for this chargeback.
            </SheetDescription>
          </SheetHeader>
          {selectedChargeback && (
            <DetailsDrawer
              data={selectedChargeback}
              fields={{
                main: [
                  { key: "id", label: "ID" },
                  { key: "bd_doc_num", label: "Document Number" },
                  { key: "customer_name", label: "Customer Name" },
                  { key: "region", label: "Region" },
                  { key: "vendor", label: "Vendor" },
                  { key: "alc", label: "ALC" },
                  { key: "customer_tas", label: "Customer TAS" },
                  { key: "org_code", label: "Org Code" },
                  { key: "chargeback_amount", label: "Chargeback Amount", type: "currency" },
                ],
                status: [
                  { key: "current_status", label: "Current Status", options: chargebackStatusOptions },
                  { key: "gsa_poc", label: "GSA POC" },
                  { key: "pfs_poc", label: "PFS POC" },
                ],
                comments: [],
              }}
              onSave={handleSaveChargeback}
              onCancel={handleCancelChargeback}
            />
          )}
        </SheetContent>
      </Sheet>
    </div>
  );
}
