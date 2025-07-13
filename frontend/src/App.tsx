import { DashboardLayout } from "./components/DashboardLayout";
import { Toaster } from "react-hot-toast";
import { Routes, Route, Outlet } from "react-router-dom";
import { ChargebacksPage } from "./pages/ChargebacksPage";
import { DelinquenciesPage } from "./pages/DelinquenciesPage";
import { DashboardPage } from "./pages/DashboardPage";
import { LandingPage } from "./pages/LandingPage";
import { AboutPage } from './pages/AboutPage';

function App() {
  const handleUploadSuccess = () => {
    // This function is called when an upload is successful.
    // The UploadReportModal handles closing itself and showing toasts.
    // Specific pages (like ChargebacksPage) will be responsible for refreshing their own data.
  };

  const AppLayout = () => (
    <DashboardLayout onUploadSuccess={handleUploadSuccess}>
      <Outlet />
    </DashboardLayout>
  );

  return (
    <>
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route element={<AppLayout />}>
          <Route path="/dashboard" element={<DashboardPage />} />
          <Route path="/chargebacks" element={<ChargebacksPage />} />
          <Route path="/delinquencies" element={<DelinquenciesPage />} />
          <Route path="/about" element={<AboutPage />} />
        </Route>
      </Routes>
      <Toaster />
    </>
  );
}

export default App;