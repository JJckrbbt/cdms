import { DashboardLayout } from "./components/DashboardLayout";
import { Toaster } from "react-hot-toast";
import { Routes, Route } from "react-router-dom";
import { ChargebacksPage } from "./pages/ChargebacksPage";
import { DelinquenciesPage } from "./pages/DelinquenciesPage";
import { DashboardPage } from "./pages/DashboardPage";

function App() {
  const handleUploadSuccess = () => {
    // This function is called when an upload is successful.
    // The UploadReportModal handles closing itself and showing toasts.
    // Specific pages (like ChargebacksPage) will be responsible for refreshing their own data.
  };

  return (
    <DashboardLayout onUploadSuccess={handleUploadSuccess}>
      <Routes>
        <Route path="/" element={<DashboardPage />} />
        <Route path="/chargebacks" element={<ChargebacksPage />} />
        <Route path="/delinquencies" element={<DelinquenciesPage />} />
      </Routes>
      <Toaster />
    </DashboardLayout>
  );
}

export default App;
