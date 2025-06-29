import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import toast from "react-hot-toast";

interface UploadReportModalProps {
  onClose: () => void;
  onUploadSuccess: () => void;
}

const ALLOWED_REPORT_TYPES = [
  "BC1300",
  "BC1048",
  "OUTSTANDING_BILLS",
  "VENDOR_CODE",
];

export function UploadReportModal({ onClose, onUploadSuccess }: UploadReportModalProps) {
  const [selectedReportType, setSelectedReportType] = useState<string | null>(null);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [isUploading, setIsUploading] = useState(false);

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files && event.target.files[0]) {
      setSelectedFile(event.target.files[0]);
    }
  };

  const handleUpload = async () => {
    if (!selectedReportType) {
      toast.error("Please select a report type.");
      return;
    }
    if (!selectedFile) {
      toast.error("Please select a file to upload.");
      return;
    }

    setIsUploading(true);
    const formData = new FormData();
    formData.append("report_file", selectedFile);

    try {
      const response = await fetch(`http://10.98.1.142:8080/api/upload/${selectedReportType}`, {
        method: "POST",
        body: formData,
      });

      if (response.ok) {
        toast.success("Upload successful!");
        onUploadSuccess();
        onClose();
      } else {
        const errorData = await response.json();
        toast.error(`Upload failed: ${errorData.message || response.statusText}`);
      }
    } catch (error: any) {
      toast.error(`Network error: ${error.message}`);
    } finally {
      setIsUploading(false);
    }
  };

  return (
    <div className="grid gap-4 py-4">
      <div className="grid grid-cols-4 items-center gap-4">
        <Label htmlFor="reportType" className="text-right">
          Report Type
        </Label>
        <Select onValueChange={setSelectedReportType} value={selectedReportType || ""}>
          <SelectTrigger className="col-span-3">
            <SelectValue placeholder="Select a report type" />
          </SelectTrigger>
          <SelectContent>
            {ALLOWED_REPORT_TYPES.map((type) => (
              <SelectItem key={type} value={type}>
                {type}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
      <div className="grid grid-cols-4 items-center gap-4">
        <Label htmlFor="file" className="text-right">
          File
        </Label>
        <Input id="file" type="file" className="col-span-3" onChange={handleFileChange} />
      </div>
      <div className="flex justify-end">
        <Button onClick={handleUpload} disabled={isUploading || !selectedReportType || !selectedFile}>
          {isUploading ? "Uploading..." : "Upload"}
        </Button>
      </div>
    </div>
  );
}
