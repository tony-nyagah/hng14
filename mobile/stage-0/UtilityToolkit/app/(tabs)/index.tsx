import MaterialIcons from "@expo/vector-icons/MaterialIcons";
import { useRouter } from "expo-router";
import React from "react";
import {
  ScrollView,
  StyleSheet,
  Text,
  TouchableOpacity,
  View,
} from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

import { Colors, Palette } from "@/constants/theme";
import { useColorScheme } from "@/hooks/use-color-scheme";

type Tool = {
  id: string;
  title: string;
  description: string;
  icon: React.ComponentProps<typeof MaterialIcons>["name"];
  accentBg: string;
  accentBgDark: string;
  route: string;
};

const TOOLS: Tool[] = [
  {
    id: "unit",
    title: "Unit Converter",
    description: "Length · Weight · Temp · Area · Speed",
    icon: "straighten",
    accentBg: "#E8EAF6",
    accentBgDark: "#1A237E",
    route: "/unit-converter",
  },
  {
    id: "currency",
    title: "Currency Converter",
    description: "20+ world currencies",
    icon: "currency-exchange",
    accentBg: "#E0F2F1",
    accentBgDark: "#004D40",
    route: "/currency",
  },
  {
    id: "bmi",
    title: "BMI Calculator",
    description: "Metric & imperial · Health category",
    icon: "fitness-center",
    accentBg: "#FBE9E7",
    accentBgDark: "#4E342E",
    route: "/bmi",
  },
];

export default function HomeScreen() {
  const colorScheme = useColorScheme() ?? "light";
  const colors = Colors[colorScheme];
  const router = useRouter();

  return (
    <SafeAreaView
      style={[styles.safeArea, { backgroundColor: colors.background }]}
    >
      <ScrollView
        contentContainerStyle={styles.scroll}
        showsVerticalScrollIndicator={false}
      >
        {/* Header */}
        <View style={styles.header}>
          <View style={[styles.iconCircle, { backgroundColor: Palette.indigo }]}>
            <MaterialIcons name="build" size={28} color="#fff" />
          </View>
          <Text style={[styles.appName, { color: colors.text }]}>
            Smart Utility Toolkit
          </Text>
          <Text style={[styles.subtitle, { color: colors.subtext }]}>
            Your everyday calculation companion
          </Text>
        </View>

        {/* Tool Cards */}
        <Text style={[styles.sectionLabel, { color: colors.subtext }]}>TOOLS</Text>
        {TOOLS.map((tool) => (
          <TouchableOpacity
            key={tool.id}
            style={[
              styles.card,
              {
                backgroundColor:
                  colorScheme === "dark" ? tool.accentBgDark : tool.accentBg,
                borderColor: colors.border,
                shadowColor: "#000",
              },
            ]}
            onPress={() => router.push(tool.route as any)}
            activeOpacity={0.75}
          >
            <View style={[styles.iconBox, { backgroundColor: Palette.indigo }]}>
              <MaterialIcons name={tool.icon} size={24} color="#fff" />
            </View>
            <View style={styles.cardText}>
              <Text style={[styles.cardTitle, { color: colors.text }]}>
                {tool.title}
              </Text>
              <Text style={[styles.cardDesc, { color: colors.subtext }]}>
                {tool.description}
              </Text>
            </View>
            <MaterialIcons name="chevron-right" size={22} color={colors.subtext} />
          </TouchableOpacity>
        ))}

        <Text style={[styles.footer, { color: colors.subtext }]}>
          3 tools · v1.0.0
        </Text>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: { flex: 1 },
  scroll: { paddingHorizontal: 20, paddingBottom: 40 },
  header: {
    alignItems: "center",
    paddingTop: 32,
    paddingBottom: 28,
  },
  iconCircle: {
    width: 64,
    height: 64,
    borderRadius: 32,
    alignItems: "center",
    justifyContent: "center",
    marginBottom: 16,
  },
  appName: {
    fontSize: 26,
    fontWeight: "700",
    textAlign: "center",
    marginBottom: 6,
  },
  subtitle: { fontSize: 14, textAlign: "center" },
  sectionLabel: {
    fontSize: 11,
    fontWeight: "700",
    letterSpacing: 1.5,
    marginBottom: 12,
  },
  card: {
    flexDirection: "row",
    alignItems: "center",
    borderRadius: 14,
    borderWidth: 1,
    padding: 16,
    marginBottom: 14,
    gap: 14,
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.08,
    shadowRadius: 6,
    elevation: 2,
  },
  iconBox: {
    width: 46,
    height: 46,
    borderRadius: 12,
    alignItems: "center",
    justifyContent: "center",
  },
  cardText: { flex: 1 },
  cardTitle: { fontSize: 16, fontWeight: "600", marginBottom: 2 },
  cardDesc: { fontSize: 12 },
  footer: { fontSize: 12, textAlign: "center", marginTop: 16 },
});
