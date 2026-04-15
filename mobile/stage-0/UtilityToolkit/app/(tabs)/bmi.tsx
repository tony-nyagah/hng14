import MaterialIcons from "@expo/vector-icons/MaterialIcons";
import React, { useState } from "react";
import {
  ScrollView,
  StyleSheet,
  Text,
  TouchableOpacity,
  View,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

import { NumberInput } from "@/components/NumberInput";
import { SectionCard } from "@/components/SectionCard";
import { Colors, Palette } from "@/constants/theme";
import { useColorScheme } from "@/hooks/use-color-scheme";

type BMICategory = {
  label: string;
  range: string;
  bg: string;
  bgDark: string;
  icon: React.ComponentProps<typeof MaterialIcons>["name"];
};

const BMI_CATEGORIES: BMICategory[] = [
  { label: "Underweight", range: "< 18.5", bg: "#BBDEFB", bgDark: "#0D47A1", icon: "trending-down" },
  { label: "Normal Weight", range: "18.5–24.9", bg: "#C8E6C9", bgDark: "#1B5E20", icon: "check-circle" },
  { label: "Overweight", range: "25–29.9", bg: "#FFF9C4", bgDark: "#F57F17", icon: "warning" },
  { label: "Obese", range: "≥ 30", bg: "#FFCDD2", bgDark: "#B71C1C", icon: "error" },
];

function getBMICategory(bmi: number): BMICategory {
  if (bmi < 18.5) return BMI_CATEGORIES[0];
  if (bmi < 25) return BMI_CATEGORIES[1];
  if (bmi < 30) return BMI_CATEGORIES[2];
  return BMI_CATEGORIES[3];
}

type UnitSystem = "metric" | "imperial";

export default function BMIScreen() {
  const colorScheme = useColorScheme() ?? "light";
  const colors = Colors[colorScheme];

  const [unitSystem, setUnitSystem] = useState<UnitSystem>("metric");
  const [weight, setWeight] = useState("");
  const [height, setHeight] = useState("");
  const [heightFt, setHeightFt] = useState("");
  const [heightIn, setHeightIn] = useState("");

  const calculateBMI = (): number | null => {
    const w = parseFloat(weight);
    if (isNaN(w) || w <= 0) return null;
    if (unitSystem === "metric") {
      const h = parseFloat(height) / 100;
      if (isNaN(h) || h <= 0) return null;
      return w / (h * h);
    } else {
      const ft = parseFloat(heightFt) || 0;
      const inches = parseFloat(heightIn) || 0;
      const totalInches = ft * 12 + inches;
      if (totalInches <= 0) return null;
      return (w / (totalInches * totalInches)) * 703;
    }
  };

  const bmi = calculateBMI();
  const category = bmi !== null ? getBMICategory(bmi) : null;
  const weightLabel = unitSystem === "metric" ? "Weight (kg)" : "Weight (lbs)";

  const reset = () => {
    setWeight("");
    setHeight("");
    setHeightFt("");
    setHeightIn("");
  };

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
          <View style={[styles.iconCircle, { backgroundColor: Palette.orange }]}>
            <MaterialIcons name="fitness-center" size={24} color="#fff" />
          </View>
          <Text style={[styles.title, { color: colors.text }]}>
            BMI Calculator
          </Text>
        </View>

        {/* Unit Toggle */}
        <View
          style={[
            styles.toggleRow,
            { backgroundColor: colors.surface, borderColor: colors.border },
          ]}
        >
          {(["metric", "imperial"] as UnitSystem[]).map((sys) => (
            <TouchableOpacity
              key={sys}
              style={[
                styles.toggleBtn,
                {
                  backgroundColor:
                    unitSystem === sys ? Palette.orange : "transparent",
                },
              ]}
              onPress={() => {
                setUnitSystem(sys);
                reset();
              }}
            >
              <Text
                style={[
                  styles.toggleText,
                  {
                    color: unitSystem === sys ? "#fff" : colors.subtext,
                    fontWeight: unitSystem === sys ? "700" : "500",
                  },
                ]}
              >
                {sys.charAt(0).toUpperCase() + sys.slice(1)}
              </Text>
            </TouchableOpacity>
          ))}
        </View>

        {/* Inputs */}
        <SectionCard title="Your Measurements">
          <NumberInput
            value={weight}
            onChangeText={setWeight}
            label={weightLabel}
            placeholder="e.g. 70"
          />
          {unitSystem === "metric" ? (
            <NumberInput
              value={height}
              onChangeText={setHeight}
              label="Height (cm)"
              placeholder="e.g. 175"
            />
          ) : (
            <View style={styles.imperialRow}>
              <View style={styles.imperialField}>
                <NumberInput
                  value={heightFt}
                  onChangeText={setHeightFt}
                  label="Feet"
                  placeholder="5"
                />
              </View>
              <View style={styles.imperialField}>
                <NumberInput
                  value={heightIn}
                  onChangeText={setHeightIn}
                  label="Inches"
                  placeholder="10"
                />
              </View>
            </View>
          )}
        </SectionCard>

        {/* Result */}
        {bmi !== null && category !== null ? (
          <View
            style={[
              styles.resultCard,
              {
                backgroundColor: colorScheme === "dark" ? category.bgDark : category.bg,
                borderColor: colors.border,
                shadowColor: "#000",
              },
            ]}
          >
            <MaterialIcons
              name={category.icon}
              size={32}
              color={colorScheme === "dark" ? "#fff" : Palette.black}
            />
            <Text style={[styles.bmiValue, { color: colorScheme === "dark" ? "#fff" : Palette.black }]}>
              {bmi.toFixed(1)}
            </Text>
            <Text style={[styles.categoryLabel, { color: colorScheme === "dark" ? "#fff" : Palette.black }]}>
              {category.label}
            </Text>
            <Text style={[styles.categoryRange, { color: colorScheme === "dark" ? "#ddd" : "#555" }]}>
              BMI range: {category.range}
            </Text>
          </View>
        ) : null}

        {/* Reference Table */}
        <SectionCard title="BMI Reference">
          {BMI_CATEGORIES.map((cat) => (
            <View
              key={cat.label}
              style={[
                styles.refRow,
                {
                  backgroundColor:
                    category?.label === cat.label
                      ? colorScheme === "dark" ? cat.bgDark : cat.bg
                      : "transparent",
                },
              ]}
            >
              <View
                style={[
                  styles.refDot,
                  { backgroundColor: colorScheme === "dark" ? cat.bgDark : cat.bg },
                ]}
              />
              <Text style={[styles.refLabel, { color: colors.text }]}>{cat.label}</Text>
              <Text style={[styles.refRange, { color: colors.subtext }]}>{cat.range}</Text>
            </View>
          ))}
        </SectionCard>

        {/* Reset */}
        {weight || height || heightFt || heightIn ? (
          <TouchableOpacity
            style={[styles.resetBtn, { borderColor: colors.border }]}
            onPress={reset}
          >
            <MaterialIcons name="refresh" size={18} color={colors.subtext} />
            <Text style={[styles.resetText, { color: colors.subtext }]}>Reset</Text>
          </TouchableOpacity>
        ) : null}

        <Text style={[styles.disclaimer, { color: colors.subtext }]}>
          ⚠ BMI is a screening tool. Consult a healthcare professional for medical advice.
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
  toggleRow: {
    flexDirection: "row",
    borderRadius: 12,
    borderWidth: 1,
    padding: 4,
    marginBottom: 16,
    gap: 4,
  },
  toggleBtn: {
    flex: 1,
    paddingVertical: 10,
    borderRadius: 9,
    alignItems: "center",
  },
  toggleText: { fontSize: 14 },
  imperialRow: { flexDirection: "row", gap: 10 },
  imperialField: { flex: 1 },
  resultCard: {
    borderRadius: 16,
    borderWidth: 1,
    padding: 28,
    alignItems: "center",
    gap: 8,
    marginBottom: 16,
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 6,
    elevation: 3,
  },
  bmiValue: { fontSize: 60, fontWeight: "700", lineHeight: 66 },
  categoryLabel: { fontSize: 18, fontWeight: "600" },
  categoryRange: { fontSize: 13 },
  refRow: {
    flexDirection: "row",
    alignItems: "center",
    gap: 10,
    paddingVertical: 8,
    paddingHorizontal: 8,
    borderRadius: 8,
    marginBottom: 4,
  },
  refDot: { width: 12, height: 12, borderRadius: 6 },
  refLabel: { flex: 1, fontSize: 13, fontWeight: "500" },
  refRange: { fontSize: 12 },
  resetBtn: {
    flexDirection: "row",
    alignItems: "center",
    gap: 6,
    alignSelf: "center",
    paddingVertical: 8,
    paddingHorizontal: 20,
    borderRadius: 20,
    borderWidth: 1,
    marginBottom: 16,
  },
  resetText: { fontSize: 13, fontWeight: "500" },
  disclaimer: { fontSize: 11, textAlign: "center", lineHeight: 16, marginTop: 8 },
});
