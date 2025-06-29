import { Button } from "@/components/ui/button";
import { Bell, Home, LineChart, Settings } from "lucide-react";
import AppLogo from "@/assets/CDMSTransparent.png";
import { Link } from "react-router-dom";

export function Sidebar() {
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

      {/* Spacer to push admin link to the bottom */}
      <div className="mt-auto" />

      {/* Admin Navigation */}
      <div className="pt-4 border-t">
        <nav className="flex flex-col gap-2">
          <Button variant="ghost" className="justify-start gap-2">
            <Settings className="h-4 w-4" />
            Admin
          </Button>
        </nav>
      </div>
    </div>
  );
}
