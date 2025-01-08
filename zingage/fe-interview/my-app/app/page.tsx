'use client';
import React, { useState } from 'react';

type CellValue = string | number;
type GridData = CellValue[][];
type CellFormula = { raw: string; value: number } | null;
interface CellState {
  value: CellValue;
  formula: CellFormula;
}

const SpreadsheetGrid: React.FC = () => {
  const [gridData, setGridData] = useState<CellState[][]>(
    Array(3).fill(null).map(() => 
      Array(3).fill(null).map(() => ({ value: '', formula: null }))
    )
  );
  const [editingCell, setEditingCell] = useState<string | null>(null);

  // type Result = [number | null, Error | null];
  const getCellReference = (ref: string): number => {
    // Convert cell reference (e.g., 'A1', 'B2') to value
    const colLetter = ref.charAt(0).toUpperCase();
    const row = parseInt(ref.slice(1)) - 1;
    const col = colLetter.charCodeAt(0) - 65; // 'A' = 0, 'B' = 1, etc.

    // Validate cell reference
    if (row < 0 || row >= gridData.length || col < 0 || col >= gridData[0].length) {
      throw new Error(`Invalid cell reference: ${ref}`);
    }

    const cellState = gridData[row][col];
    const cellValue = cellState.formula ? cellState.formula.value : Number(cellState.value) || undefined;

    if (cellValue == undefined) {
      throw new Error(`Invalid cell value: ${ref}`);
    }
    
    return cellValue;
  };

  
// First, let's define our Result type properly
type Result = [number | null, string | null];
const evaluateFormula = (formula: string): Result => {
    try {
        // Remove the '=' sign and all spaces
        const expression = formula.substring(1).trim();
        
        // Split the expression into tokens
        const tokens = expression.split(/([+\-*/()]|\s+)/).filter(token => token.trim());
        
        // Flag to track if we encountered any undefined cells
        let hasUndefinedCells = false;
        
        // Process each token
        const processedTokens = tokens.map(token => {
            token = token.trim();
            if (!token) return token;
            
            // Check if token is a cell reference (e.g., A1, B2, etc.)
            if (/^[A-Ca-c][1-3]$/.test(token)) {
                const cellValue = getCellReference(token);
                if (cellValue === null || cellValue === undefined) {
                    hasUndefinedCells = true;
                    return '0'; // Temporary placeholder to avoid syntax error
                }
                return cellValue.toString();
            }
            return token;
        });

        // If we found any undefined cells, return error
        if (hasUndefinedCells) {
            return [null, '#REF!'];
        }

        // Join tokens and evaluate
        const result = Function(`'use strict'; return (${processedTokens.join(' ')})`)();
        return [Number(result.toFixed(2)), null];
    } catch (error) {
        console.error('Formula evaluation error:', error);
        return [null, '#REF!'];
    }
};

  const handleCellChange = (rowIndex: number, colIndex: number, value: string) => {
    const newData = gridData.map((row, rIndex) =>
      rIndex === rowIndex
        ? row.map((cell, cIndex) => 
            cIndex === colIndex 
              ? { ...cell, value: value }
              : cell
          )
        : row
    );
    setGridData(newData);
  };

  const handleCellBlur = (rowIndex: number, colIndex: number, value: string) => {
    if (value.startsWith('=')) {
        const [result, error] = evaluateFormula(value);

        const newData = gridData.map((row, rIndex) =>
            rIndex === rowIndex
                ? row.map((cell, cIndex) => 
                    cIndex === colIndex 
                        ? { 
                            value: error ? error : result!,
                            formula: { 
                                raw: value, 
                                value: result ?? 0 
                            }
                        }
                        : cell
                )
                : row
        );
        setGridData(newData);
    }
    setEditingCell(null);
};

  const handleCellFocus = (rowIndex: number, colIndex: number) => {
    setEditingCell(`${rowIndex}-${colIndex}`);
  };

  const getCellDisplayValue = (cell: CellState, rowIndex: number, colIndex: number): string => {
    if (editingCell === `${rowIndex}-${colIndex}` && cell.formula) {
      return cell.formula.raw;
    }
    return cell.value.toString();
  };

  const getColumnLabel = (index: number): string => {
    return String.fromCharCode(65 + index); // A, B, C, etc.
  };

  return (
    <div className="w-full max-w-2xl mx-auto">
      {/* Column Headers */}
      <div className="grid grid-cols-3 gap-1 mb-1">
        {Array(3).fill(null).map((_, index) => (
          <div key={`header-${index}`} className="text-center text-sm font-bold text-gray-600">
            {getColumnLabel(index)}
          </div>
        ))}
      </div>
      
      {/* Grid */}
      <div className="grid grid-cols-3 gap-1 bg-gray-200 p-2 rounded-lg">
        {gridData.map((row, rowIndex) => (
          <React.Fragment key={`row-${rowIndex}`}>
            {row.map((cell, colIndex) => (
              <input
                key={`${rowIndex}-${colIndex}`}
                type="text"
                value={getCellDisplayValue(cell, rowIndex, colIndex)}
                onChange={(e) => handleCellChange(rowIndex, colIndex, e.target.value)}
                onBlur={(e) => handleCellBlur(rowIndex, colIndex, e.target.value)}
                onFocus={() => handleCellFocus(rowIndex, colIndex)}
                className="text-black w-full p-2 text-center border border-gray-300 focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none rounded"
                aria-label={`${getColumnLabel(colIndex)}${rowIndex + 1}`}
              />
            ))}
          </React.Fragment>
        ))}
      </div>
      
      {/* Instructions */}
      <div className="mt-4 text-sm text-gray-600">
        <p>Supported operations: +, -, *, /</p>
        <p>Example formulas:</p>
        <ul className="list-disc pl-5">
          <li>=A1+B1</li>
          <li>=B2*C2</li>
          <li>=A1+B2+C3</li>
        </ul>
      </div>
    </div>
  );
};

export default SpreadsheetGrid;