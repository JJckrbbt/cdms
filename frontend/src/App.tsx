import { useState, useEffect } from "react";
import { DashboardLayout } from "./components/DashboardLayout";
import { DataTable } from "./components/ui/DataTable";
import { Chargeback, columns } from "./components/chargebacks/columns";

function App() {
  const [chargebacks, setChargebacks] = useState<Chargeback[]>([]);

  useEffect(() => {
    async function fetchChargebacks() {
      try {
        const response = await fetch("http://10.98.1.142:8080/api/chargebacks");
        const responseData = await response.json();
        
        // --- THIS IS THE FIX ---
        // We now correctly access the 'data' array inside the response object.
        if (responseData && Array.isArray(responseData.data)) {
          setChargebacks(responseData.data);
        } else {
          console.error("API response did not contain a 'data' array:", responseData);
          setChargebacks([]); // Default to an empty array if the structure is wrong
        }

      } catch (error) {
        console.error("Failed to fetch chargebacks:", error);
      }
    }

    fetchChargebacks();
  }, []);

  return (
    <DashboardLayout>
      <div className="space-y-4">
        <h1 className="text-3xl font-bold">Chargebacks</h1>
        <p className="text-sm text-muted-foreground">
          A list of recent chargebacks from the live API.
        </p>
        
        <DataTable columns={columns} data={chargebacks} />

      </div>
    </DashboardLayout>
  );
}

export default App;