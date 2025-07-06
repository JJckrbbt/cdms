import { useState, useEffect } from 'react';

interface StatusHistoryEntry {
  status_history_id: number;
  status: string;
  status_date: string;
  notes: string;
  user_id: number;
  user_first_name: string;
  user_last_name: string;
  user_email: string;
}

interface StatusHistoryProps {
  id: number;
  type: 'chargeback' | 'delinquency';
}

export function StatusHistory({ id, type }: StatusHistoryProps) {
  const [history, setHistory] = useState<StatusHistoryEntry[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    console.log(`StatusHistory: useEffect triggered for ${type} ID`, id);
    const fetchHistory = async () => {
      try {
        setIsLoading(true);
        const url = `http://10.98.1.142:8080/api/${type === 'delinquency' ? 'delinquencies' : type + 's'}/history/${id}`;
        console.log('StatusHistory: Fetching from URL:', url);
        const response = await fetch(url);
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        console.log('StatusHistory: Fetched data:', data);
        setHistory(data);
      } catch (error) {
        console.error('StatusHistory: Failed to fetch status history:', error);
      } finally {
        setIsLoading(false);
      }
    };

    if (id) {
      fetchHistory();
    }
  }, [id, type]);

  if (isLoading) {
    console.log('StatusHistory: Loading state...');
    return <p>Loading status history...</p>;
  }

  if (!history || history.length === 0) {
    console.log('StatusHistory: No history data or empty array.');
    return <p>No status history available.</p>;
  }

  console.log('StatusHistory: Rendering history:', history);
  return (
    <div className="space-y-4">
      {history.map((entry) => (
        <div key={entry.status_history_id} className="p-2 border rounded-md">
          <p><strong>Status:</strong> {entry.status}</p>
          <p><strong>Date:</strong> {new Date(entry.status_date).toLocaleString()}</p>
          <p><strong>User:</strong> {entry.user_first_name} {entry.user_last_name} ({entry.user_email})</p>
          <p><strong>Notes:</strong> {entry.notes}</p>
        </div>
      ))}
    </div>
  );
}
