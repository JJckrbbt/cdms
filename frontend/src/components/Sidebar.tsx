import { Button } from "@/components/ui/button";
import { Bell, Home, LineChart, Settings, Upload, Gauge, BadgeInfo } from "lucide-react";
import AppLogo from "@/assets/CDMSTransparent.png";
import { Link } from "react-router-dom";

interface SidebarProps {
  onUploadReportClick?: () => void;
}

export function Sidebar({ onUploadReportClick }: SidebarProps) {
  return (
    <div className="flex h-full flex-col p-4">
      <div className="mb-8 flex items-center gap-2">
        {/* Use the imported logo */}
        <img src={AppLogo} alt="Application Logo" className="w-12 h-12" />
        <h2 className="text-xl font-bold">CDMS</h2>
      </div>

      {/* Main Navigation */}
      <nav className="flex flex-col gap-2">
        <Link to="/" className="w-full">
          <Button variant="ghost" className="justify-start gap-2 w-full">
            <Home className="h-4 w-4" />
            Home
          </Button>
        </Link>
        <Link to="/dashboard" className="w-full">
          <Button variant="ghost" className="justify-start gap-2 w-full">
            <Gauge className="h-4 w-4" />
            Dashboard
          </Button>
        </Link>
        <Link to="/chargebacks" className="w-full">
          <Button variant="ghost" className="justify-start gap-2 w-full">
            <LineChart className="h-4 w-4" />
            Chargebacks
          </Button>
        </Link>
        <Link to="/delinquencies" className="w-full">
          <Button variant="ghost" className="justify-start gap-2 w-full">
            <Bell className="h-4 w-4" />
            Delinquencies
          </Button>
        </Link>
      </nav>

      {/* Spacer to push admin link and upload button to the bottom */}
      <div className="mt-auto" />

      {/* Admin and Upload Navigation */}
      <div className="pt-4 border-t">
        <nav className="flex flex-col gap-2">
          <Button variant="ghost" className="justify-start gap-2" onClick={onUploadReportClick}>
            <Upload className="h-4 w-4" />
            Upload Report
          </Button>
          <Link to="/uploads" className="w-full">
            <Button variant="ghost" className="justify-start gap-2 w-full">
              <BadgeInfo className="h-4 w-4" />
              Upload Reporting
            </Button>
          </Link>
          <Button variant="ghost" className="justify-start gap-2">
            <Settings className="h-4 w-4" />
            Admin
          </Button>
        </nav>
      </div>
    </div>
  );
}