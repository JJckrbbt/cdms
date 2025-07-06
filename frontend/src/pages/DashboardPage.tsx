import { useState, useEffect } from "react";
import { StatCard } from "@/components/StatCard";
// FIXED: Removed the unused 'TrendingUp', 'TrendingDown', and 'Hourglass' imports
import { Scale, ReceiptText } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { formatCurrency } from "@/lib/utils";

// --- Data Structures from API ---
interface StatusSummary {
  current_status: string;
  status_count: string;
  total_value: string;
  percentage_of_total: string;
}

interface TimeWindowStats {
  new_items_count: number;
  new_items_value: string;
  avg_days_to_pfs: number;
  avg_days_for_pfs_complete: number;
  passed_to_pfs: number;
  completed_by_pfs: number;
}

interface CombinedChargebackStats {
  status_summary: StatusSummary[];
  time_windows: {
    "7d": TimeWindowStats;
    "14d": TimeWindowStats;
    "21d": TimeWindowStats;
    "28d": TimeWindowStats;
  };
}

type TimeWindowKey = keyof CombinedChargebackStats['time_windows'];

// --- Component ---
export function DashboardPage() {
  const [totalDelinquencies, setTotalDelinquencies] = useState<number | null>(null);
  const [totalChargebacks, setTotalChargebacks] = useState<number | null>(null);
  const [totalChargebacksValue, setTotalChargebacksValue] = useState<number | null>(null);
  const [totalDelinquenciesValue, setTotalDelinquenciesValue] = useState<number | null>(null);
  const [chargebackStats, setChargebackStats] = useState<CombinedChargebackStats | null>(null);

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

    const fetchChargebackStats = async () => {
      try {
        const response = await fetch(`http://10.98.1.142:8080/api/dashboard/chargeback-stats`);
        const data = await response.json();
        setChargebackStats(data);
      } catch (error) {
        console.error("Failed to fetch chargeback stats:", error);
      }
    };

    fetchTotals();
    fetchChargebackStats();
  }, []);

  if (!chargebackStats) {
    return <p>Loading dashboard...</p>;
  }

  const orderedStatuses = [
    'Open',
    'In Research',
    'Hold Pending Internal Action',
    'Hold Pending External Action',
    'Passed to PFS',
    'PFS Return to GSA',
    'Completed by PFS'
  ];

  const orderedStatusSummary = orderedStatuses.map(status => 
    chargebackStats.status_summary.find(item => item.current_status === status)
  ).filter((item): item is StatusSummary => item !== undefined);

  const windows = ["7d", "14d", "21d", "28d"];
  const timeWindowLabels: Record<string, string> = { "7d": "Last 7 Days", "14d": "Last 14 Days", "21d": "Last 21 Days", "28d": "Last 28 Days" };

  const tableData = [
    ["New Items", ...windows.map(w => chargebackStats.time_windows[w as TimeWindowKey].new_items_count.toLocaleString())],
    ["Value of New Items", ...windows.map(w => formatCurrency(parseFloat(chargebackStats.time_windows[w as TimeWindowKey].new_items_value)))],
    ["Passed to PFS", ...windows.map(w => chargebackStats.time_windows[w as TimeWindowKey].passed_to_pfs.toLocaleString())],
    ["Completed by PFS", ...windows.map(w => chargebackStats.time_windows[w as TimeWindowKey].completed_by_pfs.toLocaleString())],
    ["Avg Days to PFS", ...windows.map(w => `${chargebackStats.time_windows[w as TimeWindowKey].avg_days_to_pfs.toFixed(2)}`)],
    ["Avg PFS Completion", ...windows.map(w => `${chargebackStats.time_windows[w as TimeWindowKey].avg_days_for_pfs_complete.toFixed(2)}`)],
  ];

  return (
    <div className="p-6 space-y-6">
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
       <StatCard
          title="Total Chargebacks"
          value={totalChargebacks !== null ? totalChargebacks.toString() : "Loading..."}
          icon={<ReceiptText className="h-4 w-4 text-muted-foreground" />}
          description="Total number of active chargebacks"
        />
        <StatCard
          title="Total Delinquencies"
          value={totalDelinquencies !== null ? totalDelinquencies.toString() : "Loading..."}
          icon={<Scale className="h-4 w-4 text-muted-foreground" />}
          description="Total number of active delinquencies"
        />
     </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
         <Card>
          <CardHeader>
            <CardTitle>Chargebacks Status Summary</CardTitle>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Count</TableHead>
                  <TableHead className="text-right">Value</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {orderedStatusSummary.map(item => (
                  <TableRow key={item.current_status}>
                    <TableCell>{item.current_status}</TableCell>
                    <TableCell className="text-right">{parseInt(item.status_count).toLocaleString()}</TableCell>
                    <TableCell className="text-right">{formatCurrency(parseFloat(item.total_value))}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>Chargeback Trends</CardTitle>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Metric</TableHead>
                  {windows.map(w => <TableHead key={w} className="text-right">{timeWindowLabels[w]}</TableHead>)}
                </TableRow>
              </TableHeader>
              <TableBody>
                {tableData.map((row, i) => (
                  <TableRow key={i}>
                    <TableCell className="font-medium">{row[0]}</TableCell>
                    {row.slice(1).map((cell, j) => <TableCell key={j} className="text-right">{String(cell)}</TableCell>)}
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>

     </div>
    </div>
  );
}
