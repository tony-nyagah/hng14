import React from "react";
import { StyleSheet, Text, View } from "react-native";

import { Colors, Palette } from "@/constants/theme";
import { useColorScheme } from "@/hooks/use-color-scheme";

type ResultDisplayProps = {
  result: string | null;
  label?: string;
  accentBg?: string;
};

export function ResultDisplay({ result, label, accentBg }: ResultDisplayProps) {
  const colorScheme = useColorScheme() ?? "light";
  const colors = Colors[colorScheme];

  if (!result) return null;

  const bg = accentBg ?? (colorScheme === "dark" ? "#1A237E" : "#E8EAF6");

  return (
    <View
      style={[
        styles.container,
        {
          backgroundColor: bg,
          borderColor: colors.border,
          shadowColor: Palette.indigo,
        },
      ]}
    >
      {label ? (
        <Text style={[styles.label, { color: colors.subtext }]}>{label}</Text>
      ) : null}
      <Text style={[styles.result, { color: colors.text }]}>{result}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    padding: 24,
    borderRadius: 14,
    borderWidth: 1,
    alignItems: "center",
    marginTop: 20,
    marginBottom: 8,
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 6,
    elevation: 3,
  },
  label: {
    fontSize: 12,
    fontWeight: "500",
    marginBottom: 8,
    opacity: 0.7,
  },
  result: {
    fontSize: 34,
    fontWeight: "700",
    textAlign: "center",
    lineHeight: 40,
  },
});
