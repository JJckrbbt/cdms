import { useState, useEffect } from "react";
import { DataTable } from "@/components/ui/DataTable";
import { Chargeback, columns } from "@/components/chargebacks/columns";
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetDescription } from "@/components/ui/sheet";
import { DetailsDrawer } from "@/components/DetailsDrawer";
import { useAuth0 } from "@auth0/auth0-react";
import { apiClient } from "@/lib/api";

const PAGE_SIZE = 500;

const chargebackStatusOptions = [
  'Open',
  'In Research',
  'Hold Pending Internal Action',
  'Hold Pending External Action',
  'Passed to PFS',
  'PFS Return to GSA',
  'Completed by PFS',
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
  const { getAccessTokenSilently, isAuthenticated } = useAuth0();

  const fetchChargebacks = async () => {
    try {
      const token = await getAccessTokenSilently({
        authorizationParams: {
          audience: import.meta.env.VITE_AUTH0_AUDIENCE,
        },
      });

      console.log("THE TOKEN IS: ", token);
      console.log("TOKEN TYPE IS: ", typeof token);

      const responseData = await apiClient.get(`/api/chargebacks?limit=${PAGE_SIZE}&page=${page}`, token);

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
  if (isAuthenticated) {
    fetchChargebacks();
  }
}, [page, isAuthenticated]);

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
      const payload = {
        bd_doc_num: updatedData.bd_doc_num,
        customer_name: updatedData.customer_name,
        current_status: updatedData.current_status,
        region: updatedData.region,
        vendor: updatedData.vendor,
        alc: updatedData.alc,
        customer_tas: updatedData.customer_tas,
        org_code: updatedData.org_code,
        gsa_poc: updatedData.gsa_poc,
        pfs_poc: updatedData.pfs_poc,
        chargeback_amount: updatedData.chargeback_amount,
      };

      const token = await getAccessTokenSilently({
        authorizationParams: {
          audience: import.meta.env.VITE_AUTH0_AUDIENCE,
        },
      });
      await apiClient.patch(`/api/chargebacks/${updatedData.id}`, payload);

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
              id={selectedChargeback.id}
              type="chargeback"
            />
          )}
        </SheetContent>
      </Sheet>
    </div>
  );
}
