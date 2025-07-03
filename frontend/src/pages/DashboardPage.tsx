import { useState, useEffect } from "react";
import { StatCard } from "@/components/StatCard";
import { DollarSign, Scale, ReceiptText } from "lucide-react";

export function DashboardPage() {
  const [totalDelinquencies, setTotalDelinquencies] = useState<number | null>(null);
  const [totalChargebacks, setTotalChargebacks] = useState<number | null>(null);

  useEffect(() => {
    const fetchTotals = async () => {
      try {
        const [delinquenciesResponse, chargebacksResponse] = await Promise.all([
          fetch(`http://10.98.1.142:8080/api/delinquencies?limit=1&page=1`),
          fetch(`http://10.98.1.142:8080/api/chargebacks?limit=1&page=1`),
        ]);

        const delinquenciesData = await delinquenciesResponse.json();
        const chargebacksData = await chargebacksResponse.json();

        if (delinquenciesData && typeof delinquenciesData.total_count === 'number') {
          setTotalDelinquencies(delinquenciesData.total_count);
        }
        if (chargebacksData && typeof chargebacksData.total_count === 'number') {
          setTotalChargebacks(chargebacksData.total_count);
        }

      } catch (error) {
        console.error("Failed to fetch totals:", error);
      }
    };

    fetchTotals();
  }, []);

  return (
    <div className="p-6">
      <h1 className="text-3xl font-bold">Dashboard</h1>
      <p className="text-muted-foreground">Welcome to your dashboard!</p>
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4 mt-4">
        <StatCard
          title="Total Active Delinquencies"
          value={totalDelinquencies !== null ? totalDelinquencies.toString() : "Loading..."}
          description="Total number of active delinquencies"
          icon={<Scale className="h-4 w-4 text-muted-foreground" />}
        />
        <StatCard
          title="Total Active Chargebacks"
          value={totalChargebacks !== null ? totalChargebacks.toString() : "Loading..."}
          description="Total number of active chargebacks"
          icon={<ReceiptText className="h-4 w-4 text-muted-foreground" />}
        />
      </div>
    </div>
  );
}
