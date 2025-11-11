import React, { ReactNode } from 'react';

interface Column {
  key: string;
  label: string;
  render?: (value: any, row: any) => ReactNode;
  className?: string;
  mobileLabel?: string; // Optional custom label for mobile
}

interface ResponsiveTableProps {
  columns: Column[];
  data: any[];
  keyField: string;
  emptyMessage?: string;
  onRowClick?: (row: any) => void;
  className?: string;
}

const ResponsiveTable: React.FC<ResponsiveTableProps> = ({
  columns,
  data,
  keyField,
  emptyMessage = 'No data available',
  onRowClick,
  className = '',
}) => {
  if (data.length === 0) {
    return (
      <div className={`bg-white dark:bg-dark-800 rounded-lg shadow p-12 text-center ${className}`}>
        <p className="text-gray-500 dark:text-gray-400">{emptyMessage}</p>
      </div>
    );
  }

  return (
    <>
      {/* Desktop Table */}
      <div className={`hidden md:block bg-white dark:bg-dark-800 rounded-lg shadow overflow-hidden ${className}`}>
        <table className="min-w-full divide-y divide-gray-200 dark:divide-dark-700">
          <thead className="bg-gray-50 dark:bg-dark-900">
            <tr>
              {columns.map((column) => (
                <th
                  key={column.key}
                  className={`px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider ${column.className || ''}`}
                >
                  {column.label}
                </th>
              ))}
            </tr>
          </thead>
          <tbody className="bg-white dark:bg-dark-800 divide-y divide-gray-200 dark:divide-dark-700">
            {data.map((row) => (
              <tr
                key={row[keyField]}
                className={`${onRowClick ? 'cursor-pointer hover:bg-gray-50 dark:hover:bg-dark-700 transition-colors' : ''}`}
                onClick={() => onRowClick && onRowClick(row)}
              >
                {columns.map((column) => (
                  <td
                    key={column.key}
                    className={`px-6 py-4 whitespace-nowrap text-sm ${column.className || ''}`}
                  >
                    {column.render ? column.render(row[column.key], row) : row[column.key]}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Mobile Cards */}
      <div className="md:hidden space-y-4">
        {data.map((row) => (
          <div
            key={row[keyField]}
            className={`bg-white dark:bg-dark-800 rounded-lg shadow p-4 ${onRowClick ? 'cursor-pointer active:scale-95 transition-transform' : ''}`}
            onClick={() => onRowClick && onRowClick(row)}
          >
            {columns.map((column) => (
              <div key={column.key} className="flex justify-between items-start py-2 border-b border-gray-100 dark:border-dark-700 last:border-0">
                <span className="text-sm font-medium text-gray-500 dark:text-gray-400">
                  {column.mobileLabel || column.label}
                </span>
                <span className={`text-sm text-gray-900 dark:text-white text-right ml-4 ${column.className || ''}`}>
                  {column.render ? column.render(row[column.key], row) : row[column.key]}
                </span>
              </div>
            ))}
          </div>
        ))}
      </div>
    </>
  );
};

export default ResponsiveTable;
