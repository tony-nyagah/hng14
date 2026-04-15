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
import { UNIT_CATEGORIES, UnitCategory, convertUnit } from "@/constants/units";
import { useColorScheme } from "@/hooks/use-color-scheme";

export default function UnitConverterScreen() {
  const colorScheme = useColorScheme() ?? "light";
  const colors = Colors[colorScheme];

  const [activeCategoryIndex, setActiveCategoryIndex] = useState(0);
  const [inputValue, setInputValue] = useState("");
  const [fromUnitIndex, setFromUnitIndex] = useState(0);
  const [toUnitIndex, setToUnitIndex] = useState(1);

  const category: UnitCategory = UNIT_CATEGORIES[activeCategoryIndex];
  const fromUnit = category.units[fromUnitIndex];
  const toUnit = category.units[toUnitIndex];

  const result = useCallback((): string | null => {
    const num = parseFloat(inputValue);
    if (isNaN(num) || inputValue.trim() === "") return null;
    const converted = convertUnit(num, fromUnit, toUnit, category.name);
    const formatted =
      Math.abs(converted) >= 1000
        ? converted.toLocaleString("en-US", { maximumFractionDigits: 4 })
        : parseFloat(converted.toPrecision(6)).toString();
    return `${formatted} ${toUnit.symbol}`;
  }, [inputValue, fromUnit, toUnit, category.name]);

  const swapUnits = () => {
    const prev = fromUnitIndex;
    setFromUnitIndex(toUnitIndex);
    setToUnitIndex(prev);
  };

  const unitOptions = category.units.map((u) => ({
    label: `${u.label} (${u.symbol})`,
    value: u.label,
  }));

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
          <View style={[styles.iconCircle, { backgroundColor: Palette.indigo }]}>
            <MaterialIcons name="straighten" size={24} color="#fff" />
          </View>
          <Text style={[styles.title, { color: colors.text }]}>
            Unit Converter
          </Text>
        </View>

        {/* Category chips */}
        <ScrollView
          horizontal
          showsHorizontalScrollIndicator={false}
          contentContainerStyle={styles.chipsRow}
        >
          {UNIT_CATEGORIES.map((cat, idx) => (
            <TouchableOpacity
              key={cat.name}
              style={[
                styles.chip,
                {
                  backgroundColor:
                    idx === activeCategoryIndex ? Palette.indigo : colors.surface,
                  borderColor:
                    idx === activeCategoryIndex ? Palette.indigo : colors.border,
                },
              ]}
              onPress={() => {
                setActiveCategoryIndex(idx);
                setFromUnitIndex(0);
                setToUnitIndex(1);
                setInputValue("");
              }}
            >
              <Text
                style={[
                  styles.chipText,
                  { color: idx === activeCategoryIndex ? "#fff" : colors.text },
                ]}
              >
                {cat.name}
              </Text>
            </TouchableOpacity>
          ))}
        </ScrollView>

        {/* Input */}
        <SectionCard title="Enter Value">
          <NumberInput
            value={inputValue}
            onChangeText={setInputValue}
            placeholder="Type a number..."
            label={`Value in ${fromUnit.symbol}`}
          />
        </SectionCard>

        {/* From / To pickers */}
        <SectionCard title="Convert">
          <Text style={[styles.pickerLabel, { color: colors.subtext }]}>From</Text>
          <PickerSelect
            selectedValue={fromUnit.label}
            onValueChange={(val) => {
              const idx = category.units.findIndex((u) => u.label === val);
              if (idx >= 0) setFromUnitIndex(idx);
            }}
            options={unitOptions}
          />

          <TouchableOpacity
            style={[styles.swapBtn, { backgroundColor: Palette.indigo }]}
            onPress={swapUnits}
          >
            <MaterialIcons name="swap-vert" size={18} color="#fff" />
            <Text style={styles.swapText}>Swap</Text>
          </TouchableOpacity>

          <Text style={[styles.pickerLabel, { color: colors.subtext }]}>To</Text>
          <PickerSelect
            selectedValue={toUnit.label}
            onValueChange={(val) => {
              const idx = category.units.findIndex((u) => u.label === val);
              if (idx >= 0) setToUnitIndex(idx);
            }}
            options={unitOptions}
          />
        </SectionCard>

        {/* Result */}
        <ResultDisplay
          result={result()}
          label={`${inputValue || "0"} ${fromUnit.symbol} =`}
          accentBg={colorScheme === "dark" ? "#1A237E" : "#E8EAF6"}
        />
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
  chipsRow: { gap: 8, paddingBottom: 16, paddingRight: 4 },
  chip: {
    paddingHorizontal: 14,
    paddingVertical: 8,
    borderRadius: 20,
    borderWidth: 1,
  },
  chipText: { fontSize: 13, fontWeight: "600" },
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
});
