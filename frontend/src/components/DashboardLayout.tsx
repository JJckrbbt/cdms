import { useState } from "react";
import { Sidebar } from "@/components/Sidebar";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Sheet, SheetContent, SheetTrigger, SheetHeader, SheetTitle, SheetDescription } from "@/components/ui/sheet";
import { Upload, PanelLeft, Sun, Moon } from "lucide-react";
import { UploadReportModal } from "./UploadReportModal";
import { Switch } from "@/components/ui/switch";
import { useTheme } from "../hooks/useTheme";

interface DashboardLayoutProps {
  children: React.ReactNode;
  onUploadSuccess: () => void;
}

export function DashboardLayout({ children, onUploadSuccess }: DashboardLayoutProps) {
  const [isUploadModalOpen, setIsUploadModalOpen] = useState(false);
  const { theme, toggleTheme } = useTheme();

  return (
    <div className="h-screen flex">
      {/* Persistent Sidebar for Desktop */}
      <div className="hidden md:block fixed top-0 left-0 h-full w-[280px] border-r bg-muted/40">
        <Sidebar onUploadReportClick={() => setIsUploadModalOpen(true)} />
      </div>

      <div className="flex flex-col flex-1 md:ml-[280px]">
        <header className="flex h-16 items-center justify-between border-b bg-background px-6 sticky top-0 z-10">
          {/* Mobile Navigation */}
          <Sheet>
            <SheetTrigger asChild>
              <Button size="icon" variant="outline" className="md:hidden">
                <PanelLeft className="h-5 w-5" />
                <span className="sr-only">Toggle Menu</span>
              </Button>
            </SheetTrigger>
            <SheetContent side="left" className="p-0 w-64">
              <Sidebar />
            </SheetContent>
          </Sheet>
          
          {/* Header Actions */}
          <div className="flex w-full items-center gap-4 justify-end">
            <div className="flex items-center space-x-2">
              {theme === "light" ? <Sun className="h-5 w-5" /> : <Moon className="h-5 w-5" />}
              <Switch
                checked={theme === "dark"}
                onCheckedChange={toggleTheme}
              />
            </div>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button
                  variant="secondary"
                  size="icon"
                  className="rounded-full"
                >
                  <Avatar>
                    <AvatarImage src="../assets/jjckrbbt.png" alt="@jjckrbbt" />
                    <AvatarFallback>JJ</AvatarFallback>
                  </Avatar>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuLabel>My Account</DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem>Settings</DropdownMenuItem>
                <DropdownMenuItem>Support</DropdownMenuItem>
                <DropdownMenuItem>Logout</DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
          </header>

        <main className="flex-1 flex flex-col overflow-hidden">
          <div className="p-6 h-full">
            {children}
          </div>
        </main>
        <Sheet open={isUploadModalOpen} onOpenChange={setIsUploadModalOpen}>
              <SheetContent side="bottom" className="w-1/3 ml-0 mr-auto">
                <SheetHeader>
                  <SheetTitle>Upload Report</SheetTitle>
                  <SheetDescription>
                    Select a report type and upload your file.
                  </SheetDescription>
                </SheetHeader>
                <UploadReportModal onClose={() => setIsUploadModalOpen(false)} onUploadSuccess={onUploadSuccess} />
              </SheetContent>
            </Sheet>
      </div>
      
    </div>
  );
}
