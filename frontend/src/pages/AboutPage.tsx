import React from 'react';

export const AboutPage: React.FC = () => {
  return (
    <div className="flex flex-col items-center justify-center h-full bg-background text-foreground">
      <h1 className="text-5xl font-bold">About CDMS</h1>
      <p className="mt-4 text-lg text-center max-w-2xl">
        The Chargeback & Delinquency Management System (CDMS) is a comprehensive solution for managing chargebacks and delinquencies. It provides a centralized platform for tracking, managing, and resolving disputes, helping organizations  streamline their resolution process and reduce losses.
      </p>
    </div>
  );
};
