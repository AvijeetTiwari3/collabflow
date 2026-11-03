// Client-side Fractional Indexing for O(1) Drag and Drop Ordering

export const DEFAULT_INITIAL_POS = 1000.0;
export const DEFAULT_STEP = 1000.0;

/**
 * Calculates new fractional position between previous and next items
 */
export function calculatePosition(prevPos: number | null, nextPos: number | null): number {
  // Case 1: First item in empty list
  if (prevPos === null && nextPos === null) {
    return DEFAULT_INITIAL_POS;
  }

  // Case 2: Inserted at top of list
  if (prevPos === null && nextPos !== null) {
    if (nextPos <= 1.0) {
      return nextPos / 2.0;
    }
    return nextPos - DEFAULT_STEP;
  }

  // Case 3: Inserted at bottom of list
  if (prevPos !== null && nextPos === null) {
    return prevPos + DEFAULT_STEP;
  }

  // Case 4: Inserted between two items
  if (prevPos !== null && nextPos !== null) {
    return prevPos + (nextPos - prevPos) / 2.0;
  }

  return DEFAULT_INITIAL_POS;
}
