import { DashboardLayout } from "./components/DashboardLayout";
import { StatCard } from "./components/StatCard";
import { DollarSign, Users, CreditCard, Activity } from "lucide-react";

function App() {
  return (
    <DashboardLayout>
      <div className="space-y-4">
        <h1 className="text-3xl font-bold">Dashboard</h1>
        
        {/* Stat Cards Grid */}
        <div className="grid gap-4 md:grid-cols-2 md:gap-8 lg:grid-cols-4">
          <StatCard 
            title="Total Revenue"
            value="$45,231.89"
            description="+20.1% from last month"
            icon={<DollarSign className="h-4 w-4 text-muted-foreground" />}
          />
          <StatCard 
            title="Active Users"
            value="+2350"
            description="+180.1% from last month"
            icon={<Users className="h-4 w-4 text-muted-foreground" />}
          />
          <StatCard 
            title="Open Chargebacks"
            value="12"
            description="+19% from last month"
            icon={<CreditCard className="h-4 w-4 text-muted-foreground" />}
          />
          <StatCard 
            title="New Delinquencies"
            value="3"
            description="+2 since last week"
            icon={<Activity className="h-4 w-4 text-muted-foreground" />}
          />
        </div>

        {/* Placeholder for future charts/tables */}
        <div className="mt-8">
          <h2 className="text-xl font-semibold">Recent Activity</h2>
          <div className="mt-4 h-96 rounded-lg border border-dashed p-4">
            <p className="text-center text-muted-foreground">
              Charts and data tables will be rendered here.
            </p>
          </div>
        </div>

      </div>
    </DashboardLayout>
  );
}

export default App;
