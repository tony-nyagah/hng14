import { Picker } from "@react-native-picker/picker";
import React from "react";
import { Platform, StyleSheet, View } from "react-native";

import { Colors } from "@/constants/theme";
import { useColorScheme } from "@/hooks/use-color-scheme";

type Option = {
  label: string;
  value: string;
};

type PickerSelectProps = {
  selectedValue: string;
  onValueChange: (value: string) => void;
  options: Option[];
};

export function PickerSelect({
  selectedValue,
  onValueChange,
  options,
}: PickerSelectProps) {
  const colorScheme = useColorScheme() ?? "light";
  const colors = Colors[colorScheme];

  return (
    <View
      style={[
        styles.wrapper,
        {
          backgroundColor: colors.surface,
          borderColor: colors.border,
        },
      ]}
    >
      <Picker
        selectedValue={selectedValue}
        onValueChange={onValueChange}
        style={[styles.picker, { color: colors.text }]}
        dropdownIconColor={colors.subtext}
        itemStyle={{ color: colors.text, fontSize: 14 }}
      >
        {options.map((opt) => (
          <Picker.Item
            key={opt.value}
            label={opt.label}
            value={opt.value}
            color={Platform.OS === "android" ? colors.text : undefined}
          />
        ))}
      </Picker>
    </View>
  );
}

const styles = StyleSheet.create({
  wrapper: {
    borderWidth: 1,
    borderRadius: 10,
    marginVertical: 6,
    overflow: "hidden",
  },
  picker: {
    height: Platform.OS === "ios" ? 150 : 54,
    width: "100%",
  },
});
