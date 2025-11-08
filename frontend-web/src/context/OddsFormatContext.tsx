import React, { createContext, useState, useContext, ReactNode } from 'react';
import { OddsFormat } from '../types';

interface OddsFormatContextType {
  oddsFormat: OddsFormat;
  setOddsFormat: (format: OddsFormat) => void;
}

const OddsFormatContext = createContext<OddsFormatContextType | undefined>(undefined);

export const OddsFormatProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [oddsFormat, setOddsFormatState] = useState<OddsFormat>(() => {
    const saved = localStorage.getItem('oddsFormat');
    return (saved as OddsFormat) || 'american';
  });

  const setOddsFormat = (format: OddsFormat) => {
    setOddsFormatState(format);
    localStorage.setItem('oddsFormat', format);
  };

  return (
    <OddsFormatContext.Provider value={{ oddsFormat, setOddsFormat }}>
      {children}
    </OddsFormatContext.Provider>
  );
};

export const useOddsFormat = (): OddsFormatContextType => {
  const context = useContext(OddsFormatContext);
  if (context === undefined) {
    throw new Error('useOddsFormat must be used within an OddsFormatProvider');
  }
  return context;
};
