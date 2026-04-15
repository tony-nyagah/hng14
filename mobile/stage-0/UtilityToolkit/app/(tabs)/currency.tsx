import MaterialIcons from "@expo/vector-icons/MaterialIcons";
import React, { useState, useCallback } from "react";
import {
  ScrollView,
  StyleSheet,
  Text,
  TouchableOpacity,
  View,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

import { NumberInput } from "@/components/NumberInput";
import { PickerSelect } from "@/components/PickerSelect";
import { ResultDisplay } from "@/components/ResultDisplay";
import { SectionCard } from "@/components/SectionCard";
import { Colors, Palette } from "@/constants/theme";
import { CURRENCIES, convertCurrency } from "@/constants/units";
import { useColorScheme } from "@/hooks/use-color-scheme";

export default function CurrencyScreen() {
  const colorScheme = useColorScheme() ?? "light";
  const colors = Colors[colorScheme];

  const [amount, setAmount] = useState("");
  const [fromCode, setFromCode] = useState("USD");
  const [toCode, setToCode] = useState("NGN");

  const fromCurrency = CURRENCIES.find((c) => c.code === fromCode)!;
  const toCurrency = CURRENCIES.find((c) => c.code === toCode)!;

  const result = useCallback((): string | null => {
    const num = parseFloat(amount);
    if (isNaN(num) || amount.trim() === "") return null;
    const converted = convertCurrency(num, fromCurrency, toCurrency);
    return `${toCurrency.symbol} ${converted.toLocaleString("en-US", {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    })}`;
  }, [amount, fromCurrency, toCurrency]);

  const swapCurrencies = () => {
    const prev = fromCode;
    setFromCode(toCode);
    setToCode(prev);
  };

  const currencyOptions = CURRENCIES.map((c) => ({
    label: `${c.code} — ${c.name}`,
    value: c.code,
  }));

  const exchangeRate = (
    toCurrency.rateToUSD / fromCurrency.rateToUSD
  ).toLocaleString("en-US", {
    minimumFractionDigits: 4,
    maximumFractionDigits: 4,
  });

  return (
    <SafeAreaView
      style={[styles.safeArea, { backgroundColor: colors.background }]}
    >
      <ScrollView
        contentContainerStyle={styles.scroll}
        showsVerticalScrollIndicator={false}
        keyboardShouldPersistTaps="handled"
      >
        {/* Header */}
        <View style={styles.header}>
          <View style={[styles.iconCircle, { backgroundColor: Palette.teal }]}>
            <MaterialIcons name="currency-exchange" size={24} color="#fff" />
          </View>
          <Text style={[styles.title, { color: colors.text }]}>
            Currency Converter
          </Text>
        </View>

        {/* Rate badge */}
        <View
          style={[
            styles.rateBadge,
            { backgroundColor: colors.surface, borderColor: colors.border },
          ]}
        >
          <MaterialIcons name="info-outline" size={14} color={colors.subtext} />
          <Text style={[styles.rateText, { color: colors.subtext }]}>
            1 {fromCurrency.code} = {exchangeRate} {toCurrency.code} · indicative
          </Text>
        </View>

        {/* Amount Input */}
        <SectionCard title="Amount">
          <NumberInput
            value={amount}
            onChangeText={setAmount}
            placeholder="Enter amount..."
            label={`${fromCurrency.name} (${fromCurrency.symbol})`}
          />
        </SectionCard>

        {/* Currency selectors */}
        <SectionCard title="Currencies">
          <Text style={[styles.pickerLabel, { color: colors.subtext }]}>From</Text>
          <PickerSelect
            selectedValue={fromCode}
            onValueChange={setFromCode}
            options={currencyOptions}
          />

          <TouchableOpacity
            style={[styles.swapBtn, { backgroundColor: Palette.teal }]}
            onPress={swapCurrencies}
          >
            <MaterialIcons name="swap-vert" size={18} color="#fff" />
            <Text style={styles.swapText}>Swap</Text>
          </TouchableOpacity>

          <Text style={[styles.pickerLabel, { color: colors.subtext }]}>To</Text>
          <PickerSelect
            selectedValue={toCode}
            onValueChange={setToCode}
            options={currencyOptions}
          />
        </SectionCard>

        {/* Result */}
        <ResultDisplay
          result={result()}
          label={`${amount || "0"} ${fromCurrency.code} =`}
          accentBg={colorScheme === "dark" ? "#004D40" : "#E0F2F1"}
        />

        <Text style={[styles.disclaimer, { color: colors.subtext }]}>
          ⚠ Rates are hardcoded approximations (Apr 2026).
        </Text>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: { flex: 1 },
  scroll: { paddingHorizontal: 16, paddingBottom: 40 },
  header: {
    flexDirection: "row",
    alignItems: "center",
    gap: 12,
    paddingTop: 24,
    paddingBottom: 20,
  },
  iconCircle: {
    width: 42,
    height: 42,
    borderRadius: 21,
    alignItems: "center",
    justifyContent: "center",
  },
  title: { fontSize: 22, fontWeight: "700" },
  rateBadge: {
    flexDirection: "row",
    alignItems: "center",
    gap: 6,
    padding: 10,
    borderRadius: 10,
    borderWidth: 1,
    marginBottom: 16,
  },
  rateText: { fontSize: 12, fontWeight: "500", flex: 1 },
  pickerLabel: { fontSize: 12, fontWeight: "600", marginTop: 4, marginBottom: 2 },
  swapBtn: {
    flexDirection: "row",
    alignItems: "center",
    gap: 6,
    alignSelf: "center",
    paddingVertical: 8,
    paddingHorizontal: 20,
    marginVertical: 10,
    borderRadius: 20,
  },
  swapText: { fontSize: 13, fontWeight: "600", color: "#fff" },
  disclaimer: {
    fontSize: 11,
    fontWeight: "500",
    textAlign: "center",
    marginTop: 16,
    lineHeight: 16,
  },
});
