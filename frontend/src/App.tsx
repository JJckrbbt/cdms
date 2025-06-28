import { DashboardLayout } from "./components/DashboardLayout";

function App() {
  return (
    <DashboardLayout>
      <div className="space-y-4">
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <p className="text-sm text-muted-foreground">
          An overview of your chargebacks and delinquencies.
        </p>
        {/* We will add our StatCards and DataTable here in the next steps */}
      </div>
    </DashboardLayout>
  );
}

export default App;