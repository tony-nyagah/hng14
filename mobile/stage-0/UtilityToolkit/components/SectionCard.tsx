import React from "react";
import { StyleSheet, Text, View, ViewStyle } from "react-native";

import { Colors } from "@/constants/theme";
import { useColorScheme } from "@/hooks/use-color-scheme";

type SectionCardProps = {
  title?: string;
  children: React.ReactNode;
  style?: ViewStyle;
  accentColor?: string;
};

export function SectionCard({
  title,
  children,
  style,
  accentColor: _accentColor,
}: SectionCardProps) {
  const colorScheme = useColorScheme() ?? "light";
  const colors = Colors[colorScheme];

  return (
    <View
      style={[
        styles.card,
        {
          backgroundColor: colors.surface,
          borderColor: colors.border,
          shadowColor: "#000",
        },
        style,
      ]}
    >
      {title ? (
        <Text style={[styles.titleText, { color: colors.subtext }]}>
          {title}
        </Text>
      ) : null}
      <View style={styles.content}>{children}</View>
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    borderRadius: 14,
    borderWidth: 1,
    marginBottom: 16,
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.08,
    shadowRadius: 6,
    elevation: 2,
    overflow: "hidden",
  },
  titleText: {
    fontSize: 12,
    fontWeight: "600",
    paddingHorizontal: 16,
    paddingTop: 14,
    paddingBottom: 4,
    textTransform: "uppercase",
    letterSpacing: 0.8,
  },
  content: {
    padding: 16,
  },
});
