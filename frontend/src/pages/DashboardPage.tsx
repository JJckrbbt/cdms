import { useState, useEffect } from "react";
import { ReportingTable } from "@/components/ReportingTable";
import { StatCard } from "@/components/StatCard";
import { Scale } from "lucide-react";
import { formatCurrency } from "@/lib/utils";

// --- Type definitions for our API data ---
interface TimeWindowStats {
  new_items_count: number;
  new_items_value: string;
  avg_days_to_pfs: number;
  avg_days_for_pfs_complete: number;
  passed_to_pfs: number;
  completed_by_pfs: number;
}

interface ChargebackStats {
  status_summary: {
    current_status: string;
    status_count: number;
    total_value: string;
    percentage_of_total: string;
  }[];
  time_windows: {
    "7d": TimeWindowStats;
    "14d": TimeWindowStats;
    "21d": TimeWindowStats;
    "28d": TimeWindowStats;
  };
}

// --- Helper for table headers ---
const timeWindowLabels: { [key: string]: string } = {
  "7d": "Last 7 Days",
  "14d": "8-14 Days Ago",
  "21d": "15-21 Days Ago",
  "28d": "22-28 Days Ago",
};

export function DashboardPage() {
  const [totalDelinquencies, setTotalDelinquencies] = useState<number | null>(null);
  const [chargebackStats, setChargebackStats] = useState<ChargebackStats | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchDashboardData = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const [delinquenciesResponse, statsResponse] = await Promise.all([
          fetch(`http://10.98.1.142:8080/api/delinquencies?limit=1&page=1`),
          fetch(`http://10.98.1.142:8080/api/dashboard/chargeback-stats`),
        ]);

        if (!delinquenciesResponse.ok) throw new Error('Failed to fetch delinquencies');
        if (!statsResponse.ok) {
          const errorBody = await statsResponse.text();
          console.error("Chargeback stats API error body:", errorBody);
          throw new Error('Failed to fetch chargeback stats');
        }

        const delinquenciesData = await delinquenciesResponse.json();
        const statsData = await statsResponse.json();
        
        setTotalDelinquencies(delinquenciesData?.total_count ?? 0);
        setChargebackStats(statsData);

      } catch (err: any) {
        console.error("Failed to fetch dashboard data:", err);
        setError(err.message || "An unknown error occurred.");
      } finally {
        setIsLoading(false);
      }
    };

    fetchDashboardData();
  }, []);

  // --- Data Transformation functions for Rendering ---
  const getChargebackReportData = () => {
    if (!chargebackStats?.time_windows) return null;

    const windows = ["7d", "14d", "21d", "28d"];
    const headers = ["Metric", ...windows.map(w => timeWindowLabels[w])];
    const rows = [
      ["New Items", ...windows.map(w => chargebackStats.time_windows[w as keyof typeof timeWindowLabels].new_items_count.toLocaleString())],
      ["Value of New Items", ...windows.map(w => formatCurrency(parseFloat(chargebackStats.time_windows[w as keyof typeof timeWindowLabels].new_items_value)))],
      ["Passed to PFS", ...windows.map(w => chargebackStats.time_windows[w as keyof typeof timeWindowLabels].passed_to_pfs.toLocaleString())],
      ["Completed by PFS", ...windows.map(w => chargebackStats.time_windows[w as keyof typeof timeWindowLabels].completed_by_pfs.toLocaleString())],
      ["Avg Days to PFS", ...windows.map(w => `${chargebackStats.time_windows[w as keyof typeof timeWindowLabels].avg_days_to_pfs.toFixed(2)}`)],
      ["Avg PFS Completion", ...windows.map(w => `${chargebackStats.time_windows[w as keyof typeof timeWindowLabels].avg_days_for_pfs_complete.toFixed(2)}`)],
    ];
    return { headers, rows };
  };

  const getStatusSummaryData = () => {
    if (!chargebackStats?.status_summary) return null;

    const totalCount = chargebackStats.status_summary.reduce((sum, s) => sum + s.status_count, 0);
    const totalValue = chargebackStats.status_summary.reduce((sum, s) => sum + parseFloat(s.total_value), 0);

    const rows = chargebackStats.status_summary.map(s => [
        s.current_status,
        s.status_count.toLocaleString(),
        formatCurrency(parseFloat(s.total_value)),
        `${s.percentage_of_total}%`
    ]);

    rows.push([
        "Total",
        totalCount.toLocaleString(),
        formatCurrency(totalValue),
        "100%"
    ]);

    return {
        headers: ["Status", "Count", "Total Value", "% of Total"],
        rows: rows
    };
  };

  const chargebackReportData = getChargebackReportData();
  const statusSummaryData = getStatusSummaryData();

  if (isLoading) {
    return <div className="p-6">Loading dashboard...</div>;
  }

  if (error) {
    return <div className="p-6 text-destructive">Error: {error}</div>;
  }

  return (
    <div className="p-6">
      <h1 className="text-3xl font-bold">Dashboard</h1>
      <p className="text-muted-foreground">An overview of chargeback and delinquency metrics.</p>
      
      {/* NEW: Status Summary is now the first table */}
      {statusSummaryData && (
        <ReportingTable 
            title="Active Chargebacks by Status"
            headers={statusSummaryData.headers}
            rows={statusSummaryData.rows}
        />
      )}

      {/* Chargeback Trends Table */}
      {chargebackReportData && (
        <ReportingTable 
            title="Chargeback Trends"
            headers={chargebackReportData.headers}
            rows={chargebackReportData.rows}
        />
      )}

      {/* Delinquency Overview Section */}
      <h2 className="text-2xl font-semibold mt-8 mb-4">Delinquency Overview</h2>
       <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4 mt-4">
        <StatCard
          title="Total Active Delinquencies"
          value={totalDelinquencies !== null ? totalDelinquencies.toString() : "Loading..."}
          description="Total number of active delinquencies"
          icon={<Scale className="h-4 w-4 text-muted-foreground" />}
        />
      </div>
    </div>
  );
}
