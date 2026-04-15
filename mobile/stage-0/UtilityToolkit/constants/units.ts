/**
 * All unit conversion definitions for the Smart Utility Toolkit.
 * Conversions are defined relative to a base unit.
 */

export type Unit = {
  label: string;
  symbol: string;
  toBase: number; // multiply by this to get the base unit
};

export type UnitCategory = {
  name: string;
  baseUnit: string;
  units: Unit[];
};

export const UNIT_CATEGORIES: UnitCategory[] = [
  {
    name: "Length",
    baseUnit: "Meter",
    units: [
      { label: "Millimeter", symbol: "mm", toBase: 0.001 },
      { label: "Centimeter", symbol: "cm", toBase: 0.01 },
      { label: "Meter", symbol: "m", toBase: 1 },
      { label: "Kilometer", symbol: "km", toBase: 1000 },
      { label: "Inch", symbol: "in", toBase: 0.0254 },
      { label: "Foot", symbol: "ft", toBase: 0.3048 },
      { label: "Yard", symbol: "yd", toBase: 0.9144 },
      { label: "Mile", symbol: "mi", toBase: 1609.344 },
    ],
  },
  {
    name: "Weight",
    baseUnit: "Kilogram",
    units: [
      { label: "Milligram", symbol: "mg", toBase: 0.000001 },
      { label: "Gram", symbol: "g", toBase: 0.001 },
      { label: "Kilogram", symbol: "kg", toBase: 1 },
      { label: "Tonne", symbol: "t", toBase: 1000 },
      { label: "Ounce", symbol: "oz", toBase: 0.0283495 },
      { label: "Pound", symbol: "lb", toBase: 0.453592 },
      { label: "Stone", symbol: "st", toBase: 6.35029 },
    ],
  },
  {
    name: "Temperature",
    baseUnit: "Celsius",
    units: [
      { label: "Celsius", symbol: "°C", toBase: 1 },
      { label: "Fahrenheit", symbol: "°F", toBase: 1 },
      { label: "Kelvin", symbol: "K", toBase: 1 },
    ],
  },
  {
    name: "Area",
    baseUnit: "Square Meter",
    units: [
      { label: "Square Millimeter", symbol: "mm²", toBase: 0.000001 },
      { label: "Square Centimeter", symbol: "cm²", toBase: 0.0001 },
      { label: "Square Meter", symbol: "m²", toBase: 1 },
      { label: "Hectare", symbol: "ha", toBase: 10000 },
      { label: "Square Kilometer", symbol: "km²", toBase: 1000000 },
      { label: "Square Inch", symbol: "in²", toBase: 0.00064516 },
      { label: "Square Foot", symbol: "ft²", toBase: 0.092903 },
      { label: "Acre", symbol: "ac", toBase: 4046.86 },
    ],
  },
  {
    name: "Speed",
    baseUnit: "Meters per Second",
    units: [
      { label: "Meters per Second", symbol: "m/s", toBase: 1 },
      { label: "Kilometers per Hour", symbol: "km/h", toBase: 0.277778 },
      { label: "Miles per Hour", symbol: "mph", toBase: 0.44704 },
      { label: "Knot", symbol: "kn", toBase: 0.514444 },
      { label: "Feet per Second", symbol: "ft/s", toBase: 0.3048 },
    ],
  },
];

/**
 * Convert a value between two units. Handles temperature specially.
 */
export function convertUnit(
  value: number,
  fromUnit: Unit,
  toUnit: Unit,
  category: string,
): number {
  if (category === "Temperature") {
    return convertTemperature(value, fromUnit.label, toUnit.label);
  }
  // Generic: convert to base, then to target
  const inBase = value * fromUnit.toBase;
  return inBase / toUnit.toBase;
}

function convertTemperature(value: number, from: string, to: string): number {
  if (from === to) return value;
  // First convert to Celsius
  let celsius: number;
  switch (from) {
    case "Fahrenheit":
      celsius = (value - 32) * (5 / 9);
      break;
    case "Kelvin":
      celsius = value - 273.15;
      break;
    default:
      celsius = value; // already Celsius
  }
  // Then from Celsius to target
  switch (to) {
    case "Fahrenheit":
      return celsius * (9 / 5) + 32;
    case "Kelvin":
      return celsius + 273.15;
    default:
      return celsius;
  }
}

/**
 * Currency data (hardcoded rates relative to USD)
 * Rates as of April 2026 (approximate)
 */
export type Currency = {
  code: string;
  name: string;
  symbol: string;
  rateToUSD: number; // how many of this currency = 1 USD
};

export const CURRENCIES: Currency[] = [
  { code: "USD", name: "US Dollar", symbol: "$", rateToUSD: 1 },
  { code: "EUR", name: "Euro", symbol: "€", rateToUSD: 0.923 },
  { code: "GBP", name: "British Pound", symbol: "£", rateToUSD: 0.791 },
  { code: "JPY", name: "Japanese Yen", symbol: "¥", rateToUSD: 152.5 },
  { code: "CAD", name: "Canadian Dollar", symbol: "C$", rateToUSD: 1.362 },
  { code: "AUD", name: "Australian Dollar", symbol: "A$", rateToUSD: 1.541 },
  { code: "CHF", name: "Swiss Franc", symbol: "Fr", rateToUSD: 0.901 },
  { code: "CNY", name: "Chinese Yuan", symbol: "¥", rateToUSD: 7.25 },
  { code: "INR", name: "Indian Rupee", symbol: "₹", rateToUSD: 83.4 },
  { code: "NGN", name: "Nigerian Naira", symbol: "₦", rateToUSD: 1580 },
  { code: "ZAR", name: "South African Rand", symbol: "R", rateToUSD: 18.63 },
  { code: "KES", name: "Kenyan Shilling", symbol: "KSh", rateToUSD: 129.5 },
  { code: "GHS", name: "Ghanaian Cedi", symbol: "GH₵", rateToUSD: 15.2 },
  { code: "BRL", name: "Brazilian Real", symbol: "R$", rateToUSD: 5.08 },
  { code: "MXN", name: "Mexican Peso", symbol: "MX$", rateToUSD: 17.15 },
  { code: "AED", name: "UAE Dirham", symbol: "د.إ", rateToUSD: 3.673 },
  { code: "SAR", name: "Saudi Riyal", symbol: "﷼", rateToUSD: 3.75 },
  { code: "SGD", name: "Singapore Dollar", symbol: "S$", rateToUSD: 1.345 },
  { code: "HKD", name: "Hong Kong Dollar", symbol: "HK$", rateToUSD: 7.82 },
  { code: "NOK", name: "Norwegian Krone", symbol: "kr", rateToUSD: 10.62 },
];

export function convertCurrency(
  amount: number,
  from: Currency,
  to: Currency,
): number {
  // Convert to USD first, then to target
  const inUSD = amount / from.rateToUSD;
  return inUSD * to.rateToUSD;
}
